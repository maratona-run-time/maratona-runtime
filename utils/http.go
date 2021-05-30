package utils

import (
	"bytes"
	"fmt"
	"mime/multipart"
	"net/http"
)

// MakeSubmissionRequest calls path with a submission id on the request form
func MakeSubmissionRequest(path string, id string) (*http.Response, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	err := CreateFormField(writer, "id", id)
	if err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", path, buffer)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	return client.Do(req)
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
