package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	executor "github.com/maratona-run-time/Maratona-Runtime/executor/src"
	"github.com/martini-contrib/binding"
)

// FileForm define o tipo de dados esperado no POST.
// Recebe um arquivo binário e um conjunto de arquivos de entrada.
type FileForm struct {
	Binary *multipart.FileHeader   `form:"binary"`
	Inputs []*multipart.FileHeader `form:"inputs"`
}

func main() {
	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, f FileForm) []byte {
		receivedFile, rErr := f.Binary.Open()
		if rErr != nil {
			rs.WriteHeader(http.StatusBadRequest)
			rs.Write([]byte("An error occurred while trying to open the binary file named '" + f.Binary.Filename + "'"))
			return nil
		}

		binaryFile, bErr := os.Create("program.out")
		if bErr != nil {
			rs.WriteHeader(http.StatusInternalServerError)
			rs.Write([]byte("An error occurred while trying to create a local empty file"))
			return nil
		}

		exeErr := os.Chmod("program.out", 0777)
		if exeErr != nil {
			rs.WriteHeader(http.StatusInternalServerError)
			rs.Write([]byte("An error occurred while trying to give execution permission to a local empty file"))
			return nil
		}

		_, copyErr := io.Copy(binaryFile, receivedFile)

		if copyErr != nil {
			rs.WriteHeader(http.StatusInternalServerError)
			rs.Write([]byte("An error occurred while trying to copy the received binary to a local file"))
			return nil
		}

		binaryFile.Close()
		receivedFile.Close()

		os.Mkdir("inputs", 0700)

		for _, file := range f.Inputs {
			if file == nil {
				rs.WriteHeader(http.StatusBadRequest)
				rs.Write([]byte("Received nil input file on the executor"))
				return nil
			}
			testFileName := fmt.Sprintf("inputs/%s", file.Filename)
			testFile, testFileErr := os.Create(testFileName)
			if testFileErr != nil {
				rs.WriteHeader(http.StatusBadRequest)
				rs.Write([]byte("An error occurred while trying to create a local file named '" + file.Filename + "' on 'inputs/' folder"))
				return nil
			}
			defer testFile.Close()
			receivedTestFile, rfErr := file.Open()
			if rfErr != nil {
				rs.WriteHeader(http.StatusBadRequest)
				rs.Write([]byte("An error occurred while trying to open the received test file named '" + file.Filename + "'"))
				return nil
			}
			defer receivedTestFile.Close()
			_, copyErr := io.Copy(testFile, receivedTestFile)
			if copyErr != nil {
				rs.WriteHeader(http.StatusInternalServerError)
				rs.Write([]byte("An error occurred while trying to copy the received test to a local file named '" + file.Filename + "' on 'inputs/' folder"))
				return nil
			}
		}

		res := executor.Execute("program.out", "inputs", 2.)
		jsonResult, convertErr := json.Marshal(res)
		if convertErr != nil {
			rs.WriteHeader(http.StatusInternalServerError)
			rs.Write([]byte("An error occurred while trying to convert the execution result into a json format"))
			return nil
		}
		return jsonResult
	})
	m.RunOnAddr(":8080")
}
