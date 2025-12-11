package generatedebate

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
)

// ClaudeClient handles communication with the Anthropic Claude API
type ClaudeClient struct {
	client anthropic.Client
}

// NewClaudeClient creates a new Claude API client
func NewClaudeClient() (*ClaudeClient, error) {
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, errors.New("ANTHROPIC_API_KEY environment variable not set")
	}

	return &ClaudeClient{
		client: anthropic.NewClient(
			option.WithAPIKey(apiKey),
		),
	}, nil
}

// GenerateDebate streams a debate between the selected panelists
func (c *ClaudeClient) GenerateDebate(ctx context.Context, req *DebateRequest, writer io.Writer) error {
	// Build the debate prompt
	prompt := c.buildDebatePrompt(req)

	// Create streaming request
	stream := c.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeSonnet4_5,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	// Stream the response
	return c.streamResponse(stream, writer)
}

// buildDebatePrompt creates the prompt for Claude to generate a debate
func (c *ClaudeClient) buildDebatePrompt(req *DebateRequest) string {
	var prompt strings.Builder

	prompt.WriteString("You are a neutral moderator orchestrating a theological/philosophical debate between historical figures.\n\n")
	prompt.WriteString(fmt.Sprintf("Topic: %s\n\n", req.Topic))
	prompt.WriteString("Panelists:\n")

	for i, panelist := range req.SelectedPanelists {
		prompt.WriteString(fmt.Sprintf("%d. %s (ID: %s)\n", i+1, panelist.Name, panelist.ID))
		if panelist.Bio != "" {
			prompt.WriteString(fmt.Sprintf("   Bio: %s\n", panelist.Bio))
		}
		if panelist.Position != "" {
			prompt.WriteString(fmt.Sprintf("   Position: %s\n", panelist.Position))
		}
	}

	prompt.WriteString("\nGenerate a moderated debate with the following structure:\n")
	prompt.WriteString("1. FIRST MESSAGE MUST BE: [moderator]: (introducing the topic and panelists)\n")
	prompt.WriteString("2. Include 12-18 exchanges between panelists\n")
	prompt.WriteString("3. The moderator may occasionally intervene between panelist exchanges to:\n")
	prompt.WriteString("   - Redirect the conversation\n")
	prompt.WriteString("   - Ask clarifying questions\n")
	prompt.WriteString("   - Highlight contrasting viewpoints\n")
	prompt.WriteString("   - Summarize progress\n")
	prompt.WriteString("4. LAST MESSAGE MUST BE: [moderator]: (providing a concluding summary that synthesizes the key points, acknowledges different perspectives, and gracefully ends the debate - 3-5 sentences)\n\n")
	prompt.WriteString("CRITICAL FORMAT REQUIREMENTS:\n")
	prompt.WriteString("- Each response MUST start on a new line with the exact format: [ID]: text\n")
	prompt.WriteString("- Use [moderator]: for moderator messages\n")
	prompt.WriteString("- Use [PANELIST_ID]: for panelist messages (IDs listed above)\n")
	prompt.WriteString("- NO extra text before the [ID]: marker\n")
	prompt.WriteString("- Start your response immediately with [moderator]:\n\n")
	prompt.WriteString("Guidelines:\n")
	prompt.WriteString("- Moderator responses: 1-3 sentences, neutral and facilitating\n")
	prompt.WriteString("- Panelist responses: 2-4 sentences (50-100 words)\n")
	prompt.WriteString("- Maintain each panelist's historical perspective and known positions\n")
	prompt.WriteString("- Create engaging exchanges with direct responses and counter-arguments\n")
	prompt.WriteString("- Let panelists speak to each other directly, not just to the moderator\n")
	prompt.WriteString("- Moderator should intervene naturally, not after every exchange\n")
	prompt.WriteString("- Ensure philosophical depth while remaining accessible\n\n")
	prompt.WriteString("Begin the debate:")

	return prompt.String()
}

// streamResponse processes the SDK stream and writes formatted chunks
func (c *ClaudeClient) streamResponse(stream *ssestream.Stream[anthropic.MessageStreamEventUnion], writer io.Writer) error {
	var patternBuffer strings.Builder // Buffer only for incomplete [handle]: patterns
	var currentSpeaker string
	inPattern := false
	flusher, _ := writer.(http.Flusher)

	sendChunk := func(speaker, text string) {
		if text == "" {
			return
		}
		chunk := StreamChunk{
			Type:       "message",
			PanelistID: speaker,
			Text:       text,
			Done:       false,
		}
		json.NewEncoder(writer).Encode(chunk)
		if flusher != nil {
			flusher.Flush()
		}
	}

	for stream.Next() {
		event := stream.Current()

		if event.Delta.Text == "" {
			continue
		}

		text := event.Delta.Text

		// Process character by character (runes, not bytes - handles UTF-8 correctly)
		for _, char := range text {
			if char == '[' && !inPattern {
				// Start of potential pattern - flush any accumulated text first
				if patternBuffer.Len() > 0 && currentSpeaker != "" {
					sendChunk(currentSpeaker, patternBuffer.String())
					patternBuffer.Reset()
				}
				inPattern = true
				patternBuffer.WriteRune(char)
			} else if inPattern {
				patternBuffer.WriteRune(char)

				// Check if we have a complete pattern [handle]:
				bufStr := patternBuffer.String()
				if strings.HasSuffix(bufStr, "]: ") || strings.HasSuffix(bufStr, "]:") {
					// Extract speaker ID
					var newSpeaker string
					if idx := strings.Index(bufStr, "]: "); idx != -1 {
						newSpeaker = bufStr[1:idx] // Skip opening [
					} else if idx := strings.Index(bufStr, "]:"); idx != -1 {
						newSpeaker = bufStr[1:idx]
					}

					if newSpeaker != "" {
						currentSpeaker = newSpeaker
						patternBuffer.Reset()
						inPattern = false
					}
				}
			} else {
				// Normal text - accumulate and send immediately
				patternBuffer.WriteRune(char)
				
				// Send chunks frequently for responsiveness (every 10 chars or at word boundaries)
				if patternBuffer.Len() >= 10 || char == ' ' || char == '\n' {
					if currentSpeaker != "" {
						sendChunk(currentSpeaker, patternBuffer.String())
						patternBuffer.Reset()
					}
				}
			}
		}

		// Flush any remaining buffered text at end of each delta
		if patternBuffer.Len() > 0 && currentSpeaker != "" && !inPattern {
			sendChunk(currentSpeaker, patternBuffer.String())
			patternBuffer.Reset()
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	// Send any final buffered text
	if patternBuffer.Len() > 0 && currentSpeaker != "" {
		sendChunk(currentSpeaker, patternBuffer.String())
	}

	// Send done signal
	chunk := StreamChunk{
		Type: "done",
		Done: true,
	}
	json.NewEncoder(writer).Encode(chunk)
	if flusher != nil {
		flusher.Flush()
	}

	return nil
}
