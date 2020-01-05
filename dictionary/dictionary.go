package dictionary

import (
	"bytes"
	"sort"
	"strings"
)

// keys are alphabetized words, values are the words
type Dictionary map[string][]string

// Build fills in a Dictionary from a buffer, which
// should constitute all the lines from a file of words,
// one word per line.
func Build(buffer []byte) Dictionary {
	dict := Dictionary(make(map[string][]string))
	lines := bytes.Fields(buffer)
	for _, word := range lines {
		if w, a, saveit := Alphabetizer([]rune(string(word))); saveit {
			if _, ok := dict[a]; ok {
				dict[a] = append(dict[a], string(w))
				continue
			}
			words := make([]string, 1)
			words[0] = string(w)
			dict[a] = words
		}
	}
	return dict
}

// Dedupe eliminates duplicate words from list of words that
// alphabetize to a given string.
func (dict Dictionary) Dedupe() (int, int) {
	maxbefore := 0
	max := 0
	for alphabetized, words := range dict {
		if n := len(words); n > maxbefore {
			maxbefore = n
		}
		d := make(map[string]bool)
		for _, word := range words {
			d[word] = true
		}
		var newwords []string
		for word, _ := range d {
			newwords = append(newwords, word)
		}
		if n := len(newwords); n > max {
			max = n
		}
		dict[alphabetized] = newwords
	}
	return maxbefore, max
}

type RuneArray []rune

func (r RuneArray) Len() int           { return len(r) }
func (r RuneArray) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RuneArray) Less(i, j int) bool { return r[i] < r[j] }

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
	sort.Sort(RuneArray(alphabetized))
	return word, string(alphabetized), true
}
