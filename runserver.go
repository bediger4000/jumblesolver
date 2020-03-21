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
	dictionaryFileName := flag.String("d", "/usr/share/dict/words", "file name full of words")
	portString := flag.String("p", "8012", "TCP port on which to listen")
	debug := flag.Bool("v", false, "verbose output per request")
	flag.Parse()

	buffer, err := ioutil.ReadFile(*dictionaryFileName)
	if err != nil {
		log.Fatalf("Problem reading dictionary file %q: %v\n", *dictionaryFileName, err)
	}

	before := time.Now()
	dict := dictionary.Build(buffer)
	dict.Dedupe()
	fmt.Printf("Dictionary file %s has %d keys, construction %v\n", *dictionaryFileName, len(dict), time.Since(before))

	srv := &srvr.Srvr{
		Router:    http.NewServeMux(),
		FindWords: dict,
		Debug:     *debug,
	}

	srv.Routes()

	s := &http.Server{
		Addr:           ":" + *portString,
		Handler:        srv.Router,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Jumbled-word solver listening on TCP port %s\n", *portString)

	log.Fatal(s.ListenAndServe())

}
