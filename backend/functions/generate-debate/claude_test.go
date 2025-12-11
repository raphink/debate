package main

import (
	"testing"
)

func TestParseMessage(t *testing.T) {
	client := &ClaudeClient{}

	tests := []struct {
		name            string
		input           string
		expectedID      string
		expectedMessage string
	}{
		{
			name:            "Standard format with space",
			input:           "[moderator]: Welcome to the debate",
			expectedID:      "moderator",
			expectedMessage: "Welcome to the debate",
		},
		{
			name:            "Standard format without space",
			input:           "[moderator]:Welcome to the debate",
			expectedID:      "moderator",
			expectedMessage: "Welcome to the debate",
		},
		{
			name:            "Panelist ID with numbers",
			input:           "[Augustine354]: I believe divine law supersedes human law.",
			expectedID:      "Augustine354",
			expectedMessage: "I believe divine law supersedes human law.",
		},
		{
			name:            "Message with extra whitespace",
			input:           "[moderator]:   Thank you both   ",
			expectedID:      "moderator",
			expectedMessage: "Thank you both",
		},
		{
			name:            "Incomplete pattern - no closing bracket",
			input:           "[moderator",
			expectedID:      "",
			expectedMessage: "",
		},
		{
			name:            "Incomplete pattern - no colon",
			input:           "[moderator]",
			expectedID:      "",
			expectedMessage: "",
		},
		{
			name:            "Partial bracket only",
			input:           "[",
			expectedID:      "",
			expectedMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, message := client.parseMessage(tt.input)
			
			if id != tt.expectedID {
				t.Errorf("parseMessage() id = %q, want %q", id, tt.expectedID)
			}
			
			if message != tt.expectedMessage {
				t.Errorf("parseMessage() message = %q, want %q", message, tt.expectedMessage)
			}
		})
	}
}
