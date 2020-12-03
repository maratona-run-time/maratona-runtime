package executor

import (
	"github.com/rs/zerolog"
	"io/ioutil"
	"strings"
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

func TestExecuteOK(t *testing.T) {
	logger := initLogger()
	status := Execute("./tests/program", "./tests", 1.0, logger)
	if status[0].Status != "OK" {
		t.Errorf("Expected status OK but got %s", status[0].Status)
	}
	if strings.EqualFold(status[0].Message, "3") {
		t.Errorf("Expected output 3 but got %s", status[0].Message)
	}
}

func TestExecuteTLE(t *testing.T) {
	logger := initLogger()
	status := Execute("./tests/programLento", "./tests", 1.0, logger)
	if status[0].Status != "TLE" {
		t.Errorf("Expected status TLE but got %s", status[0].Status)
	}
}

func TestExecuteRTE(t *testing.T) {
	logger := initLogger()
	status := Execute("./]tests/programRuntimeError", "./tests", 1.0, logger)
	if status[0].Status != "RTE" {
		t.Errorf("Expected status RLE but got %s", status[0].Status)
	}
}
