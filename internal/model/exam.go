package model

import (
	"encoding/json"
	"time"
)

type SubmitExamRequest struct {
	Subject    string        `json:"subject"`
	ExamName   string        `json:"exam_name"`
	SchoolName string        `json:"school_name"`
	Extend     string        `json:"extend"`
	IsPublic   bool          `json:"is_public"`
	Questions  []QuestionDTO `json:"questions"`
}

type QuestionDTO struct {
	Content       string           `json:"content"`
	Level         string           `json:"level"`
	Type          string           `json:"type"` // single | multiple | essay
	Options       json.RawMessage  `json:"options"`
	OrderNo       int              `json:"order"`
	ParentOrder   *int             `json:"parent_order"`
	CorrectAnswer *json.RawMessage `json:"correct_answer"`
	UserAnswer    interface{}      `json:"user_answer"`
}
type ExamSet struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	IsPublic  bool   `json:"is_public"`
	CreatedBy int    `json:"created_by"`
}
type ExamSetItem struct {
	ID         int64     `json:"id"`
	Name       string    `json:"name"`
	SchoolName string    `json:"school_name"`
	Extend     string    `json:"extend"`
	IsPublic   bool      `json:"is_public"`
	CreatedAt  time.Time `json:"created_at"`
}
type ExamAttemptItem struct {
	AttemptID   int64      `json:"attempt_id"`
	ExamSetID   int64      `json:"exam_set_id"`
	ExamName    string     `json:"exam_name"`
	IsPublic    bool       `json:"is_public"`
	StartedAt   time.Time  `json:"started_at"`
	SubmittedAt *time.Time `json:"submitted_at,omitempty"`
}
