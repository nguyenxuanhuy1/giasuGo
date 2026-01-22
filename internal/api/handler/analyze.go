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

	if err := c.Request.ParseMultipartForm(20 << 20); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid multipart form"})
		return
	}

	files := c.Request.MultipartForm.File["images"]
	if len(files) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "images are required"})
		return
	}
	mode := c.PostForm("mode")
	customPrompt := c.PostForm("prompt")

	finalPrompt := customPrompt
	if finalPrompt == "" {
		if mode == "exam" {
			finalPrompt = prompt.ExamPrompt
		} else {
			finalPrompt = "Extract all visible text from the image(s)"
		}
	}

	var images []service.ImageInput

	for _, file := range files {
		f, err := file.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot open image"})
			return
		}

		data, err := io.ReadAll(f)
		f.Close()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "cannot read image"})
			return
		}

		images = append(images, service.ImageInput{
			Data:     data,
			MimeType: file.Header.Get("Content-Type"),
		})
	}

	// giới hạn câu hỏi
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

	jsonStr, err := service.AnalyzeImagesWithGemini(
		ctx,
		images,
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
