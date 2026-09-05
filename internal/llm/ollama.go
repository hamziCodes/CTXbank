package llm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// OllamaClient communicates with a local Ollama HTTP daemon.
type OllamaClient struct {
	BaseURL    string
	Model      string
	HTTPClient *http.Client
}

// NewOllamaClient initializes an Ollama client with a 3-second connection timeout.
func NewOllamaClient(baseURL, model string) *OllamaClient {
	if baseURL == "" {
		baseURL = "http://localhost:11434"
	}
	if model == "" {
		model = "llama3.2"
	}
	return &OllamaClient{
		BaseURL: strings.TrimRight(baseURL, "/"),
		Model:   model,
		HTTPClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// IsAvailable checks if the Ollama server is currently responding.
func (c *OllamaClient) IsAvailable() bool {
	resp, err := c.HTTPClient.Get(c.BaseURL + "/api/tags")
	if err != nil || resp.StatusCode != http.StatusOK {
		return false
	}
	_ = resp.Body.Close()
	return true
}

// ClassifyChunk sends text to Ollama to categorize into a memory bank target document.
func (c *OllamaClient) ClassifyChunk(text string) (TargetDoc, error) {
	if !c.IsAvailable() {
		return DocProductContext, ErrLLMUnavailable
	}

	prompt := fmt.Sprintf(`Classify this text into exactly ONE category from:
[projectbrief, productContext, systemPatterns, techContext, discard]

Text:
"""
%s
"""

Respond with ONLY the category name and nothing else.`, text)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":  c.Model,
		"prompt": prompt,
		"stream": false,
		"options": map[string]interface{}{
			"temperature": 0.0,
		},
	})

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return DocProductContext, fmt.Errorf("%w: %v", ErrLLMUnavailable, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return DocProductContext, fmt.Errorf("%w: status %d", ErrLLMUnavailable, resp.StatusCode)
	}

	var parsed struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return DocProductContext, err
	}

	raw := strings.ToLower(strings.TrimSpace(parsed.Response))
	switch {
	case strings.Contains(raw, "projectbrief"):
		return DocProjectBrief, nil
	case strings.Contains(raw, "productcontext"):
		return DocProductContext, nil
	case strings.Contains(raw, "systempatterns"):
		return DocSystemPatterns, nil
	case strings.Contains(raw, "techcontext"):
		return DocTechContext, nil
	case strings.Contains(raw, "discard"):
		return DocDiscard, nil
	default:
		return DocProductContext, nil
	}
}

// SummarizeNarrative turns structured audit data into concise product narrative.
func (c *OllamaClient) SummarizeNarrative(contextData string) (string, error) {
	if !c.IsAvailable() {
		return "", ErrLLMUnavailable
	}

	prompt := fmt.Sprintf(`Given this repository summary, write a concise 2-sentence explanation of what this product accomplishes:
%s

Respond only with the 2-sentence summary.`, contextData)

	reqBody, _ := json.Marshal(map[string]interface{}{
		"model":  c.Model,
		"prompt": prompt,
		"stream": false,
	})

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/generate", "application/json", bytes.NewReader(reqBody))
	if err != nil {
		return "", fmt.Errorf("%w: %v", ErrLLMUnavailable, err)
	}
	defer resp.Body.Close()

	var parsed struct {
		Response string `json:"response"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return "", err
	}

	return strings.TrimSpace(parsed.Response), nil
}
