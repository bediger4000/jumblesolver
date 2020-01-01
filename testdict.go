package main

import (
	"jumble/dictionary"
	"time"

	"fmt"
	"io/ioutil"
	"log"
	"os"
)

func main() {
	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Problem reading %q: %vn", os.Args[1], err)
	}
	before := time.Now()
	dict := dictionary.Build(buffer)
	m1, m2 := dict.Dedupe()
	fmt.Printf("%d words in dictionary in %v\n", len(dict), time.Since(before))
	fmt.Printf("Max before dedupe %d, max after %d\n", m1, m2)
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
