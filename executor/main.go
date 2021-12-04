package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	graphql "github.com/hasura/go-graphql-client"
	executor "github.com/maratona-run-time/Maratona-Runtime/executor/src"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
	"github.com/martini-contrib/binding"
	"github.com/rs/zerolog"
)

// FileForm receives a submission ID
type FileForm struct {
	ID string `form:"id"`
}

type (
	Input struct {
		FileName string
		Content  []byte
	}
	Challenge struct {
		TimeLimit float32
		Inputs    []Input
	}
	Submission struct {
		Challenge Challenge
	}
	Info struct {
		Submission Submission `graphql:"submission(id: $id)"`
	}
)

func createExecutorServer(client utils.QueryClient, logger zerolog.Logger) *martini.ClassicMartini {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, req FileForm) []byte {
		var info Info
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
			return nil
		}

		os.Mkdir("inputs", 0700)

		for _, input := range info.Submission.Challenge.Inputs {
			testFileName := fmt.Sprintf("inputs/%s", input.FileName)
			writeErr := ioutil.WriteFile(testFileName, input.Content, 0777)
			if writeErr != nil {
				msg := "An error occurred while trying to write the received test to a local file named '" + testFileName
				utils.WriteResponse(rs, http.StatusInternalServerError, msg, writeErr)
				logger.Error().
					Err(writeErr).
					Msg(msg)
				return nil
			}
		}

		res := executor.Execute("/var/program.out", "inputs", info.Submission.Challenge.TimeLimit, logger)
		jsonResult, convertErr := json.Marshal(res)
		if convertErr != nil {
			msg := "An error occurred while trying to convert the execution result into a json format"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, convertErr)
			logger.Error().
				Err(convertErr).
				Msg(msg)
			return nil
		}
		return jsonResult
	})
	return m
}

func main() {
	logger, logFile := utils.InitLogger("executor")
	defer logFile.Close()
	client := graphql.NewClient("http://orm:8084/graphql", nil)
	m := createExecutorServer(client, logger)
	m.RunOnAddr(":8082")
}
