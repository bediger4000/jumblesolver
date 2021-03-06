package stringnumber

// Digit - a single, string-valued digit in a number analog
type Digit struct {
	values       []string
	currentIndex int
	maxIndex     int
}

// Number - a multi-digit number analog, allowing iteration
// through all combinations of Digits
type Number []*Digit

// NewDigit creates a string-valued digit-analog
func NewDigit(vals []string) *Digit {
	return &Digit{
		values:   vals,
		maxIndex: len(vals) - 1,
	}
}

// Increment a string-valued digit, return true if it incremented
// past the last digit, i.e. it carries the incrementation to
// the next digit in the "number"
func (d *Digit) Increment() (carry bool) {
	d.currentIndex++
	carry = false
	if d.currentIndex > d.maxIndex {
		carry = true
		d.currentIndex = 0
	}
	return
}

// Current gives back the current value of the digit analog
func (d *Digit) Current() string {
	return d.values[d.currentIndex]
}

// Next gives back a number-analog and a bool,
// which is false while the number-analog hasn't
// gone through all possible internal states.
func (n *Number) Next() ([]string, bool) {
	return n.Current(), n.Increment()
}

// Reset a Number to its original starting state
func (n *Number) Reset() {
	for _, d := range *n {
		d.currentIndex = 0
	}
}

// Current gives back the current digits in a []string,
// least significant to most significant, left to right
func (n *Number) Current() []string {
	answer := make([]string, len(*n))

	for i, d := range *n {
		answer[i] = d.Current()
	}

	return answer
}

// Increment a number analog, doing carries as necessary,
// returning false until the most significant digit needs
// to carry, which means the iteration through "numbers" is done
func (n *Number) Increment() (done bool) {
	carry := true
	digitIndex := 0
	digitCount := len(*n)

	for carry && digitIndex < digitCount {
		carry = (*n)[digitIndex].Increment()
		digitIndex++
	}

	return carry
}
