package main

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	orm "github.com/maratona-run-time/Maratona-Runtime/orm/src"
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

func setChallengeRoutes(m *martini.ClassicMartini) {
	m.Get("/challenge", func(rs http.ResponseWriter, rq *http.Request) {
		challenges, err := orm.FindAllChallenges()
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find all challenges", err)
			return
		}
		writeJSONResponse(rs, challenges)
	})

	m.Get("/challenge/:id", func(rs http.ResponseWriter, rq *http.Request, params martini.Params) {
		id, err := strconv.ParseUint(params["id"], 10, 64)
		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, "Challenge ID "+fmt.Sprint(id)+" must be a number", err)
			return
		}
		challenge, err := orm.FindChallenge(uint(id))
		if err != nil {
			msg := fmt.Sprintf("Database error trying to find challenge with id %v", fmt.Sprint(id))
			utils.WriteResponse(rs, http.StatusNotFound, msg, challengeNotFound)
			return
		}
		writeJSONResponse(rs, challenge)
	})

	m.Put("/challenge/:id", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm, params martini.Params) {
		id, err := strconv.ParseUint(params["id"], 10, 64)
		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, fmt.Sprintf("Challenge ID %v must be a number", id), err)
			return
		}
		challenge, err := orm.FindChallenge(uint(id))

		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, fmt.Sprintf("Database error trying to find challenge with id %v to update", id), err)
			return
		}

		challenge.Title = f.Title
		challenge.TimeLimit = f.TimeLimit
		challenge.MemoryLimit = f.MemoryLimit

		files, err := parseTestFiles(f.Inputs)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access input files", err)
			return
		}
		challenge.Inputs = (model.TestFileArray)(files).InputFiles()

		files, err = parseTestFiles(f.Outputs)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access output files", err)
			return
		}
		challenge.Outputs = (model.TestFileArray)(files).OutputFiles()

		err = orm.UpdateChallenge(challenge)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to update challenge", err)
			return
		}
		writeJSONResponse(rs, challenge)
	})

	m.Patch("/challenge/:id", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm, params martini.Params) {
		id, err := strconv.ParseUint(params["id"], 10, 64)
		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, "Challenge ID "+fmt.Sprint(id)+" must be a number", err)
			return
		}
		challenge, err := orm.FindChallenge(uint(id))

		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+string(rune(id))+" to update", err)
			return
		}

		if len(f.Title) > 0 {
			challenge.Title = f.Title
		}
		if f.TimeLimit > 0 {
			challenge.TimeLimit = f.TimeLimit
		}
		if f.MemoryLimit > 0 {
			challenge.MemoryLimit = f.MemoryLimit
		}

		if len(f.Inputs) > 0 {
			inputsArray, err := parseTestFiles(f.Inputs)
			if err != nil {
				utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access input files", err)
				return
			}
			challenge.Inputs = (model.TestFileArray)(inputsArray).InputFiles()
		}

		if len(f.Outputs) > 0 {
			outputsArray, err := parseTestFiles(f.Outputs)
			if err != nil {
				utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access output files", err)
				return
			}
			challenge.Outputs = (model.TestFileArray)(outputsArray).OutputFiles()
		}

		err = orm.UpdateChallenge(challenge)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to patch challenge", err)
			return
		}
		writeJSONResponse(rs, challenge)
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
		inputs := (model.TestFileArray)(inputsArray).InputFiles()
		outputs := (model.TestFileArray)(outputsArray).OutputFiles()
		challenge := model.Challenge{Title: f.Title, TimeLimit: f.TimeLimit, MemoryLimit: f.MemoryLimit, Inputs: inputs, Outputs: outputs}
		err = orm.CreateChallenge(&challenge)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to create challenge", err)
			return
		}
		writeJSONResponse(rs, challenge)
	})

	m.Delete("/challenge/:id", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm, params martini.Params) {
		id, err := strconv.ParseUint(params["id"], 10, 64)
		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, "Challenge ID "+fmt.Sprint(id)+" must be a number", err)
			return
		}
		err = orm.DeleteChallenge(uint(id))
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, fmt.Sprintf("Database error trying to find challenge with id %v to update", id), err)
			return
		}
		utils.WriteResponse(rs, http.StatusNoContent, "", nil)
	})

}
