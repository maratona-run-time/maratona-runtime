package utils

import (
	"io/ioutil"
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// InitDummyLogger creates a dummy logger to be used in tests.
func InitDummyLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: ioutil.Discard}
	logger := zerolog.New(consoleWriter)
	return logger
}

// InitLogger creates and configures a logger to be used in a specific container.
func InitLogger(containerName string) (zerolog.Logger, *os.File) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logFile, errLog := os.OpenFile(containerName+".log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	if errLog != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Fatal().Err(errLog).Msg("Could not create log file")
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	logger := zerolog.
		New(multi).
		With().
		Timestamp().
		Str("MaRT", containerName).
		Logger().
		Level(zerolog.DebugLevel)
	return logger, logFile
}
