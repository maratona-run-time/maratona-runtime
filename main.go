package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/maratona-run-time/Maratona-Runtime/comparator"
	"github.com/maratona-run-time/Maratona-Runtime/compiler"
	"github.com/maratona-run-time/Maratona-Runtime/executor"
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

	statusChan := make(chan []string)
	timeout := float32(2)

	go executor.Execute(path, inputFileName, timeout, statusChan)

	status := <-statusChan

	switch status[0] {
	case "TLE":
		fmt.Println("TLE")
	case "RTE":
		fmt.Println("RTE")
		fmt.Println(status[1])
	case "OK":
		expectedData, _ := ioutil.ReadFile(outputFileName)
		expectedOutput := string(expectedData)
		programOutput := status[1]
		if comparator.Compare(expectedOutput, programOutput) {
			fmt.Println("AC")
		} else {
			fmt.Println("WA")
		}
	}
}
