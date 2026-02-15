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
	schoolName string,
	extend string,
	userID int,
	isPublic int,
) (int64, error) {

	var id int64
	err := tx.QueryRow(`
		insert into exam_sets (name, school_name, extend, created_by, is_public)
		values ($1, $2, $3, $4, $5)
		returning id
	`, name, schoolName, extend, userID, isPublic).Scan(&id)

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
func (r *ExamRepo) GetExamSetByID(
	tx *sql.Tx,
	examSetID int64,
) (*model.ExamSet, error) {

	var e model.ExamSet
	err := tx.QueryRow(`
		select
			id,
			name,
			is_public,
			created_by
		from exam_sets
		where id = $1
	`, examSetID).Scan(
		&e.ID,
		&e.Name,
		&e.IsPublic,
		&e.CreatedBy,
	)

	return &e, err
}

func (r *ExamRepo) GetQuestionsByExamSet(
	tx *sql.Tx,
	examSetID int64,
) ([]model.QuestionDTO, error) {

	rows, err := tx.Query(`
		select
			content,
			level,
			type,
			options,
			order_no,
			parent_order,
			correct_answer
		from exam_questions
		where exam_set_id = $1
		order by order_no
	`, examSetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var questions []model.QuestionDTO
	for rows.Next() {
		var q model.QuestionDTO
		if err := rows.Scan(
			&q.Content,
			&q.Level,
			&q.Type,
			&q.Options,
			&q.OrderNo,
			&q.ParentOrder,
			&q.CorrectAnswer,
		); err != nil {
			return nil, err
		}
		questions = append(questions, q)
	}

	return questions, nil
}
func (r *ExamRepo) GetMyExamSets(
	userID int,
) ([]model.ExamSetItem, error) {

	rows, err := r.DB.Query(`
		select
			id,
			name,
			school_name,
			extend,
			is_public,
			created_at
		from exam_sets
		where created_by = $1
		order by created_at desc
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []model.ExamSetItem
	for rows.Next() {
		var e model.ExamSetItem
		if err := rows.Scan(
			&e.ID,
			&e.Name,
			&e.SchoolName,
			&e.Extend,
			&e.IsPublic,
			&e.CreatedAt,
		); err != nil {
			return nil, err
		}
		result = append(result, e)
	}

	return result, nil
}
