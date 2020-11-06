package executor

import "testing"

func TestExecuteOK(t *testing.T) {
	status := Execute("./tests/program", "./tests", 1.0)
	if status[0][1] != "OK" && status[0][2] == "3" {
		t.Errorf("Execução não deu OK")
	}
}

func TestExecuteTLE(t *testing.T) {
	status := Execute("./tests/programLento", "./tests", 1.0)
	if status[0][1] != "TLE" {
		t.Errorf("Execução não excedeu tempo limite")
	}
}

func TestExecuteRTE(t *testing.T) {
	status := Execute("./]tests/programRuntimeError", "./tests", 1.0)
	if status[0][1] != "RTE" {
		t.Errorf("Execução não causou erro de runtime")
	}
}
