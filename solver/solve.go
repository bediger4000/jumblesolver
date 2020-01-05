package solver

import "jumble/dictionary"

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
