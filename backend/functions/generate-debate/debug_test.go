package generatedebate

import (
	"strings"
	"testing"
)

func TestParseMessageDebug(t *testing.T) {
	client := &ClaudeClient{}

	tests := []struct {
		name    string
		input   string
		wantID  string
		wantMsg string
	}{
		{
			name:    "moderator with space",
			input:   "[moderator]: Welcome to our debate",
			wantID:  "moderator",
			wantMsg: "Welcome to our debate",
		},
		{
			name:    "moderator no space",
			input:   "[moderator]:Welcome to our debate",
			wantID:  "moderator",
			wantMsg: "Welcome to our debate",
		},
		{
			name:    "panelist with space",
			input:   "[augustine]: I believe that...",
			wantID:  "augustine",
			wantMsg: "I believe that...",
		},
		{
			name:    "partial pattern",
			input:   "[moderator",
			wantID:  "",
			wantMsg: "",
		},
		{
			name:    "partial pattern with bracket",
			input:   "[moderator]",
			wantID:  "",
			wantMsg: "",
		},
		{
			name:    "multiple messages",
			input:   "[moderator]: Welcome\n[augustine]: Thank you",
			wantID:  "moderator",
			wantMsg: "Welcome\n[augustine]: Thank you",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotID, gotMsg := client.parseMessage(tt.input)
			if gotID != tt.wantID {
				t.Errorf("parseMessage() gotID = %v, want %v", gotID, tt.wantID)
			}
			if gotMsg != tt.wantMsg {
				t.Errorf("parseMessage() gotMsg = %v, want %v", gotMsg, tt.wantMsg)
			}
		})
	}
}

func TestStripTrailingPatterns(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "no pattern",
			input: "This is a complete message.",
			want:  "This is a complete message.",
		},
		{
			name:  "trailing pattern with space",
			input: "This is a complete message. [moderator]: ",
			want:  "This is a complete message.",
		},
		{
			name:  "trailing pattern no space",
			input: "This is a complete message. [moderator]:",
			want:  "This is a complete message.",
		},
		{
			name:  "partial trailing pattern",
			input: "This is a complete message. [mod",
			want:  "This is a complete message. [mod",
		},
		{
			name:  "pattern in middle",
			input: "Start [moderator]: Middle text",
			want:  "Start",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripTrailingPatterns(tt.input)
			if got != tt.want {
				t.Errorf("stripTrailingPatterns() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestStreamingSimulation(t *testing.T) {
	// Simulate how Claude might stream a debate
	chunks := []string{
		"[modera",
		"tor]: Welc",
		"ome to our debate on ethics. ",
		"Today we have Augustine and ",
		"Aquinas.\n\n[augustine",
		"]: Thank you for ",
		"having me. I believe",
		" that...\n\n[aquin",
		"as]: I respectfully ",
		"disagree because...",
	}

	var buffer strings.Builder // Simulates the buffer in streamResponse
	var currentPanelist string
	var currentMessage strings.Builder
	client := &ClaudeClient{}

	t.Log("=== Simulating streaming chunks with buffer-based approach ===")
	for i, chunk := range chunks {
		buffer.WriteString(chunk)
		bufferText := buffer.String()

		t.Logf("\n--- Chunk %d: %q ---", i, chunk)
		t.Logf("Buffer: %q", bufferText)

		id, msg := client.parseMessage(bufferText)
		if id != "" {
			if currentPanelist != "" {
				t.Logf("✓ SEND message for [%s]: %q", currentPanelist, currentMessage.String())
			}
			t.Logf("→ NEW speaker detected: [%s]", id)
			currentPanelist = id
			currentMessage.Reset()
			currentMessage.WriteString(msg)
			buffer.Reset()
			t.Logf("Message starts: %q", msg)
		} else if currentPanelist != "" {
			t.Logf("+ Append to [%s]", currentPanelist)
			currentMessage.WriteString(bufferText)
			buffer.Reset()
		} else {
			t.Logf("⏳ Waiting for pattern, keeping in buffer")
		}
	}

	if currentPanelist != "" {
		t.Logf("\n✓ FINAL SEND for [%s]: %q", currentPanelist, currentMessage.String())
	}
}
