package main

import (
	"Maratona-Runtime/comparator"
	"Maratona-Runtime/executor"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	executable := "a.out"
	inputFileName := "in"
	outputFileName := "out"
	if len(os.Args[1:]) > 0 {
		executable = os.Args[1]
	}
	if len(os.Args[2:]) > 0 {
		inputFileName = os.Args[2]
	}
	if len(os.Args[3:]) > 0 {
		outputFileName = os.Args[3]
	}
	executablePath := make(chan string)

	ctx, _ := timerContext()
	errorOutput := make(chan error)
	output := make(chan []byte)
	go executor.Execute(ctx, executablePath, inputFileName, output, errorOutput)
	executablePath <- executable // Compiler
	select {
	case <-ctx.Done():
		fmt.Println("TLE")
	case err := <-errorOutput:
		fmt.Println("RTE")
		fmt.Println(err)
	case out := <-output:
		expectedData, _ := ioutil.ReadFile(outputFileName)
		expectedOutput := string(expectedData)
		programOutput := string(out)
		if comparator.Compare(expectedOutput, programOutput) {
			fmt.Println("AC")
		} else {
			fmt.Println("WA")
		}
	}
}

func timerContext() (context.Context, context.CancelFunc) {
	return context.WithDeadline(context.Background(), time.Now().Add(time.Second*2))
}
