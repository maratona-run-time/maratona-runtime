package comparator

import (
	"strings"
)

func Compare(expectedOutput string, programOutput string) bool {
	return strings.EqualFold(programOutput, expectedOutput)
}
