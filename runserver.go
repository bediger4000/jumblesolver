package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"jumble/dictionary"
	"jumble/srvr"
	"log"
	"net/http"
	"time"
)

func main() {
	dictionaryFileName := flag.String("d", "", "file name full of words")
	debug := flag.Bool("v", false, "verbose output per request")
	flag.Parse()
	buffer, err := ioutil.ReadFile(*dictionaryFileName)
	if err != nil {
		log.Fatalf("Problem reading %q: %v\n", *dictionaryFileName, err)
	}

	before := time.Now()
	dict := dictionary.Build(buffer)
	dict.Dedupe()
	fmt.Printf("Dictionary %d keys, construction %v\n", len(dict), time.Since(before))

	srv := &srvr.Srvr{
		Router:    http.NewServeMux(),
		FindWords: dict,
		Debug:     *debug,
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
