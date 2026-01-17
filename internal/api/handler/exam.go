package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"traingolang/internal/model"
	"traingolang/internal/service"
)

func SubmitExamHandler(examService *service.ExamService) gin.HandlerFunc {
	return func(c *gin.Context) {

		// JWT middleware của bạn chắc đã set
		userID := c.GetInt("user_id")
		role := c.GetString("role")

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
