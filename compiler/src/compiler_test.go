package compiler

import (
	"fmt"
	"github.com/rs/zerolog"
	"io/ioutil"
	"os/exec"
	"testing"
)

func initLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: ioutil.Discard}
	logger := zerolog.
		New(consoleWriter).
		With().
		Logger()
	return logger
}

func setup(lang string, file string) {
	exec.Command("cp", "tests/"+file, "./"+file).Output()
}

func teardown(file string) {
	exec.Command("rm", "program.out", file).Output()
}

func TestCompilation(t *testing.T) {
	logger := initLogger()
	tests := []struct {
		lang string
		file string
	}{
		{"C", "program.c"},
		{"Python", "program.py"},
		{"C++11", "program.cpp"},
		{"Go", "program.go"},
	}
	for _, test := range tests {
		t.Run(fmt.Sprintf("%s/%s", test.lang, test.file), func(t *testing.T) {
			setup(test.lang, test.file)
			_, err := Compile(test.lang, test.file, logger)
			if err != nil {
				t.Errorf("Compilation failed!")
			}
			teardown(test.file)
		})
	}
}
