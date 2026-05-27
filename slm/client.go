package slm

import (
	"bayar-woy-project/config"
	"bytes"
	"encoding/json"
	"net/http"
	"time"
)

type classifyResponse struct {
	Category            string  `json:"category"`
	SecondaryCategory   string  `json:"secondary_category"`
	TransactionType     string  `json:"transaction_type"`
	Confidence          float64 `json:"confidence"`
}

// ClassifyResult holds category and type prediction from SLM.
type ClassifyResult struct {
	Category            string
	SecondaryCategory   string
	TransactionType     string
}

// Classify calls the SLM service and returns category + transaction type.
// Returns empty strings if SLM is unreachable — caller still proceeds.
func Classify(title string) ClassifyResult {
	baseURL := config.GetEnv("SLM_URL")

	payload, _ := json.Marshal(map[string]string{"title": title})
	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Post(baseURL+"/classify", "application/json", bytes.NewReader(payload))
	if err != nil {
		return ClassifyResult{}
	}
	defer resp.Body.Close()

	var result classifyResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return ClassifyResult{}
	}
	return ClassifyResult{
		Category:            result.Category,
		SecondaryCategory:   result.SecondaryCategory,
		TransactionType:     result.TransactionType,
	}
}

// ClassifyTitle is kept for backward compatibility.
func ClassifyTitle(title string) string {
	return Classify(title).Category
}

// Feedback sends a confirmed or corrected classification to the SLM for reinforcement learning.
// Best-effort — errors are silently ignored (same policy as Classify).
func Feedback(title, category, secondaryCategory string) {
	baseURL := config.GetEnv("SLM_URL")

	payload, _ := json.Marshal(map[string]string{
		"title":                      title,
		"correct_category":           category,
		"correct_secondary_category": secondaryCategory,
	})
	client := &http.Client{Timeout: 2 * time.Second}

	resp, err := client.Post(baseURL+"/feedback", "application/json", bytes.NewReader(payload))
	if err != nil {
		return
	}
	resp.Body.Close()
}
