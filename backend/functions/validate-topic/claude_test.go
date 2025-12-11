package validatetopic

import (
	"context"
	"os"
	"testing"
)

func TestNewClaudeClient(t *testing.T) {
	// Save original env var
	originalKey := os.Getenv("ANTHROPIC_API_KEY")
	defer func() {
		if originalKey != "" {
			os.Setenv("ANTHROPIC_API_KEY", originalKey)
		} else {
			os.Unsetenv("ANTHROPIC_API_KEY")
		}
	}()

	t.Run("missing API key", func(t *testing.T) {
		os.Unsetenv("ANTHROPIC_API_KEY")
		_, err := NewClaudeClient()
		if err == nil {
			t.Error("NewClaudeClient() expected error for missing API key, got nil")
		}
	})

	t.Run("valid API key", func(t *testing.T) {
		os.Setenv("ANTHROPIC_API_KEY", "test-key")
		client, err := NewClaudeClient()
		if err != nil {
			t.Errorf("NewClaudeClient() unexpected error: %v", err)
		}
		if client == nil {
			t.Error("NewClaudeClient() returned nil client")
		}
		if client.apiKey != "test-key" {
			t.Errorf("NewClaudeClient() apiKey = %v, want test-key", client.apiKey)
		}
		if client.model == "" {
			t.Error("NewClaudeClient() model is empty")
		}
	})
}

func TestParseValidationResponse(t *testing.T) {
	client := &ClaudeClient{
		apiKey: "test-key",
		model:  "test-model",
	}

	tests := []struct {
		name         string
		response     string
		wantRelevant bool
		wantMessage  string
		wantErr      bool
	}{
		{
			name:         "valid JSON response",
			response:     `{"isRelevant": true, "message": "This is a theological topic."}`,
			wantRelevant: true,
			wantMessage:  "This is a theological topic.",
			wantErr:      false,
		},
		{
			name:         "valid JSON with extra text",
			response:     `Some preamble {"isRelevant": false, "message": "Not relevant."} some postamble`,
			wantRelevant: false,
			wantMessage:  "Not relevant.",
			wantErr:      false,
		},
		{
			name:     "no JSON in response",
			response: "This is not a JSON response",
			wantErr:  true,
		},
		{
			name:     "invalid JSON",
			response: `{"isRelevant": true, "message":`,
			wantErr:  true,
		},
		{
			name:         "message truncation",
			response:     `{"isRelevant": true, "message": "` + string(make([]byte, 250)) + `"}`,
			wantRelevant: true,
			wantMessage:  "",
			wantErr:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotRelevant, gotMessage, err := client.parseValidationResponse(tt.response)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValidationResponse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err != nil {
				return
			}
			if gotRelevant != tt.wantRelevant {
				t.Errorf("parseValidationResponse() gotRelevant = %v, want %v", gotRelevant, tt.wantRelevant)
			}
			if tt.wantMessage != "" && gotMessage != tt.wantMessage {
				t.Errorf("parseValidationResponse() gotMessage = %v, want %v", gotMessage, tt.wantMessage)
			}
			if len(gotMessage) > 200 {
				t.Errorf("parseValidationResponse() message length = %d, want <= 200", len(gotMessage))
			}
		})
	}
}

func TestValidateTopicRelevance(t *testing.T) {
	// Skip this test if ANTHROPIC_API_KEY is not set (integration test)
	if os.Getenv("ANTHROPIC_API_KEY") == "" {
		t.Skip("Skipping integration test: ANTHROPIC_API_KEY not set")
	}

	client, err := NewClaudeClient()
	if err != nil {
		t.Fatalf("Failed to create Claude client: %v", err)
	}

	tests := []struct {
		name  string
		topic string
	}{
		{
			name:  "theological topic",
			topic: "Should Christians defy authorities when the law is unfair?",
		},
		{
			name:  "philosophical topic",
			topic: "What is the nature of consciousness?",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := context.Background()
			isRelevant, message, err := client.ValidateTopicRelevance(ctx, tt.topic)
			if err != nil {
				t.Errorf("ValidateTopicRelevance() error = %v", err)
				return
			}
			if message == "" {
				t.Error("ValidateTopicRelevance() returned empty message")
			}
			if len(message) > 200 {
				t.Errorf("ValidateTopicRelevance() message length = %d, want <= 200", len(message))
			}
			t.Logf("Topic: %s\nRelevant: %v\nMessage: %s", tt.topic, isRelevant, message)
		})
	}
}
