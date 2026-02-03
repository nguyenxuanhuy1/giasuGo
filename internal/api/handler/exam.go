package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"strconv"
	"traingolang/internal/auth"
	"traingolang/internal/model"
	"traingolang/internal/service"
)

func SubmitExamHandler(examService *service.ExamService) gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, ok := auth.GetCurrentUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userID := int(claims.UserID)
		role := claims.Role

		var req model.SubmitExamRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid json",
			})
			return
		}

		examSetID, attemptID, err := examService.SubmitExam(
			userID,
			role,
			req,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusCreated, gin.H{
			"exam_set_id": examSetID,
			"attempt_id":  attemptID,
		})
	}
}
func RedoExamHandler(examService *service.ExamService) gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, ok := auth.GetCurrentUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userID := int(claims.UserID)

		examSetIDParam := c.Param("id")
		examSetID, err := strconv.ParseInt(examSetIDParam, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid exam_set_id",
			})
			return
		}

		result, err := examService.RedoExam(examSetID, userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
func GetMyExamSetsHandler(examService *service.ExamService) gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, ok := auth.GetCurrentUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userID := int(claims.UserID)

		exams, err := examService.GetMyExamSets(userID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"data": exams,
		})
	}
}

func GetExamQuestionsHandler(
	examService *service.ExamService,
) gin.HandlerFunc {
	return func(c *gin.Context) {

		claims, ok := auth.GetCurrentUser(c)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "unauthorized",
			})
			return
		}

		userID := int(claims.UserID)

		examSetID, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid exam_set_id",
			})
			return
		}

		questions, err := examService.GetExamQuestionsForUser(
			examSetID,
			userID,
		)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"error": err.Error(),
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"exam_set_id": examSetID,
			"questions":   questions,
		})
	}
}
