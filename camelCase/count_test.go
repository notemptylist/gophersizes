package camelCase_test

import (
	"fmt"
	"testing"

	"github.com/notemptylist/gophersizes/camelCase"
)

func TestCount(t *testing.T) {

	var tests = []struct {
		s    string
		want int
	}{
		{"saveChangesInTheEditor", 5},
		{"camelCaseStringWithSixWords", 6},
	}
	for _, tt := range tests {
		testname := fmt.Sprintf("%s, %d", tt.s, tt.want)
		t.Run(testname, func(t *testing.T) {
			ans := camelCase.CountCamel(tt.s)
			if ans != tt.want {
				t.Errorf("got %d want %d", ans, tt.want)
			}
		})
	}
}
