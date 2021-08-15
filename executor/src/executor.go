package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/maratona-run-time/Maratona-Runtime/model"

	"github.com/rs/zerolog"
)

// Execute executes the compiled code (at path) of a given submission against a set of inputs (contained in inputsFolder).
// After executing all submissions the function returns an array containing the results of every execution.
// If the execution takes longer than `timeout` seconds, "TLE" is returned.
// If a run time error occurs, "RTE" is returned.
func Execute(path string,
	inputsFolder string,
	timeout float32,
	logger zerolog.Logger) []model.ExecutionResult {

	var files []string

	root := inputsFolder
	err := filepath.Walk(root, func(inputPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		basename := info.Name()
		if filepath.Ext(basename) == ".in" {
			files = append(files, inputPath)
		}
		return nil
	})
	if err != nil {
		logger.Error().Err(err)
	}

	var res []model.ExecutionResult

	for _, inputFileName := range files {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		defer cancel()

		errorOutput := make(chan error)
		output := make(chan []byte)

		executable := fmt.Sprintf("%s", path)

		file, fileErr := os.Open(inputFileName)
		if fileErr != nil {
			logger.Error().Err(fileErr)
		}
		defer file.Close()

		go execute(ctx, executable, file, output, errorOutput)

		select {
		case <-ctx.Done():
			res = append(res, model.ExecutionResult{TestName: inputFileName, Status: model.TIME_LIMIT_EXCEEDED, Message: "Time limit exceeded"})
			logger.Debug().Msg("Time limit exceeded")
			return res
		case err := <-errorOutput:
			res = append(res, model.ExecutionResult{TestName: inputFileName, Status: model.RUNTIME_ERROR, Message: err.Error()})
			logger.Debug().Msg("Runtime error")
			return res
		case out := <-output:
			res = append(res, model.ExecutionResult{TestName: inputFileName, Status: "OK", Message: string(out)})
		}
	}
	logger.Debug().Msg("Executions finished")
	return res
}

func execute(ctx context.Context, executable string, inputFile *os.File, output chan<- []byte, errorOutput chan<- error) {
	cmd := exec.CommandContext(ctx, executable)
	cmd.Stdin = inputFile
	programOutput, err := cmd.Output()
	if err != nil {
		errorOutput <- err
		return
	}
	output <- programOutput
}
