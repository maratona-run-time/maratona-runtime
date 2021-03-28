package main

import (
	"encoding/json"
	"github.com/go-martini/martini"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	orm "github.com/maratona-run-time/Maratona-Runtime/orm/src"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
	"github.com/martini-contrib/binding"
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

func setChallengeRoutes(m *martini.ClassicMartini) {
	m.Get("/challenge/:id", func(rs http.ResponseWriter, rq *http.Request, params martini.Params) {
		id := params["id"]
		challenge, err := orm.FindChallenge(id)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+string(id), err)
			return
		}
		writeChallenge(rs, challenge)
	})

	m.Put("/challenge/:id", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm, params martini.Params) {
		id := params["id"]
		challenge, err := orm.FindChallenge(id)

		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+string(id)+" to update", err)
			return
		}

		challenge.Title = f.Title
		challenge.TimeLimit = f.TimeLimit
		challenge.MemoryLimit = f.MemoryLimit

		files, err := parseRequestFiles(f.Inputs)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access input files", err)
			return
		}
		challenge.Inputs = (model.TestFileArray)(files).InputFiles()

		files, err = parseRequestFiles(f.Outputs)
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
		writeChallenge(rs, challenge)
	})

	m.Patch("/challenge/:id", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm, params martini.Params) {
		id := params["id"]
		challenge, err := orm.FindChallenge(id)

		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+string(id)+" to update", err)
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
			inputsArray, err := parseRequestFiles(f.Inputs)
			if err != nil {
				utils.WriteResponse(rs, http.StatusInternalServerError, "Error trying to access input files", err)
				return
			}
			challenge.Inputs = (model.TestFileArray)(inputsArray).InputFiles()
		}

		if len(f.Outputs) > 0 {
			outputsArray, err := parseRequestFiles(f.Outputs)
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
		inputs := (model.TestFileArray)(inputsArray).InputFiles()
		outputs := (model.TestFileArray)(outputsArray).OutputFiles()
		challenge := model.Challenge{Title: f.Title, TimeLimit: f.TimeLimit, MemoryLimit: f.MemoryLimit, Inputs: inputs, Outputs: outputs}
		err = orm.CreateChallenge(&challenge)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to create challenge", err)
			return
		}
		writeChallenge(rs, challenge)
	})

	m.Delete("/challenge/:id", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm, params martini.Params) {
		id := params["id"]
		err := orm.DeleteChallenge(id)
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, "Database error trying to find challenge with id "+string(id)+" to update", err)
			return
		}
		utils.WriteResponse(rs, http.StatusNoContent, "", nil)
	})

}
