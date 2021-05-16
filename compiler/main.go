package main

import (
	"context"
	"io/ioutil"
	"net/http"
	"os"

	compiler "github.com/maratona-run-time/Maratona-Runtime/compiler/src"
	"github.com/maratona-run-time/Maratona-Runtime/utils"

	"github.com/go-martini/martini"
	graphql "github.com/hasura/go-graphql-client"
	"github.com/martini-contrib/binding"
	"github.com/rs/zerolog"
)

// FileForm receives a submission ID
type FileForm struct {
	ID string `form:"id"`
}

var sourceFileName = map[string]string{
	"C":      "program.c",
	"C++":    "program.cpp",
	"C++11":  "program.cpp",
	"Python": "program.py",
	"Go":     "program.go",
}

type submission struct {
	Submission struct {
		Language string
		Source   []byte
	} `graphql:"submission(id: $id)"`
}

func createCompilerServer(client utils.QueryClient, logger zerolog.Logger) *martini.ClassicMartini {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, req FileForm) {
		var info submission
		variables := map[string]interface{}{
			"id": graphql.ID(req.ID),
		}
		graphqlErr := client.Query(context.Background(), &info, variables)
		if graphqlErr != nil {
			msg := "An error occurred while trying to fetch submission '" + req.ID + "' details"
			logger.Error().
				Err(graphqlErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, graphqlErr)
			return
		}
		fileName := sourceFileName[info.Submission.Language]
		f, createErr := os.Create(fileName)
		if createErr != nil {
			msg := "An error occurred while trying to create a file named '" + fileName + "'"
			logger.Error().
				Err(createErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, createErr)
			return
		}
		f.Write(info.Submission.Source)
		f.Close()
		ret, compilerErr := compiler.Compile(info.Submission.Language, fileName, logger)
		err := os.Remove(fileName)
		if err != nil {
			msg := "Could not remove source file"
			logger.Error().
				Err(err).
				Msg(msg)
		}
		if compilerErr != nil {
			msg := "An error occurred while trying to compile program in language '" + info.Submission.Language + "'"
			logger.Error().
				Err(compilerErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, compilerErr)
			return
		}
		binary, readErr := ioutil.ReadFile(ret)
		if readErr != nil {
			msg := "An error occurred while trying to read binary file"
			logger.Error().
				Err(readErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, readErr)
			return
		}
		writeErr := ioutil.WriteFile("/var/program.out", binary, 0777)
		if writeErr != nil {
			msg := "An error occurred while trying to write binary file to shared volume"
			logger.Error().
				Err(writeErr).
				Msg(msg)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, writeErr)
			return
		}
		http.ServeFile(rs, rq, ret)
	})
	return m
}

func main() {
	logger, logFile := utils.InitLogger("compiler")
	defer logFile.Close()
	client := graphql.NewClient("http://orm:8084/graphql", nil)
	m := createCompilerServer(client, logger)
	m.RunOnAddr(":8081")
}
