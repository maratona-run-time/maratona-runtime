package verdict

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/maratona-run-time/Maratona-Runtime/comparator/src"
	"github.com/maratona-run-time/Maratona-Runtime/executor/src"
)

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

	res := executor.Execute(executablePath, inputFilesFolder, timeout)

	for _, executionResult := range res {
		_, fileName := path.Split(executionResult[0])
		testResult := executionResult[1]
		programOutput := executionResult[2]
		switch testResult {
		case "TLE":
			result <- "TLE"
			return
		case "RTE":
			result <- "RTE"
			return
		case "OK":
			testName := strings.TrimSuffix(fileName, filepath.Ext(fileName))
			outputFileName := filepath.Join(outputFilesFolder, testName+".out")
			expectedData, err := ioutil.ReadFile(outputFileName)
			if err != nil {
				fmt.Println(err)
			}
			expectedOutput := string(expectedData)
			if !comparator.Compare(expectedOutput, programOutput) {
				result <- "WA"
				return
			}
		}
	}
	result <- "AC"
	return
}
