package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/notemptylist/gophersizes/linkparser/pkg/linkparse"
)

var HTMLfiles = [...]string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"}

func main() {

	fname := flag.String("inputfile", "", "HTML input file to parse.")
	flag.Parse()
	file, err := os.Open(*fname)
	if err != nil {
		panic(err)
	}
	parser, err := linkparse.New(file)
	if err != nil {
		panic(err)
	}
	for _, link := range parser.EmitLinks() {
		fmt.Printf("'%s' -> '%s'\n", link.Text, link.Href)
	}
}
