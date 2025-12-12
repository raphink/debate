package autocompletetopics

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	"github.com/raphink/debate/backend/shared/firebase"
	"google.golang.org/api/iterator"
)

// AutocompleteTopicsHandler handles GET requests for topic autocomplete
func AutocompleteTopicsHandler(w http.ResponseWriter, r *http.Request) {
	// CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle OPTIONS preflight
	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	// Parse query parameter
	query := r.URL.Query().Get("q")
	if len(query) < 3 {
		http.Error(w, `{"error":"Query must be at least 3 characters"}`, http.StatusBadRequest)
		return
	}

	// Parse limit parameter
	limitStr := r.URL.Query().Get("limit")
	limit := 10
	if limitStr != "" {
		if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 && parsed <= 50 {
			limit = parsed
		}
	}

	// Initialize Firestore (if not already initialized)
	ctx := r.Context()
	if firebase.GetClient() == nil {
		if err := firebase.InitFirestore(ctx); err != nil {
			log.Printf("Failed to initialize Firestore: %v", err)
			http.Error(w, `{"error":"Database connection failed"}`, http.StatusInternalServerError)
			return
		}
	}

	// Query Firestore
	results, err := queryDebates(ctx, query, limit)
	if err != nil {
		log.Printf("Failed to query debates: %v", err)
		http.Error(w, `{"error":"Failed to fetch debates"}`, http.StatusInternalServerError)
		return
	}

	// Return response
	response := AutocompleteResponse{
		Debates: results,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Failed to encode response: %v", err)
		http.Error(w, `{"error":"Failed to encode response"}`, http.StatusInternalServerError)
		return
	}
}

// queryDebates performs Firestore query for topic autocomplete
func queryDebates(ctx context.Context, queryStr string, limit int) ([]DebateSummary, error) {
	client := firebase.GetClient()

	// Normalize query for case-insensitive search
	q := strings.ToLower(queryStr)

	// Fetch recent debates (up to 100) and filter client-side for substring matching
	// This is more effective than Firestore's limited text search capabilities
	iter := client.Collection("debates").
		OrderBy("createdAt", firestore.Desc).
		Limit(100).
		Documents(ctx)

	defer iter.Stop()

	var results []DebateSummary
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating debate documents: %v", err)
			continue
		}

		var debate firebase.DebateDocument
		if err := doc.DataTo(&debate); err != nil {
			log.Printf("Error parsing debate document: %v", err)
			continue
		}

		// Filter by substring match (case-insensitive)
		topicLower := strings.ToLower(debate.Topic.Text)
		if strings.Contains(topicLower, q) {
			summary := transformToSummary(doc.Ref.ID, &debate)
			results = append(results, summary)

			// Stop once we have enough results
			if len(results) >= limit {
				break
			}
		}
	}

	return results, nil
}

// transformToSummary converts DebateDocument to DebateSummary
func transformToSummary(id string, debate *firebase.DebateDocument) DebateSummary {
	panelists := make([]PanelistSummary, len(debate.Panelists))
	for i, p := range debate.Panelists {
		// Generate slug from ID if not present
		slug := p.ID
		if slug == "" {
			slug = strings.ToLower(strings.ReplaceAll(p.Name, " ", "-"))
		}

		panelists[i] = PanelistSummary{
			ID:   p.ID,
			Name: p.Name,
			Slug: slug,
		}
	}

	// Format timestamp
	createdAt := ""
	if !debate.StartedAt.IsZero() {
		createdAt = debate.StartedAt.Format(time.RFC3339)
	}

	return DebateSummary{
		ID:            id,
		Topic:         debate.Topic.Text,
		Panelists:     panelists,
		PanelistCount: len(panelists),
		CreatedAt:     createdAt,
	}
}
