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
		namesSection = "\n\nUser has suggested considering these individuals (if they have relevant documented views on this topic):\n"
		for _, name := range suggestedNames {
			if name != "" {
				namesSection += fmt.Sprintf("- %s\n", name)
			}
		}
		namesSection += "\nYou may include these individuals in your suggestions if they have known, documented positions relevant to this topic. If not, you may omit them.\n"
	}

	// Build the combined prompt for Claude
	prompt := fmt.Sprintf(`You are an expert in theology and philosophy. Your task is to:
1. Determine if the following topic is suitable for a theological or philosophical debate
2. If suitable, suggest 8-20 historical or contemporary figures who would make excellent panelists

Topic: "%s"%s

First, evaluate whether this topic relates to:
- Theology (study of God, religion, faith, sacred texts)
- Philosophy (ethics, morality, epistemology, metaphysics, logic)
- Religious practice and doctrine
- Moral and ethical questions
- Existential questions

If the topic IS relevant, suggest 8-20 diverse panelists with:
- Different theological/philosophical positions on this topic
- IMPORTANT: Mix of time periods across the last 2000 years (e.g., early church fathers, medieval scholars, reformation thinkers, enlightenment philosophers, modern theologians, contemporary figures)
- Aim for roughly equal representation: 25%% ancient/early church (0-500 AD), 25%% medieval/reformation (500-1700 AD), 25%% modern (1700-1950 AD), 25%% contemporary (1950-present)
- Different traditions (Catholic, Protestant, Orthodox, Jewish, Islamic, secular, etc.)
- Mix of perspectives (theist/atheist, conservative/progressive, different schools of thought)
- Historical or contemporary figures with known, documented views on related topics

Respond with a JSON object:
{
  "isRelevant": boolean,
  "message": "brief explanation (max 200 chars)",
  "panelists": [
    {
      "id": "unique-kebab-case-id",
      "name": "Full Name",
      "tagline": "One-line description with era (max 60 chars, e.g., '4th Century Church Father')",
      "bio": "Brief bio explaining their relevant expertise and historical period (max 300 chars)",
      "avatarUrl": "placeholder-avatar.svg",
      "position": "Brief statement of their likely position on this topic (max 100 chars)"
    }
  ]
}

If not relevant, set panelists to an empty array.
Format your response as valid JSON only, no other text.`, topic, namesSection)

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

// streamPanelistResponse processes the stream and emits validation + panelists incrementally
func (c *ClaudeClient) streamPanelistResponse(stream *ssestream.Stream[anthropic.MessageStreamEventUnion], writer io.Writer) error {
	var buffer strings.Builder
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

	var validationSent bool
	var inPanelistsArray bool
	var panelistBuffer strings.Builder
	var braceDepth int
	var inString bool
	var escapeNext bool

	// Process stream character by character to detect complete panelist objects
	for stream.Next() {
		event := stream.Current()
		if event.Delta.Text == "" {
			continue
		}

		text := event.Delta.Text
		buffer.WriteString(text)

		// Parse incrementally to find complete panelist objects
		for _, char := range text {
			// Track string context to ignore braces inside strings
			if escapeNext {
				escapeNext = false
				panelistBuffer.WriteRune(char)
				continue
			}
			if char == '\\' && inString {
				escapeNext = true
				panelistBuffer.WriteRune(char)
				continue
			}
			if char == '"' {
				inString = !inString
			}

			// Send validation once we find "isRelevant"
			if !validationSent && strings.Contains(buffer.String(), `"isRelevant"`) {
				// Try to parse validation fields
				currentJSON := buffer.String()
				if startIdx := strings.Index(currentJSON, "{"); startIdx != -1 {
					// Look for both isRelevant and message
					if strings.Contains(currentJSON, `"message"`) {
						// Extract up to the panelists array
						endIdx := strings.Index(currentJSON, `"panelists"`)
						if endIdx == -1 {
							endIdx = len(currentJSON)
						}
						
						partialJSON := currentJSON[startIdx:endIdx] + `}`
						
						var validation struct {
							IsRelevant bool   `json:"isRelevant"`
							Message    string `json:"message"`
						}
						
						if err := json.Unmarshal([]byte(partialJSON), &validation); err == nil {
							validationData, _ := json.Marshal(map[string]interface{}{
								"isRelevant": validation.IsRelevant,
								"message":    validation.Message,
							})
							sendChunk("validation", string(validationData))
							validationSent = true
						}
					}
				}
			}

			// Detect when we enter the panelists array
			if !inPanelistsArray && strings.Contains(buffer.String(), `"panelists"`) {
				if char == '[' && !inString {
					inPanelistsArray = true
					braceDepth = 0
					panelistBuffer.Reset()
					continue
				}
			}

			// If we're in the panelists array, track complete objects
			if inPanelistsArray {
				if !inString {
					if char == '{' {
						if braceDepth == 0 {
							panelistBuffer.Reset()
						}
						braceDepth++
					} else if char == '}' {
						braceDepth--
						panelistBuffer.WriteRune(char)
						
						// Complete panelist object found
						if braceDepth == 0 {
							panelistJSON := panelistBuffer.String()
							
							var panelist Panelist
							if err := json.Unmarshal([]byte(panelistJSON), &panelist); err == nil {
								// Validate and sanitize
								if panelist.Name != "" && panelist.ID != "" {
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
							
							panelistBuffer.Reset()
							continue
						}
					}
				}
				
				if braceDepth > 0 {
					panelistBuffer.WriteRune(char)
				}
			}
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	// Send done signal
	sendChunk("done", "")

	return nil
}
