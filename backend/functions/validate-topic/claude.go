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
		model:  "claude-sonnet-4-5", // Latest Claude model
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}, nil
}

// ValidateTopicAndSuggestPanelists validates topic relevance and suggests panelists in one API call
func (c *ClaudeClient) ValidateTopicAndSuggestPanelists(ctx context.Context, topic string) (bool, string, []Panelist, error) {
	// Build the combined prompt for Claude
	prompt := fmt.Sprintf(`You are an expert in theology and philosophy. Your task is to:
1. Determine if the following topic is suitable for a theological or philosophical debate
2. If suitable, suggest 5-8 historical or contemporary figures who would make excellent panelists

Topic: "%s"

First, evaluate whether this topic relates to:
- Theology (study of God, religion, faith, sacred texts)
- Philosophy (ethics, morality, epistemology, metaphysics, logic)
- Religious practice and doctrine
- Moral and ethical questions
- Existential questions

If the topic IS relevant, suggest 8-20 diverse panelists with:
- Different theological/philosophical positions on this topic
- Historical or contemporary figures with known views
- Mix of perspectives (theist/atheist, conservative/progressive, different traditions)

Respond with a JSON object:
{
  "isRelevant": boolean,
  "message": "brief explanation (max 200 chars)",
  "panelists": [
    {
      "id": "unique-kebab-case-id",
      "name": "Full Name",
      "tagline": "One-line description (max 60 chars)",
      "bio": "Brief bio explaining their relevant expertise (max 300 chars)",
      "avatarUrl": "/avatars/placeholder-avatar.png",
      "position": "Brief statement of their likely position on this topic (max 100 chars)"
    }
  ]
}

If not relevant, set panelists to an empty array.
Format your response as valid JSON only, no other text.`, topic)

	// Create the request body with higher token limit for panelist suggestions
	requestBody := map[string]interface{}{
		"model":      c.model,
		"max_tokens": 2000, // Increased for panelist data
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return false, "", nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Call Claude API
	response, err := c.callClaudeAPI(ctx, requestJSON)
	if err != nil {
		return false, "", nil, err
	}

	// Parse the response
	return c.parseValidationAndPanelistResponse(response)
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

// parseValidationAndPanelistResponse parses Claude's JSON response with panelists
func (c *ClaudeClient) parseValidationAndPanelistResponse(response string) (bool, string, []Panelist, error) {
	// Extract JSON from response (Claude might include it in content)
	var result struct {
		IsRelevant bool       `json:"isRelevant"`
		Message    string     `json:"message"`
		Panelists  []Panelist `json:"panelists"`
	}

	// Try to find JSON in the response
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")

	if startIdx == -1 || endIdx == -1 {
		return false, "", nil, errors.New("no JSON found in Claude response")
	}

	jsonStr := response[startIdx : endIdx+1]

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return false, "", nil, fmt.Errorf("failed to parse Claude response: %w", err)
	}

	// Limit message length to 200 characters
	if len(result.Message) > 200 {
		result.Message = result.Message[:197] + "..."
	}

	// If not relevant, panelists should be empty
	if !result.IsRelevant {
		return false, result.Message, []Panelist{}, nil
	}

	// Validate and sanitize panelist data
	validPanelists := make([]Panelist, 0, len(result.Panelists))
	for _, p := range result.Panelists {
		// Basic validation
		if p.Name == "" || p.ID == "" {
			continue
		}
		
		// Limit field lengths
		if len(p.Tagline) > 60 {
			p.Tagline = p.Tagline[:57] + "..."
		}
		if len(p.Bio) > 300 {
			p.Bio = p.Bio[:297] + "..."
		}
		if len(p.Position) > 100 {
			p.Position = p.Position[:97] + "..."
		}
		
		validPanelists = append(validPanelists, p)
	}

	return result.IsRelevant, result.Message, validPanelists, nil
}

