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

func TestParseMessageMultiplePanelists(t *testing.T) {
	client := &ClaudeClient{}

	tests := []struct {
		name     string
		input    string
		wantID   string
		wantText string
		desc     string
	}{
		{
			name: "Single complete message",
			input: `[moderator]: Welcome to the debate`,
			wantID:   "moderator",
			wantText: "Welcome to the debate",
			desc: "Should extract first message",
		},
		{
			name: "Two complete messages - returns first with rest of text",
			input: `[moderator]: Welcome to the debate
[Augustine354]: I believe divine law supersedes human law.`,
			wantID:   "moderator",
			wantText: "Welcome to the debate\n[Augustine354]: I believe divine law supersedes human law.",
			desc: "parseMessage extracts first ID and returns ALL remaining text (including next [ID]:)",
		},
		{
			name: "Three panelists - returns first with all text",
			input: `[moderator]: Let's begin
[Augustine354]: My position is clear
[MLKJr]: I respectfully disagree`,
			wantID:   "moderator",
			wantText: "Let's begin\n[Augustine354]: My position is clear\n[MLKJr]: I respectfully disagree",
			desc: "parseMessage returns first ID and all text including subsequent patterns",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			id, text := client.parseMessage(tt.input)
			
			if id != tt.wantID {
				t.Errorf("parseMessage() id = %q, want %q (%s)", id, tt.wantID, tt.desc)
			}
			
			if text != tt.wantText {
				t.Errorf("parseMessage() text = %q, want %q (%s)", text, tt.wantText, tt.desc)
			}
		})
	}
}

func TestFindNextPanelistInText(t *testing.T) {
	client := &ClaudeClient{}
	
	// Test that we can detect a second [ID]: pattern within already-extracted text
	tests := []struct {
		name         string
		currentID    string
		messageText  string
		expectNextID string
		expectSplit  bool
	}{
		{
			name:         "No second pattern",
			currentID:    "moderator",
			messageText:  "Welcome to the debate",
			expectNextID: "",
			expectSplit:  false,
		},
		{
			name:         "Contains second pattern with space",
			currentID:    "moderator",
			messageText:  "Welcome\n[Augustine354]: Thank you",
			expectNextID: "Augustine354",
			expectSplit:  true,
		},
		{
			name:         "Contains second pattern without space",
			currentID:    "moderator",
			messageText:  "Welcome\n[Augustine354]:Thank you",
			expectNextID: "Augustine354",
			expectSplit:  true,
		},
		{
			name:         "Contains two additional patterns",
			currentID:    "moderator",
			messageText:  "Welcome\n[Augustine354]: Thanks\n[MLKJr]: Indeed",
			expectNextID: "Augustine354",
			expectSplit:  true,
		},
		{
			name:         "Same ID appears again - should not split",
			currentID:    "moderator",
			messageText:  "Welcome\n[moderator]: Continuing",
			expectNextID: "moderator",
			expectSplit:  false, // Same ID, not a new message
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nextID, _ := client.parseMessage(tt.messageText)
			
			if tt.expectSplit {
				if nextID != tt.expectNextID {
					t.Errorf("Expected to find next ID %q in text, got %q", tt.expectNextID, nextID)
				}
			} else {
				// When not expecting split, either no pattern found or same ID
				if nextID != "" && nextID != tt.currentID && nextID != tt.expectNextID {
					t.Errorf("Expected no different ID, but got %q", nextID)
				}
			}
		})
	}
}

func TestParseMessageStreamingBehavior(t *testing.T) {
	client := &ClaudeClient{}

	// The key insight: parseMessage only finds the FIRST [ID]: pattern
	// The streaming loop handles detecting when a NEW pattern starts
	// and sends the accumulated message before starting the next one
	
	t.Run("Returns only first message from block", func(t *testing.T) {
		input := `[moderator]: Welcome
[Augustine354]: Thank you`
		
		id, text := client.parseMessage(input)
		
		if id != "moderator" {
			t.Errorf("Expected first panelist 'moderator', got %q", id)
		}
		
		if text != "Welcome\n[Augustine354]: Thank you" {
			t.Errorf("Expected full text after first ID, got %q", text)
		}
	})
}
