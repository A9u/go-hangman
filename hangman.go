package main

import (
	"fmt"
	"strings"
)

func getEntries(entries map[string]bool) (keys []string) {
	for k, _ := range entries {
		keys = append(keys, k)
	}

	return
}

func main() {
	word := "elephant"

	// lookup for entries made by the user.
	entries := make(map[string]bool)

	// list of "_" corrosponding to the number of letters in the word. [ _ _ _ _ _ ]
	placeholder := make([]string, len(word), len(word))

	for k := range word {
		placeholder[k] = "_"
	}

	chances := len(word)
	for {

		correctEntries := strings.Join(placeholder, "")

		// evaluate a loss! If user guesses a wrong letter or the wrong word, they lose a chance.

		if chances == 0 && correctEntries != word {
			fmt.Println("You lose")
			break
		}

		// evaluate a win!
		if correctEntries == word {
			fmt.Println("You won")
			break
		}

		// Console display
		fmt.Println("\n")
		fmt.Println(placeholder)         // render the placeholder
		fmt.Println(chances)             // render the chances left
		fmt.Println(getEntries(entries)) // show the letters or words guessed till now.
		fmt.Printf("Guess a letter or the word: ")

		// take the input
		str := ""
		fmt.Scanln(&str)

		// compare and update entries, placeholder and chances.
		if _, ok := entries[str]; ok {
			fmt.Println("You have already entered")
			continue
		}

		if len(str) > 2 {
			if str == word {
				fmt.Println("You have guessed correctly")
				break
			} else {
				entries[str] = true
				chances = chances - 1
				continue
			}
		}

		entries[str] = true

		if strings.Contains(word, str) {
			for k, v := range word {
				guessedString := string(v)
				if guessedString == str {
					placeholder[k] = guessedString
				}
			}
		} else {
			chances = chances - 1
		}
	}
}
