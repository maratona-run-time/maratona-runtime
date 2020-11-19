package main

import (
	"fmt"
	"os"

	compiler "github.com/maratona-run-time/Maratona-Runtime/compiler/src"
	"github.com/maratona-run-time/Maratona-Runtime/verdict/src"
)

func main() {
	language := "C"
	inputFileName := "in"
	outputFileName := "out"
	sourceFileName := "program.c"
	if len(os.Args[1:]) > 0 {
		language = os.Args[1]
	}
	if len(os.Args[2:]) > 0 {
		inputFileName = os.Args[2]
	}
	if len(os.Args[3:]) > 0 {
		outputFileName = os.Args[3]
	}
	if len(os.Args[4:]) > 0 {
		sourceFileName = os.Args[4]
	}
	path, errCompiler := compiler.Compile(language, sourceFileName)
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
