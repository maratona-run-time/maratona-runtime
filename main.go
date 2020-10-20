package main

import (
	"fmt"
	"os"

	"github.com/maratona-run-time/Maratona-Runtime/compiler"
	"github.com/maratona-run-time/Maratona-Runtime/verdict"
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

	timeout := float32(2)
	resultChan := make(chan string)
	go verdict.Verdict(timeout, path, inputFileName, outputFileName, resultChan)

	result := <-resultChan
	fmt.Println(result)
}
