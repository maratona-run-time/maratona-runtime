package main

import (
	"bytes"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"testing"

	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

func createRequestForm(writer *multipart.Writer, language, filepath string) error {

	languageField, err := writer.CreateFormField("language")
	if err != nil {
		return err
	}
	_, err = languageField.Write([]byte(language))
	if err != nil {
		return err
	}

	content, err := ioutil.ReadFile(filepath)
	if err != nil {
		return err
	}
	field, err := writer.CreateFormFile("program", "source")
	if err != nil {
		return err
	}
	_, err = field.Write(content)
	return err
}

func createRequest(t *testing.T, language, filepath string) *http.Request {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	err := createRequestForm(writer, language, filepath)
	if err != nil {
		t.Error("could not create request form")
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/", buffer)
	if err != nil {
		t.Error("could not create request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

func TestCompilerServer(t *testing.T) {
	tests := []struct {
		name           string
		language       string
		filepath       string
		expectedOutput string
		expectedStatus int
	}{
		{
			name:           "C++/OK",
			language:       "C++11",
			filepath:       "src/tests/program.cpp",
			expectedOutput: "Hello World!\n",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Go/OK",
			language:       "Go",
			filepath:       "src/tests/program.go",
			expectedOutput: "Hello World!\n",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Python/OK",
			language:       "Python",
			filepath:       "src/tests/program.py",
			expectedOutput: "Hello World!\n",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "C/OK",
			language:       "C",
			filepath:       "src/tests/program.c",
			expectedOutput: "Hello World!",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "C++/CE",
			language:       "C++11",
			filepath:       "src/tests/compilation_error.cpp",
			expectedOutput: "Hello World!\n",
			expectedStatus: http.StatusBadRequest,
		},
	}

	logger := utils.InitDummyLogger()
	m := createCompilerServer(logger)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := createRequest(t, test.language, test.filepath)
			res := httptest.NewRecorder()
			m.ServeHTTP(res, req)
			if res.Code != test.expectedStatus {
				t.Errorf("expected status %v, got %v", test.expectedStatus, res.Code)
			}
			if res.Code != http.StatusOK {
				return
			}

			binary, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("could not read request response: %v", err.Error())
			}
			err = ioutil.WriteFile("executable", binary, 0777)
			if err != nil {
				t.Error("could not create executable file")
			}

			r, w, err := os.Pipe()
			if err != nil {
				t.Error("could not create pipe")
			}

			cmd := &exec.Cmd{
				Path:   "./executable",
				Args:   []string{},
				Stdout: w,
			}
			err = cmd.Run()
			if err != nil {
				t.Error("could not run executable")
			}
			w.Close()
			out, err := ioutil.ReadAll(r)
			if err != nil {
				t.Error("could not read executable output")
			}

			if string(out) != test.expectedOutput {
				t.Errorf("expected output to be %v, got %v", test.expectedOutput, string(out))
			}
		})
	}
}
