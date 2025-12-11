package getportrait

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
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
	
	// Create request with proper User-Agent header (required by Wikipedia)
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		fmt.Printf("[WIKIMEDIA] Error creating request for %s: %v\n", personName, err)
		return ""
	}
	req.Header.Set("User-Agent", "DebateApp/1.0 (https://github.com/raphink/debate; debate@example.com)")
	
	// Make the request
	resp, err := w.client.Do(req)
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
