package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"traingolang/internal/prompt"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
)

const (
	RequestTimeout = 120 * time.Second
	MaxConcurrent  = 10
)

var geminiSemaphore = make(chan struct{}, MaxConcurrent)

func AnalyzeImage(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), RequestTimeout)
	defer cancel()

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image is required"})
		return
	}

	mode := c.PostForm("mode")
	customPrompt := c.PostForm("prompt")

	finalPrompt := customPrompt
	if finalPrompt == "" {
		if mode == "exam" {
			finalPrompt = prompt.ExamPrompt
		} else {
			finalPrompt = "Trích xuất toàn bộ nội dung văn bản trong ảnh"
		}
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open image"})
		return
	}
	defer f.Close()

	imageBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read image"})
		return
	}

	// giới hạn số request gọi
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

	jsonStr, err := service.AnalyzeImageWithGemini(
		ctx,
		imageBytes,
		file.Header.Get("Content-Type"),
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
