package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type Hangman struct {
	Entries     map[string]bool
	Placeholder []string
	Chances     int
}

func getEntries(entries map[string]bool) (keys []string) {
	for k, _ := range entries {
		keys = append(keys, k)
	}

	return
}

func play(h Hangman, word string, result chan<- string) {
	for {

		correctEntries := strings.Join(h.Placeholder, "")

		// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.

		if h.Chances == 0 && correctEntries != word {
			result <- `You lose
			Correct word` + word
			break
		}

		// evaluate a win!
		if correctEntries == word {
			result <- "You guessed it right:" + word
			break
		}

		// Console display
		fmt.Println("\n")
		fmt.Println(h.Placeholder)         // render the placeholder
		fmt.Println(h.Chances)             // render the chances left
		fmt.Println(getEntries(h.Entries)) // show the letters or words guessed till now.
		fmt.Printf("Guess a letter or the word: ")

		// take the input
		str := ""
		fmt.Scanln(&str)

		// compare and update entries, placeholder and chances.
		if _, ok := h.Entries[str]; ok {
			fmt.Println("You have already entered")
			continue
		}

		if len(str) > 2 {
			if str == word {
				result <- "You have guessed correctly"
				break
			} else {
				h.Entries[str] = true
				h.Chances = h.Chances - 1
				continue
			}
		}

		h.Entries[str] = true

		if strings.Contains(word, str) {
			for k, v := range word {
				guessedString := string(v)
				if guessedString == str {
					h.Placeholder[k] = guessedString
				}
			}
		} else {
			h.Chances = h.Chances - 1
		}
	}

}

func getWord(dev bool) string {
	defaultWord := "elephant"

	if dev {
		return defaultWord
	}

	response, err := http.Get("https://random-word-api.herokuapp.com/word")

	if err != nil {
		return defaultWord
	}

	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		return defaultWord
	}

	var words []string

	err = json.Unmarshal(body, &words)

	if err != nil {
		return defaultWord
	}

	return words[0]
}

func main() {
	var dev bool

	flag.BoolVar(&dev, "dev", false, "use default word")
	flag.Parse()

	word := getWord(dev)

	h := Hangman{Entries: make(map[string]bool), Placeholder: make([]string, len(word), len(word)), Chances: len(word)}

	for k := range word {
		h.Placeholder[k] = "_"
	}

	result := make(chan string)
	go play(h, word, result)

	ticker := time.NewTicker(1 * time.Minute)

	for {
		select {
		case output := <-result:
			fmt.Println(output)
			return

		case <-ticker.C:
			fmt.Println("\nTimeout")
			fmt.Printf("Word:%v\n", word)
			return
		}
	}
}
