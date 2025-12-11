package main

import (
	"bytes"
	"encoding/json"
	"strings"
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
			name:     "Single complete message",
			input:    `[moderator]: Welcome to the debate`,
			wantID:   "moderator",
			wantText: "Welcome to the debate",
			desc:     "Should extract first message",
		},
		{
			name: "Two complete messages - returns first with rest of text",
			input: `[moderator]: Welcome to the debate
[Augustine354]: I believe divine law supersedes human law.`,
			wantID:   "moderator",
			wantText: "Welcome to the debate\n[Augustine354]: I believe divine law supersedes human law.",
			desc:     "parseMessage extracts first ID and returns ALL remaining text (including next [ID]:)",
		},
		{
			name: "Three panelists - returns first with all text",
			input: `[moderator]: Let's begin
[Augustine354]: My position is clear
[MLKJr]: I respectfully disagree`,
			wantID:   "moderator",
			wantText: "Let's begin\n[Augustine354]: My position is clear\n[MLKJr]: I respectfully disagree",
			desc:     "parseMessage returns first ID and all text including subsequent patterns",
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

func TestStreamResponseMultipleSpeakers(t *testing.T) {
	client := &ClaudeClient{}

	tests := []struct {
		name           string
		mockSSEData    string
		expectedChunks []StreamChunk
	}{
		{
			name: "Single speaker in one chunk",
			// This simulates Claude sending: [moderator]: Welcome to the debate
			// in a SINGLE SSE event
			mockSSEData: `data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"[moderator]: Welcome to the debate"}}

`,
			// We expect ONE output chunk (one speech bubble)
			expectedChunks: []StreamChunk{
				{Type: "message", PanelistID: "moderator", Text: "Welcome to the debate", Done: false},
			},
		},
		{
			name: "Two speakers in one chunk",
			// ⭐ KEY TEST: Claude sends TWO speaker patterns in ONE SSE event:
			//    "[moderator]: Welcome\n[Augustine354]: Thank you"
			// This is the edge case we're testing!
			mockSSEData: `data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"[moderator]: Welcome\n[Augustine354]: Thank you"}}

`,
			// ⭐ We expect TWO output chunks (two separate speech bubbles)
			// even though it came in ONE SSE chunk from Claude
			expectedChunks: []StreamChunk{
				{Type: "message", PanelistID: "moderator", Text: "Welcome", Done: false},
				{Type: "message", PanelistID: "Augustine354", Text: "Thank you", Done: false},
			},
		},
		{
			name: "Three speakers in one chunk",
			// ⭐ EXTREME CASE: THREE speakers in ONE SSE chunk
			//    "[moderator]: Let's begin\n[Augustine354]: I believe in divine law\n[MLKJr]: I advocate nonviolence"
			mockSSEData: `data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"[moderator]: Let's begin\n[Augustine354]: I believe in divine law\n[MLKJr]: I advocate nonviolence"}}

`,
			// ⭐ We expect THREE output chunks (three separate speech bubbles)
			// The loop in streamResponse must iterate to find ALL patterns
			expectedChunks: []StreamChunk{
				{Type: "message", PanelistID: "moderator", Text: "Let's begin", Done: false},
				{Type: "message", PanelistID: "Augustine354", Text: "I believe in divine law", Done: false},
				{Type: "message", PanelistID: "MLKJr", Text: "I advocate nonviolence", Done: false},
			},
		},
		{
			name: "Two speakers with empty line between them",
			// Real-world case: Claude sometimes adds blank lines between speakers
			mockSSEData: `data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"[moderator]: Welcome to today's debate.\n\n[john-macarthur]: This tithe is important."}}

`,
			expectedChunks: []StreamChunk{
				{Type: "message", PanelistID: "moderator", Text: "Welcome to today's debate.", Done: false},
				{Type: "message", PanelistID: "john-macarthur", Text: "This tithe is important.", Done: false},
			},
		},
		{
			name: "Speaker change across chunks",
			// This tests speaker changes across MULTIPLE SSE events (normal case)
			// Event 1: "[moderator]: Welcome to "
			// Event 2: "the debate\n[Augustine354]: "
			// Event 3: "Thank you moderator"
			mockSSEData: `data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"[moderator]: Welcome to "}}

data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"the debate\n[Augustine354]: "}}

data: {"type":"content_block_delta","delta":{"type":"text_delta","text":"Thank you moderator"}}

`,
			expectedChunks: []StreamChunk{
				{Type: "message", PanelistID: "moderator", Text: "Welcome to", Done: false}, // Sent when Augustine starts
				{Type: "message", PanelistID: "Augustine354", Text: "Thank you moderator", Done: false},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a mock reader with SSE data
			reader := strings.NewReader(tt.mockSSEData)

			// Create a buffer to capture output
			// streamResponse will write JSON-encoded StreamChunk objects here,
			// one per line, for each separate message it detects
			var output bytes.Buffer

			// Call streamResponse - this is where the magic happens!
			// It reads the SSE data, detects speaker patterns, and writes
			// separate StreamChunk objects for each speaker
			// Call streamResponse - this is where the magic happens!
			// It reads the SSE data, detects speaker patterns, and writes
			// separate StreamChunk objects for each speaker
			err := client.streamResponse(reader, &output, []Panelist{})
			if err != nil {
				t.Fatalf("streamResponse failed: %v", err)
			}

			// Parse output into chunks
			// The output buffer now contains JSON lines like:
			// {"type":"message","panelistId":"moderator","text":"Welcome","done":false}
			// {"type":"message","panelistId":"Augustine354","text":"Thank you","done":false}
			// Each line represents ONE speech bubble that should appear in the UI
			lines := strings.Split(output.String(), "\n")
			var gotChunks []StreamChunk

			for _, line := range lines {
				if line == "" {
					continue
				}

				var chunk StreamChunk
				if err := json.Unmarshal([]byte(line), &chunk); err != nil {
					continue
				}

				// Skip "done" chunks for this test
				if chunk.Type == "done" {
					continue
				}

				gotChunks = append(gotChunks, chunk)
			}

			// ⭐ THE KEY ASSERTION: Compare number of chunks
			// If we sent "[moderator]: A\n[Augustine354]: B" (TWO speakers in ONE SSE event)
			// we should get TWO StreamChunks in the output (TWO speech bubbles)
			if len(gotChunks) != len(tt.expectedChunks) {
				t.Errorf("Expected %d chunks, got %d", len(tt.expectedChunks), len(gotChunks))
				t.Logf("Got chunks: %+v", gotChunks)
				return
			}

			// Verify each chunk has correct speaker ID and text
			for i, expected := range tt.expectedChunks {
				got := gotChunks[i]

				if got.PanelistID != expected.PanelistID {
					t.Errorf("Chunk %d: expected PanelistID=%q, got %q", i, expected.PanelistID, got.PanelistID)
				}

				if got.Text != expected.Text {
					t.Errorf("Chunk %d: expected Text=%q, got %q", i, expected.Text, got.Text)
				}
			}
		})
	}
}
