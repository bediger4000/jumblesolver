package dictionary

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
)

// Dictionary keys are alphabetized words, values are the words
type Dictionary map[string][]string

// Build fills in a Dictionary from a buffer, which
// should constitute all the lines from a file of words,
// one word per line.
func Build(buffer []byte, stopWords map[string]bool) Dictionary {
	dict := Dictionary(make(map[string][]string))
	lines := bytes.Fields(buffer)
	for _, word := range lines {
		w, a, saveit := Alphabetizer([]rune(string(word)))
		if !saveit {
			continue
		}
		if stopWords[w] {
			continue
		}
		if _, ok := dict[a]; ok {
			dict[a] = append(dict[a], string(w))
			continue
		}
		dict[a] = []string{string(w)}
	}
	return dict
}

// Dump method allows human-readable output of alphabetized-word-keys,
// and the corresponding values, word lists. There's a command line flag
// that ends up invoking this.
func (dict Dictionary) Dump(w io.Writer) {
	for alphabetized, word := range dict {
		fmt.Fprintf(w, "%s\t%s\n", alphabetized, word)
	}
}

// Dedupe eliminates duplicate words from list of words that
// alphabetize to a given string.
func (dict Dictionary) Dedupe() {
	for alphabetizedKey, words := range dict {
		d := make(map[string]bool)
		for _, word := range words {
			d[word] = true
		}
		var newwords []string
		for word := range d {
			newwords = append(newwords, word)
		}
		dict[alphabetizedKey] = newwords
	}
}

// runeArray does duplicate type RuneArray in package solver,
// but if I use solver.RuneArray here, I get a circular import.
// The answer is a separate package, but too much work.
type runeArray []rune

func (r runeArray) Len() int           { return len(r) }
func (r runeArray) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r runeArray) Less(i, j int) bool { return r[i] < r[j] }

// Alphabetizer receives a word as an array of bytes, returns
// a lowercase word, the word with all runes alphabetized,
// a bool indicating whether to save this word or not.
// Don't save words with apostrophes, commas, etc in them.
func Alphabetizer(runes []rune) (string, string, bool) {
	word := strings.ToLower(string(runes))
	if strings.ContainsAny(word, `!@#$%^&*()-_+={[}]:;"'<,>.?/`) {
		return "", "", false
	}
	alphabetized := make([]rune, len(runes))
	copy(alphabetized, runes)
	sort.Sort(runeArray(alphabetized))
	return word, string(alphabetized), true
}
