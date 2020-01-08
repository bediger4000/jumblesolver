package solver

import (
	"jumble/runenumber"

	"jumble/dictionary"
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

func CreateKeys(alternates [][]rune, length int) []string {

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

	var strings []string
	for str, _ := range combostrings {
		strings = append(strings, str)
	}

	return strings
}
