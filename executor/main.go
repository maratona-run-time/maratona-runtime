package main

import (
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/maratona-run-time/Maratona-Runtime/utils"

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
	logger := utils.InitLogger("executor")

	m := martini.Classic()
	m.Post("/", binding.MultipartForm(FileForm{}), func(rs http.ResponseWriter, rq *http.Request, f FileForm) []byte {
		receivedFile, rErr := f.Binary.Open()
		if rErr != nil {
			msg := "An error occurred while trying to open the binary file named '" + f.Binary.Filename + "'"
			utils.WriteResponse(rs, http.StatusBadRequest, msg, rErr)
			logger.Error().
				Err(rErr).
				Msg(msg)
			return nil
		}

		binaryFile, bErr := os.Create("program.out")
		if bErr != nil {
			msg := "An error occurred while trying to create a local empty file"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, bErr)
			logger.Error().
				Err(bErr).
				Msg(msg)
			return nil
		}

		exeErr := os.Chmod("program.out", 0777)
		if exeErr != nil {
			msg := "An error occurred while trying to give execution permission to a local empty file"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, exeErr)
			logger.Error().
				Err(exeErr).
				Msg(msg)
			return nil
		}

		_, copyErr := io.Copy(binaryFile, receivedFile)
		if copyErr != nil {
			msg := "An error occurred while trying to copy the received binary to a local file"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, copyErr)
			logger.Error().
				Err(copyErr).
				Msg(msg)
			return nil
		}

		binaryFile.Close()
		receivedFile.Close()

		os.Mkdir("inputs", 0700)

		for _, file := range f.Inputs {
			if file == nil {
				msg := "Received nil input file on the executor"
				utils.WriteResponse(rs, http.StatusBadRequest, msg, nil)
				logger.Error().
					Msg(msg)
				return nil
			}
			testFileName := fmt.Sprintf("inputs/%s", file.Filename)
			testFile, testFileErr := os.Create(testFileName)
			if testFileErr != nil {
				msg := "An error occurred while trying to create a local file named '" + file.Filename + "' on 'inputs/' folder"
				utils.WriteResponse(rs, http.StatusBadRequest, msg, testFileErr)
				logger.Error().
					Err(testFileErr).
					Msg(msg)
				return nil
			}
			defer testFile.Close()
			receivedTestFile, rfErr := file.Open()
			if rfErr != nil {
				msg := "An error occurred while trying to open the received test file named '" + file.Filename + "'"
				utils.WriteResponse(rs, http.StatusBadRequest, msg, rfErr)
				logger.Error().
					Err(rfErr).
					Msg(msg)
				return nil
			}
			defer receivedTestFile.Close()
			_, copyErr := io.Copy(testFile, receivedTestFile)
			if copyErr != nil {
				msg := "An error occurred while trying to copy the received test to a local file named '" + file.Filename + "' on 'inputs/' folder"
				utils.WriteResponse(rs, http.StatusInternalServerError, msg, copyErr)
				logger.Error().
					Err(copyErr).
					Msg(msg)
				return nil
			}
		}

		res := executor.Execute("program.out", "inputs", 2., logger)
		jsonResult, convertErr := json.Marshal(res)
		if convertErr != nil {
			msg := "An error occurred while trying to convert the execution result into a json format"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, convertErr)
			logger.Error().
				Err(convertErr).
				Msg(msg)
			return nil
		}
		return jsonResult
	})
	m.RunOnAddr(":8080")
}
