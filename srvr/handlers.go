package srvr

import (
	"fmt"
	"net/http"
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
		fmt.Printf("Raw form value:\n%s\n", r.FormValue("word"))
	}
}
