package dictionary

import (
	"bytes"
	"io/ioutil"
	"sort"
	"strings"
)

// keys are alphabetized words, values are the words
type Dictionary map[string][]string

func Build(filename string) (Dictionary, error) {
	buffer, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	dict := Dictionary(make(map[string][]string))
	lines := bytes.Fields(buffer)
	for _, word := range lines {
		if w, a, saveit := alphabetizer(word); saveit {
			if _, ok := dict[a]; ok {
				dict[a] = append(dict[a], string(w))
				continue
			}
			words := make([]string, 1)
			words[0] = string(w)
			dict[a] = words
		}
	}
	// eliminate duplicate words from list of words that
	// alphabetize to a given string.
	for alphabetized, words := range dict {
		d := make(map[string]bool)
		for _, word := range words {
			d[word] = true
		}
		var newwords []string
		for word, _ := range d {
			newwords = append(newwords, word)
		}
		dict[alphabetized] = newwords
	}
	return dict, nil
}

type RuneArray []rune

func (r RuneArray) Len() int           { return len(r) }
func (r RuneArray) Swap(i, j int)      { r[i], r[j] = r[j], r[i] }
func (r RuneArray) Less(i, j int) bool { return r[i] < r[j] }

// alphabetizer receives a word as an array of bytes, returns
// a lowercase word, the word with all runes alphabetized,
// a bool indicating whether to save this word or not.
// Don't save words with apostrophes, commas, etc in them.
func alphabetizer(rawWord []byte) (string, string, bool) {
	word := strings.ToLower(string(rawWord))
	runes := []rune(word)
	if strings.ContainsAny(word, `'"!,:;.?/`) {
		return "", "", false
	}
	sort.Sort(RuneArray(runes))
	return word, string(runes), true
}
