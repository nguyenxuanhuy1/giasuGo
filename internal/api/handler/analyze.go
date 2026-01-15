package handler

import (
	"encoding/json"
	"io"
	"net/http"
	"traingolang/internal/prompt"
	"traingolang/internal/service"

	"github.com/gin-gonic/gin"
)

func AnalyzeImage(c *gin.Context) {
	file, err := c.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "image is required",
		})
		return
	}

	mode := c.PostForm("mode") // cv | exam
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

	result, err := service.AnalyzeImageWithGemini(
		imageBytes,
		file.Header.Get("Content-Type"),
		finalPrompt,
	)
	if err != nil {
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
