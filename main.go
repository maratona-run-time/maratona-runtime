package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"
	"Maratona-Runtime/comparator"
	"Maratona-Runtime/executor"
)

func main() {
	executable := "a.out"
	if len(os.Args[1:]) > 0 {
		executable = os.Args[1]
	}
	executablePath := make(chan string)

	// TODO: get input and output file names from command line
	ctx, _ := timerContext()
	errorOutput := make(chan error)
	output := make(chan []byte)
	go executor.Execute(ctx, executablePath, "in", output, errorOutput)
	executablePath <- executable // Compiler
	select {
	case <-ctx.Done():
		fmt.Println("TLE")
	case err := <-errorOutput:
		fmt.Println("RTE")
		fmt.Println(err)
	case out := <-output:
		expectedData, _ := ioutil.ReadFile("out")
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
