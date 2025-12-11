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

	// Accumulate the full response
	for stream.Next() {
		event := stream.Current()
		if event.Delta.Text != "" {
			buffer.WriteString(event.Delta.Text)
		}
	}

	// Check for stream errors
	if err := stream.Err(); err != nil {
		return fmt.Errorf("stream error: %w", err)
	}

	// Parse the complete JSON response
	response := buffer.String()
	
	// Extract JSON from response
	startIdx := strings.Index(response, "{")
	endIdx := strings.LastIndex(response, "}")
	if startIdx == -1 || endIdx == -1 {
		return fmt.Errorf("no JSON found in response")
	}

	jsonStr := response[startIdx : endIdx+1]
	
	var result struct {
		IsRelevant bool       `json:"isRelevant"`
		Message    string     `json:"message"`
		Panelists  []Panelist `json:"panelists"`
	}

	if err := json.Unmarshal([]byte(jsonStr), &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	// Send validation result first
	validationData, _ := json.Marshal(map[string]interface{}{
		"isRelevant": result.IsRelevant,
		"message":    result.Message,
	})
	sendChunk("validation", string(validationData))

	// Stream panelists one by one
	for _, panelist := range result.Panelists {
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

	// Send done signal
	sendChunk("done", "")

	return nil
}
