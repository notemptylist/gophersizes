package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"

	"github.com/notemptylist/gophersizes/linkparser/pkg/linkparse"
)

var HTMLfiles = [...]string{"ex1.html", "ex2.html", "ex3.html", "ex4.html"}

func main() {

	fname := flag.String("inputfile", "", "HTML input file to parse.")
	flag.Parse()
	if *fname == "" {
		fmt.Println("Using a random default file.")
		fname = &HTMLfiles[rand.Intn(len(HTMLfiles))]
	}
	file, err := os.Open(*fname)
	if err != nil {
		panic(err)
	}
	for _, link := range linkparse.ParseLinks(file) {
		fmt.Printf("'%s' -> '%s'\n", link.Text, link.Href)
	}
}
