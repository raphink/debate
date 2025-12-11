package validatetopic

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// WikimediaAPI client for fetching portrait images
type WikimediaAPI struct {
	client *http.Client
}

// NewWikimediaAPI creates a new Wikimedia API client
func NewWikimediaAPI() *WikimediaAPI {
	return &WikimediaAPI{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

// WikimediaImageResponse represents the API response structure
type WikimediaImageResponse struct {
	Query struct {
		Pages map[string]struct {
			PageID   int    `json:"pageid"`
			Title    string `json:"title"`
			Thumbnail struct {
				Source string `json:"source"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"thumbnail,omitempty"`
			Original struct {
				Source string `json:"source"`
				Width  int    `json:"width"`
				Height int    `json:"height"`
			} `json:"original,omitempty"`
		} `json:"pages"`
	} `json:"query"`
}

// FetchPortraitURL fetches a portrait image URL from Wikimedia Commons for a given person
// Returns empty string if no suitable image is found
func (w *WikimediaAPI) FetchPortraitURL(personName string) string {
	// Construct Wikipedia API URL
	// We'll use the pageimages API which provides the main image for a page
	baseURL := "https://en.wikipedia.org/w/api.php"
	
	params := url.Values{}
	params.Add("action", "query")
	params.Add("titles", personName)
	params.Add("prop", "pageimages")
	params.Add("pithumbsize", "300") // Request 300px thumbnail
	params.Add("format", "json")
	params.Add("formatversion", "2")
	
	apiURL := baseURL + "?" + params.Encode()
	
	// Make the request
	resp, err := w.client.Get(apiURL)
	if err != nil {
		fmt.Printf("[WIKIMEDIA] Error fetching portrait for %s: %v\n", personName, err)
		return ""
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("[WIKIMEDIA] Non-200 status for %s: %d\n", personName, resp.StatusCode)
		return ""
	}
	
	// Parse response - using simpler structure
	var result struct {
		Query struct {
			Pages []struct {
				PageID   int    `json:"pageid"`
				Title    string `json:"title"`
				Thumbnail *struct {
					Source string `json:"source"`
					Width  int    `json:"width"`
					Height int    `json:"height"`
				} `json:"thumbnail,omitempty"`
			} `json:"pages"`
		} `json:"query"`
	}
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Printf("[WIKIMEDIA] Error decoding response for %s: %v\n", personName, err)
		return ""
	}
	
	// Extract image URL from first page
	if len(result.Query.Pages) > 0 && result.Query.Pages[0].Thumbnail != nil {
		imageURL := result.Query.Pages[0].Thumbnail.Source
		fmt.Printf("[WIKIMEDIA] Found portrait for %s: %s\n", personName, imageURL)
		return imageURL
	}
	
	fmt.Printf("[WIKIMEDIA] No portrait found for %s\n", personName)
	return ""
}

// EnrichPanelistsWithPortraits fetches portrait URLs for panelists and updates their avatarUrl
// Falls back to placeholder-avatar.svg if no portrait is found
func EnrichPanelistsWithPortraits(panelists []Panelist) []Panelist {
	wiki := NewWikimediaAPI()
	
	for i := range panelists {
		// Only fetch if using placeholder
		if panelists[i].AvatarURL == "placeholder-avatar.svg" || panelists[i].AvatarURL == "" {
			portraitURL := wiki.FetchPortraitURL(panelists[i].Name)
			if portraitURL != "" {
				panelists[i].AvatarURL = portraitURL
			} else {
				// Keep placeholder as fallback
				panelists[i].AvatarURL = "placeholder-avatar.svg"
			}
		}
	}
	
	return panelists
}

// CleanPersonName removes common suffixes and titles to improve Wikipedia search
func CleanPersonName(name string) string {
	// Remove common suffixes in parentheses
	if idx := strings.Index(name, "("); idx > 0 {
		name = strings.TrimSpace(name[:idx])
	}
	
	// Remove titles
	name = strings.TrimPrefix(name, "Saint ")
	name = strings.TrimPrefix(name, "St. ")
	
	return strings.TrimSpace(name)
}
