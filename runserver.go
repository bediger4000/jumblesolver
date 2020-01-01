package main

import (
	"io/ioutil"
	"jumble/dictionary"
	"jumble/srvr"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	buffer, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatalf("Problem reading %q: %vn", os.Args[1], err)
	}

	srv := &srvr.Srvr{
		Router:    http.NewServeMux(),
		FindWords: dictionary.Build(buffer),
	}
	srv.Routes()
	s := &http.Server{
		Addr:           ":8012",
		Handler:        srv.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Fatal(s.ListenAndServe())

}
