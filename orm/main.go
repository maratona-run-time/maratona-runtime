package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/graphql-go/handler"
	"github.com/martini-contrib/binding"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	orm "github.com/maratona-run-time/Maratona-Runtime/orm/src"
	"github.com/maratona-run-time/Maratona-Runtime/queue"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

// ChallengeForm is a struct for receiving input and output files for a challenge via HTTP.
type ChallengeForm struct {
	Title       string                  `form:"title"`
	TimeLimit   float32                 `form:"timeLimit"`
	MemoryLimit int                     `form:"memoryLimit"`
	Inputs      []*multipart.FileHeader `form:"inputs"`
	Outputs     []*multipart.FileHeader `form:"outputs"`
}

var challengeNotFound = errors.New("Challenge not found")

func writeChallenge(rs http.ResponseWriter, challenge model.Challenge) {
	jsonChallenge, err := json.Marshal(challenge)
	if err != nil {
		utils.WriteResponse(rs, http.StatusInternalServerError, "Error parsing challenge to JSON", err)
		return
	}
	rs.Header().Set("Content-Type", "application/json")
	rs.Write(jsonChallenge)
}

func writeSubmission(rs http.ResponseWriter, submission model.Submission) {
	jsonSubmission, err := json.Marshal(submission)
	if err != nil {
		utils.WriteResponse(rs, http.StatusInternalServerError, "Error parsing submission to JSON", err)
		return
	}
	rs.Header().Set("Content-Type", "application/json")
	rs.Write(jsonSubmission)
}

func parseRequestFile(file *multipart.FileHeader) (string, []byte, error) {
	content, err := file.Open()
	if err != nil {
		return "", nil, err
	}
	defer content.Close()
	byteInput, err := ioutil.ReadAll(content)
	if err != nil {
		return "", nil, err
	}
	return file.Filename, byteInput, nil
}

func parseTestFiles(files []*multipart.FileHeader) ([]model.TestFile, error) {
	array := make([]model.TestFile, len(files))
	for i, file := range files {
		name, content, err := parseRequestFile(file)
		if err != nil {
			return nil, err
		}
		array[i] = model.TestFile{Filename: name, Content: content}
	}
	return array, nil
}

func challengeExists(challengeID uint) bool {
	_, err := orm.FindChallenge(challengeID)
	return err == nil
}

func createOrmServer() *martini.ClassicMartini {
	m := martini.Classic()

	m.Get("/challenge/:id", func(rs http.ResponseWriter, rq *http.Request, params martini.Params) {
		id, err := strconv.ParseUint(params["id"], 10, 64)
		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, "Challenge ID "+fmt.Sprint(id)+" must be a number", err)
			return
		}
		challenge, err := orm.FindChallenge(uint(id))
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+fmt.Sprint(id), err)
			return
		}
		writeChallenge(rs, challenge)
	})

	m.Post("/challenge", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm) {
		inputsArray, err := parseTestFiles(f.Inputs)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access input files", err)
			return
		}
		outputsArray, err := parseTestFiles(f.Outputs)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access output files", err)
			return
		}
		challenge := model.Challenge{Title: f.Title, TimeLimit: f.TimeLimit, MemoryLimit: f.MemoryLimit, Inputs: inputsArray, Outputs: outputsArray}
		err = orm.CreateChallenge(&challenge)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to create challenge", err)
			return
		}
		writeChallenge(rs, challenge)
	})

	m.Post("/submit", binding.MultipartForm(model.SubmissionForm{}), func(rs http.ResponseWriter, rq *http.Request, form model.SubmissionForm) {
		if !challengeExists(form.ChallengeID) {
			msg := fmt.Sprintf("Could not find challenge %v", form.ChallengeID)
			utils.WriteResponse(rs, http.StatusNotFound, msg, challengeNotFound)
		}
		_, content, err := parseRequestFile(form.Source)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access source file", err)
			return
		}
		submission := model.Submission{Language: form.Language, Source: content, ChallengeID: form.ChallengeID}
		err = orm.CreateSubmission(&submission)
		if err != nil {
			msg := fmt.Sprintf("Could not save submission")
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, err)
		}
		err = queue.Submit(fmt.Sprint(submission.ID))
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to queue submission ID", err)
			return
		}
		writeSubmission(rs, submission)
	})

	m.Get("/submission/:id", func(rs http.ResponseWriter, rq *http.Request, params martini.Params) {
		id, err := strconv.ParseUint(params["id"], 10, 64)
		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, "Submission ID "+fmt.Sprint(id)+" must be a number", err)
			return
		}
		submission, err := orm.FindSubmission(uint(id))
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find submission with id "+fmt.Sprint(id), err)
			return
		}
		writeSubmission(rs, submission)
	})
	return m
}

func main() {
	m := createOrmServer()
	h := handler.New(&handler.Config{
		Schema:   &orm.Schema,
		Pretty:   true,
		GraphiQL: true,
	})
	m.Any("/graphql", h.ServeHTTP)
	m.RunOnAddr(":8084")
}
