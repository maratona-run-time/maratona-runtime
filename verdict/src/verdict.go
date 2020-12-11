package verdict

import (
	"fmt"
	"io/ioutil"
	"mime/multipart"
	"strings"

	"github.com/maratona-run-time/Maratona-Runtime/model"

	"github.com/rs/zerolog"
)

func compare(expectedOutput string, programOutput string) bool {
	return strings.EqualFold(programOutput, expectedOutput)
}

// Judge compares the execution output of a submission and the expected output values.
// Returns a string with the judgment status ("AC", "WA", "TLE" or "RTE") and, possibly, an error.
func Judge(result []model.ExecutionResult, outputs map[string]*multipart.FileHeader, logger zerolog.Logger) (string, error) {
	for _, testExecution := range result {
		if testExecution.Status != "OK" {
			logger.Info().Msg("Judgment finished sentence " + testExecution.Status + " " + testExecution.TestName)
			return testExecution.Status + " " + testExecution.TestName, nil
		}

		testName := testExecution.TestName[len("inputs/") : len(testExecution.TestName)-len(".in")]
		expectedOutputContent, err := outputs[testName].Open()
		if err != nil {
			msg := "Failed judgment\nAn error occurred while trying to open the output file named '" + testName + "'"
			logger.Error().
				Err(err).
				Msg(msg)
			return "", fmt.Errorf(msg)
		}
		defer expectedOutputContent.Close()
		byteExpectedOutput, err := ioutil.ReadAll(expectedOutputContent)
		expectedOutput := string(byteExpectedOutput)
		if compare(testExecution.Message, expectedOutput) == false {
			logger.Info().Msg("Judgment finished sentence Wrong Answer")
			return "WA" + " " + testExecution.TestName, nil
		}
	}
	logger.Info().Msg("Judgment finished sentence Accepted")
	return "AC", nil
}
