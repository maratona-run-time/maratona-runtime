package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"context"

	"github.com/go-martini/martini"
	executor "github.com/maratona-run-time/Maratona-Runtime/executor/src"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
	"github.com/martini-contrib/binding"
	"github.com/rs/zerolog"
	graphql "github.com/hasura/go-graphql-client"
)

// FileForm receives a submission ID
type FileForm struct {
	ID string `form:"id"`
}

func createExecutorServer(logger zerolog.Logger) *martini.ClassicMartini {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, req FileForm) []byte {
		fmt.Println(req.ID)
		client := graphql.NewClient("http://orm:8084/graphql", nil)
		var info struct {
			Submission struct {
				Challenge struct {
					TimeLimit float32
					Inputs    []struct {
						FileName string
						Content  []byte
					}
				}
			} `graphql:"submission(id: $id)"`
		}
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
	m := createExecutorServer(logger)
	m.RunOnAddr(":8082")
}
