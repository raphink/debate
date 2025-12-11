package main

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	claudeAPIURL = "https://api.anthropic.com/v1/messages"
	apiVersion   = "2023-06-01"
)

// ClaudeClient handles communication with the Anthropic Claude API
type ClaudeClient struct {
	apiKey     string
	model      string
	httpClient *http.Client
}

// NewClaudeClient creates a new Claude API client
func NewClaudeClient() (*ClaudeClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, errors.New("ANTHROPIC_API_KEY environment variable not set")
	}

	return &ClaudeClient{
		apiKey: apiKey,
		model:  "claude-3-5-sonnet-20241022", // Latest Claude model
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// ValidateTopicRelevance asks Claude to determine if the topic is relevant for theological/philosophical debate
func (c *ClaudeClient) ValidateTopicRelevance(ctx context.Context, topic string) (bool, string, error) {
	// Build the prompt for Claude
	prompt := fmt.Sprintf(`You are an expert in theology and philosophy. Your task is to determine if the following topic is suitable for a theological or philosophical debate.

Topic: "%s"

Evaluate whether this topic relates to:
- Theology (study of God, religion, faith, sacred texts)
- Philosophy (ethics, morality, epistemology, metaphysics, logic)
- Religious practice and doctrine
- Moral and ethical questions
- Existential questions

Respond with a JSON object containing:
1. "isRelevant": true if the topic is suitable for theological/philosophical debate, false otherwise
2. "message": A brief explanation (max 200 characters) of why the topic is or isn't suitable

Format your response as valid JSON only, no other text:
{"isRelevant": boolean, "message": "explanation"}`, topic)

	// Create the request body
	requestBody := map[string]interface{}{
		"model":      c.model,
		"max_tokens": 200,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return false, "", fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call Claude API
	response, err := c.callClaudeAPI(ctx, requestJSON)
	if err != nil {
		return false, "", err
	}

	// Parse the response
	return c.parseValidationResponse(response)
}

// ClaudeResponse represents the response from Claude API
type ClaudeResponse struct {
	ID      string `json:"id"`
	Type    string `json:"type"`
	Role    string `json:"role"`
	Content []struct {
		Type string `json:"type"`
		Text string `json:"text"`
	} `json:"content"`
	Model        string `json:"model"`
	StopReason   string `json:"stop_reason"`
	StopSequence string `json:"stop_sequence"`
	Usage        struct {
		InputTokens  int `json:"input_tokens"`
		OutputTokens int `json:"output_tokens"`
	} `json:"usage"`
}

// callClaudeAPI makes the HTTP request to Claude API
func (c *ClaudeClient) callClaudeAPI(ctx context.Context, requestBody []byte) (string, error) {
	// Create HTTP request
	req, err := http.NewRequestWithContext(ctx, "POST", claudeAPIURL, bytes.NewBuffer(requestBody))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	req.Header.Set("anthropic-version", apiVersion)

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to call Claude API: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Claude API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Parse Claude response
	var claudeResp ClaudeResponse
	if err := json.Unmarshal(body, &claudeResp); err != nil {
		return "", fmt.Errorf("failed to parse Claude response: %w", err)
	}

	// Extract text content
	if len(claudeResp.Content) == 0 {
		return "", errors.New("no content in Claude response")
	}

	return claudeResp.Content[0].Text, nil
}

// parseValidationResponse parses Claude's JSON response
func (c *ClaudeClient) parseValidationResponse(response string) (bool, string, error) {
	// Extract JSON from response (Claude might include it in content)
	var result struct {
		IsRelevant bool   `json:"isRelevant"`
		Message    string `json:"message"`
	}

	// Try to find JSON in the response
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")
	
	if startIdx == -1 || endIdx == -1 {
		return false, "", errors.New("no JSON found in Claude response")
	}

	jsonStr := response[startIdx : endIdx+1]
	
	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return false, "", fmt.Errorf("failed to parse Claude response: %w", err)
	}

	// Limit message length to 200 characters
	if len(result.Message) > 200 {
		result.Message = result.Message[:197] + "..."
	}

	return result.IsRelevant, result.Message, nil
}
