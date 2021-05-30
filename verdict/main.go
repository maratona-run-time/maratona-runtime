package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-martini/martini"
	graphql "github.com/hasura/go-graphql-client"
	model "github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
	verdict "github.com/maratona-run-time/Maratona-Runtime/verdict/src"
	"github.com/martini-contrib/binding"
)

var compilationError = errors.New("Compilation Error")

// VerdictForm receives a submission ID
type VerdictForm struct {
	ID string `form:"id"`
}

func handleCompiling(id string) error {
	res, err := utils.MakeSubmissionRequest("http://localhost:8081", id)
	if err != nil {
		return err
	}

	if res.StatusCode != http.StatusOK {
		return compilationError
	}
	return nil
}

func handleExecute(id string) ([]model.ExecutionResult, error) {
	res, err := utils.MakeSubmissionRequest("http://localhost:8082", id)
	if err != nil {
		return nil, err
	}

	executionResult := new([]model.ExecutionResult)
	err = json.NewDecoder(res.Body).Decode(executionResult)
	if err != nil {
		return nil, err
	}
	return *executionResult, nil
}

func main() {
	logger, logFile := utils.InitLogger("verdict")
	defer logFile.Close()

	m := martini.Classic()
	m.Post("/", binding.MultipartForm(VerdictForm{}), func(rs http.ResponseWriter, rq *http.Request, req VerdictForm) string {
		client := graphql.NewClient("http://orm:8084/graphql", nil)
		var info struct {
			Submission struct {
				Challenge struct {
					Outputs []struct {
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
			return ""
		}

		compilerErr := handleCompiling(req.ID)
		if errors.Is(compilerErr, compilationError) {
			rs.WriteHeader(http.StatusOK)
			logger.Debug().
				Msg("Compilation Error")
			return "CE"
		}
		if compilerErr != nil {
			msg := "Failed Judgment\nAn error occurred while trying to compile the source file"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, compilerErr)
			logger.Error().
				Err(compilerErr).
				Msg(msg)
			return ""
		}

		results, executorErr := handleExecute(req.ID)
		if executorErr != nil {
			msg := "Failed judgment\nAn error occurred while trying to execute the program with the received input files"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, executorErr)
			logger.Error().
				Err(executorErr).
				Msg(msg)
			return ""
		}

		outputs := map[string]string{}

		for _, output := range info.Submission.Challenge.Outputs {
			outputName := output.FileName[:len(output.FileName)-len(".out")]
			outputs[outputName] = string(output.Content)
		}

		result, err := verdict.Judge(results, outputs, logger)

		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, "", err)
			return ""
		}
		return result
	})
	m.RunOnAddr(":8083")
}
