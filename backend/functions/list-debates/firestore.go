package listdebates

import (
	"context"
	"fmt"
	"time"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
)

// queryDebates fetches debates from Firestore with pagination
func queryDebates(ctx context.Context, client *firestore.Client, limit, offset int) ([]DebateSummary, int, error) {
	// Query debates ordered by startedAt descending
	query := client.Collection("debates").
		OrderBy("startedAt", firestore.Desc).
		Limit(limit).
		Offset(offset)

	// Execute query
	iter := query.Documents(ctx)
	defer iter.Stop()

	var debates []DebateSummary

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, 0, fmt.Errorf("failed to iterate debates: %w", err)
		}

		// Parse document data
		var data map[string]interface{}
		if err := doc.DataTo(&data); err != nil {
			return nil, 0, fmt.Errorf("failed to parse debate data: %w", err)
		}

		// Extract debate summary






























































}	return ""	}		return val	if val, ok := m[key].(string); ok {func getString(m map[string]interface{}, key string) string {// getString safely extracts a string from a map}	return count, nil	}		count++		}			return 0, err		if err != nil {		}			break		if err == iterator.Done {		_, err := iter.Next()	for {	count := 0	defer iter.Stop()	iter := client.Collection("debates").Documents(ctx)	// Count all documents in debates collectionfunc getTotalDebateCount(ctx context.Context, client *firestore.Client) (int, error) {// getTotalDebateCount returns the total number of debates in Firestore}	return debates, total, nil	}		return nil, 0, fmt.Errorf("failed to get total count: %w", err)	if err != nil {	total, err := getTotalDebateCount(ctx, client)	// Get total count	}		debates = append(debates, debate)		}			debate.StartedAt = startedAt		if startedAt, ok := data["startedAt"].(time.Time); ok {		// Extract timestamp		}			}				}					})						Name: getString(panelistMap, "name"),						ID:   getString(panelistMap, "id"),					debate.Panelists = append(debate.Panelists, PanelistInfo{				if panelistMap, ok := p.(map[string]interface{}); ok {			for _, p := range panelists {		if panelists, ok := data["panelists"].([]interface{}); ok {		// Extract panelists		}			Topic: getString(data, "topic"),			ID:    doc.Ref.ID,		debate := DebateSummary{