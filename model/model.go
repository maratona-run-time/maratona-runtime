package model

import (
	"mime/multipart"

	"gorm.io/gorm"
)

// ExecutionResult represents the result of executing a given submission against one test case.
// TestName is the name of the test case, Status is the status of the test (i.e. "OK", "TLE" nad "RTE")
// and Message is used to store any relevant information, such as error messages.
type ExecutionResult struct {
	TestName string `json:"testName"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

// SubmissionForm is used to represent a submission.
// It contains the source code, a string identifiying the language of the source code and the ID of the challenge associated with this submission.
type SubmissionForm struct {
	Language    string                `form:"language"`
	Source      *multipart.FileHeader `form:"source"`
	ChallengeID uint                  `form:"challengeID"`
}

// Challenge is a representation for the ORM of our challenges and their relevant information such as
// the challenge title, the memory and time constraints for the submissions and the set of input and outputs that represent the test cases.
type Challenge struct {
	gorm.Model
	Title       string
	TimeLimit   int
	MemoryLimit int
	Inputs      []InputFile  `gorm:"ForeignKey:ChallengeID"`
	Outputs     []OutputFile `gorm:"ForeignKey:ChallengeID"`
}

type Submission struct {
	gorm.Model
	Language    string
	Source      []byte
	ChallengeID uint
}

// TestFile represents the input of a single test case for a given challenge.
type TestFile struct {
	gorm.Model
	Filename    string
	Content     []byte
	ChallengeID uint
}

type TestFileArray []TestFile

func (files TestFileArray) InputFiles() []InputFile {
	inputs := make([]InputFile, len(files))
	for index, file := range files {
		inputs[index] = InputFile{
			TestFile: file,
		}
	}
	return inputs
}

func (files TestFileArray) OutputFiles() []OutputFile {
	outputs := make([]OutputFile, len(files))
	for index, file := range files {
		outputs[index] = OutputFile{
			TestFile: file,
		}
	}
	return outputs
}

type InputFile struct {
	TestFile
}

type InputsArray []InputFile

func (files InputsArray) TestFiles() []TestFile {
	testFiles := make([]TestFile, len(files))
	for index, file := range files {
		testFiles[index] = file.TestFile
	}
	return testFiles
}

type OutputFile struct {
	TestFile
}

type OutputsArray []OutputFile

func (files OutputsArray) TestFiles() []TestFile {
	testFiles := make([]TestFile, len(files))
	for index, file := range files {
		testFiles[index] = file.TestFile
	}
	return testFiles
}
