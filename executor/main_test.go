package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"path"
	"strconv"
	"testing"

	"github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
)

func resultEqual(a, b model.ExecutionResult) bool {
	if a.Status == "OK" {
		return a.Status == b.Status && a.TestName == b.TestName && a.Message == b.Message
	} else {
		return a.Status == b.Status && a.TestName == b.TestName
	}
}

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
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var inputs []Input
			for _, filepath := range test.inputPaths {
				binary, err := ioutil.ReadFile(filepath)
				if err != nil {
					panic("Could not read testfile from " + filepath)
				}
				inputs = append(inputs, Input{
					FileName: path.Base(filepath),
					Content:  binary,
				})
			}
			binary, err := ioutil.ReadFile(test.filePath)
			if err != nil {
				panic("Could not read executable file from " + test.filePath)
			}
			err = ioutil.WriteFile("/var/program.out", binary, 0777)
			if err != nil {
				panic("Could not write executable file to /var/program.out")
			}
			id := strconv.Itoa(rand.Int())
			var client utils.GraphqlMock = utils.GraphqlMock{
				Test: t,
				Object: Info{
					Submission: Submission{
						Challenge: Challenge{
							TimeLimit: 1.0,
							Inputs:    inputs,
						},
					},
				},
				Variables: map[string]interface{}{
					"id": id,
				},
			}
			m := createExecutorServer(client, logger)
			req := createRequest(t, id)
			res := httptest.NewRecorder()
			m.ServeHTTP(res, req)
			if res.Code != test.expectedStatus {
				t.Errorf("expected status %v, got %v", test.expectedStatus, res.Code)
			}
			if res.Code != http.StatusOK {
				return
			}

			jsonResult, err := ioutil.ReadAll(res.Body)
			if err != nil {
				t.Errorf("could not read request response: %v", err.Error())
			}
			var results []model.ExecutionResult
			err = json.Unmarshal(jsonResult, &results)
			if err != nil {
				t.Errorf("could not unmarshall execution result: %v", err.Error())
			}
			for i, result := range results {
				if !resultEqual(result, test.expectedResult[i]) {
					t.Errorf("expected %v, got %v", test.expectedResult[i], result)
				}
			}
		})
	}
}
