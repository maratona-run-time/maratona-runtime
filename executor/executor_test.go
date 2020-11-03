package executor

import "testing"

func TestExecuteOK(t *testing.T) {
	ver := make(chan []string)
	go Execute("./tests/program", "./tests/in", 1.0, ver)

	status := <-ver
	if status[0] != "OK" && status[1] == "3" {
		t.Errorf("Execução não deu OK")
	}
}

func TestExecuteTLE(t *testing.T) {
	ver := make(chan []string)
	go Execute("./tests/programLento", "./tests/in", 1.0, ver)

	status := <-ver
	if status[0] != "TLE" {
		t.Errorf("Execução não excedeu tempo limite")
	}
}

func TestExecuteRTE(t *testing.T) {
	ver := make(chan []string)
	go Execute("./tests/programRuntimeError", "./tests/in", 1.0, ver)

	status := <-ver
	if status[0] != "RTE" {
		t.Errorf("Execução não causou erro de runtime")
	}
}
