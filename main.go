package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"
)

func compare(expectedOutput string, programOutput string) bool {
	return strings.EqualFold(programOutput, expectedOutput)
}

func verdict(executable []string, inputFileName string, outputFileName string) (string, error) {
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
		if compare(expectedOutput, programOutput) {
			return "AC", nil
		} else {
			return "WA", nil
		}
	}
}

func main() {
	executable := "a.out"
	if len(os.Args[1:]) > 0 {
		executable = os.Args[1]
	}
	file := []string{fmt.Sprintf("./%s", executable)}

	// TODO: get input and output file names from command line
	status, err := verdict(file, "in", "out")
	fmt.Println(status)
	if err != nil {
		fmt.Println(err)
	}
}

func execute(ctx context.Context, executable []string, inputFile *os.File, output chan<- []byte, errorOutput chan<- error) {
	cmd := exec.CommandContext(ctx, executable[0], executable[1:]...)
	cmd.Stdin = inputFile
	fmt.Println("Pegando o output..")
	programOutput, err := cmd.Output() // Nao ta conseguindo pegar o output de um arquivo que tenha dado RTE, tirar o case da linha 23 resolve
	fmt.Println("pegou")
	fmt.Println(err)
	fmt.Println(programOutput)
	if err != nil {
		errorOutput <- err
		return
	}
	output <- programOutput
}

func timerContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
}
