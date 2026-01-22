package model

import (
	"encoding/json"
	"time"
)

type SubmitExamRequest struct {
	Subject   string        `json:"subject"`
	ExamName  string        `json:"exam_name"`
	IsPublic  bool          `json:"is_public"`
	Questions []QuestionDTO `json:"questions"`
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
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	IsPublic bool   `json:"is_public"`
}
type ExamSetItem struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	IsPublic  bool      `json:"is_public"`
	CreatedAt time.Time `json:"created_at"`
}
