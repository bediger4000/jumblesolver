package solver

import (
	"fmt"
	"sort"
)

func keyCandidates(instring []rune, firstWord int) [][][]rune {
	switch firstWord {
	case 1:
		return keyCandidates1(instring)
	case 2:
		return keyCandidates2(instring)
	case 3:
		return keyCandidates3(instring)
	case 4:
		return keyCandidates4(instring)
	case 5:
		return keyCandidates5(instring)
	case 6:
		return keyCandidates6(instring)
	case 7:
	}
	fmt.Printf("exit keyCandidates for %q, the wrong way\n", string(instring))
	return nil
}

func keyCandidates1(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n; mark1++ {

		part1 := RuneArray{letters[mark1]}

		var otherpart RuneArray
		for idx := 0; idx < n; idx++ {
			if idx != mark1 {
				otherpart = append(otherpart, letters[idx])
			}
		}

		uniques[string(part1)+string(otherpart)] = [][]rune{part1, otherpart}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}

func keyCandidates2(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-1; mark1++ {
		for mark2 := mark1 + 1; mark2 < n; mark2++ {

			part1 := RuneArray{
				letters[mark1], letters[mark2],
			}

			var otherpart RuneArray
			for idx := 0; idx < n; idx++ {
				if idx != mark1 && idx != mark2 {
					otherpart = append(otherpart, letters[idx])
				}
			}

			uniques[string(part1)+string(otherpart)] = [][]rune{part1, otherpart}
		}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}

func keyCandidates3(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-2; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-1; mark2++ {
			for mark3 := mark2 + 1; mark3 < n; mark3++ {

				part1 := RuneArray{
					letters[mark1], letters[mark2],
					letters[mark3],
				}

				var otherpart RuneArray
				for idx := 0; idx < n; idx++ {
					if idx != mark1 && idx != mark2 && idx != mark3 {
						otherpart = append(otherpart, letters[idx])
					}
				}

				uniques[string(part1)+string(otherpart)] = [][]rune{part1, otherpart}
			}
		}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}

func keyCandidates4(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-3; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-2; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-1; mark3++ {
				for mark4 := mark3 + 1; mark4 < n; mark4++ {

					part1 := RuneArray{
						letters[mark1], letters[mark2],
						letters[mark3], letters[mark4],
					}

					var otherpart RuneArray
					for idx := 0; idx < n; idx++ {
						if idx != mark1 && idx != mark2 && idx != mark3 && idx != mark4 {
							otherpart = append(otherpart, letters[idx])
						}
					}

					uniques[string(part1)+string(otherpart)] = [][]rune{part1, otherpart}
				}
			}
		}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}

func keyCandidates5(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-4; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-3; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-2; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-1; mark4++ {
					for mark5 := mark4 + 1; mark5 < n; mark5++ {

						part1 := RuneArray{
							letters[mark1], letters[mark2],
							letters[mark3], letters[mark4],
							letters[mark5],
						}

						var otherpart RuneArray
						for idx := 0; idx < n; idx++ {
							if idx != mark1 && idx != mark2 && idx != mark3 && idx != mark4 && idx != mark5 {
								otherpart = append(otherpart, letters[idx])
							}
						}

						uniques[string(part1)+string(otherpart)] = [][]rune{part1, otherpart}
					}
				}
			}
		}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}

func keyCandidates6(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-5; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-4; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-3; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-2; mark4++ {
					for mark5 := mark4 + 1; mark5 < n-1; mark5++ {
						for mark6 := mark5 + 1; mark6 < n; mark6++ {

							part1 := RuneArray{
								letters[mark1], letters[mark2],
								letters[mark3], letters[mark4],
								letters[mark5], letters[mark6],
							}

							var otherpart RuneArray
							for idx := 0; idx < n; idx++ {
								if idx != mark1 && idx != mark2 && idx != mark3 && idx != mark4 && idx != mark5 && idx != mark6 {
									otherpart = append(otherpart, letters[idx])
								}
							}

							uniques[string(part1)+string(otherpart)] = [][]rune{part1, otherpart}
						}
					}
				}
			}
		}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}
