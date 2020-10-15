package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
)

func Execute(ctx context.Context, path string, inputFileName string, output chan<- []byte, errorOutput chan<- error) {
	inputFile, _ := os.Open(inputFileName)

	executable := fmt.Sprintf("./%s", path)

	go execute(ctx, executable, inputFile, output, errorOutput)
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
