package validatetopic

import (
	"context"
	"encoding/json"
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
		return nil, fmt.Errorf("ANTHROPIC_API_KEY environment variable not set")
	}

	return &ClaudeClient{
		client: anthropic.NewClient(
			option.WithAPIKey(apiKey),
		),
	}, nil
}

// ValidateTopicAndSuggestPanelists validates topic and streams panelist suggestions
func (c *ClaudeClient) ValidateTopicAndSuggestPanelists(ctx context.Context, topic string, suggestedNames []string, writer io.Writer) error {
	// Build user-suggested names section
	namesSection := ""
	if len(suggestedNames) > 0 {
		namesSection = "\n\nIMPORTANT - User has specifically requested these panelists:\n"
		for _, name := range suggestedNames {
			if name != "" {
				namesSection += fmt.Sprintf("- %s\n", name)
			}
		}
		namesSection += `
PRIORITY REQUIREMENT: You MUST include these individuals in your panelist list if they meet ANY of these criteria:
1. They have ANY documented historical views or writings (even if not directly on this topic)
2. Their philosophical/theological tradition would give them a position on this topic
3. They are recognized figures in theology, philosophy, or related fields

Only exclude a user-suggested panelist if:
- They are completely fictional/non-existent
- They have absolutely no connection to theology, philosophy, or intellectual discourse
- Their inclusion would be factually impossible (e.g., purely a fictional character)

For user-suggested panelists who qualify, infer their likely position based on their known works and tradition, even if they never directly addressed this specific topic.
`
	}

	// Build the combined prompt for Claude
	prompt := fmt.Sprintf(`You are an expert in theology and philosophy. Your task is to evaluate if a topic is suitable for a theological or philosophical debate, and if so, suggest panelists.

Topic: "%s"%s

First, evaluate whether this topic relates to:
- Theology (study of God, religion, faith, sacred texts)
- Philosophy (ethics, morality, epistemology, metaphysics, logic)
- Religious practice and doctrine
- Moral and ethical questions
- Existential questions

If the topic is NOT relevant for theological/philosophical debate:
Return ONLY this JSON structure:
{
  "type": "rejection",
  "message": "brief explanation why this topic is not suitable (max 200 chars)"
}

If the topic IS relevant:
Return 8-20 panelist objects in this format, one per line:
{"type":"panelist","data":{"id":"unique-kebab-case-id","name":"Full Name","tagline":"One-line with era (max 60 chars)","bio":"Brief bio (max 300 chars)","avatarUrl":"placeholder-avatar.svg","position":"Their position (max 100 chars)"}}

Requirements for panelists:
- Different theological/philosophical positions on this topic
- Mix of time periods: 25%% ancient/early church (0-500 AD), 25%% medieval/reformation (500-1700 AD), 25%% modern (1700-1950 AD), 25%% contemporary (1950-present)
- Different traditions (Catholic, Protestant, Orthodox, Jewish, Islamic, secular, etc.)
- Mix of perspectives (theist/atheist, conservative/progressive, different schools of thought)
- Only include historical/contemporary figures with known, documented views on related topics

Format: Each panelist on its own line as shown above. No other text.`, topic, namesSection)

	// Create streaming request
	stream := c.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: 4096,
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.NewTextBlock(prompt)),
		},
	})

	// Stream the response
	return c.streamPanelistResponse(stream, writer)
}

// streamPanelistResponse processes the stream and emits panelists or rejection incrementally
func (c *ClaudeClient) streamPanelistResponse(stream *ssestream.Stream[anthropic.MessageStreamEventUnion], writer io.Writer) error {
	flusher, _ := writer.(http.Flusher)

	sendChunk := func(chunkType, data string) {
		chunk := map[string]string{
			"type": chunkType,
			"data": data,
		}
		json.NewEncoder(writer).Encode(chunk)
		if flusher != nil {
			flusher.Flush()
		}
	}

	var lineBuffer strings.Builder
	var fullBuffer strings.Builder

	// Process stream incrementally, emitting complete lines as they arrive
	for stream.Next() {
		event := stream.Current()
		if event.Delta.Text == "" {
			continue
		}

		text := event.Delta.Text
		fullBuffer.WriteString(text)
		fmt.Printf("[DEBUG] Received text chunk: %s\n", text)
		
		// Process character by character to detect complete lines
		for _, char := range text {
			lineBuffer.WriteRune(char)
			
			// Check if we have a complete JSON object (must start with { and end with }})
			currentLine := lineBuffer.String()
			if char == '\n' && strings.TrimSpace(currentLine) != "" {
				line := strings.TrimSpace(currentLine)
				fmt.Printf("[DEBUG] Newline detected, line: %s\n", line)
				lineBuffer.Reset()
				
				// Skip empty lines or standalone closing braces
				if line == "" || line == "}" || line == "{" || !strings.HasPrefix(line, "{") {
					continue
				}

				// Try to parse as a chunk
				var chunk struct {
					Type    string          `json:"type"`
					Data    json.RawMessage `json:"data"`
					Message string          `json:"message"`
				}

				if err := json.Unmarshal([]byte(line), &chunk); err != nil {
					fmt.Printf("[DEBUG] Failed to parse line as chunk: %v (line was: %s)\n", err, line)
					continue // Skip malformed lines
				}

				fmt.Printf("[DEBUG] Parsed chunk type: %s\n", chunk.Type)

				if chunk.Type == "rejection" {
					// Send rejection message
					rejectionData, _ := json.Marshal(map[string]interface{}{
						"isRelevant": false,
						"message":    chunk.Message,
					})
					sendChunk("validation", string(rejectionData))
					// Continue processing in case there's more
				} else if chunk.Type == "panelist" {
					// Parse and send panelist immediately
					var panelist Panelist
					if err := json.Unmarshal(chunk.Data, &panelist); err != nil {
						fmt.Printf("[DEBUG] Failed to parse panelist data: %v\n", err)
						continue
					}

					fmt.Printf("[DEBUG] Parsed panelist: %s\n", panelist.Name)

					// Validate and sanitize
					if panelist.Name == "" || panelist.ID == "" {
						continue
					}
					if len(panelist.Tagline) > 60 {
						panelist.Tagline = panelist.Tagline[:57] + "..."
					}
					if len(panelist.Bio) > 300 {
						panelist.Bio = panelist.Bio[:297] + "..."
					}
					if len(panelist.Position) > 100 {
						panelist.Position = panelist.Position[:97] + "..."
					}

					panelistData, _ := json.Marshal(panelist)
					sendChunk("panelist", string(panelistData))
				}
			}
		}
	}

	// After stream completes, check if we accumulated text but didn't emit anything via line parsing
	// This handles the case where Claude returns the old format instead of line-delimited
	fullText := fullBuffer.String()
	fmt.Printf("[DEBUG] Full accumulated text: %s\n", fullText)
	
	if fullText != "" {
		// Try to parse as old format (single JSON object with isRelevant, message, panelists array)
		var oldFormat struct {
			IsRelevant bool       `json:"isRelevant"`
			Message    string     `json:"message"`
			Panelists  []Panelist `json:"panelists"`
		}

		if err := json.Unmarshal([]byte(fullText), &oldFormat); err == nil {
			fmt.Printf("[DEBUG] Parsed as old format - isRelevant: %v, panelists: %d\n", oldFormat.IsRelevant, len(oldFormat.Panelists))
			
			// Send validation result
			validationData, _ := json.Marshal(map[string]interface{}{
				"isRelevant": oldFormat.IsRelevant,
				"message":    oldFormat.Message,
			})
			sendChunk("validation", string(validationData))

			// Send each panelist
			for _, panelist := range oldFormat.Panelists {
				panelistJSON, _ := json.Marshal(panelist)
				sendChunk("panelist", string(panelistJSON))
			}
		} else {
			fmt.Printf("[DEBUG] Failed to parse as old format: %v\n", err)
		}
	}

	// Send done signal
	sendChunk("done", "")

	return nil
}
