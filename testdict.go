package main

import (
	"jumble/dictionary"

	"fmt"
	"log"
	"os"
)

func main() {
	dict, err := dictionary.Build(os.Args[1])
	if err != nil {
		log.Fatalf("Problem with %q: %v\n", os.Args[1], err)
	}
	fmt.Printf("%d words in dictionary\n", len(dict))
	if len(os.Args) > 2 {
		if _, alphabetized, valid := dictionary.Alphabetizer([]byte(os.Args[2])); valid {
			if matches, ok := dict[alphabetized]; ok {
				spacer := ""
				for _, word := range matches {
					fmt.Printf("%s%s", spacer, word)
					spacer = ", "
				}
				fmt.Printf("\n")
				return
			}
			fmt.Printf("Did not find %q in dictionary\n", os.Args[2])
			return
		}
	}
	i := 0
	for a, words := range dict {
		i++
		if i > 200 {
			break
		}
		fmt.Printf("%q\t", a)
		spacer := ""
		for _, word := range words {
			fmt.Printf("%s%q", spacer, word)
			spacer = ", "
		}
		fmt.Printf("\n")
	}
}
