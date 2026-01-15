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
	RequestTimeout = 95 * time.Second
	MaxConcurrent  = 10 // Số request đồng thời tối đa
)

// Semaphore để giới hạn số request đồng thời xử lý Gemini
var geminiSemaphore = make(chan struct{}, MaxConcurrent)

func AnalyzeImage(c *gin.Context) {
	// Thêm timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), RequestTimeout)
	defer cancel()

	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image is required",
		})
		return
	}

	mode := c.PostForm("mode")
	customPrompt := c.PostForm("prompt")

	var finalPrompt string
	if customPrompt != "" {
		finalPrompt = customPrompt
	} else {
		switch mode {
		case "exam":
			finalPrompt = prompt.ExamPrompt
		default:
			finalPrompt = "Trích xuất toàn bộ nội dung văn bản trong ảnh"
		}
	}

	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot open image",
		})
		return
	}
	defer f.Close()

	imageBytes, err := io.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "cannot read image",
		})
		return
	}

	// Thử acquire semaphore - giới hạn số request gọi Gemini đồng thời
	select {
	case geminiSemaphore <- struct{}{}:
		// Acquired, tiếp tục xử lý
		defer func() { <-geminiSemaphore }()
	case <-ctx.Done():
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "server đang xử lý quá nhiều request, vui lòng thử lại sau",
		})
		return
	case <-time.After(5 * time.Second):
		// Nếu chờ quá 5s vẫn không có slot → reject
		c.JSON(http.StatusServiceUnavailable, gin.H{
			"error": "server đang bận, vui lòng thử lại sau",
		})
		return
	}

	result, err := service.AnalyzeImageWithGemini(
		ctx,
		imageBytes,
		file.Header.Get("Content-Type"),
		finalPrompt,
	)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			c.JSON(http.StatusRequestTimeout, gin.H{
				"error": "request timeout",
			})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	var parsed any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		c.JSON(500, gin.H{
			"error": "invalid json",
			"raw":   result,
		})
		return
	}

	c.JSON(200, parsed)
}
