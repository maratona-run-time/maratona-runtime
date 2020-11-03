package verdict

import (
	"fmt"
	"testing"
)

func TestVerdict(t *testing.T) {
	tests := []struct {
		ver  string
		file string
	}{
		{"AC", "ac"},
		{"WA", "wa"},
		{"TLE", "tle"},
		{"RTE", "rte"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s", test.ver), func(t *testing.T) {
			res := make(chan string)
			go Verdict(1., "./tests/sum/"+test.file, "./tests/sum/inputs/", "./tests/sum/outputs/", res)
			status := <-res
			if status != test.ver {
				t.Errorf("Programa gerou status %v ao invés de %v", status, test.ver)
			}
		})
	}
}
