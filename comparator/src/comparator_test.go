package comparator

import (
	"testing"
)

func TestCompare(t *testing.T) {
	expected, program := "YES", "yes"
	if Compare(expected, program) == false {
		t.Errorf("\"%s\" should be equal to \"%s\"", expected, program)
	}
	expected, program = "yes", "no"
	if Compare(expected, program) == true {
		t.Errorf("\"%s\" should be different to \"%s\"", expected, program)
	}
}
