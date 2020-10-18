package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"jumble/dictionary"
	"jumble/srvr"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	dictionaryFileName := flag.String("d", "/usr/share/dict/words", "file name full of words")
	stopWordsFileName := flag.String("s", "", "file name full of words to ignore in dictionaries")
	portString := flag.String("p", "8012", "TCP port on which to listen")
	debug := flag.Bool("v", false, "verbose output per request")
	dump := flag.Bool("D", false, "dump final dictionary on stdout")
	flag.Parse()

	stopWords := readStopWords(*stopWordsFileName)

	buffer, err := ioutil.ReadFile(*dictionaryFileName)
	if err != nil {
		log.Fatalf("Problem reading dictionary file %q: %v\n", *dictionaryFileName, err)
	}

	before := time.Now()
	dict := dictionary.Build(buffer, stopWords)
	dict.Dedupe()

	if *dump {
		dict.Dump(os.Stdout)
		return
	}

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

// readStopWords has the name of a file as its formal argument,
// returns a map of strings to bool, if a string appears,
// don't put it in the dictionary.
func readStopWords(stopWordsFileName string) map[string]bool {

	if stopWordsFileName == "" {
		return make(map[string]bool)
	}

	buffer, err := ioutil.ReadFile(stopWordsFileName)
	if err != nil {
		log.Fatalf("Problem reading stop words file %q: %v\n", stopWordsFileName, err)
	}

	lines := bytes.Split(buffer, []byte{'\n'})

	removes := make(map[string]bool)
	for i := range lines {
		if len(lines[i]) == 0 {
			continue
		}
		removes[string(lines[i])] = true
	}

	return removes
}
