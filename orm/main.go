package main

import (
	"encoding/json"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/go-martini/martini"
	"github.com/maratona-run-time/Maratona-Runtime/orm/src"
	"github.com/martini-contrib/binding"
)

type SubmissionForm struct {
	Language  string                `form:"language"`
	Source    *multipart.FileHeader `form:"source"`
	ProblemId int                   `form:"problemId"`
}

type ChallengeForm struct {
	Title       string `form:"title"`
	TimeLimit   int    `form:"timeLimit"`
	MemoryLimit int    `form:"memoryLimit"`
	//SolutionBinary // Lang
	//ComparatorBinary
	//Inputs      []*multipart.FileHeader `form:"inputs"`
	//Outputs     []*multipart.FileHeader `form:"outputs"`
}

func createOrmServer() *martini.ClassicMartini {
	m := martini.Classic()

	m.Get("/challenge/:id", func(params martini.Params) []byte {
		id := params["id"]
		c := orm.FindChallenge(id)
		jsonChallenge, _ := json.Marshal(c)
		return jsonChallenge
	})

	m.Post("/challenge", binding.MultipartForm(ChallengeForm{}), func(rs http.ResponseWriter, rq *http.Request, f ChallengeForm) string {
		challenge := orm.Challenge{Title: f.Title, TimeLimit: f.TimeLimit, MemoryLimit: f.MemoryLimit}
		orm.CreateChallenge(&challenge)
		return strconv.Itoa(int(challenge.ID))
	})

	m.Get("/test", func() {
		orm.Test()
	})

	return m
}

func main() {
	m := createOrmServer()
	m.RunOnAddr(":8085")
}
