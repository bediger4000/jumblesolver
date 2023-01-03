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
		return keyCandidates7(instring)
	case 8:
		return keyCandidates8(instring)
	case 9:
		return keyCandidates9(instring)
	case 10:
		return keyCandidates10(instring)
	case 11:
		return keyCandidates11(instring)
	case 12:
		return keyCandidates12(instring)
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

func keyCandidates7(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-5; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-4; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-3; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-2; mark4++ {
					for mark5 := mark4 + 1; mark5 < n-1; mark5++ {
						for mark6 := mark5 + 1; mark6 < n; mark6++ {
							for mark7 := mark6 + 1; mark7 < n; mark7++ {

								part1 := RuneArray{
									letters[mark1], letters[mark2],
									letters[mark3], letters[mark4],
									letters[mark5], letters[mark6],
									letters[mark7],
								}

								var otherpart RuneArray
								for idx := 0; idx < n; idx++ {
									if idx != mark1 && idx != mark2 && idx != mark3 &&
										idx != mark4 && idx != mark5 && idx != mark6 &&
										idx != mark7 {
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
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}

func keyCandidates8(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-5; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-4; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-3; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-2; mark4++ {
					for mark5 := mark4 + 1; mark5 < n-1; mark5++ {
						for mark6 := mark5 + 1; mark6 < n; mark6++ {
							for mark7 := mark6 + 1; mark7 < n; mark7++ {
								for mark8 := mark7 + 1; mark8 < n; mark8++ {

									part1 := RuneArray{
										letters[mark1], letters[mark2],
										letters[mark3], letters[mark4],
										letters[mark5], letters[mark6],
										letters[mark7], letters[mark8],
									}

									var otherpart RuneArray
									for idx := 0; idx < n; idx++ {
										if idx != mark1 && idx != mark2 && idx != mark3 &&
											idx != mark4 && idx != mark5 && idx != mark6 &&
											idx != mark7 && idx != mark8 {
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
		}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}
func keyCandidates9(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-5; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-4; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-3; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-2; mark4++ {
					for mark5 := mark4 + 1; mark5 < n-1; mark5++ {
						for mark6 := mark5 + 1; mark6 < n; mark6++ {
							for mark7 := mark6 + 1; mark7 < n; mark7++ {
								for mark8 := mark7 + 1; mark8 < n; mark8++ {
									for mark9 := mark8 + 1; mark9 < n; mark9++ {

										part1 := RuneArray{
											letters[mark1], letters[mark2],
											letters[mark3], letters[mark4],
											letters[mark5], letters[mark6],
											letters[mark7], letters[mark8], letters[mark9],
										}

										var otherpart RuneArray
										for idx := 0; idx < n; idx++ {
											if idx != mark1 && idx != mark2 && idx != mark3 &&
												idx != mark4 && idx != mark5 && idx != mark6 &&
												idx != mark7 && idx != mark8 && idx != mark9 {
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
			}
		}
	}

	var words [][][]rune

	for _, wordpair := range uniques {
		words = append(words, wordpair)
	}

	return words
}
func keyCandidates10(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-5; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-4; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-3; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-2; mark4++ {
					for mark5 := mark4 + 1; mark5 < n-1; mark5++ {
						for mark6 := mark5 + 1; mark6 < n; mark6++ {
							for mark7 := mark6 + 1; mark7 < n; mark7++ {
								for mark8 := mark7 + 1; mark8 < n; mark8++ {
									for mark9 := mark8 + 1; mark9 < n; mark9++ {
										for mark10 := mark9 + 1; mark10 < n; mark10++ {

											part1 := RuneArray{
												letters[mark1], letters[mark2],
												letters[mark3], letters[mark4],
												letters[mark5], letters[mark6],
												letters[mark7], letters[mark8],
												letters[mark9], letters[mark10],
											}

											var otherpart RuneArray
											for idx := 0; idx < n; idx++ {
												if idx != mark1 && idx != mark2 && idx != mark3 &&
													idx != mark4 && idx != mark5 && idx != mark6 &&
													idx != mark7 && idx != mark8 && idx != mark9 && idx != mark10 {
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
func keyCandidates11(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-5; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-4; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-3; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-2; mark4++ {
					for mark5 := mark4 + 1; mark5 < n-1; mark5++ {
						for mark6 := mark5 + 1; mark6 < n; mark6++ {
							for mark7 := mark6 + 1; mark7 < n; mark7++ {
								for mark8 := mark7 + 1; mark8 < n; mark8++ {
									for mark9 := mark8 + 1; mark9 < n; mark9++ {
										for mark10 := mark9 + 1; mark10 < n; mark10++ {
											for mark11 := mark10 + 1; mark11 < n; mark11++ {

												part1 := RuneArray{
													letters[mark1], letters[mark2],
													letters[mark3], letters[mark4],
													letters[mark5], letters[mark6],
													letters[mark7], letters[mark8],
													letters[mark9], letters[mark10],
													letters[mark11],
												}

												var otherpart RuneArray
												for idx := 0; idx < n; idx++ {
													if idx != mark1 && idx != mark2 && idx != mark3 &&
														idx != mark4 && idx != mark5 && idx != mark6 &&
														idx != mark7 && idx != mark8 && idx != mark9 &&
														idx != mark10 && idx != mark11 {
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
func keyCandidates12(letters []rune) [][][]rune {

	n := len(letters)
	sort.Sort(RuneArray(letters))

	uniques := make(map[string][][]rune)

	for mark1 := 0; mark1 < n-5; mark1++ {
		for mark2 := mark1 + 1; mark2 < n-4; mark2++ {
			for mark3 := mark2 + 1; mark3 < n-3; mark3++ {
				for mark4 := mark3 + 1; mark4 < n-2; mark4++ {
					for mark5 := mark4 + 1; mark5 < n-1; mark5++ {
						for mark6 := mark5 + 1; mark6 < n; mark6++ {
							for mark7 := mark6 + 1; mark7 < n; mark7++ {
								for mark8 := mark7 + 1; mark8 < n; mark8++ {
									for mark9 := mark8 + 1; mark9 < n; mark9++ {
										for mark10 := mark9 + 1; mark10 < n; mark10++ {
											for mark11 := mark10 + 1; mark11 < n; mark11++ {
												for mark12 := mark11 + 1; mark12 < n; mark12++ {

													part1 := RuneArray{
														letters[mark1], letters[mark2],
														letters[mark3], letters[mark4],
														letters[mark5], letters[mark6],
														letters[mark7], letters[mark8],
														letters[mark9], letters[mark10],
														letters[mark11], letters[mark12],
													}

													var otherpart RuneArray
													for idx := 0; idx < n; idx++ {
														if idx != mark1 && idx != mark2 && idx != mark3 &&
															idx != mark4 && idx != mark5 && idx != mark6 &&
															idx != mark7 && idx != mark8 && idx != mark9 &&
															idx != mark10 && idx != mark11 && idx != mark12 {
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
