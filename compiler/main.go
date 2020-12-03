package main

import (
	"io"
	"mime/multipart"
	"net/http"

	"github.com/maratona-run-time/Maratona-Runtime/errors"

	"github.com/go-martini/martini"
	compiler "github.com/maratona-run-time/Maratona-Runtime/compiler/src"
	"github.com/martini-contrib/binding"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type FileForm struct {
	Language string                `form:"language"`
	Program  *multipart.FileHeader `form:"program"`
}

var sourceFileName = map[string]string{
	"C":      "program.c",
	"C++":    "program.cpp",
	"C++11":  "program.cpp",
	"Python": "program.py",
	"Go":     "program.go",
}

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logFile, logErr := os.OpenFile("compiler.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	defer logFile.Close()
	if logErr != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Fatal().Err(logErr).Msg("Could not create log file")
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	logger := zerolog.
		New(multi).
		With().
		Timestamp().
		Str("MaRT", "compiler").
		Logger().
		Level(zerolog.DebugLevel)

	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, req FileForm) {
		fileName := sourceFileName[req.Language]
		f, createErr := os.Create(fileName)
		if createErr != nil {
			msg := "An error occurred while trying to create a file named '" + fileName + "'"
			logger.Error().
				Err(createErr).
				Msg(msg)
			errors.WriteResponse(rs, http.StatusBadRequest, msg, createErr)
			return
		}
		program, pErr := req.Program.Open()
		if pErr != nil {
			msg := "An error occurred while trying to open the received program"
			logger.Error().
				Err(pErr).
				Msg(msg)
			errors.WriteResponse(rs, http.StatusBadRequest, msg, pErr)
			return
		}
		io.Copy(f, program)
		f.Close()
		program.Close()
		ret, compilerErr := compiler.Compile(req.Language, fileName, logger)
		if compilerErr != nil {
			msg := "An error occurred while trying compile program in language '" + req.Language + "'"
			logger.Error().
				Err(compilerErr).
				Msg(msg)
			errors.WriteResponse(rs, http.StatusBadRequest, msg, compilerErr)
			return
		}
		http.ServeFile(rs, rq, ret)
	})
	m.RunOnAddr(":8080")
}
