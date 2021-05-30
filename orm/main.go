package main

import (
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/go-martini/martini"

	orm "github.com/maratona-run-time/Maratona-Runtime/orm/src"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

func writeJSONResponse(rs http.ResponseWriter, response interface{}) {
	jsonResponse, err := json.Marshal(response)
	if err != nil {
		utils.WriteResponse(rs, http.StatusInternalServerError, "Error parsing response to JSON", err)
		return
	}
	rs.Header().Set("Content-Type", "application/json")
	rs.Write(jsonResponse)
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

func challengeExists(challengeID uint) bool {
	_, err := orm.FindChallenge(challengeID)
	return err == nil
}

func createOrmServer() *martini.ClassicMartini {
	m := martini.Classic()
	setChallengeRoutes(m)
	setSubmissionRoutes(m)
	return m
}

func main() {
	m := createOrmServer()
	m.RunOnAddr(":8084")
}
