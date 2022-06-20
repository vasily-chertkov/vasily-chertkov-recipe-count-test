package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

var jsonfast = jsoniter.ConfigCompatibleWithStandardLibrary

type stringSet []string

func (s *stringSet) String() string {
	return strings.Join(*s, ", ")
}

func (s *stringSet) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func main() {
	var postcode string
	var fromStr string
	var toStr string
	var filters stringSet

	fSrv := flag.NewFlagSet("stats-counter", flag.ExitOnError)
	fSrv.StringVar(&postcode, "postcode", "10120", "postcode to count deliveries")
	fSrv.StringVar(&fromStr, "from", "10AM", "start time to count deliveries")
	fSrv.StringVar(&toStr, "to", "3PM", "end time to count deliveries")
	fSrv.Var(&filters, "f", "Filters for recipes (default 'Potato, Veggie, Mushroom', allowed use multiple values)")

	if err := fSrv.Parse(os.Args[1:]); err != nil {
		fSrv.PrintDefaults()
		log.Fatalf("flag parsing arguments. Error: %+v", err)
	}

	from, err := time.Parse("3PM", fromStr)
	if err != nil {
		log.Fatal(err)
	}
	to, err := time.Parse("3PM", toStr)
	if err != nil {
		log.Fatal(err)
	}

	if len(filters) == 0 {
		filters = append(filters, []string{"Potato", "Veggie", "Mushroom"}...)
	}

	now := time.Now()
	t := now

	path := "/input.json"

	f, err := os.Open(path)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	stats := NewStatsIter(f, postcode, from, to, filters)

	enc := jsonfast.NewEncoder(os.Stdout)
	enc.SetIndent("", "    ")
	if err := enc.Encode(&stats); err != nil {
		log.Fatal(err)
	}

	now = time.Now()
	fmt.Println("Processing time", now.Sub(t))
	t = now
}
