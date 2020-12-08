package utils

import (
	"io/ioutil"

	"github.com/rs/zerolog"
)

func InitDummyLogger() zerolog.Logger {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: ioutil.Discard}
	logger := zerolog.New(consoleWriter)
	return logger
}
