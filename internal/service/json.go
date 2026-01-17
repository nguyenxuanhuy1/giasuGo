package service

import (
	"errors"
	"strings"
)

func ExtractJSON(s string) (string, error) {
	start := strings.Index(s, "{")
	end := strings.LastIndex(s, "}")

	if start == -1 || end == -1 || end <= start {
		return "", errors.New("no valid json block")
	}

	return s[start : end+1], nil
}
