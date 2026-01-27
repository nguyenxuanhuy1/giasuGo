package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"traingolang/internal/prompt"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
)

type QuestionRequest struct {
	Mode    string          `json:"mode"`
	Content json.RawMessage `json:"content" binding:"required"`
	Prompt  string          `json:"prompt"`
}

func AnalyzeQuestion(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), RequestTimeout)
	defer cancel()

	var req QuestionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid json body",
		})
		return
	}

	finalPrompt := req.Prompt
	if finalPrompt == "" {
		switch req.Mode {
		case "exam":
			finalPrompt = prompt.ExamQuestionPrompt
		default:
			finalPrompt = "hãy nói bạn yêu tôi"
		}
	}

	questionJSON, err := json.MarshalIndent(req.Content, "", "  ")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid question content",
		})
		return
	}

	select {
	case geminiSemaphore <- struct{}{}:
		defer func() { <-geminiSemaphore }()
	case <-ctx.Done():
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "server busy"})
		return
	case <-time.After(5 * time.Second):
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": "server busy"})
		return
	}

	jsonStr, err := service.AnalyzeTextWithGemini(
		ctx,
		string(questionJSON),
		finalPrompt,
	)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusRequestTimeout, gin.H{"error": "request timeout"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	var parsed any
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "invalid json after cleaning",
		})
		return
	}

	c.JSON(http.StatusOK, parsed)
}
