package main

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/binding"
)

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

	binary, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return binary, nil
}

func handleExecute(binary string, inputs []*multipart.FileHeader) ([][]string, error) {
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

	stringResult := new([][]string)
	err = json.NewDecoder(res.Body).Decode(stringResult)
	if err != nil {
		return nil, err
	}
	return *stringResult, nil
}

func compare(expectedOutput string, programOutput string) bool {
	return strings.EqualFold(programOutput, expectedOutput)
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(VerdictForm{}), func(f VerdictForm) string {
		binary, compilerErr := handleCompiling(f.Language, f.Source)
		if compilerErr != nil {
			panic(compilerErr)
		}
		writeErr := ioutil.WriteFile("binary", binary, 0777)
		if writeErr != nil {
			panic(writeErr)
		}
		result, executorErr := handleExecute("binary", f.Inputs)
		if executorErr != nil {
			panic(executorErr)
		}

		outputs := map[string]*multipart.FileHeader{}

		for _, out := range f.Outputs {
			outputName := out.Filename[:len(out.Filename)-len(".out")]
			outputs[outputName] = out
		}

		for _, testExecution := range result {
			if testExecution[1] != "OK" {
				return testExecution[1] + " " + testExecution[0]
			}

			testName := testExecution[0][len("inputs/") : len(testExecution[0])-len(".in")]
			expectedOutputContent, err := outputs[testName].Open()
			if err != nil {
				panic(err)
			}
			defer expectedOutputContent.Close()
			byteExpectedOutput, err := ioutil.ReadAll(expectedOutputContent)
			expectedOutput := string(byteExpectedOutput)
			if compare(testExecution[2], expectedOutput) == false {
				return "WA" + " " + testExecution[0]
			}
		}
		return "AC"
	})
	m.RunOnAddr(":8080")
}
