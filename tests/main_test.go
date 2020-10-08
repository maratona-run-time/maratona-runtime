package main

import (
	"Maratona-Runtime/executor"
	"testing"
)

func TestCompare(t *testing.T) {
	expected, program := "YES", "yes"
	if executor.Compare(expected, program) == false {
		t.Errorf("\"%s\" should be equal to \"%s\"", expected, program)
	}
	expected, program = "yes", "no"
	if executor.Compare(expected, program) == true {
		t.Errorf("\"%s\" should be different to \"%s\"", expected, program)
	}
}
