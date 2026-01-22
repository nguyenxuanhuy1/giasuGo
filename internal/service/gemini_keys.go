package service

import (
	"os"
	"strings"
	"sync/atomic"
)

var (
	geminiKeys     []string
	geminiKeyIndex uint32
)

func loadGeminiKeysOnce() {
	if len(geminiKeys) > 0 {
		return
	}

	raw := os.Getenv("GEMINI_API_KEYS")
	if raw == "" {
		panic("GEMINI_API_KEYS is not set")
	}

	geminiKeys = strings.Split(raw, ",")
}

func nextGeminiKey() string {
	loadGeminiKeysOnce()
	i := atomic.AddUint32(&geminiKeyIndex, 1)
	return geminiKeys[int(i)%len(geminiKeys)]
}

func geminiKeyCount() int {
	loadGeminiKeysOnce()
	return len(geminiKeys)
}
