package comparator

import (
	"strings"
)

//Compare checks if the program output and the expected output are exactly the same.
func Compare(expectedOutput string, programOutput string) bool {
	return strings.EqualFold(programOutput, expectedOutput)
}
