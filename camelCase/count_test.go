package main

import (
	"testing"
)

const check = "\u2713"
const cross = "\u2717"

func TestCount(t *testing.T) {

	var tests = []struct {
		s    string
		want int
	}{
		{"saveChangesInTheEditor", 5},
		{"camelCaseStringWithSixWords", 6},
		{"hump", 1},
		{"", 0},
	}
	for _, tt := range tests {
		ans := countCamel(tt.s)
		if ans != tt.want {
			t.Errorf("%s got %d want %d", cross, ans, tt.want)
		} else {
			t.Logf("%s got %d want %d", check, ans, tt.want)
		}
	}
}
