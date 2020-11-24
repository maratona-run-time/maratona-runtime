package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	model "github.com/maratona-run-time/Maratona-Runtime/model"
)

func Execute(path string,
	inputsFolder string,
	timeout float32) []model.ExecutionResult {

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
		fmt.Println(err)
	}

	var res []model.ExecutionResult

	for _, inputFileName := range files {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
		defer cancel()

		errorOutput := make(chan error)
		output := make(chan []byte)

		executable := fmt.Sprintf("./%s", path)

		file, fileErr := os.Open(inputFileName)
		if fileErr != nil {
			panic(fileErr)
		}
		defer file.Close()

		go execute(ctx, executable, file, output, errorOutput)

		select {
		case <-ctx.Done():
			res = append(res, model.ExecutionResult{inputFileName, "TLE", "Tempo limite excedido"})
			return res
		case err := <-errorOutput:
			res = append(res, model.ExecutionResult{inputFileName, "RTE", err.Error()})
			return res
		case out := <-output:
			res = append(res, model.ExecutionResult{inputFileName, "OK", string(out)})
		}
	}

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
