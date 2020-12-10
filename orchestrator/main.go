package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"net/http"

	"github.com/go-martini/martini"
	model "github.com/maratona-run-time/Maratona-Runtime/model"
	utils "github.com/maratona-run-time/Maratona-Runtime/utils"
	"github.com/martini-contrib/binding"
)

var getChallengeError = errors.New("Error getting challenge")
var verdictResponseError = errors.New("Error on verdict response")

func createTestFileField(writer *multipart.Writer, fieldName string, files []model.TestFile) error {
	for _, file := range files {
		err := utils.CreateFormFileFromContent(writer, fieldName, file.Content, file.Filename)
		if err != nil {
			return err
		}
	}
	return nil
}

func getChallengeInfo(challengeID string) (model.Challenge, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://orm:8080/challenge/%v", challengeID), new(bytes.Buffer))
	if err != nil {
		return model.Challenge{}, err
	}
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return model.Challenge{}, err
	}

	if res.StatusCode != http.StatusOK {
		return model.Challenge{}, getChallengeError
	}

	var challenge model.Challenge
	binary, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return model.Challenge{}, err
	}
	err = json.Unmarshal(binary, &challenge)
	return challenge, err
}

func callVerdict(challenge model.Challenge, form model.SubmissionForm) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	fieldName := "language"
	err := utils.CreateFormField(writer, fieldName, form.Language)
	if err != nil {
		return nil, err
	}

	err = utils.CreateFormFileFromFileHeader(writer, "source", form.Source)
	if err != nil {
		return nil, err
	}

	err = createTestFileField(writer, "inputs", challenge.Inputs)
	if err != nil {
		return nil, err
	}

	err = createTestFileField(writer, "outputs", challenge.Outputs)
	if err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://verdict:8080", buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode != http.StatusOK {
		return nil, verdictResponseError
	}

	binary, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(model.SubmissionForm{}), func(rs http.ResponseWriter, rq *http.Request, form model.SubmissionForm) {
		challenge, err := getChallengeInfo(form.ChallengeID)
		if err != nil {
			msg := fmt.Sprintf("Could not find challenge %v", form.ChallengeID)
			utils.WriteResponse(rs, http.StatusBadRequest, msg, err)
		}

		verdictResponse, err := callVerdict(challenge, form)
		if err != nil {
			msg := "Could not get response from verdict"
			utils.WriteResponse(rs, http.StatusBadRequest, msg, err)
			return
		}

		rs.Write(verdictResponse)
	})
	m.RunOnAddr(":8080")
}
