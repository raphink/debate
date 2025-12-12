package firebase

import (
	"context"
	"strings"
	"time"
)

// Topic represents the debate topic
type Topic struct {
	Text              string   `firestore:"text" json:"text"`
	SuggestedNames    []string `firestore:"suggestedNames,omitempty" json:"suggestedNames,omitempty"`
	IsRelevant        bool     `firestore:"isRelevant" json:"isRelevant"`
	ValidationMessage string   `firestore:"validationMessage,omitempty" json:"validationMessage,omitempty"`
}

// Panelist represents a debate participant
type Panelist struct {
	ID        string `firestore:"id" json:"id"`
	Name      string `firestore:"name" json:"name"`
	Tagline   string `firestore:"tagline" json:"tagline"`
	Biography string `firestore:"biography" json:"biography"`
	AvatarURL string `firestore:"avatarUrl" json:"avatarUrl"`
	Position  string `firestore:"position,omitempty" json:"position,omitempty"`
}

// Message represents a single debate contribution
type Message struct {
	ID           string    `firestore:"id" json:"id"`
	PanelistID   string    `firestore:"panelistId" json:"panelistId"`
	PanelistName string    `firestore:"panelistName" json:"panelistName"`
	AvatarURL    string    `firestore:"avatarUrl" json:"avatarUrl"`
	Text         string    `firestore:"text" json:"text"`
	Timestamp    time.Time `firestore:"timestamp" json:"timestamp"`
	Sequence     int       `firestore:"sequence" json:"sequence"`
	IsComplete   bool      `firestore:"isComplete" json:"isComplete"`
}

// Metadata contains debate metadata
type Metadata struct {
	CreatedBy   string `firestore:"createdBy" json:"createdBy"`
	UserAgent   string `firestore:"userAgent,omitempty" json:"userAgent,omitempty"`
	Version     string `firestore:"version" json:"version"`
	GeneratedBy string `firestore:"generatedBy" json:"generatedBy"`
}

// DebateDocument represents a complete debate stored in Firestore
type DebateDocument struct {
	ID             string     `firestore:"id" json:"id"`
	Topic          Topic      `firestore:"topic" json:"topic"`
	TopicLowercase string     `firestore:"topic_lowercase" json:"-"`
	Panelists      []Panelist `firestore:"panelists" json:"panelists"`
	Messages       []Message  `firestore:"messages" json:"messages"`
	Status         string     `firestore:"status" json:"status"`
	StartedAt      time.Time  `firestore:"startedAt" json:"startedAt"`
	CompletedAt    time.Time  `firestore:"completedAt" json:"completedAt"`
	Metadata       Metadata   `firestore:"metadata" json:"metadata"`
	CreatedAt      time.Time  `firestore:"createdAt" json:"createdAt"`
}

// SaveDebate saves a debate document to Firestore
func SaveDebate(ctx context.Context, uuid string, debate *DebateDocument) error {
	// Add lowercase topic for efficient autocomplete queries
	debate.TopicLowercase = strings.ToLower(debate.Topic.Text)
	
	// Set createdAt if not already set
	if debate.CreatedAt.IsZero() {
		debate.CreatedAt = time.Now()
	}
	
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
