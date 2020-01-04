package srvr

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"jumble/dictionary"
	"jumble/solver"
)

var indexHTML string = `
<html>
<head>
</head>
<body>
<form name="f" method="post" action="/form">
<input name="word" />
<input type="submit" />
</form>
</body>
</html>
`

var formHTML string = `
<html>
<head>
</head>
<body>
<form name="f" method="post" action="/form">
<p>%s</p>
<input name="word" />
<input type="submit" />
</form>
</body>
</html>
`

var explainHTML string = `<!DOCTYPE html>
<html>
    <head>
    <meta charset="UTF-8">
	</head>
	<body>
`

var errorHTML string = `<!DOCTYPE html>
<html>
    <head>
    <meta charset="UTF-8">
	</head>
	<body>
		<p>I had a problem.</p>
	</body>
</html>
`

func (s *Srvr) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Enter handleIndex closure\n")
		defer fmt.Printf("Exit handleIndex closure\n")
		w.Write([]byte(indexHTML))
	}
}

func (s *Srvr) handleJumble() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Enter handleJumble closure\n")
		defer fmt.Printf("Exit handleJumble closure\n")

		w.Header().Set("Content-Type", "text/html")

		wordCount, err := strconv.Atoi(strings.TrimSpace(r.FormValue("wordcount")))
		if err != nil {
			log.Printf("Finding value of wordcount: %v\n", err)
			w.Write([]byte(errorHTML))
			return
		}
		fmt.Printf("wordCount %d\n", wordCount)

		var words []solver.Word
		for wordNumber := 0; wordNumber < wordCount; wordNumber++ {
			var marks []int
			var word []rune
			for charNumber := 0; charNumber < 10; charNumber++ {
				wordCode := fmt.Sprintf("w%dc%d", wordNumber, charNumber)
				wordChar := strings.TrimSpace(r.FormValue(wordCode))
				if wordChar != "" {
					word = append(word, []rune(wordChar)[0])
				}

				markCode := wordCode + "forward"
				m := strings.TrimSpace(r.FormValue(markCode))
				if m == "on" {
					marks = append(marks, charNumber)
				}

			}
			asIsCode := fmt.Sprintf("w%dasis", wordNumber)
			aic := strings.TrimSpace(r.FormValue(asIsCode))
			asIs := false
			if aic == "on" {
				asIs = true
			}
			wrd := solver.Word{
				Word:        word,
				MarkedChars: marks,
				AsIs:        asIs,
			}
			words = append(words, wrd)
		}
		w.Write([]byte(explainHTML))
		for idx, word := range words {
			w.Write([]byte(fmt.Sprintf("<h1>Word %d</h1>\n", idx)))
			w.Write([]byte(fmt.Sprintf("<p>%q</p>\n", string(word.Word))))
			w.Write([]byte(fmt.Sprintf("<p>Marks: %v</p>\n", word.MarkedChars)))
			w.Write([]byte(fmt.Sprintf("<p>As-Is: %v</p>\n", word.AsIs)))
		}
		w.Write([]byte("</body></html>\n"))
	}
}

func (s *Srvr) handleForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Enter handleForm closure\n")
		defer fmt.Printf("Exit handleForm closure\n")

		w.Header().Set("Content-Type", "text/html")
		x := strings.TrimSpace(r.FormValue("word"))
		fmt.Printf("Form word value:\n%s\n", x)
		text := ""

		if _, alphabetized, valid := dictionary.Alphabetizer([]byte(x)); valid {
			if matches, ok := s.FindWords[alphabetized]; ok {
				spacer := ""
				for _, word := range matches {
					text += fmt.Sprintf("%s%s", spacer, word)
					spacer = ", "
				}
			}
		}
		w.Write([]byte(fmt.Sprintf(formHTML, text)))

	}
}
