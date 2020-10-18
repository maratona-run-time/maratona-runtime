package main

import (
	"testing"

	"github.com/maratona-run-time/Maratona-Runtime/comparator"
)

func TestCompare(t *testing.T) {
	expected, program := "YES", "yes"
	if comparator.Compare(expected, program) == false {
		t.Errorf("\"%s\" should be equal to \"%s\"", expected, program)
	}
	expected, program = "yes", "no"
	if comparator.Compare(expected, program) == true {
		t.Errorf("\"%s\" should be different to \"%s\"", expected, program)
	}
}
