package getportrait

import (
	"sync"
)

// PortraitCache provides thread-safe in-memory caching of portrait URLs
type PortraitCache struct {
	cache map[string]string
	mu    sync.RWMutex
}

// NewPortraitCache creates a new portrait cache
func NewPortraitCache() *PortraitCache {
	return &PortraitCache{
		cache: make(map[string]string),
	}
}

// Get retrieves a portrait URL from cache
// Returns empty string if not found
func (pc *PortraitCache) Get(panelistID string) (string, bool) {
	pc.mu.RLock()
	defer pc.mu.RUnlock()

	url, found := pc.cache[panelistID]
	return url, found
}

// Set stores a portrait URL in cache
func (pc *PortraitCache) Set(panelistID, portraitURL string) {
	pc.mu.Lock()
	defer pc.mu.Unlock()

	pc.cache[panelistID] = portraitURL
}

// Global cache instance (persists across function invocations in same instance)
var portraitCache = NewPortraitCache()
