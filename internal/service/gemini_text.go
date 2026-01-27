package service

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

func AnalyzeTextWithGemini(
	ctx context.Context,
	question string,
	prompt string,
) (string, error) {

	if strings.TrimSpace(question) == "" {
		return "", errors.New("empty question")
	}

	parts := []map[string]interface{}{
		{"text": prompt},
		{"text": question},
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

		if resp.StatusCode == http.StatusTooManyRequests ||
			resp.StatusCode == http.StatusForbidden {
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

		jsonStr = strings.ReplaceAll(jsonStr, `\`, `\\`)

		return jsonStr, nil
	}

	return "", errors.New("all Gemini API keys exhausted: " + lastErr.Error())
}
