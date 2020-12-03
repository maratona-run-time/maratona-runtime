package main

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"

	"github.com/maratona-run-time/Maratona-Runtime/errors"
	"github.com/maratona-run-time/Maratona-Runtime/orm/src"
)

type SubmissionForm struct {
	Language  string                `form:"language"`
	Source    *multipart.FileHeader `form:"source"`
	ProblemId int                   `form:"problemId"`
}

type ChallengeForm struct {
	Title       string                  `form:"title"`
	TimeLimit   int                     `form:"timeLimit"`
	MemoryLimit int                     `form:"memoryLimit"`
	Inputs      []*multipart.FileHeader `form:"inputs"`
}

func writeChallenge(rs http.ResponseWriter, challenge orm.Challenge) {
	jsonChallenge, err := json.Marshal(challenge)
	if err != nil {
		errors.WriteResponse(rs, http.StatusInternalServerError, "Error parsing challenge to JSON", err)
		return
	}
	rs.Write(jsonChallenge)
}

func createOrmServer() *martini.ClassicMartini {
	m := martini.Classic()

	m.Get("/challenge/:id", func(rs http.ResponseWriter, rq *http.Request, params martini.Params) {
		id := params["id"]
		challenge, err := orm.FindChallenge(id)
		if err != nil {
			errors.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+string(id), err)
			return
		}
		writeChallenge(rs, challenge)
	})

	m.Post("/challenge", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm) {
		inputsArray := make([]orm.Input, len(f.Inputs))
		for i, input := range f.Inputs {
			inputContent, err := input.Open()
			if err != nil {
				errors.WriteResponse(rs, http.StatusInternalServerError, "Error trying to open input files", err)
				return
			}
			defer inputContent.Close()
			byteInput, err := ioutil.ReadAll(inputContent)
			if err != nil {
				errors.WriteResponse(rs, http.StatusInternalServerError, "Error trying to read input files", err)
				return
			}
			inputsArray[i] = orm.Input{Filename: input.Filename, Content: byteInput}
		}
		challenge := orm.Challenge{Title: f.Title, TimeLimit: f.TimeLimit, MemoryLimit: f.MemoryLimit, Inputs: inputsArray}
		err := orm.CreateChallenge(&challenge)
		if err != nil {
			errors.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to create challenge", err)
			return
		}
		writeChallenge(rs, challenge)
	})
	return m
}

func main() {
	m := createOrmServer()
	m.RunOnAddr(":8080")
}
