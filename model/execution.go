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
	ChallengeID string                `form:"challengeID"`
}

// Challenge is a representation for the ORM of our challenges and their relevant information such as
// the challenge title, the memory and time constraints for the submissions and the set of input and outputs that represent the test cases.
type Challenge struct {
	gorm.Model
	Title       string
	TimeLimit   int
	MemoryLimit int
	Inputs      []TestFile `gorm:"ForeignKey:ChallengeID"`
	Outputs     []TestFile `gorm:"ForeignKey:ChallengeID"`
}

// TestFile represents the input of a single test case for a given challenge.
type TestFile struct {
	gorm.Model
	Filename    string
	Content     []byte
	ChallengeID uint
}
