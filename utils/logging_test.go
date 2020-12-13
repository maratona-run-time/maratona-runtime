package utils

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"

	"github.com/rs/zerolog/log"
)

func TestInitLogger(t *testing.T) {
	type logEntry struct {
		Level, MaRT, Message string
		Time                 float32
	}

	// Logging
	containerName := "utils_test"
	message := "This is a test"
	logger, logFile := InitLogger(containerName)
	logger.Debug().Msg(message)
	logFile.Close()

	// Reading Log
	file, errOpen := ioutil.ReadFile(containerName + ".log")
	if errOpen != nil {
		log.Error().
			Err(errOpen).
			Msg("Error opening log file")
	}

	// Validating Log
	var data logEntry
	json.Unmarshal(file, &data)
	if data.MaRT != containerName {
		t.Errorf("Expected %s as name of container but got %s\n", containerName, data.MaRT)
	}

	// Cleanup
	errRem := os.Remove(containerName + ".log")
	if errRem != nil {
		log.Error().
			Err(errRem).
			Msg("Error removing log file")
	}

}
