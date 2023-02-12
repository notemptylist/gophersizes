package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/notemptylist/gophersizes/cyoa"
)

func main() {
	file := flag.String("file", "gopher.json", "Story content in JSON format")
	port := flag.Int("port", 9000, "the port to bind to")
	flag.Parse()
	fmt.Printf("Parsing the story in %s\n", *file)
	story, err := cyoa.ParseFile(*file)
	if err != nil {
		log.Fatal(err)
	}

	h := cyoa.NewHandler(story)
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Running web server on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}
