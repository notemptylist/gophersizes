package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/notemptylist/gophersizes/cyoa"
)

func parseFile(fname string) cyoa.Story {
	f, err := os.Open(fname)
	if err != nil {
		log.Fatal(err)
	}

	d := json.NewDecoder(f)
	var story cyoa.Story
	if err = d.Decode(&story); err != nil {
		log.Fatal(err)
	}

	return story
}
func main() {
	file := flag.String("file", "gopher.json", "Story content in JSON format")
	flag.Parse()
	fmt.Printf("Parsing the story in %s\n", *file)
	fmt.Printf("%+v\n", parseFile(*file))

}
