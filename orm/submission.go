package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	orm "github.com/maratona-run-time/Maratona-Runtime/orm/src"
	"github.com/maratona-run-time/Maratona-Runtime/queue"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

func setSubmissionRoutes(m *martini.ClassicMartini) {
	m.Get("/submission/:id", func(rs http.ResponseWriter, rq *http.Request, params martini.Params) {
		id, err := strconv.ParseUint(params["id"], 10, 64)
		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, fmt.Sprintf("Submission ID %v must be a number", id), err)
			return
		}
		submission, err := orm.FindSubmission(uint(id))
		if err != nil {
			utils.WriteResponse(rs, http.StatusInternalServerError, fmt.Sprintf("Database error trying to find submission with id %v", id), err)
			return
		}
		writeJSONResponse(rs, submission)
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
		writeJSONResponse(rs, submission)
	})
}
