package srvr

import (
	"fmt"
	"jumble/dictionary"
	"net/http"
	"strings"
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

func (s *Srvr) handleIndex() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Enter handleIndex closure\n")
		defer fmt.Printf("Exit handleIndex closure\n")
		w.Write([]byte(indexHTML))
	}
}

func (s *Srvr) handleForm() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Enter handleForm closure\n")
		defer fmt.Printf("Exit handleForm closure\n")

		w.Header().Set("Content-Type", "text/html")
		x := strings.TrimSpace(r.FormValue("word"))
		fmt.Printf("Form word value:\n%s\n", x)

		if _, alphabetized, valid := dictionary.Alphabetizer([]byte(x)); valid {
			if matches, ok := s.FindWords[alphabetized]; ok {
				spacer := ""
				text := ""
				for _, word := range matches {
					text += fmt.Sprintf("%s%s", spacer, word)
					spacer = ", "
				}
				w.Write([]byte(fmt.Sprintf(formHTML, text)))
			}
		}

	}
}
