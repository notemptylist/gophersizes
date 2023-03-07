package camelCase

import (
	"flag"
	"fmt"
	"os"
)

const upperMin = 65 // A
const upperMax = 90 // Z
// countCamel counts the number of words encoded in the provided string.
func CountCamel(s string) int {
	if len(s) < 1 {
		return 0
	}
	var n int
	for _, ch := range s {
		if int(ch) >= upperMin && int(ch) <= upperMax {
			n++
		}
	}
	return n + 1
}

func main() {
	camelString := flag.String("input", "", "camelCase string to process.")
	flag.Parse()
	if *camelString == "" {
		flag.Usage()
		os.Exit(1)
	}

	fmt.Printf("%s contains %d words\n", *camelString, CountCamel(*camelString))
}
