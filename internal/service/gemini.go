package service

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

type ImageInput struct {
	Data     []byte
	MimeType string
}

type geminiResponse struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
		} `json:"content"`
	} `json:"candidates"`
}

func AnalyzeImagesWithGemini(
	ctx context.Context,
	images []ImageInput,
	prompt string,
) (string, error) {

	if len(images) == 0 {
		return "", errors.New("no images provided")
	}

	// Build parts: prompt + multiple images
	parts := []map[string]interface{}{
		{"text": prompt},
	}

	for _, img := range images {
		if len(img.Data) == 0 {
			return "", errors.New("empty image")
		}
		if !strings.HasPrefix(img.MimeType, "image/") {
			return "", errors.New("invalid image mime type")
		}

		encodedImage := base64.StdEncoding.EncodeToString(img.Data)

		parts = append(parts, map[string]interface{}{
			"inline_data": map[string]string{
				"mime_type": img.MimeType,
				"data":      encodedImage,
			},
		})
	}

	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": parts,
			},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 120 * time.Second}
	url := os.Getenv("GEMINI_URL")

	var lastErr error

	for i := 0; i < geminiKeyCount(); i++ {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			url,
			bytes.NewBuffer(body),
		)
		if err != nil {
			return "", err
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-goog-api-key", nextGeminiKey())

		resp, err := client.Do(req)
		if err != nil {
			lastErr = err
			continue
		}

		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusForbidden {
			lastErr = errors.New(string(respBody))
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", errors.New(string(respBody))
		}

		var geminiResp geminiResponse
		if err := json.Unmarshal(respBody, &geminiResp); err != nil {
			return "", err
		}

		if len(geminiResp.Candidates) == 0 ||
			len(geminiResp.Candidates[0].Content.Parts) == 0 {
			return "", errors.New("no response from gemini")
		}

		raw := strings.TrimSpace(
			geminiResp.Candidates[0].Content.Parts[0].Text,
		)

		jsonStr, err := ExtractJSON(raw)
		if err != nil {
			return "", err
		}

		// Escape single backslash for Go JSON safety
		jsonStr = strings.ReplaceAll(jsonStr, `\`, `\\`)

		return jsonStr, nil
	}

	return "", errors.New("all Gemini API keys exhausted: " + lastErr.Error())
}
