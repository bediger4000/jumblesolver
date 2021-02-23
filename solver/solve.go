package solver

import (
	"fmt"
	"jumblesolver/runenumber"
	"jumblesolver/stringnumber"

	"jumblesolver/dictionary"
)

// LookupWords finds a []string for each solver.Word in its input,
// based on the contents of dict.
// It puts a zero-length []string for use-as-is, and for words that
// don't qualify, which dictionary.Alphabetizer can decide to say
// about an input Word.
func LookupWords(dict dictionary.Dictionary, words []Word) [][]string {
	var matches [][]string
	for _, word := range words {
		var matched []string
		if word.AsIs {
			matched = []string{string(word.Word)}
		} else {
			_, alphabetized, useit := dictionary.Alphabetizer(word.Word)
			if useit {
				if m, ok := dict[alphabetized]; ok {
					matched = m
				}
			}
		}
		matches = append(matches, matched)
	}
	return matches
}

// CreateKeys makes all the unique (alphabetized) dictionary keys
// that can be calculated from an array of "alternates". Alternates
// come from when a jumbled word has more than one match in the dictionary,
// and so the promoted letters might end up different.
// For example, key: aelm, first and last letters promoted would end up
// matching unjumbled words "meal" and "lame", which would promote 'm','l"
// in one case, 'e', 'l' in the other. That leads to alternate sets of promoted
// letters for the final solution.
func CreateKeys(alternates [][]rune) [][]rune {

	var combos runenumber.Number

	for _, alt := range alternates {
		combos = append(combos, runenumber.NewDigit(alt))
	}

	combostrings := make(map[string]bool)

	done := false
	var current []rune
	for !done {
		current, done = combos.Next()
		_, alpha, _ := dictionary.Alphabetizer(current)
		combostrings[alpha] = true
	}

	var strings [][]rune
	for str := range combostrings {
		strings = append(strings, []rune(str))
	}

	return strings
}

// FindUniqueMatches uses a dictionary.Dict to find potential matches
// for the arrary of rune-arrays named alternates.
// alternates is of length of the answer string.
// Each array position is one or more runes that are specified as
// carried forward from the unjumbled words. Since a jumbled word could
// match two or more unjumbled words ("amet" for example), each position
// carried forward into the answer could have two or more possible runes.
// For example:
// amet - jumbled word
//
// meat - unjumbled word
// team - unjumbled word
// meta - unjumbled word
//
// If the first unjumbled position is carried forward, it can have ['m', 't'].
// If the 4th unjumbled position is carried forward, it can have ['a', 'm', 't']
func FindUniqueMatches(alternates [][]rune, dict dictionary.Dictionary) map[string][]string {

	// Here's where the alternates in a given position might end up
	// as more strings. keys should contain an array of len(alternates) strings,
	// and each should be unique.
	keys := CreateKeys(alternates)

	finalmatches := make(map[string][]string)

	for _, key := range keys {
		uniquematches := make(map[string]bool)
		if matches, ok := dict[string(key)]; ok {
			for _, match := range matches {
				uniquematches[match] = true
			}
		}
		var allmatched []string
		for match := range uniquematches {
			allmatched = append(allmatched, match)
		}
		finalmatches[string(key)] = allmatched
	}

	return finalmatches
}

// GenerateKeyCombos creates all the key (character-alphabetized string) combinations,
// based on the sets of letters promoted from the unjumbled words, the number and
// lengths of words in the final solution.
func GenerateKeyCombos(debug bool, alternates [][]rune, finalcount int, finalsizes []int) [][][]rune {

	if debug {
		fmt.Printf("Enter GenerateKeyCombos\n")
		defer fmt.Printf("Exit GenerateKeyCombos\n")
	}

	// func generateKeys uses len of finalsizes to know when to
	// stop recursing.
	finalsizes = fixSizes(len(alternates), finalsizes)
	fmt.Printf("finalsizes = %d\n", finalsizes)

	var keyCombos [][][]rune

	keyList := CreateKeys(alternates)
	for _, key := range keyList {
		keycandidates := generateKeyCombos(debug, key, finalsizes)
		keyCombos = append(keyCombos, keycandidates...)
	}

	return keyCombos
}

func generateKeyCombos(debug bool, wholeKey []rune, sizes []int) [][][]rune {

	if len(sizes) == 1 {
		return [][][]rune{[][]rune{wholeKey}}
	}

	n := sizes[0]

	keycandidates := keyCandidates(wholeKey, n)

	var allCandidates [][][]rune

	for _, candidate := range keycandidates {
		// candidate is [][]rune, but really a pair of []rune.
		// len(candidate[0]) == n, len(candidate[1]) == len(wholeKey) - n

		// subcandidates is an array of []rune, each []rune is a piece of wholeKey
		// made out of candidate[1] base on the rest of sizes[]
		subcandidates := generateKeyCombos(debug, candidate[1], sizes[1:])

		// join candidate[0] and each []rune in subcandidate - keyCandidates()
		// should have only given us unique candidate[] sub-arrays
		for _, subcandidate := range subcandidates {
			t := append([][]rune{candidate[0]}, subcandidate...)
			allCandidates = append(allCandidates, t)
		}
	}

	return allCandidates
}

func fixSizes(desired int, array []int) []int {
	sum := sumarray(array)
	for sum > desired {
		array = array[:len(array)-1]
		sum = sumarray(array)
	}
	if sum < desired {
		array = append(array, desired-sum)
	}
	return array
}

func sumarray(array []int) int {
	sum := 0
	for _, v := range array {
		sum += v
	}
	return sum
}

// SolutionsFromKeyCombos does the work of finding multi-word solutions
// from the keyCombos, which are alphabetized keys that fit the lengths
// and number of final solution words.
func SolutionsFromKeyCombos(debug bool, keyCombos [][][]rune, dict dictionary.Dictionary) [][]string {
	if debug {
		fmt.Printf("Enter SolutionsFromKeyCombos\n")
		defer fmt.Printf("Exit SolutionsFromKeyCombos\n")
	}

	var solutions [][]string

	for _, candidate := range keyCombos {
		// candidate is an array of []rune, each []rune can match a set of words in dict
		matchesForAll := true
		var matchingWords stringnumber.Number
		for _, candidatekey := range candidate {
			if matchedWords, ok := dict[string(candidatekey)]; ok {
				matchingWords = append(matchingWords, stringnumber.NewDigit(matchedWords))
			} else {
				matchesForAll = false
				break
			}
		}
		if matchesForAll {
			// word(s) existed for each key in the key combination
			done := false
			for !done {
				var match []string
				match, done = matchingWords.Next()
				solutions = append(solutions, match)
			}
		}
	}

	return solutions
}
