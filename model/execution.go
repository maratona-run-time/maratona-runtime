package model

import (
	"gorm.io/gorm"
	"mime/multipart"
)

type ExecutionResult struct {
	TestName string `json:"testName"`
	Status   string `json:"status"`
	Message  string `json:"message"`
}

type SubmissionForm struct {
	Language    string                `form:"language"`
	Source      *multipart.FileHeader `form:"source"`
	ChallengeID string                `form:"challengeID"`
}

type Challenge struct {
	gorm.Model
	Title       string
	TimeLimit   int
	MemoryLimit int
	Inputs      []TestFile `gorm:"ForeignKey:ChallengeID"`
	Outputs     []TestFile `gorm:"ForeignKey:ChallengeID"`
}

type TestFile struct {
	gorm.Model
	Filename    string
	Content     []byte
	ChallengeID uint
}
