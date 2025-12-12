package generatedebate

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/raphink/debate/shared/firebase"
)

// DebateMessage represents an accumulated message
type DebateMessage struct {
	ID           string
	PanelistID   string
	PanelistName string
	AvatarURL    string
	Text         string
	Sequence     int
	Timestamp    time.Time
}

// DebateAccumulator accumulates debate messages during streaming
type DebateAccumulator struct {
	DebateID        string
	Topic           string
	Panelists       []Panelist
	Messages        []DebateMessage
	PanelistMap     map[string]Panelist
	CurrentSequence int
	StartedAt       time.Time
}

// NewDebateAccumulator creates a new accumulator
func NewDebateAccumulator(debateID, topic string, panelists []Panelist) *DebateAccumulator {
	panelistMap := make(map[string]Panelist)
	for _, p := range panelists {
		panelistMap[p.ID] = p
	}
	// Add moderator
	panelistMap["moderator"] = Panelist{
		ID:        "moderator",
		Name:      "Moderator",
		Tagline:   "Neutral Facilitator",
		Bio:       "Guiding the conversation",
		AvatarURL: "/avatars/moderator-avatar.png",
	}

	return &DebateAccumulator{
		DebateID:        debateID,
		Topic:           topic,
		Panelists:       panelists,
		Messages:        make([]DebateMessage, 0),
		PanelistMap:     panelistMap,
		CurrentSequence: 0,
		StartedAt:       time.Now(),
	}
}

// AddMessage accumulates a message chunk
func (acc *DebateAccumulator) AddMessage(panelistID, text string) {
	// Find if we're continuing the last message or starting a new one
	if len(acc.Messages) > 0 {
		lastMsg := &acc.Messages[len(acc.Messages)-1]
		if lastMsg.PanelistID == panelistID {
			// Continue existing message
			lastMsg.Text += text
			return
		}
	}

	// New message
	panelist, ok := acc.PanelistMap[panelistID]
	if !ok {
		log.Printf("Unknown panelist ID: %s", panelistID)
		panelist = Panelist{
			ID:        panelistID,
			Name:      panelistID,
			AvatarURL: "/avatars/placeholder-avatar.png",
		}
	}

	msg := DebateMessage{
		ID:           panelistID + "-" + string(rune(acc.CurrentSequence)),
		PanelistID:   panelistID,
		PanelistName: panelist.Name,
		AvatarURL:    panelist.AvatarURL,
		Text:         text,
		Sequence:     acc.CurrentSequence,
		Timestamp:    time.Now(),
	}

	acc.Messages = append(acc.Messages, msg)
	acc.CurrentSequence++
}

// AccumulatingWriter wraps an http.ResponseWriter to accumulate messages
type AccumulatingWriter struct {
	writer      http.ResponseWriter
	accumulator *DebateAccumulator
}

// Write intercepts writes to accumulate message chunks
func (aw *AccumulatingWriter) Write(p []byte) (n int, err error) {
	// Try to parse as StreamChunk
	var chunk StreamChunk
	if err := json.Unmarshal(p, &chunk); err == nil {
		if chunk.Type == "message" && chunk.PanelistID != "" && chunk.Text != "" {
			// Accumulate this message
			aw.accumulator.AddMessage(chunk.PanelistID, chunk.Text)
		}
	}

	// Pass through to original writer
	return aw.writer.Write(p)
}

// Header passes through to original writer
func (aw *AccumulatingWriter) Header() http.Header {
	return aw.writer.Header()
}

// WriteHeader passes through to original writer
func (aw *AccumulatingWriter) WriteHeader(statusCode int) {
	aw.writer.WriteHeader(statusCode)
}

// Flush passes through flush if available
func (aw *AccumulatingWriter) Flush() {
	if flusher, ok := aw.writer.(http.Flusher); ok {
		flusher.Flush()
	}
}

// saveDebateToFirestore saves the completed debate asynchronously
func saveDebateToFirestore(ctx context.Context, acc *DebateAccumulator, userAgent string) {
	if firebase.GetClient() == nil {
		log.Println("Firestore not initialized, skipping debate save")
		return
	}

	// Clean up message text (trim whitespace)
	messages := make([]firebase.Message, len(acc.Messages))
	for i, msg := range acc.Messages {
		messages[i] = firebase.Message{
			ID:           msg.ID,
			PanelistID:   msg.PanelistID,
			PanelistName: msg.PanelistName,
			AvatarURL:    msg.AvatarURL,
			Text:         strings.TrimSpace(msg.Text),
			Timestamp:    msg.Timestamp,
			Sequence:     msg.Sequence,
			IsComplete:   true,
		}
	}

	// Convert panelists
	panelists := make([]firebase.Panelist, len(acc.Panelists))
	for i, p := range acc.Panelists {
		panelists[i] = firebase.Panelist{
			ID:        p.ID,
			Name:      p.Name,
			Tagline:   p.Tagline,
			Biography: p.Bio,
			AvatarURL: p.AvatarURL,
			Position:  p.Position,
		}
	}

	// Create debate document
	debate := firebase.DebateDocument{
		ID: acc.DebateID,
		Topic: firebase.Topic{
			Text:       acc.Topic,
			IsRelevant: true,
		},
		Panelists:   panelists,
		Messages:    messages,
		Status:      "complete",
		StartedAt:   acc.StartedAt,
		CompletedAt: time.Now(),
		Metadata: firebase.Metadata{
			CreatedBy:   "anonymous",
			UserAgent:   userAgent,
			Version:     "1.0",
			GeneratedBy: "backend",
		},
	}

	// Save to Firestore
	if err := firebase.SaveDebate(ctx, acc.DebateID, &debate); err != nil {
		log.Printf("Failed to save debate to Firestore (debate still succeeded): %v", err)
	} else {
		log.Printf("Successfully saved debate %s to Firestore", acc.DebateID)
	}
}
