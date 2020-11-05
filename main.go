package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	comparator "github.com/maratona-run-time/Maratona-Runtime/comparator/src"
	compiler "github.com/maratona-run-time/Maratona-Runtime/compiler/src"
	executor "github.com/maratona-run-time/Maratona-Runtime/executor/src"
)

func main() {
	language := "C"
	inputFileName := "in"
	outputFileName := "out"
	if len(os.Args[1:]) > 0 {
		language = os.Args[1]
	}
	if len(os.Args[2:]) > 0 {
		inputFileName = os.Args[2]
	}
	if len(os.Args[3:]) > 0 {
		outputFileName = os.Args[3]
	}
	path, errCompiler := compiler.Compile(language)
	if errCompiler != nil {
		fmt.Println("compiler error:", errCompiler)
		return
	}

	errorOutput := make(chan error)
	output := make(chan []byte)

	ctx, cancel := timerContext()
	defer cancel()

	go executor.Execute(ctx, path, inputFileName, output, errorOutput)
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
