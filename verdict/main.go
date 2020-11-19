package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/go-martini/martini"
	model "github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/martini-contrib/binding"
)

var compilationError = errors.New("Compilation Error")

type VerdictForm struct {
	Language string                  `form:"language"`
	Source   *multipart.FileHeader   `form:"source"`
	Inputs   []*multipart.FileHeader `form:"inputs"`
	Outputs  []*multipart.FileHeader `form:"outputs"`
}

func createFileField(writer *multipart.Writer, fieldName string, file *multipart.FileHeader) error {
	field, err := writer.CreateFormFile(fieldName, file.Filename)
	if err != nil {
		return err
	}
	content, err := file.Open()
	if err != nil {
		return err
	}
	io.Copy(field, content)
	defer content.Close()
	return nil
}

func handleCompiling(language string, source *multipart.FileHeader) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	languageField, err := writer.CreateFormField("language")
	if err != nil {
		return nil, err
	}
	languageField.Write([]byte(language))

	err = createFileField(writer, "program", source)
	if err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://compiler:8080", buffer)
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
		return nil, compilationError
	}

	binary, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func handleExecute(binary string, inputs []*multipart.FileHeader) ([]model.ExecutionResult, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	binaryField, err := writer.CreateFormFile("binary", binary)
	if err != nil {
		return nil, err
	}
	binaryFile, err := os.Open(binary)
	if err != nil {
		return nil, err
	}
	defer binaryFile.Close()
	io.Copy(binaryField, binaryFile)

	for _, input := range inputs {
		err = createFileField(writer, "inputs", input)
		if err != nil {
			return nil, err
		}
	}

	writer.Close()

	req, err := http.NewRequest("POST", "http://executor:8080", buffer)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	executionResult := new([]model.ExecutionResult)
	err = json.NewDecoder(res.Body).Decode(executionResult)
	if err != nil {
		return nil, err
	}
	return *executionResult, nil
}

func compare(expectedOutput string, programOutput string) bool {
	return strings.EqualFold(programOutput, expectedOutput)
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(VerdictForm{}), func(rs http.ResponseWriter, rq *http.Request, f VerdictForm) string {
		binary, compilerErr := handleCompiling(f.Language, f.Source)
		if errors.Is(compilerErr, compilationError) {
			rs.WriteHeader(http.StatusOK)
			return "CE" // Compilation Error
		}
		if compilerErr != nil {
			rs.WriteHeader(http.StatusInternalServerError)
			rs.Write([]byte("An error occurred while trying to compile the file '" + f.Source.Filename + "' on the language '" + f.Language + "'\n"))
			return "Failed judgment"
		}
		writeErr := ioutil.WriteFile("binary", binary, 0777)
		if writeErr != nil {
			rs.WriteHeader(http.StatusInternalServerError)
			rs.Write([]byte("An error occurred while trying to create a local copy of the binary compilation of '" + f.Source.Filename + "'\n"))
			return "Failed judgment"
		}
		result, executorErr := handleExecute("binary", f.Inputs)
		if executorErr != nil {
			rs.WriteHeader(http.StatusInternalServerError)
			rs.Write([]byte("An error occurred while trying to execute the program with the received input files\n"))
			return "Failed judgment"
		}

		outputs := map[string]*multipart.FileHeader{}

		for _, out := range f.Outputs {
			outputName := out.Filename[:len(out.Filename)-len(".out")]
			outputs[outputName] = out
		}

		for _, testExecution := range result {
			if testExecution.Status != "OK" {
				return testExecution.Status + " " + testExecution.TestName
			}

			testName := testExecution.TestName[len("inputs/") : len(testExecution.TestName)-len(".in")]
			expectedOutputContent, err := outputs[testName].Open()
			if err != nil {
				rs.WriteHeader(http.StatusBadRequest)
				rs.Write([]byte("An error occurred while trying to open the output file named '" + testName + "'\n"))
				return "Failed judgment"
			}
			defer expectedOutputContent.Close()
			byteExpectedOutput, err := ioutil.ReadAll(expectedOutputContent)
			expectedOutput := string(byteExpectedOutput)
			if compare(testExecution.Message, expectedOutput) == false {
				return "WA" + " " + testExecution.TestName
			}
		}
		return "AC"
	})
	m.RunOnAddr(":8080")
}
