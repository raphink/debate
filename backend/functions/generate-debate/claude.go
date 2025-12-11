package generatedebate

import (
	"bufio"
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
		model:  "claude-sonnet-4-5",
		httpClient: &http.Client{
			Timeout: 120 * time.Second, // Longer timeout for streaming debates
		},
	}, nil
}

// GenerateDebate streams a debate between the selected panelists
func (c *ClaudeClient) GenerateDebate(ctx context.Context, req *DebateRequest, writer io.Writer) error {
	// Build the debate prompt
	prompt := c.buildDebatePrompt(req)

	// Create the request body with streaming enabled
	requestBody := map[string]interface{}{
		"model":      c.model,
		"max_tokens": 4096,
		"stream":     true,
		"messages": []map[string]string{
			{
				"role":    "user",
				"content": prompt,
			},
		},
	}

	requestJSON, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequestWithContext(ctx, "POST", claudeAPIURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("x-api-key", c.apiKey)
	httpReq.Header.Set("anthropic-version", apiVersion)

	// Make the request
	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to call Claude API: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("Claude API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Stream the response
	return c.streamResponse(resp.Body, writer, req.SelectedPanelists)
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

// streamResponse parses the streaming SSE response and writes formatted chunks
func (c *ClaudeClient) streamResponse(reader io.Reader, writer io.Writer, panelists []Panelist) error {
	scanner := bufio.NewScanner(reader)
	var buffer strings.Builder // Buffer for incomplete patterns only
	var currentPanelistID string
	var currentMessage strings.Builder
	flusher, _ := writer.(http.Flusher)

	fmt.Printf("[DEBUG] Starting stream parsing\n")

	for scanner.Scan() {
		line := scanner.Text()

		// Skip empty lines and event type lines
		if line == "" || strings.HasPrefix(line, "event:") {
			continue
		}

		// Parse data lines
		if strings.HasPrefix(line, "data: ") {
			data := strings.TrimPrefix(line, "data: ")

			// Check for stream end
			if data == "[DONE]" {
				fmt.Printf("[DEBUG] Stream end detected\n")
				// Send any remaining message
				if currentPanelistID != "" && currentMessage.Len() > 0 {
					finalText := currentMessage.String()
					fmt.Printf("[DEBUG] Final message for %s: %q\n", currentPanelistID, finalText)
					chunk := StreamChunk{
						Type:       "message",
						PanelistID: currentPanelistID,
						Text:       finalText,
						Done:       false,
					}
					json.NewEncoder(writer).Encode(chunk)
					if flusher != nil {
						flusher.Flush()
					}
				}

				chunk := StreamChunk{
					Type: "done",
					Done: true,
				}
				json.NewEncoder(writer).Encode(chunk)
				if flusher != nil {
					flusher.Flush()
				}
				break
			}

			// Parse JSON data
			var eventData map[string]interface{}
			if err := json.Unmarshal([]byte(data), &eventData); err != nil {
				fmt.Printf("[DEBUG] JSON parse error: %v\n", err)
				continue
			}

			// Extract text from content delta
			if delta, ok := eventData["delta"].(map[string]interface{}); ok {
				if text, ok := delta["text"].(string); ok {
					fmt.Printf("[DEBUG] Received text chunk: %q\n", text)
					
					// Add to buffer
					buffer.WriteString(text)
					bufferText := buffer.String()
					fmt.Printf("[DEBUG] Buffer contains: %q\n", bufferText)
					
					// Process ALL complete patterns in buffer
					for {
						panelistID, messageText := c.parseMessage(bufferText)
						if panelistID == "" {
							// No complete pattern found
							if currentPanelistID != "" {
								// We have an active speaker, add buffer content to their message
								fmt.Printf("[DEBUG] Appending buffer to %s's message: %q\n", currentPanelistID, bufferText)
								currentMessage.WriteString(bufferText)
								buffer.Reset()
							} else {
								// No active speaker yet, keep in buffer for next iteration
								fmt.Printf("[DEBUG] No speaker yet, keeping in buffer\n")
							}
							break // Exit loop, wait for more data
						}

						fmt.Printf("[DEBUG] Detected new speaker: %s, message start: %q\n", panelistID, messageText[:min(50, len(messageText))])
						
						// Send previous message if exists
						if currentPanelistID != "" && currentMessage.Len() > 0 {
							finalText := strings.TrimSpace(currentMessage.String())
							fmt.Printf("[DEBUG] Sending complete message for %s: %q\n", currentPanelistID, finalText)
							chunk := StreamChunk{
								Type:       "message",
								PanelistID: currentPanelistID,
								Text:       finalText,
								Done:       false,
							}
							json.NewEncoder(writer).Encode(chunk)
							if flusher != nil {
								flusher.Flush()
							}
						}

					// Start new message
					fmt.Printf("[DEBUG] Starting new message for %s\n", panelistID)
					currentPanelistID = panelistID
					currentMessage.Reset()
					
					// Check if messageText itself contains another pattern
					pos, nextID, _ := c.findNextPattern(messageText, 0)
					if pos > 0 && nextID != "" {
						// Multiple patterns in this chunk! 
						// Current message is everything before the next pattern
						currentMessage.WriteString(strings.TrimSpace(messageText[:pos]))
						// Update buffer to continue processing from next pattern
						bufferText = messageText[pos:]
						buffer.Reset()
						buffer.WriteString(bufferText)
						// Continue loop to process next pattern
					} else {
						// No additional pattern, this is the message content
						currentMessage.WriteString(messageText)
						buffer.Reset()
							break // Exit loop, wait for more chunks
						}
					}
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading stream: %w", err)
	}

	// Send any remaining message when stream ends
	if currentPanelistID != "" && currentMessage.Len() > 0 {
		chunk := StreamChunk{
			Type:       "message",
			PanelistID: currentPanelistID,
			Text:       stripTrailingPatterns(currentMessage.String()),
			Done:       false,
		}
		json.NewEncoder(writer).Encode(chunk)
		if flusher != nil {
			flusher.Flush()
		}
	}

	return nil
}

// findNextPattern searches for a [ID]: pattern starting from position start+1
// Returns the position of '[', the panelist ID, and the message text after the pattern
func (c *ClaudeClient) findNextPattern(text string, start int) (pos int, panelistID, messageText string) {
	searchText := text[start:]

	// Look for '[' character
	for i := 1; i < len(searchText); i++ {
		if searchText[i] == '[' {
			// Try to parse from this position
			testText := searchText[i:]
			if id, msg := c.parseMessage(testText); id != "" {
				return start + i, id, msg
			}
		}
	}

	return -1, "", ""
}

// stripTrailingPatterns removes any [ID]: patterns from the end of the text
// This handles cases where we've accumulated text that includes the start of the next speaker's message
func stripTrailingPatterns(text string) string {
	// Look for any [ID]: pattern in the text
	for i := 0; i < len(text); i++ {
		if text[i] == '[' {
			// Check if this looks like a pattern
			testText := text[i:]
			// Try to parse from this position
			if idx := strings.Index(testText, "]: "); idx != -1 {
				if startIdx := strings.LastIndex(testText[:idx], "["); startIdx == 0 {
					// Found a pattern! Return everything before it
					return strings.TrimSpace(text[:i])
				}
			}
			if idx := strings.Index(testText, "]:"); idx != -1 {
				if startIdx := strings.LastIndex(testText[:idx], "["); startIdx == 0 {
					// Found a pattern! Return everything before it
					return strings.TrimSpace(text[:i])
				}
			}
		}
	}
	return text
}

// parseMessage extracts panelist ID and message text from formatted response
func (c *ClaudeClient) parseMessage(text string) (panelistID, messageText string) {
	// Look for pattern: [PANELIST_ID]: Message text or [PANELIST_ID]:Message text
	// Try with space after colon first
	if idx := strings.Index(text, "]: "); idx != -1 {
		if startIdx := strings.LastIndex(text[:idx], "["); startIdx != -1 {
			panelistID = text[startIdx+1 : idx]
			messageText = strings.TrimSpace(text[idx+3:]) // Skip ]: and space
			fmt.Printf("[DEBUG] parseMessage found (with space): ID=%s, msgStart=%q\n", panelistID, messageText[:min(50, len(messageText))])
			return panelistID, messageText
		}
	}

	// Fallback to no space after colon
	if idx := strings.Index(text, "]:"); idx != -1 {
		if startIdx := strings.LastIndex(text[:idx], "["); startIdx != -1 {
			panelistID = text[startIdx+1 : idx]
			messageText = strings.TrimSpace(text[idx+2:]) // Skip ]:
			fmt.Printf("[DEBUG] parseMessage found (no space): ID=%s, msgStart=%q\n", panelistID, messageText[:min(50, len(messageText))])
			return panelistID, messageText
		}
	}

	return "", ""
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
