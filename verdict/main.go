package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"

	"github.com/go-martini/martini"
	model "github.com/maratona-run-time/Maratona-Runtime/model"
	"github.com/maratona-run-time/Maratona-Runtime/utils"
	"github.com/maratona-run-time/Maratona-Runtime/verdict/src"
	"github.com/martini-contrib/binding"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var compilationError = errors.New("Compilation Error")

type VerdictForm struct {
	Language string                  `form:"language"`
	Source   *multipart.FileHeader   `form:"source"`
	Inputs   []*multipart.FileHeader `form:"inputs"`
	Outputs  []*multipart.FileHeader `form:"outputs"`
}

func handleCompiling(language string, source *multipart.FileHeader) ([]byte, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)

	fieldName := "language"
	err := utils.CreateFormField(writer, fieldName, language)
	if err != nil {
		return nil, err
	}

	fieldName = "source"
	err = utils.CreateFormFileFromFileHeader(writer, fieldName, source)
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

func handleExecute(binaryFilePath string, inputs []*multipart.FileHeader) ([]model.ExecutionResult, error) {
	buffer := new(bytes.Buffer)
	writer := multipart.NewWriter(buffer)
	fieldName := "binary"
	fileName := "binary"
	err := utils.CreateFormFileFromFilePath(writer, fieldName, fileName, binaryFilePath)
	if err != nil {
		return nil, err
	}

	for _, input := range inputs {
		err = utils.CreateFormFileFromFileHeader(writer, "inputs", input)
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

func main() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	consoleWriter := zerolog.ConsoleWriter{Out: os.Stdout}
	logFile, logErr := os.OpenFile("verdict.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
	defer logFile.Close()
	if logErr != nil {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
		log.Fatal().Err(logErr).Msg("Could not create log file")
	}
	multi := zerolog.MultiLevelWriter(consoleWriter, logFile)
	logger := zerolog.
		New(multi).
		With().
		Timestamp().
		Str("MaRT", "verdict").
		Logger().
		Level(zerolog.DebugLevel)

	m := martini.Classic()
	m.Post("/", binding.MultipartForm(VerdictForm{}), func(rs http.ResponseWriter, rq *http.Request, f VerdictForm) string {
		binary, compilerErr := handleCompiling(f.Language, f.Source)
		if errors.Is(compilerErr, compilationError) {
			rs.WriteHeader(http.StatusOK)
			logger.Debug().
				Msg("Compilation Error")
			return "CE"
		}
		if compilerErr != nil {
			msg := "Failed Judgment\nAn error occurred while trying to compile the file '" + f.Source.Filename + "' on the language '" + f.Language + "'"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, compilerErr)
			logger.Error().
				Err(compilerErr).
				Msg(msg)
			return ""
		}

		const binaryFileName = "binary"
		writeErr := ioutil.WriteFile(binaryFileName, binary, 0777)
		if writeErr != nil {
			msg := "Failed judgment\nAn error occurred while trying to create a local copy of the binary compilation of '" + f.Source.Filename + "'"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, writeErr)
			logger.Error().
				Err(writeErr).
				Msg(msg)
			return ""
		}
		results, executorErr := handleExecute(binaryFileName, f.Inputs)
		if executorErr != nil {
			msg := "Failed judgment\nAn error occurred while trying to execute the program with the received input files"
			utils.WriteResponse(rs, http.StatusInternalServerError, msg, executorErr)
			logger.Error().
				Err(executorErr).
				Msg(msg)
			return ""
		}

		outputs := map[string]*multipart.FileHeader{}

		for _, out := range f.Outputs {
			outputName := out.Filename[:len(out.Filename)-len(".out")]
			outputs[outputName] = out
		}

		result, err := verdict.Judge(results, outputs, logger)

		if err != nil {
			utils.WriteResponse(rs, http.StatusBadRequest, "", err)
			return ""
		}
		return result
	})
	m.RunOnAddr(":8080")
}
