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
	template := flag.String("template", "", "Alternative template to use")
	port := flag.Int("port", 9000, "the port to bind to")
	flag.Parse()
	fmt.Printf("Parsing the story in %s\n", *file)
	story, err := cyoa.ParseFile(*file)
	if err != nil {
		log.Fatal(err)
	}
	var opts []cyoa.HandlerOption
	if *template != "" {
		opts = append(opts, cyoa.WithTemplate(*template))
	}

	h := cyoa.NewHandler(story, opts...)
	addr := fmt.Sprintf(":%d", *port)
	fmt.Printf("Running web server on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, h))
}
