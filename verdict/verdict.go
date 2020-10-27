package verdict

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/maratona-run-time/Maratona-Runtime/comparator"
	"github.com/maratona-run-time/Maratona-Runtime/executor"
)

func runTest(timeout float32, executablePath string, inputFileName string, outputFileName string, result chan<- string) {
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

func Verdict(timeout float32, executablePath string, inputFilesFolder string, outputFilesFolder string, result chan<- string) {
	var files [][2]string

	root := inputFilesFolder
	err := filepath.Walk(root, func(inputPath string, info os.FileInfo, err error) error {
		if info.IsDir() {
			return nil
		}
		basename := info.Name()
		if filepath.Ext(basename) == ".in" {
			testName := strings.TrimSuffix(basename, filepath.Ext(basename))
			outputPath := filepath.Join(outputFilesFolder, testName+".out")
			files = append(files, [2]string{inputPath, outputPath})
		}
		return nil
	})
	if err != nil {
		fmt.Println(err)
	}
	for _, testFiles := range files {
		testResult := make(chan string)
		go runTest(timeout, executablePath, testFiles[0], testFiles[1], testResult)
		switch <-testResult {
		case "TLE":
			result <- "TLE"
			return
		case "RTE":
			result <- "RTE"
			return
		case "WA":
			result <- "WA"
			return
		}
	}
	result <- "AC"
	return
}
