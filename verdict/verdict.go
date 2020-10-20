package verdict

import (
	"io/ioutil"

	"github.com/maratona-run-time/Maratona-Runtime/comparator"
	"github.com/maratona-run-time/Maratona-Runtime/executor"
)

func Verdict(timeout float32, executablePath string, inputFileName string, outputFileName string, result chan<- string) {
	statusChan := make(chan []string)

	go executor.Execute(executablePath, inputFileName, timeout, statusChan)

	status := <-statusChan

	switch status[0] {
	case "TLE":
		result <- "TLE"
	case "RTE":
		result <- "RTE"
	case "OK":
		expectedData, _ := ioutil.ReadFile(outputFileName)
		expectedOutput := string(expectedData)
		programOutput := status[1]
		if comparator.Compare(expectedOutput, programOutput) {
			result <- "AC"
		} else {
			result <- "WA"
		}
	}
}
