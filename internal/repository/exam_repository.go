package repository

import (
	"database/sql"

	"traingolang/internal/model"
	"traingolang/internal/util"
)

type ExamRepo struct {
	DB *sql.DB
}

func NewExamRepo(db *sql.DB) *ExamRepo {
	return &ExamRepo{DB: db}
}

func (r *ExamRepo) CreateExamSet(
	tx *sql.Tx,
	name string,
	userID int,
	isPublic bool,
) (int64, error) {

	var id int64
	err := tx.QueryRow(`
		insert into exam_sets (name, created_by, is_public)
		values ($1,$2,$3)
		returning id
	`, name, userID, isPublic).Scan(&id)

	return id, err
}

func (r *ExamRepo) CreateAttempt(
	tx *sql.Tx,
	examSetID int64,
	userID int,
) (int64, error) {

	var id int64
	err := tx.QueryRow(`
		insert into exam_attempts (exam_set_id, user_id)
		values ($1,$2)
		returning id
	`, examSetID, userID).Scan(&id)

	return id, err
}

func (r *ExamRepo) CreateQuestion(
	tx *sql.Tx,
	examSetID int64,
	q model.QuestionDTO,
	correctAnswer interface{},
) (int64, error) {

	var id int64
	err := tx.QueryRow(`
		insert into exam_questions
		(exam_set_id, content, type, level, options, correct_answer, order_no, parent_order)
		values ($1,$2,$3,$4,$5,$6,$7,$8)
		returning id
	`,
		examSetID,
		q.Content,
		q.Type,
		q.Level,
		util.ToJSONB(q.Options),
		util.ToJSONB(correctAnswer),
		q.OrderNo,
		q.ParentOrder,
	).Scan(&id)

	return id, err
}

func (r *ExamRepo) CreateAnswer(
	tx *sql.Tx,
	attemptID int64,
	questionID int64,
	answer interface{},
) error {

	_, err := tx.Exec(`
		insert into exam_answers (attempt_id, question_id, answer)
		values ($1,$2,$3)
	`, attemptID, questionID, util.ToJSONB(answer))

	return err
}
