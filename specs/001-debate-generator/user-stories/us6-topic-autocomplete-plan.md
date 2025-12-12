# US6 Topic Autocomplete - Technical Implementation Plan

## Architecture Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Flow                                │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  1. Type "eth..." in Home.jsx topic input                        │
│  2. useTopicAutocomplete hook debounces (300ms)                  │
│  3. GET /api/autocomplete-topics?q=eth&limit=10                  │
│  4. TopicAutocompleteDropdown shows results                      │
│  5. Select "Ethics of AI in healthcare" (3 panelists)            │
│  6. Navigate to PanelistSelection with state:                    │
│     { debateId, topic, panelists, skipValidation: true }         │
│  7. PanelistSelection pre-fills chips, shows "Modify" button     │
│  8. User clicks "Generate" without changes                       │
│  9. Cache detection: isCacheHit(topic, panelists) → true         │
│ 10. Redirect to /d/{debateId} (load from Firestore)              │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Backend Architecture                          │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Cloud Function: autocomplete-topics                             │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  1. Parse query params: q, limit (default 10)        │       │
│  │  2. Validate q.length >= 3                            │       │
│  │  3. Query Firestore:                                  │       │
│  │     - OrderBy startedAt DESC                          │       │
│  │     - Limit(100) recent debates                       │       │
│  │  4. Filter results client-side:                       │       │
│  │     - strings.Contains(topic.lower, query.lower)      │       │
│  │     - Return first <limit> matches                    │       │
│  │  5. Transform results:                                │       │
│  │     - Map to: {id, topic, panelists, count, created}  │       │
│  │  6. Return JSON with CORS headers                     │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
│  Firestore Schema (No Changes Required):                         │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  debates/{id}:                                        │       │
│  │    topic: string (existing)                           │       │
│  │    panelists: array (existing)                        │       │
│  │    startedAt: timestamp (existing)                    │       │
│  │    ...other fields                                    │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
│  Index Required: None (uses existing startedAt index)            │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘

┌─────────────────────────────────────────────────────────────────┐
│                    Frontend Architecture                         │
├─────────────────────────────────────────────────────────────────┤
│                                                                   │
│  Hook: useTopicAutocomplete(query)                               │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  State: suggestions[], loading, error                 │       │
│  │  Effect: Debounce 300ms, min 3 chars                  │       │
│  │  API Call: api.autocompleteTopics(query, limit=10)    │       │
│  │  Return: { suggestions, loading, error }              │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
│  Component: TopicAutocompleteDropdown                            │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  Props: suggestions, onSelect, loading, visible       │       │
│  │  State: selectedIndex (keyboard nav)                  │       │
│  │  Handlers:                                            │       │
│  │    - onKeyDown (ArrowUp/Down, Enter, Escape)          │       │
│  │    - onClick (select item)                            │       │
│  │    - useOnClickOutside (close dropdown)               │       │
│  │  Render:                                              │       │
│  │    - Absolute positioned below input                  │       │
│  │    - List items: topic + badge ("N panelists")        │       │
│  │    - Highlight selected index                         │       │
│  │    - Loading spinner if loading && >300ms             │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
│  Page: Home.jsx Integration                                      │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  1. Add state: topicQuery (controlled input)          │       │
│  │  2. Use hook: useTopicAutocomplete(topicQuery)        │       │
│  │  3. Render dropdown below input                       │       │
│  │  4. Handle selection:                                 │       │
│  │     - Set topic state                                 │       │
│  │     - Store debate metadata in ref/state              │       │
│  │     - Call onSubmit with skipValidation=true          │       │
│  │  5. Navigate with state:                              │       │
│  │     { debateId, topic, panelists, skipValidation }    │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
│  Page: PanelistSelection.jsx Enhancements                        │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  1. Read navigation state for pre-filled data         │       │
│  │  2. If panelists provided:                            │       │
│  │     - Pre-populate selectedPanelists state            │       │
│  │     - Set mode: 'locked' (chips disabled)             │       │
│  │     - Show "Modify Panelists" button                  │       │
│  │  3. On "Modify" click:                                │       │
│  │     - Set mode: 'editable'                            │       │
│  │     - Enable chip add/remove                          │       │
│  │  4. On "Generate Debate" click:                       │       │
│  │     - Call isCacheHit(debateId, topic, panelists)     │       │
│  │     - If true: navigate to /d/{debateId}              │       │
│  │     - If false: proceed with normal generation        │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
│  Utility: cacheDetection.js                                      │
│  ┌──────────────────────────────────────────────────────┐       │
│  │  export const isCacheHit = (                          │       │
│  │    originalDebate,                                    │       │
│  │    currentTopic,                                      │       │
│  │    currentPanelists                                   │       │
│  │  ) => {                                               │       │
│  │    if (!originalDebate) return false;                 │       │
│  │    if (originalDebate.topic !== currentTopic) {       │       │
│  │      return false;                                    │       │
│  │    }                                                  │       │
│  │    const origIds = originalDebate.panelists           │       │
│  │      .map(p => p.id).sort();                          │       │
│  │    const currIds = currentPanelists                   │       │
│  │      .map(p => p.id).sort();                          │       │
│  │    return JSON.stringify(origIds) ===                 │       │
│  │           JSON.stringify(currIds);                    │       │
│  │  };                                                   │       │
│  └──────────────────────────────────────────────────────┘       │
│                                                                   │
└─────────────────────────────────────────────────────────────────┘
```

## Implementation Phases

### Phase 1: Backend Infrastructure (Tasks T144-T150)

**Goal**: Create autocomplete-topics Cloud Function with Firestore query capability

#### T144-T145: Function Scaffolding

**Files**:
- `backend/functions/autocomplete-topics/main.go`
- `backend/functions/autocomplete-topics/handler.go`
- `backend/functions/autocomplete-topics/go.mod`
- `backend/functions/autocomplete-topics/cmd/main.go`

**Structure**:
```go
// handler.go
package autocompletetopics

import (
    "encoding/json"
    "log"
    "net/http"
    "strings"
    "cloud.google.com/go/firestore"
    "github.com/raphink/debate/backend/shared/firebase"
    "github.com/raphink/debate/backend/shared/errors"
)

type AutocompleteRequest struct {
    Query string
    Limit int
}

type DebateSummary struct {
    ID            string              `json:"id"`
    Topic         string              `json:"topic"`
    Panelists     []PanelistSummary   `json:"panelists"`
    PanelistCount int                 `json:"panelistCount"`
    CreatedAt     string              `json:"createdAt"`
}

type PanelistSummary struct {
    ID   string `json:"id"`
    Name string `json:"name"`
    Slug string `json:"slug"`
}

func AutocompleteTopicsHandler(w http.ResponseWriter, r *http.Request) {
    // CORS headers
    // Parse query params
    // Validate input
    // Query Firestore
    // Transform results
    // Return JSON
}
```

#### T146-T148: Query Implementation

**Firestore Query Logic**:
```go
// In handler.go
func queryDebates(ctx context.Context, client *firestore.Client, query string, limit int) ([]DebateSummary, error) {
    // Normalize query
    q := strings.ToLower(query)
    
    // Query with range for substring search
    iter := client.Collection("debates").
        Where("topic_lowercase", ">=", q).
        Where("topic_lowercase", "<", q+"\uf8ff"). // Unicode max char
        OrderBy("topic_lowercase", firestore.Asc).
        OrderBy("createdAt", firestore.Desc).
        Limit(limit).
        Documents(ctx)
    
    defer iter.Stop()
    
    var results []DebateSummary
    for {
        doc, err := iter.Next()
        if err == iterator.Done {
            break
        }
        if err != nil {
            return nil, err
        }
        
        var debate firebase.DebateDocument
        if err := doc.DataTo(&debate); err != nil {
            continue
        }
        
        // Transform to summary
        summary := transformToSummary(doc.Ref.ID, &debate)
        results = append(results, summary)
    }
    
    return results, nil
}

func transformToSummary(id string, debate *firebase.DebateDocument) DebateSummary {
    panelists := make([]PanelistSummary, len(debate.Panelists))
    for i, p := range debate.Panelists {
        panelists[i] = PanelistSummary{
            ID:   p.ID,
            Name: p.Name,
            Slug: p.Slug,
        }
    }
    
    return DebateSummary{
        ID:            id,
        Topic:         debate.Topic.Text,
        Panelists:     panelists,
        PanelistCount: len(panelists),
        CreatedAt:     debate.Metadata.CreatedAt.Format(time.RFC3339),
    }
}
```

**CORS and Error Handling**:
```go
func AutocompleteTopicsHandler(w http.ResponseWriter, r *http.Request) {
    // CORS
    w.Header().Set("Access-Control-Allow-Origin", "*")
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
    w.Header().Set("Content-Type", "application/json")
    
    if r.Method == http.MethodOptions {
        w.WriteHeader(http.StatusOK)
        return
    }
    
    if r.Method != http.MethodGet {
        errors.SendError(w, http.StatusMethodNotAllowed, "Method not allowed")
        return
    }
    
    // Parse params
    query := r.URL.Query().Get("q")
    if len(query) < 3 {
        errors.SendError(w, http.StatusBadRequest, "Query must be at least 3 characters")
        return
    }
    
    limitStr := r.URL.Query().Get("limit")
    limit := 10
    if limitStr != "" {
        if parsed, err := strconv.Atoi(limitStr); err == nil && parsed > 0 {
            limit = parsed
        }
    }
    
    // Query
    ctx := r.Context()
    client, err := firebase.GetClient(ctx)
    if err != nil {
        log.Printf("Failed to get Firestore client: %v", err)
        errors.SendError(w, http.StatusInternalServerError, "Internal server error")
        return
    }
    
    results, err := queryDebates(ctx, client, query, limit)
    if err != nil {
        log.Printf("Failed to query debates: %v", err)
        errors.SendError(w, http.StatusInternalServerError, "Failed to fetch debates")
        return
    }
    
    response := map[string]interface{}{
        "debates": results,
    }
    
    json.NewEncoder(w).Encode(response)
}
```

#### T149: Dockerfile

**File**: `backend/functions/autocomplete-topics/Dockerfile`

```dockerfile
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o autocomplete-topics ./cmd

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/autocomplete-topics .
EXPOSE 8080
CMD ["./autocomplete-topics"]
```

#### T150: Deployment Integration

**File**: `deploy.sh` (add to existing script)

```bash
# Add to deploy.sh after other function deployments
echo "Deploying autocomplete-topics function..."
gcloud functions deploy autocomplete-topics \
  --gen2 \
  --runtime=go124 \
  --region=us-central1 \
  --source=./backend/functions/autocomplete-topics \
  --entry-point=AutocompleteTopicsHandler \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars FIRESTORE_PROJECT_ID=${PROJECT_ID}
```

#### Firestore Schema Update

**File**: `backend/shared/firebase/debates.go` (update SaveDebate)

```go
// Add to DebateDocument struct
type DebateDocument struct {
    Topic         Topic       `firestore:"topic" json:"topic"`
    TopicLowercase string     `firestore:"topic_lowercase" json:"-"` // NEW
    // ... existing fields
}

// Update SaveDebate function
func SaveDebate(ctx context.Context, client *firestore.Client, debate *DebateDocument) (string, error) {
    // Add lowercase topic for querying
    debate.TopicLowercase = strings.ToLower(debate.Topic.Text)
    
    // ... rest of existing save logic
}
```

**Firestore Index**:
Create `firestore.indexes.json`:
```json
{
  "indexes": [
    {
      "collectionGroup": "debates",
      "queryScope": "COLLECTION",
      "fields": [
        {
          "fieldPath": "topic_lowercase",
          "order": "ASCENDING"
        },
        {
          "fieldPath": "createdAt",
          "order": "DESCENDING"
        }
      ]
    }
  ]
}
```

Deploy index:
```bash
gcloud firestore indexes create --database=(default) \
  --collection-group=debates \
  --field=topic_lowercase,ASC \
  --field=createdAt,DESC
```

### Phase 2: Frontend Autocomplete UI (Tasks T151-T158)

**Goal**: Create autocomplete dropdown with debouncing and keyboard navigation

#### T151: API Service Method

**File**: `frontend/src/services/api.js`

```javascript
// Add to existing api.js
export const autocompleteTopics = async (query, limit = 10) => {
  if (!query || query.length < 3) {
    return { debates: [] };
  }

  const params = new URLSearchParams({
    q: query,
    limit: limit.toString(),
  });

  const response = await fetch(`${API_BASE_URL}/autocomplete-topics?${params}`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!response.ok) {
    throw new Error(`Autocomplete failed: ${response.status}`);
  }

  return response.json();
};
```

#### T152: useTopicAutocomplete Hook

**File**: `frontend/src/hooks/useTopicAutocomplete.js`

```javascript
import { useState, useEffect } from 'react';
import { autocompleteTopics } from '../services/api';

export const useTopicAutocomplete = (query, enabled = true) => {
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  useEffect(() => {
    // Reset if query too short or disabled
    if (!enabled || !query || query.length < 3) {
      setSuggestions([]);
      setLoading(false);
      return;
    }

    // Debounce 300ms
    const timer = setTimeout(async () => {
      setLoading(true);
      setError(null);

      try {
        const data = await autocompleteTopics(query, 10);
        setSuggestions(data.debates || []);
      } catch (err) {
        console.error('Autocomplete error:', err);
        setError(err);
        setSuggestions([]);
      } finally {
        setLoading(false);
      }
    }, 300);

    return () => clearTimeout(timer);
  }, [query, enabled]);

  return { suggestions, loading, error };
};
```

#### T153: TopicAutocompleteDropdown Component

**File**: `frontend/src/components/TopicInput/TopicAutocompleteDropdown.jsx`

```javascript
import React, { useState, useEffect, useRef } from 'react';
import PropTypes from 'prop-types';
import styles from './TopicAutocompleteDropdown.module.css';

export const TopicAutocompleteDropdown = ({
  suggestions,
  onSelect,
  loading,
  visible,
  onClose,
}) => {
  const [selectedIndex, setSelectedIndex] = useState(0);
  const dropdownRef = useRef(null);

  // Reset selection when suggestions change
  useEffect(() => {
    setSelectedIndex(0);
  }, [suggestions]);

  // Keyboard navigation
  useEffect(() => {
    if (!visible) return;

    const handleKeyDown = (e) => {
      if (!suggestions.length) return;

      switch (e.key) {
        case 'ArrowDown':
          e.preventDefault();
          setSelectedIndex((prev) =>
            prev < suggestions.length - 1 ? prev + 1 : 0
          );
          break;
        case 'ArrowUp':
          e.preventDefault();
          setSelectedIndex((prev) =>
            prev > 0 ? prev - 1 : suggestions.length - 1
          );
          break;
        case 'Enter':
          e.preventDefault();
          if (suggestions[selectedIndex]) {
            onSelect(suggestions[selectedIndex]);
          }
          break;
        case 'Escape':
          e.preventDefault();
          onClose();
          break;
        default:
          break;
      }
    };

    document.addEventListener('keydown', handleKeyDown);
    return () => document.removeEventListener('keydown', handleKeyDown);
  }, [visible, suggestions, selectedIndex, onSelect, onClose]);

  // Click outside to close
  useEffect(() => {
    if (!visible) return;

    const handleClickOutside = (e) => {
      if (dropdownRef.current && !dropdownRef.current.contains(e.target)) {
        onClose();
      }
    };

    document.addEventListener('mousedown', handleClickOutside);
    return () => document.removeEventListener('mousedown', handleClickOutside);
  }, [visible, onClose]);

  if (!visible || (!loading && !suggestions.length)) {
    return null;
  }

  return (
    <div className={styles.dropdown} ref={dropdownRef}>
      {loading && <div className={styles.loading}>Searching...</div>}
      {!loading && suggestions.length > 0 && (
        <ul className={styles.list}>
          {suggestions.map((debate, index) => (
            <li
              key={debate.id}
              className={`${styles.item} ${
                index === selectedIndex ? styles.selected : ''
              }`}
              onClick={() => onSelect(debate)}
              onMouseEnter={() => setSelectedIndex(index)}
            >
              <div className={styles.topic}>{debate.topic}</div>
              <div className={styles.badge}>
                {debate.panelistCount} panelist{debate.panelistCount !== 1 ? 's' : ''}
              </div>
            </li>
          ))}
        </ul>
      )}
    </div>
  );
};

TopicAutocompleteDropdown.propTypes = {
  suggestions: PropTypes.arrayOf(
    PropTypes.shape({
      id: PropTypes.string.isRequired,
      topic: PropTypes.string.isRequired,
      panelists: PropTypes.array.isRequired,
      panelistCount: PropTypes.number.isRequired,
      createdAt: PropTypes.string.isRequired,
    })
  ).isRequired,
  onSelect: PropTypes.func.isRequired,
  loading: PropTypes.bool,
  visible: PropTypes.bool.isRequired,
  onClose: PropTypes.func.isRequired,
};
```

**File**: `frontend/src/components/TopicInput/TopicAutocompleteDropdown.module.css`

```css
.dropdown {
  position: absolute;
  top: 100%;
  left: 0;
  right: 0;
  margin-top: 8px;
  background: white;
  border: 1px solid #e0e0e0;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  max-height: 400px;
  overflow-y: auto;
  z-index: 1000;
}

.loading {
  padding: 16px;
  text-align: center;
  color: #666;
  font-size: 14px;
}

.list {
  list-style: none;
  margin: 0;
  padding: 0;
}

.item {
  padding: 12px 16px;
  cursor: pointer;
  border-bottom: 1px solid #f0f0f0;
  display: flex;
  justify-content: space-between;
  align-items: center;
  transition: background-color 0.15s ease;
}

.item:last-child {
  border-bottom: none;
}

.item:hover,
.item.selected {
  background-color: #f5f5f5;
}

.topic {
  flex: 1;
  font-size: 14px;
  color: #333;
  margin-right: 12px;
}

.badge {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  font-weight: 500;
  white-space: nowrap;
}
```

#### T154-T158: Home.jsx Integration

**File**: `frontend/src/pages/Home.jsx` (update existing file)

```javascript
import { TopicAutocompleteDropdown } from '../components/TopicInput/TopicAutocompleteDropdown';
import { useTopicAutocomplete } from '../hooks/useTopicAutocomplete';

// Inside Home component
const [topic, setTopic] = useState('');
const [showAutocomplete, setShowAutocomplete] = useState(false);
const [selectedDebateMetadata, setSelectedDebateMetadata] = useState(null);

const { suggestions, loading, error } = useTopicAutocomplete(topic, showAutocomplete);

const handleTopicChange = (e) => {
  setTopic(e.target.value);
  setShowAutocomplete(true);
  setSelectedDebateMetadata(null);
};

const handleAutocompleteSelect = (debate) => {
  setTopic(debate.topic);
  setSelectedDebateMetadata(debate);
  setShowAutocomplete(false);
  
  // Navigate with metadata
  navigate('/panelist-selection', {
    state: {
      debateId: debate.id,
      topic: debate.topic,
      panelists: debate.panelists,
      skipValidation: true,
    },
  });
};

// In JSX:
<div style={{ position: 'relative' }}>
  <input
    type="text"
    value={topic}
    onChange={handleTopicChange}
    onFocus={() => setShowAutocomplete(true)}
    placeholder="Enter debate topic..."
  />
  <TopicAutocompleteDropdown
    suggestions={suggestions}
    onSelect={handleAutocompleteSelect}
    loading={loading}
    visible={showAutocomplete}
    onClose={() => setShowAutocomplete(false)}
  />
</div>
```

### Phase 3: Cache Detection and Modified Flow (Tasks T159-T164)

**Goal**: Implement cache hit detection and "Modify Panelists" workflow

#### T159: Cache Detection Utility

**File**: `frontend/src/utils/cacheDetection.js`

```javascript
/**
 * Determine if current topic + panelist combination matches a historical debate
 * @param {Object} originalDebate - The debate from autocomplete/navigation state
 * @param {string} currentTopic - Current topic text
 * @param {Array} currentPanelists - Current panelist array
 * @returns {boolean} True if cache hit (exact match)
 */
export const isCacheHit = (originalDebate, currentTopic, currentPanelists) => {
  if (!originalDebate) return false;

  // Topic must match exactly (case-sensitive)
  if (originalDebate.topic !== currentTopic) {
    return false;
  }

  // Panelist count must match
  if (originalDebate.panelists.length !== currentPanelists.length) {
    return false;
  }

  // Deep equality check on panelist IDs (order-independent)
  const originalIds = originalDebate.panelists
    .map((p) => p.id)
    .sort()
    .join(',');
  const currentIds = currentPanelists
    .map((p) => p.id)
    .sort()
    .join(',');

  return originalIds === currentIds;
};
```

#### T160-T164: PanelistSelection.jsx Enhancements

**File**: `frontend/src/pages/PanelistSelection.jsx` (major updates)

```javascript
import { useLocation, useNavigate } from 'react-router-dom';
import { isCacheHit } from '../utils/cacheDetection';
import { useState, useEffect } from 'react';

export const PanelistSelection = () => {
  const location = useLocation();
  const navigate = useNavigate();
  
  const {
    debateId: originalDebateId,
    topic,
    panelists: preFillPanelists,
    skipValidation,
  } = location.state || {};

  const [selectedPanelists, setSelectedPanelists] = useState([]);
  const [isLocked, setIsLocked] = useState(false);
  const [showCacheIndicator, setShowCacheIndicator] = useState(false);

  // Pre-fill panelists on mount if provided
  useEffect(() => {
    if (preFillPanelists && preFillPanelists.length > 0) {
      setSelectedPanelists(preFillPanelists);
      setIsLocked(true);
      setShowCacheIndicator(true);
    }
  }, [preFillPanelists]);

  const handleModifyPanelists = () => {
    setIsLocked(false);
    setShowCacheIndicator(false);
  };

  const handleGenerateDebate = async () => {
    // Check for cache hit
    if (originalDebateId && preFillPanelists) {
      const cacheHit = isCacheHit(
        { id: originalDebateId, topic, panelists: preFillPanelists },
        topic,
        selectedPanelists
      );

      if (cacheHit) {
        // Load cached debate
        console.log('Cache hit detected, loading debate:', originalDebateId);
        navigate(`/d/${originalDebateId}`);
        return;
      }
    }

    // No cache hit, proceed with normal generation
    console.log('Generating new debate with modified panelists');
    // ... existing debate generation logic
  };

  return (
    <div className={styles.container}>
      <h1>Select Panelists</h1>
      <p>{topic}</p>

      {showCacheIndicator && (
        <div className={styles.cacheIndicator}>
          <span>✓ Using cached debate</span>
        </div>
      )}

      {/* Panelist chips */}
      <div className={styles.panelistList}>
        {selectedPanelists.map((panelist) => (
          <div
            key={panelist.id}
            className={`${styles.chip} ${isLocked ? styles.locked : ''}`}
          >
            {panelist.name}
            {!isLocked && (
              <button onClick={() => handleRemovePanelist(panelist.id)}>×</button>
            )}
          </div>
        ))}
      </div>

      {/* Modify button when locked */}
      {isLocked && (
        <button onClick={handleModifyPanelists} className={styles.modifyButton}>
          Modify Panelists
        </button>
      )}

      {/* Generate button */}
      <button onClick={handleGenerateDebate} className={styles.generateButton}>
        {showCacheIndicator ? 'Load Debate' : 'Generate Debate'}
      </button>
    </div>
  );
};
```

**CSS Updates**: `frontend/src/pages/PanelistSelection.module.css`

```css
.cacheIndicator {
  background: #e8f5e9;
  border: 1px solid #4caf50;
  color: #2e7d32;
  padding: 12px 16px;
  border-radius: 8px;
  margin-bottom: 16px;
  display: flex;
  align-items: center;
  gap: 8px;
}

.chip.locked {
  opacity: 0.7;
  cursor: not-allowed;
  background: #e0e0e0;
}

.modifyButton {
  background: #ff9800;
  color: white;
  border: none;
  padding: 10px 20px;
  border-radius: 6px;
  cursor: pointer;
  margin-bottom: 16px;
}

.modifyButton:hover {
  background: #fb8c00;
}
```

### Phase 4: Testing (Tasks T165-T171)

**Goal**: Comprehensive testing of backend and frontend features

#### T165-T166: Backend Tests

**File**: `backend/functions/autocomplete-topics/handler_test.go`

```go
package autocompletetopics

import (
    "context"
    "net/http/httptest"
    "testing"
)

func TestAutocompleteTopicsHandler_ValidQuery(t *testing.T) {
    // Setup test Firestore
    // Insert test debates
    // Create test request: GET /autocomplete-topics?q=ethics&limit=10
    // Call handler
    // Assert response contains matching debates
    // Assert ordered by createdAt DESC
}

func TestAutocompleteTopicsHandler_LimitRespected(t *testing.T) {
    // Insert 20 test debates
    // Request with limit=5
    // Assert max 5 results returned
}

func TestAutocompleteTopicsHandler_QueryTooShort(t *testing.T) {
    // Request with q=ab (2 chars)
    // Assert 400 Bad Request
}
```

#### T167-T171: Frontend Tests

**File**: `frontend/src/hooks/useTopicAutocomplete.test.js`

```javascript
import { renderHook, waitFor } from '@testing-library/react';
import { useTopicAutocomplete } from './useTopicAutocomplete';
import * as api from '../services/api';

jest.mock('../services/api');

describe('useTopicAutocomplete', () => {
  it('should debounce API calls by 300ms', async () => {
    api.autocompleteTopics.mockResolvedValue({ debates: [] });

    const { rerender } = renderHook(({ query }) => useTopicAutocomplete(query), {
      initialProps: { query: 'eth' },
    });

    // Change query rapidly
    rerender({ query: 'ethi' });
    rerender({ query: 'ethic' });

    // Should only call once after debounce
    await waitFor(() => {
      expect(api.autocompleteTopics).toHaveBeenCalledTimes(1);
    });
  });

  it('should not call API if query < 3 chars', () => {
    renderHook(() => useTopicAutocomplete('ab'));
    expect(api.autocompleteTopics).not.toHaveBeenCalled();
  });
});
```

**File**: `frontend/src/components/TopicInput/TopicAutocompleteDropdown.test.jsx`

```javascript
import { render, screen, fireEvent } from '@testing-library/react';
import { TopicAutocompleteDropdown } from './TopicAutocompleteDropdown';

describe('TopicAutocompleteDropdown', () => {
  const mockSuggestions = [
    { id: '1', topic: 'Ethics of AI', panelists: [], panelistCount: 3, createdAt: '2025-01-15' },
    { id: '2', topic: 'AI Safety', panelists: [], panelistCount: 2, createdAt: '2025-01-14' },
  ];

  it('should render suggestions', () => {
    render(
      <TopicAutocompleteDropdown
        suggestions={mockSuggestions}
        onSelect={jest.fn()}
        visible={true}
        onClose={jest.fn()}
      />
    );

    expect(screen.getByText('Ethics of AI')).toBeInTheDocument();
    expect(screen.getByText('3 panelists')).toBeInTheDocument();
  });

  it('should call onSelect when item clicked', () => {
    const handleSelect = jest.fn();
    render(
      <TopicAutocompleteDropdown
        suggestions={mockSuggestions}
        onSelect={handleSelect}
        visible={true}
        onClose={jest.fn()}
      />
    );

    fireEvent.click(screen.getByText('Ethics of AI'));
    expect(handleSelect).toHaveBeenCalledWith(mockSuggestions[0]);
  });

  it('should support keyboard navigation', () => {
    const handleSelect = jest.fn();
    render(
      <TopicAutocompleteDropdown
        suggestions={mockSuggestions}
        onSelect={handleSelect}
        visible={true}
        onClose={jest.fn()}
      />
    );

    // Arrow down, then Enter
    fireEvent.keyDown(document, { key: 'ArrowDown' });
    fireEvent.keyDown(document, { key: 'Enter' });

    expect(handleSelect).toHaveBeenCalled();
  });
});
```

## File Changes Summary

### Backend
1. **New Function**: `backend/functions/autocomplete-topics/`
   - `main.go`, `handler.go`, `go.mod`, `cmd/main.go`, `Dockerfile`
2. **Schema Update**: `backend/shared/firebase/debates.go`
   - Add `TopicLowercase` field
   - Update `SaveDebate()` to populate it
3. **Deployment**: Update `deploy.sh` to deploy autocomplete-topics
4. **Index**: Create `firestore.indexes.json` and deploy index

### Frontend
1. **New Hook**: `frontend/src/hooks/useTopicAutocomplete.js`
2. **New Component**: `frontend/src/components/TopicInput/TopicAutocompleteDropdown.jsx`
3. **New CSS Module**: `frontend/src/components/TopicInput/TopicAutocompleteDropdown.module.css`
4. **New Utility**: `frontend/src/utils/cacheDetection.js`
5. **API Service**: Update `frontend/src/services/api.js`
6. **Page Updates**: 
   - `frontend/src/pages/Home.jsx` (integrate dropdown)
   - `frontend/src/pages/PanelistSelection.jsx` (pre-fill, modify, cache detection)
   - `frontend/src/pages/PanelistSelection.module.css` (locked/cache styles)

### Testing
1. **Backend**: `backend/functions/autocomplete-topics/handler_test.go`
2. **Frontend**: 
   - `frontend/src/hooks/useTopicAutocomplete.test.js`
   - `frontend/src/components/TopicInput/TopicAutocompleteDropdown.test.jsx`

## Deployment Checklist

- [ ] Deploy Firestore index (may take several minutes)
- [ ] Deploy autocomplete-topics Cloud Function
- [ ] Verify CORS headers in production
- [ ] Test autocomplete with real debates in Firestore
- [ ] Monitor function logs for errors
- [ ] Test cache hit detection with actual data
- [ ] Verify loading states and error handling
- [ ] Test keyboard navigation across browsers

## Success Criteria

- Autocomplete appears within 500ms of typing 3+ characters
- Dropdown shows accurate historical debates ordered by recency
- Cache hit correctly redirects to existing debate
- Modified panelists trigger new debate generation
- Graceful degradation if Firestore/API fails
- No blocking or janky typing experience
- All tests passing
