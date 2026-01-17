package service

import (
	"database/sql"

	"traingolang/internal/model"
	"traingolang/internal/repository"
)

type ExamService struct {
	DB   *sql.DB
	Repo *repository.ExamRepo
}

func NewExamService(db *sql.DB) *ExamService {
	return &ExamService{
		DB:   db,
		Repo: repository.NewExamRepo(db),
	}
}

func (s *ExamService) SubmitExam(
	userID int,
	role string,
	req model.SubmitExamRequest,
) (int64, int64, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return 0, 0, err
	}
	defer tx.Rollback()

	examSetID, err := s.Repo.CreateExamSet(
		tx,
		req.ExamName,
		userID,
		req.IsPublic,
	)
	if err != nil {
		return 0, 0, err
	}

	attemptID, err := s.Repo.CreateAttempt(tx, examSetID, userID)
	if err != nil {
		return 0, 0, err
	}

	for _, q := range req.Questions {

		var correct interface{} = nil
		if role == "admin" || role == "author" {
			correct = q.CorrectAnswer
		}

		questionID, err := s.Repo.CreateQuestion(
			tx,
			examSetID,
			q,
			correct,
		)
		if err != nil {
			return 0, 0, err
		}

		if err := s.Repo.CreateAnswer(
			tx,
			attemptID,
			questionID,
			q.UserAnswer,
		); err != nil {
			return 0, 0, err
		}
	}

	if err := tx.Commit(); err != nil {
		return 0, 0, err
	}

	return examSetID, attemptID, nil
}
