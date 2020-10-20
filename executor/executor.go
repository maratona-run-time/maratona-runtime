package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"
)

func Execute(path string,
	inputFileName string,
	timeout float32,
	status chan<- []string) {
	inputFile, errInFile := os.Open(inputFileName)
	if errInFile != nil {
		fmt.Println(errInFile)
		return
	}

	executable := fmt.Sprintf("./%s", path)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeout))
	defer cancel()

	errorOutput := make(chan error)
	output := make(chan []byte)

	go execute(ctx, executable, inputFile, output, errorOutput)

	select {
	case <-ctx.Done():
		status <- []string{"TLE"}
	case err := <-errorOutput:
		status <- []string{"RTE", err.Error()}
	case out := <-output:
		status <- []string{"OK", string(out)}
	}
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
