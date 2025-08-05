package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"go-code/pkg/models"
)

const (
	GroqAPIURL = "https://api.groq.com/openai/v1/chat/completions"
)

// GroqClient represents the Groq API client
type GroqClient struct {
	APIKey     string
	HTTPClient *http.Client
	BaseURL    string
}

// NewGroqClient creates a new Groq API client
func NewGroqClient(apiKey string) *GroqClient {
	return &GroqClient{
		APIKey:  apiKey,
		BaseURL: GroqAPIURL,
		HTTPClient: &http.Client{
			Timeout: 60 * time.Second,
		},
	}
}

// ChatRequest represents a chat completion request
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float32   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream,omitempty"`
}

// Message represents a chat message
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a chat completion response
type ChatResponse struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
	Usage   Usage    `json:"usage"`
}

// Choice represents a completion choice
type Choice struct {
	Index        int     `json:"index"`
	Message      Message `json:"message"`
	FinishReason string  `json:"finish_reason"`
}

// Usage represents token usage information
type Usage struct {
	PromptTokens     int `json:"prompt_tokens"`
	CompletionTokens int `json:"completion_tokens"`
	TotalTokens      int `json:"total_tokens"`
}

// GroqError represents an error from the Groq API
type GroqError struct {
	ErrorInfo struct {
		Message string `json:"message"`
		Type    string `json:"type"`
		Code    string `json:"code"`
	} `json:"error"`
}

func (e GroqError) Error() string {
	return fmt.Sprintf("Groq API error: %s (type: %s, code: %s)", 
		e.ErrorInfo.Message, e.ErrorInfo.Type, e.ErrorInfo.Code)
}

// SendChatRequest sends a chat completion request to Groq API
func (c *GroqClient) SendChatRequest(req ChatRequest) (*ChatResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	httpReq, err := http.NewRequest("POST", c.BaseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var groqErr GroqError
		if err := json.Unmarshal(body, &groqErr); err != nil {
			return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(body))
		}
		return nil, groqErr
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &chatResp, nil
}

// GetAvailableModels returns a list of available Groq models
func GetAvailableModels() []string {
	return []string{
		"llama-3.1-70b-versatile",
		"llama-3.1-8b-instant",
		"llama-3.2-90b-text-preview",
		"llama-3.2-11b-text-preview",
		"llama-3.2-3b-preview",
		"llama-3.2-1b-preview",
		"mixtral-8x7b-32768",
		"gemma2-9b-it",
		"gemma-7b-it",
		"qwen2.5-coder-32b-instruct",
		"llama3-groq-70b-8192-tool-use-preview",
		"llama3-groq-8b-8192-tool-use-preview",
		"openai/gpt-oss-120b",
	}
}

// IsValidModel checks if a model is available
func IsValidModel(model string) bool {
	models := GetAvailableModels()
	for _, m := range models {
		if strings.EqualFold(m, model) {
			return true
		}
	}
	return false
}

// ProcessAgentRequest processes a request using the specified agent configuration
func (c *GroqClient) ProcessAgentRequest(agentType models.AgentType, systemPrompt, userMessage string, config models.AgentConfig) (*models.Response, error) {
	req := ChatRequest{
		Model: config.Model,
		Messages: []Message{
			{Role: "system", Content: systemPrompt},
			{Role: "user", Content: userMessage},
		},
		Temperature: config.Temperature,
		MaxTokens:   config.MaxTokens,
	}

	resp, err := c.SendChatRequest(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send chat request: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned from API")
	}

	return &models.Response{
		Content:    resp.Choices[0].Message.Content,
		TokensUsed: resp.Usage.TotalTokens,
		Model:      resp.Model,
		Agent:      agentType,
	}, nil
}