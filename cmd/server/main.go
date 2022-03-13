package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	// https://yalantis.com/blog/speed-up-json-encoding-decoding/
	jsoniter "github.com/json-iterator/go"
	log "github.com/sirupsen/logrus"
)

var jsonfast = jsoniter.ConfigCompatibleWithStandardLibrary

var exists = struct{}{}

func main() {
	path := "/input.json"

	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(fmt.Errorf("could not open file: %s", err))
	}

	storage, err := NewStorage(content)
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	stats := &Stats{VMCount: len(storage.Map)}

	mux.Handle("/api/v1/attack", stats.Wrap("/api/v1/attack", storage))
	mux.Handle("/api/v1/stats", stats.Wrap("/api/v1/stats", stats))

	log.Println("Listening on :8080...")
	err = http.ListenAndServe(":8080", mux)
	log.Fatal(err)
}
