package main

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	orm "github.com/maratona-run-time/Maratona-Runtime/orm/src"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

// ChallengeForm is a struct for receiving input and output files for a challenge via HTTP.
type ChallengeForm struct {
	Title       string                  `form:"title"`
	TimeLimit   int                     `form:"timeLimit"`
	MemoryLimit int                     `form:"memoryLimit"`
	Inputs      []*multipart.FileHeader `form:"inputs"`
	Outputs     []*multipart.FileHeader `form:"outputs"`
}

func writeChallenge(rs http.ResponseWriter, challenge model.Challenge) {
	jsonChallenge, err := json.Marshal(challenge)
	if err != nil {
		utils.WriteResponse(rs, http.StatusInternalServerError, "Error parsing challenge to JSON", err)
		return
	}
	rs.Header().Set("Content-Type", "application/json")
	rs.Write(jsonChallenge)
}

func parseRequestFiles(files []*multipart.FileHeader) ([]model.TestFile, error) {
	array := make([]model.TestFile, len(files))
	for i, file := range files {
		content, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer content.Close()
		byteInput, err := ioutil.ReadAll(content)
		if err != nil {
			return nil, err
		}
		array[i] = model.TestFile{Filename: file.Filename, Content: byteInput}
	}
	return array, nil
}

func createOrmServer() *martini.ClassicMartini {
	m := martini.Classic()

	m.Get("/challenge/:id", func(rs http.ResponseWriter, rq *http.Request, params martini.Params) {
		id := params["id"]
		challenge, err := orm.FindChallenge(id)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+string(id), err)
			return
		}
		writeChallenge(rs, challenge)
	})

	m.Post("/challenge", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm) {
		inputsArray, err := parseRequestFiles(f.Inputs)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access input files", err)
			return
		}
		outputsArray, err := parseRequestFiles(f.Outputs)
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
	return m
}

func main() {
	m := createOrmServer()
	m.RunOnAddr(":8080")
}
