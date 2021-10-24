package main

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"

	graphql "github.com/hasura/go-graphql-client"
	model "github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
	verdict "github.com/maratona-run-time/Maratona-Runtime/verdict/src"
)

var compilationError = errors.New("Compilation Error")

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

	submissionID := os.Getenv("SUBMISSION_ID")
	var status struct {
		verdict, message string
	}
	client := graphql.NewClient("http://orm:8084/graphql", nil)
	defer func() {
		err := utils.SaveSubmissionStatus(client, submissionID, status.verdict, status.message)
		if err != nil {
			msg := "An error occurred while trying to save submission '" + submissionID + "' verdict"
			logger.Error().
				Err(err).
				Msg(msg)
		}
	}()
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
		"id": graphql.ID(submissionID),
	}
	graphqlErr := client.Query(context.Background(), &info, variables)
	if graphqlErr != nil {
		msg := "An error occurred while trying to fetch submission '" + submissionID + "' details"
		logger.Error().
			Err(graphqlErr).
			Msg(msg)

		status.verdict = model.REJECTED
		status.message = msg
		return
	}

	compilerErr := handleCompiling(submissionID)
	if errors.Is(compilerErr, compilationError) {
		logger.Debug().
			Msg("Compilation Error")
		status.verdict = model.COMPILATION_ERROR
		return
	}
	if compilerErr != nil {
		msg := "Failed Judgment\nAn error occurred while trying to compile the source file"
		logger.Error().
			Err(compilerErr).
			Msg(msg)
		status.verdict = model.REJECTED
		status.message = msg
		return
	}

	results, executorErr := handleExecute(submissionID)
	if executorErr != nil {
		msg := "Failed judgment\nAn error occurred while trying to execute the program with the received input files"
		logger.Error().
			Err(executorErr).
			Msg(msg)
		status.verdict = model.REJECTED
		status.message = msg
		return
	}

	outputs := map[string]string{}
	for _, output := range info.Submission.Challenge.Outputs {
		outputName := output.FileName[:len(output.FileName)-len(".out")]
		outputs[outputName] = string(output.Content)
	}

	var err error
	status.verdict, status.message, err = verdict.Judge(results, outputs, logger)
	if err != nil {
		msg := "Failed judgment\nAn error occurred while trying to check the program output with the expected answers."
		logger.Error().
			Err(executorErr).
			Msg(msg)
		status.verdict = model.REJECTED
		status.message = msg
		return
	}
}
