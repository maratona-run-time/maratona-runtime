package executor

import (
	"context"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
	"Maratona-Runtime/comparator"
)

func Verdict(executable []string, inputFileName string, outputFileName string) (string, error) {
	actualOutput := make(chan []byte)
	errorOutput := make(chan error)
	ctx, _ := timerContext()

	inputFile, _ := os.Open(inputFileName)

	go execute(ctx, executable, inputFile, actualOutput, errorOutput)
	select {
	case <-ctx.Done():
		return "TLE", nil
	case err := <-errorOutput:
		return "RTE", err
	case out := <-actualOutput:
		expectedData, _ := ioutil.ReadFile(outputFileName)
		expectedOutput := string(expectedData)
		programOutput := string(out)
		if comparator.Compare(expectedOutput, programOutput) {
			return "AC", nil
		} else {
			return "WA", nil
		}
	}
}

func execute(ctx context.Context, executable []string, inputFile *os.File, output chan<- []byte, errorOutput chan<- error) {
	cmd := exec.CommandContext(ctx, executable[0], executable[1:]...)
	cmd.Stdin = inputFile
	programOutput, err := cmd.Output()
	if err != nil {
		errorOutput <- err
		return
	}
	output <- programOutput
}

func timerContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
}
