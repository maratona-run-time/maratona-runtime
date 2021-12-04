package utils

import (
	"bytes"
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"time"

	"github.com/hasura/go-graphql-client"
)

const REQUEST_RETRIES = 10
const RETRY_INTERVAL = 2 * time.Second

// MakeSubmissionRequest calls path with a submission id on the request form
func MakeSubmissionRequest(path string, id string) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	err := CreateFormField(writer, "id", id)
	if err != nil {
		return nil, err
	}
	writer.Close()

	retry_number := 0
	for retry_number < REQUEST_RETRIES {
		var res *http.Response
		var req *http.Request
		req, err = http.NewRequest("POST", path, buffer)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", writer.FormDataContentType())
		client := &http.Client{}
		res, err = client.Do(req)
		if err == nil {
			return res, err
		}
		time.Sleep(RETRY_INTERVAL)
		retry_number++
	}
	return nil, err
}

func SaveSubmissionStatus(client *graphql.Client, id, verdict, message string) error {
	var judgeMutation struct {
		Judge struct {
			Verdict string
		} `graphql:"judge(submissionID: $id, verdict: $verdict, message: $message)"`
	}
	variables := map[string]interface{}{
		"id":      id,
		"verdict": graphql.String(verdict),
		"message": graphql.String(message),
	}
	return client.Mutate(context.Background(), &judgeMutation, variables)
}

// WriteResponse is used to write a HTTP response status in case of an error.
func WriteResponse(rs http.ResponseWriter, status int, msg string, err error) {
	rs.WriteHeader(status)
	errorMsg := msg
	if err != nil {
		errorMsg = fmt.Sprintf("%v:\n%v", msg, err.Error())
	}
	rs.Write([]byte(errorMsg))
}
