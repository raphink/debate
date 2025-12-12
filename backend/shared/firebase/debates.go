package firebase

import (
	"context"
	"time"
)

// Topic represents the debate topic
type Topic struct {
	Text              string   `firestore:"text"`
	SuggestedNames    []string `firestore:"suggestedNames,omitempty"`
	IsRelevant        bool     `firestore:"isRelevant"`
	ValidationMessage string   `firestore:"validationMessage,omitempty"`
}

// Panelist represents a debate participant
type Panelist struct {
	ID        string `firestore:"id"`
	Name      string `firestore:"name"`
	Tagline   string `firestore:"tagline"`
	Biography string `firestore:"biography"`
	AvatarURL string `firestore:"avatarUrl"`
	Position  string `firestore:"position,omitempty"`
}

// Message represents a single debate contribution
type Message struct {
	ID           string    `firestore:"id"`
	PanelistID   string    `firestore:"panelistId"`
	PanelistName string    `firestore:"panelistName"`
	AvatarURL    string    `firestore:"avatarUrl"`
	Text         string    `firestore:"text"`
	Timestamp    time.Time `firestore:"timestamp"`
	Sequence     int       `firestore:"sequence"`
	IsComplete   bool      `firestore:"isComplete"`
}

// Metadata contains debate metadata
type Metadata struct {
	CreatedBy   string `firestore:"createdBy"`
	UserAgent   string `firestore:"userAgent,omitempty"`
	Version     string `firestore:"version"`
	GeneratedBy string `firestore:"generatedBy"`
}

// DebateDocument represents a complete debate stored in Firestore
type DebateDocument struct {
	ID          string     `firestore:"id"`
	Topic       Topic      `firestore:"topic"`
	Panelists   []Panelist `firestore:"panelists"`
	Messages    []Message  `firestore:"messages"`
	Status      string     `firestore:"status"`
	StartedAt   time.Time  `firestore:"startedAt"`
	CompletedAt time.Time  `firestore:"completedAt"`
	Metadata    Metadata   `firestore:"metadata"`
}

// SaveDebate saves a debate document to Firestore
func SaveDebate(ctx context.Context, uuid string, debate *DebateDocument) error {
	_, err := GetClient().Collection("debates").Doc(uuid).Set(ctx, debate)
	return err
}

// GetDebate retrieves a debate document from Firestore by UUID
func GetDebate(ctx context.Context, uuid string) (*DebateDocument, error) {
	doc, err := GetClient().Collection("debates").Doc(uuid).Get(ctx)
	if err != nil {
		return nil, err
	}

	var debate DebateDocument
	if err := doc.DataTo(&debate); err != nil {
		return nil, err
	}

	return &debate, nil
}
