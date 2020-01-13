package solver

type Word struct {
	Word        []rune
	MarkedChars []int
	AsIs        bool
}

type RuneArray []rune

func (n RuneArray) Len() int           { return len(n) }
func (n RuneArray) Less(i, j int) bool { return n[i] < n[j] }
func (n RuneArray) Swap(i, j int)      { n[i], n[j] = n[j], n[i] }
