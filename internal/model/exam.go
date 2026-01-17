package model

type SubmitExamRequest struct {
	Subject   string        `json:"subject"`
	ExamName  string        `json:"exam_name"`
	IsPublic  bool          `json:"is_public"`
	Questions []QuestionDTO `json:"questions"`
}

type QuestionDTO struct {
	Content       string      `json:"content"`
	Level         string      `json:"level"`
	Type          string      `json:"type"` // single | multiple | essay
	Options       interface{} `json:"options"`
	OrderNo       int         `json:"order"`
	ParentOrder   *int        `json:"parent_order"`
	CorrectAnswer interface{} `json:"correct_answer"`
	UserAnswer    interface{} `json:"user_answer"`
}
