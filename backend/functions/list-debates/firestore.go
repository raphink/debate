package listdebates

import (
	"context"
	"fmt"
	"sort"
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
		debate := DebateSummary{
			ID:    doc.Ref.ID,
			Topic: getTopicText(data),
		}

		// Extract panelists
		if panelists, ok := data["panelists"].([]interface{}); ok {
			for _, p := range panelists {
				if panelistMap, ok := p.(map[string]interface{}); ok {
					debate.Panelists = append(debate.Panelists, PanelistInfo{
						ID:        getString(panelistMap, "id"),
						Name:      getString(panelistMap, "name"),
						AvatarURL: getString(panelistMap, "avatarUrl"),
						Tagline:   getString(panelistMap, "tagline"),
						Bio:       getString(panelistMap, "biography"),
					})
				}
			}
		}

		// Extract timestamp
		if startedAt, ok := data["startedAt"].(time.Time); ok {
			debate.StartedAt = startedAt
		}

		debates = append(debates, debate)
	}

	// Get total count
	total, err := getTotalDebateCount(ctx, client)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get total count: %w", err)
	}

	return debates, total, nil
}

// getTotalDebateCount returns the total number of debates in Firestore
func getTotalDebateCount(ctx context.Context, client *firestore.Client) (int, error) {
	// Count all documents in debates collection
	iter := client.Collection("debates").Documents(ctx)
	defer iter.Stop()

	count := 0
	for {
		_, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return 0, err
		}
		count++
	}

	return count, nil
}

// getString safely extracts a string from a map
func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

// getTopicText extracts the text field from the nested topic object
func getTopicText(data map[string]interface{}) string {
	if topicObj, ok := data["topic"].(map[string]interface{}); ok {
		if text, ok := topicObj["text"].(string); ok {
			return text
		}
	}
	return ""
}

// autocompleteDebates fetches recent debates and filters using normalized token matching
// Returns up to 10 matching debates ordered by match weight (DESC), then startedAt (DESC)
func autocompleteDebates(ctx context.Context, client *firestore.Client, query string) ([]DebateSummary, error) {
	// Normalize and tokenize the query
	queryTokens := NormalizeAndTokenize(query)
	
	// If query has no significant tokens (all words <3 chars), return empty results
	if len(queryTokens) == 0 {
		return []DebateSummary{}, nil
	}

	// Fetch last 50 debates ordered by startedAt DESC
	dbQuery := client.Collection("debates").
		OrderBy("startedAt", firestore.Desc).
		Limit(50)

	iter := dbQuery.Documents(ctx)
	defer iter.Stop()

	// matchWithWeight holds a debate and its match score
	type matchWithWeight struct {
		debate DebateSummary
		weight int
	}

	var matches []matchWithWeight

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("failed to iterate debates: %w", err)
		}

		// Parse document data
		var data map[string]interface{}
		if err := doc.DataTo(&data); err != nil {
			return nil, fmt.Errorf("failed to parse debate data: %w", err)
		}

		// Extract topic text and tokenize
		topicText := getTopicText(data)
		topicTokens := NormalizeAndTokenize(topicText)
		
		// Count matching tokens (bag-of-words)
		weight := CountMatchingTokens(queryTokens, topicTokens)
		
		// Skip if not all query tokens found
		if weight == 0 {
			continue
		}

		// Build debate summary
		debate := DebateSummary{
			ID:    doc.Ref.ID,
			Topic: topicText,
		}

		// Extract panelists
		if panelists, ok := data["panelists"].([]interface{}); ok {
			debate.PanelistCount = len(panelists)

			for _, p := range panelists {
				if panelistMap, ok := p.(map[string]interface{}); ok {
					debate.Panelists = append(debate.Panelists, PanelistInfo{
						ID:        getString(panelistMap, "id"),
						Name:      getString(panelistMap, "name"),
						AvatarURL: getString(panelistMap, "avatarUrl"),
						Tagline:   getString(panelistMap, "tagline"),
						Bio:       getString(panelistMap, "biography"),
					})
				}
			}
		}

		// Extract timestamp
		if startedAt, ok := data["startedAt"].(time.Time); ok {
			debate.StartedAt = startedAt
		}

		matches = append(matches, matchWithWeight{debate, weight})
	}

	// Sort by weight (DESC), then by startedAt (DESC) for ties
	sort.Slice(matches, func(i, j int) bool {
		if matches[i].weight != matches[j].weight {
			return matches[i].weight > matches[j].weight
		}
		return matches[i].debate.StartedAt.After(matches[j].debate.StartedAt)
	})

	// Extract top 10 debates
	results := make([]DebateSummary, 0, min(10, len(matches)))
	for i := 0; i < len(matches) && i < 10; i++ {
		results = append(results, matches[i].debate)
	}

	return results, nil
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
