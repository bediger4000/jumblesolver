package solver

// Word holds a jumbled word from user input,
// along with "marks" for the letters that get promoted
// into the final jumbled word(s), and a mark to use it
// "as is", because it's something like a proper name that
// isn't in the dictionary.
type Word struct {
	Word        []rune
	MarkedChars []int
	AsIs        bool
}

// RuneArray holds a sortable, sometimes sorted, array of runes
type RuneArray []rune

func (n RuneArray) Len() int           { return len(n) }
func (n RuneArray) Less(i, j int) bool { return n[i] < n[j] }
func (n RuneArray) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
