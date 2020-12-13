package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path"
	"testing"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
	"github.com/rs/zerolog/log"
)

func resultEqual(a, b model.ExecutionResult) bool {
	if a.Status == "OK" {
		return a.Status == b.Status && a.TestName == b.TestName && a.Message == b.Message
	} else {
		return a.Status == b.Status && a.TestName == b.TestName
	}
}

func createRequestForm(writer *multipart.Writer, filePath string, inputPaths []string) error {
	fieldName := "binary"
	fileName := path.Base(filePath)
	err := utils.CreateFormFileFromFilePath(writer, fieldName, fileName, filePath)
	if err != nil {
		return err
	}
	for _, inputPath := range inputPaths {
		fieldName = "inputs"
		fileName = path.Base(inputPath)
		err = utils.CreateFormFileFromFilePath(writer, fieldName, fileName, inputPath)
		if err != nil {
			return err
		}
	}
	return nil
}

func createRequest(t *testing.T, filePath string, inputPaths []string) *http.Request {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	err := createRequestForm(writer, filePath, inputPaths)
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

func cleanUp() {

	dir, err := ioutil.ReadDir("inputs")
	if err != nil {
		log.Error().
			Err(err).
			Msg("Error reading the contents of the 'inputs' directory")
	}
	for _, d := range dir {
		err = os.Remove(path.Join([]string{"inputs", d.Name()}...))
		if err != nil {
			log.Error().
				Err(err).
				Msg("Error deleting file " + d.Name())
		}
	}
}

func TestExecutorServer(t *testing.T) {
	tests := []struct {
		name           string
		filePath       string
		inputPaths     []string
		expectedResult []model.ExecutionResult
		expectedStatus int
	}{
		{
			name:     "Sum/OK",
			filePath: "src/tests/ok.sh",
			inputPaths: []string{
				"../verdict/src/tests/sum/inputs/1.in",
				"../verdict/src/tests/sum/inputs/2.in",
				"../verdict/src/tests/sum/inputs/3.in",
				"../verdict/src/tests/sum/inputs/4.in",
				"../verdict/src/tests/sum/inputs/5.in",
			},
			expectedResult: []model.ExecutionResult{
				{
					TestName: "inputs/1.in",
					Status:   "OK",
					Message:  "2\n",
				},
				{
					TestName: "inputs/2.in",
					Status:   "OK",
					Message:  "3\n",
				},
				{
					TestName: "inputs/3.in",
					Status:   "OK",
					Message:  "20\n",
				},
				{
					TestName: "inputs/4.in",
					Status:   "OK",
					Message:  "1001000000\n",
				},
				{
					TestName: "inputs/5.in",
					Status:   "OK",
					Message:  "20000000000\n",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Sum/TLE",
			filePath: "src/tests/tle.sh",
			inputPaths: []string{
				"../verdict/src/tests/sum/inputs/1.in",
				"../verdict/src/tests/sum/inputs/2.in",
				"../verdict/src/tests/sum/inputs/3.in",
				"../verdict/src/tests/sum/inputs/4.in",
				"../verdict/src/tests/sum/inputs/5.in",
			},
			expectedResult: []model.ExecutionResult{
				{
					TestName: "inputs/1.in",
					Status:   "OK",
					Message:  "2\n",
				},
				{
					TestName: "inputs/2.in",
					Status:   "OK",
					Message:  "3\n",
				},
				{
					TestName: "inputs/3.in",
					Status:   "OK",
					Message:  "20\n",
				},
				{
					TestName: "inputs/4.in",
					Status:   "TLE",
					Message:  "Tempo limite excedido",
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "Sum/RTE",
			filePath: "src/tests/rte.sh",
			inputPaths: []string{
				"../verdict/src/tests/sum/inputs/1.in",
			},
			expectedResult: []model.ExecutionResult{
				{
					TestName: "inputs/1.in",
					Status:   "RTE",
					Message:  "exit status 1",
				},
			},
			expectedStatus: http.StatusOK,
		},
	}

	logger := utils.InitDummyLogger()
	m := createExecutorServer(logger)
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			req := createRequest(t, test.filePath, test.inputPaths)
			res := httptest.NewRecorder()
			m.ServeHTTP(res, req)
			if res.Code != test.expectedStatus {
				cleanUp()
				t.Errorf("expected status %v, got %v", test.expectedStatus, res.Code)
			}
			if res.Code != http.StatusOK {
				cleanUp()
				return
			}

			jsonResult, err := ioutil.ReadAll(res.Body)
			if err != nil {
				cleanUp()
				t.Errorf("could not read request response: %v", err.Error())
			}
			var results []model.ExecutionResult
			err = json.Unmarshal(jsonResult, &results)
			if err != nil {
				cleanUp()
				t.Errorf("could not unmarshall execution result: %v", err.Error())
			}
			for i, result := range results {
				if !resultEqual(result, test.expectedResult[i]) {
					cleanUp()
					t.Errorf("expected %v, got %v", test.expectedResult[i], result)
				}
			}
		})
	}
	cleanUp()
}
