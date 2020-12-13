package main

import (
	"io"
	"mime/multipart"
	"net/http"
	"os"

	compiler "github.com/maratona-run-time/Maratona-Runtime/compiler/src"
	"github.com/maratona-run-time/Maratona-Runtime/utils"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
	"github.com/rs/zerolog"
)

// FileForm is a struct containing the source code of a submission and string identifiying it language.
type FileForm struct {
	Language string                `form:"language"`
	Source   *multipart.FileHeader `form:"source"`
}

var sourceFileName = map[string]string{
	"C":      "program.c",
	"C++":    "program.cpp",
	"C++11":  "program.cpp",
	"Python": "program.py",
	"Go":     "program.go",
}

func createCompilerServer(logger zerolog.Logger) *martini.ClassicMartini {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, req FileForm) {
		fileName := sourceFileName[req.Language]
		f, createErr := os.Create(fileName)
		if createErr != nil {
			msg := "An error occurred while trying to create a file named '" + fileName + "'"
			logger.Error().
				Err(createErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, createErr)
			return
		}
		program, pErr := req.Source.Open()
		if pErr != nil {
			msg := "An error occurred while trying to open the received program"
			logger.Error().
				Err(pErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, pErr)
			return
		}
		io.Copy(f, program)
		f.Close()
		program.Close()
		ret, compilerErr := compiler.Compile(req.Language, fileName, logger)
		err := os.Remove(fileName)
		if err != nil {
			msg := "Could not remove source file"
			logger.Error().
				Err(err).
				Msg(msg)
		}
		if compilerErr != nil {
			msg := "An error occurred while trying compile program in language '" + req.Language + "'"
			logger.Error().
				Err(compilerErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, compilerErr)
			return
		}
		http.ServeFile(rs, rq, ret)
	})
	return m
}

func main() {
	logger, logFile := utils.InitLogger("compiler")
	defer logFile.Close()
	m := createCompilerServer(logger)
	m.RunOnAddr(":8080")
}
