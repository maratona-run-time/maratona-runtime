package main

import (
	"bytes"
	"context"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"testing"

	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

func createRequest(t *testing.T, id string) *http.Request {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	err := utils.CreateFormField(writer, "id", id)
	if err != nil {
		panic("Could not create request form")
	}
	writer.Close()

	req, err := http.NewRequest("POST", "/", buffer)
	if err != nil {
		panic("Could not create request")
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	return req
}

type graphqlMock struct {
	test      *testing.T
	object    interface{}
	variables map[string]interface{}
}

func (gm graphqlMock) Query(ctx context.Context, info interface{}, variables map[string]interface{}) error {
	reflect.ValueOf(info).Elem().Set(reflect.ValueOf(gm.object))
	if !reflect.DeepEqual(gm.variables, variables) {
		gm.test.Errorf("Expect request variables to be %v, received %v", gm.variables, variables)
	}
	return nil
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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var client graphqlMock
			client.test = t
			binary, readErr := ioutil.ReadFile(test.filepath)
			if readErr != nil {
				panic("Could not read testfile from " + test.filepath)
			}
			client.object = submission{
				Submission: struct {
					Language string
					Source   []byte
				}{
					Language: test.language,
					Source:   binary,
				},
			}
			id := strconv.Itoa(rand.Int())
			client.variables = map[string]interface{}{
				"id": id,
			}
			m := createCompilerServer(client, logger)
			req := createRequest(t, id)
			res := httptest.NewRecorder()
			m.ServeHTTP(res, req)
			if res.Code != test.expectedStatus {
				t.Logf("request body: %v", res.Body)
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
