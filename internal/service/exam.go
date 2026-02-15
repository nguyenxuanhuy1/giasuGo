package service

import (
	"database/sql"
	"traingolang/internal/model"
	"traingolang/internal/repository"
	"traingolang/internal/util"
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
		req.SchoolName,
		req.Extend,
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
		correct = q.CorrectAnswer

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
func (s *ExamService) RedoExam(
	examSetID int64,
	userID int,
) (map[string]interface{}, error) {

	tx, err := s.DB.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	exam, err := s.Repo.GetExamSetByID(tx, examSetID)
	if err != nil {
		return nil, err
	}
	questions, err := s.Repo.GetQuestionsByExamSet(tx, examSetID)
	if err != nil {
		return nil, err
	}
	attemptID, err := s.Repo.CreateAttempt(tx, examSetID, userID)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	for i := range questions {
		questions[i].UserAnswer = nil
	}

	return map[string]interface{}{
		"exam_set_id": exam.ID,
		"exam_name":   exam.Name,
		"is_public":   exam.IsPublic,
		"attempt_id":  attemptID,
		"questions":   questions,
	}, nil
}
func (s *ExamService) GetMyExamSets(
	userID int,
) ([]model.ExamSetItem, error) {

	return s.Repo.GetMyExamSets(userID)
}
func (s *ExamService) GetPublicExamSets(
	req model.PublicExamListRequest,
) (*util.PaginatedResponse[model.ExamSetItem], error) {

	offset, limit := util.NewPagination(req.Page, req.PageSize)

	return s.Repo.GetPublicExamSetsPaginated(
		req.Search,
		offset,
		limit,
	)
}
