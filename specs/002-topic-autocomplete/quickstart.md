# US6: Topic Discovery - Quick Start Guide

## Prerequisites

Before implementing US6, ensure:

1. **US5 (Debate Caching & Sharing) is complete**
   - Firestore debates collection exists and is populated
   - backend/shared/firebase module is functional
   - Debate documents have required fields: id, topic, panelists, createdAt

2. **Development environment ready**
   - Docker & Docker Compose installed
   - GCP project configured with Firestore enabled
   - ALLOWED_ORIGIN environment variable set

3. **Sample data available**
   - Generate 5-10 test debates via existing US1-US3 flow
   - Verify debates saved to Firestore (check GCP Console)

## Architecture Overview

```
┌─────────────────┐
│   Home.jsx      │  User types topic
│  (TopicInput)   │  ↓
└────────┬────────┘  
         │ ≥3 chars typed
         ↓
┌─────────────────────────────┐
│ useTopicAutocomplete hook   │  Debounce 300ms
│ ↓ autocompleteTopics(query) │  
└──────────┬──────────────────┘
           │ GET /api/autocomplete-topics?q=...
           ↓
┌───────────────────────────────┐
│ autocomplete-topics Function  │
│ ├─ Sanitize query             │
│ ├─ Query Firestore            │
│ │   WHERE topicLowercase >= q │
│ │   ORDER BY createdAt DESC   │
│ │   LIMIT 10                  │
│ └─ Return debates metadata    │
└───────────┬───────────────────┘
            │ JSON response
            ↓
┌────────────────────────────┐
│ TopicAutocompleteDropdown  │  Display results
│ ├─ Topic + avatars + date  │  
│ └─ Keyboard navigation     │
└────────┬───────────────────┘
         │ User selects debate
         ↓
┌──────────────────────────┐
│ navigate('/panelist-...  │
│   state: {               │
│     source: 'autocomplete'
│     debateId: uuid       │
│     preFilled: [...]     │
│   }                      │
└────────┬─────────────────┘
         │
         ↓
┌─────────────────────────────┐
│  PanelistSelection.jsx      │
│  ├─ Check state.source      │
│  ├─ Pre-fill panelists      │
│  ├─ Run cache detection     │
│  └─ Show "View Debate" or   │
│      "Modify Panelists"     │
└─────────────────────────────┘
```

## Implementation Steps

### Step 1: Firestore Preparation (5 min)

Create composite index for efficient querying:

```bash
gcloud firestore indexes composite create \
  --collection-group=debates \
  --field-config=field-path=topicLowercase,order=ascending \
  --field-config=field-path=createdAt,order=descending \
  --project=${GCP_PROJECT_ID}
```

Modify `backend/shared/firebase/debates.go` to auto-populate `topicLowercase`:

```go
func SaveDebate(ctx context.Context, debate *DebateDocument) error {
    // Add lowercase field for autocomplete
    debate.TopicLowercase = strings.ToLower(debate.Topic)
    
    // Rest of SaveDebate logic...
}
```

### Step 2: Backend Implementation (1 hour)

Create `backend/functions/autocomplete-topics/`:

**handler.go**:
```go
package autocompletetopics

import (
    "encoding/json"
    "net/http"
    "os"
    "strconv"
    "strings"
    
    "cloud.google.com/go/firestore"
    "github.com/raphink/debate/backend/shared/firebase"
    "github.com/raphink/debate/backend/shared/sanitize"
)

func HandleAutocompleteTopics(w http.ResponseWriter, r *http.Request) {
    // CORS headers
    origin := os.Getenv("ALLOWED_ORIGIN")
    w.Header().Set("Access-Control-Allow-Origin", origin)
    w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
    
    if r.Method == "OPTIONS" {
        w.WriteHeader(http.StatusNoContent)
        return
    }
    
    // Parse and validate query
    query := r.URL.Query().Get("q")
    if len(query) < 3 {
        http.Error(w, `{"error":"Query must be at least 3 characters"}`, http.StatusBadRequest)
        return
    }
    
    // Sanitize and normalize
    query = sanitize.StripHTML(query)
    query = strings.ToLower(query)
    
    // Parse limit
    limit := 10
    if l := r.URL.Query().Get("limit"); l != "" {
        if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 && parsed <= 10 {
            limit = parsed
        }
    }
    
    // Query Firestore (see plan.md for full implementation)
    // ...
}
```

**main.go**:
```go
package autocompletetopics

import (
    "github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

func init() {
    functions.HTTP("AutocompleteTopics", HandleAutocompleteTopics)
}
```

### Step 3: Frontend Hook (30 min)

Create `frontend/src/hooks/useTopicAutocomplete.js`:

```javascript
import { useState, useEffect, useRef } from 'react';
import { autocompleteTopics } from '../services/api';

export const useTopicAutocomplete = (query, minLength = 3) => {
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const debounceTimer = useRef(null);

  useEffect(() => {
    if (debounceTimer.current) clearTimeout(debounceTimer.current);
    
    if (query.length < minLength) {
      setSuggestions([]);
      return;
    }

    setLoading(true);
    debounceTimer.current = setTimeout(async () => {
      try {
        const data = await autocompleteTopics(query, 10);
        setSuggestions(data.debates || []);
      } catch (err) {
        console.error('Autocomplete error:', err);
        setSuggestions([]);
      } finally {
        setLoading(false);
      }
    }, 300);

    return () => clearTimeout(debounceTimer.current);
  }, [query]);

  return { suggestions, loading };
};
```

### Step 4: Dropdown Component (1 hour)

Create `frontend/src/components/TopicInput/TopicAutocompleteDropdown.jsx`:

Key features:
- Absolute positioning below input
- Display: topic, avatars (24px circular), count badge, date
- Keyboard navigation (↑↓ arrows, Enter, Escape)
- Click-outside detection
- ARIA labels for accessibility

(See plan.md for full component code)

### Step 5: Integration (1 hour)

**Modify `TopicInput.jsx`**:
```javascript
const [query, setQuery] = useState('');
const [showDropdown, setShowDropdown] = useState(false);
const { suggestions, loading } = useTopicAutocomplete(query);

const handleSelect = (debate) => {
  setShowDropdown(false);
  onAutocompleteSelect(debate); // Callback to Home.jsx
};

return (
  <>
    <input value={query} onChange={(e) => setQuery(e.target.value)} />
    {showDropdown && suggestions.length > 0 && (
      <TopicAutocompleteDropdown 
        suggestions={suggestions}
        loading={loading}
        onSelect={handleSelect}
        onClose={() => setShowDropdown(false)}
      />
    )}
  </>
);
```

**Modify `Home.jsx`**:
```javascript
const handleAutocompleteSelect = (debate) => {
  navigate('/panelist-selection', {
    state: {
      source: 'autocomplete',
      debateId: debate.id,
      topic: debate.topic,
      preFilled: debate.panelists,
    }
  });
};
```

**Modify `PanelistSelection.jsx`**:
```javascript
const { state } = useLocation();
const isFromAutocomplete = state?.source === 'autocomplete';

useEffect(() => {
  if (isFromAutocomplete) {
    setSelectedPanelists(state.preFilled);
    const cacheHit = isCacheHit(state.debateId, state.preFilled);
    setShowViewDebateButton(cacheHit);
  }
}, []);
```

### Step 6: Cache Detection Utility (30 min)

Create `frontend/src/utils/cacheDetection.js`:

```javascript
export function isCacheHit(originalDebateId, currentPanelists, originalPanelists) {
  // Deep comparison: topic exact match + panelist array (order-independent)
  const normalize = (p) => ({ id: p.id, name: p.name });
  
  const orig = originalPanelists.map(normalize).sort((a,b) => a.id.localeCompare(b.id));
  const curr = currentPanelists.map(normalize).sort((a,b) => a.id.localeCompare(b.id));
  
  return JSON.stringify(orig) === JSON.stringify(curr);
}
```

## Local Development

1. **Start services**:
```bash
docker-compose up --build
```

2. **Generate test debates**:
- Open http://localhost:3000
- Create 5-10 debates with varied topics
- Verify saved to Firestore (check GCP Console or emulator UI)

3. **Test autocomplete**:
- Type topic containing substring from previous debates
- Verify dropdown appears with suggestions
- Test keyboard navigation (↑↓, Enter, Escape)
- Select debate, verify navigation to PanelistSelection with pre-filled panelists

4. **Test cache detection**:
- Select debate from autocomplete
- Don't modify panelists
- Verify "View Debate" button appears
- Click "Modify Panelists", change list
- Verify button changes to "Generate New Debate"

## Deployment

1. **Deploy Firestore index**:
```bash
gcloud firestore indexes composite create ...
```

2. **Deploy Cloud Function**:
```bash
cd backend/functions/autocomplete-topics
gcloud functions deploy autocomplete-topics \
  --gen2 \
  --runtime=go124 \
  --region=us-central1 \
  --entry-point=HandleAutocompleteTopics \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars ALLOWED_ORIGIN=https://raphink.github.io,GCP_PROJECT_ID=${GCP_PROJECT_ID}
```

3. **Deploy frontend**:
```bash
cd frontend
npm run build
# Deploy to GitHub Pages (existing workflow)
```

4. **Update API base URL** in `frontend/src/services/api.js`:
```javascript
const AUTOCOMPLETE_URL = process.env.REACT_APP_AUTOCOMPLETE_URL || 
  'https://us-central1-${PROJECT}.cloudfunctions.net/autocomplete-topics';
```

## Testing Checklist

- [ ] Backend: Query parameter validation (min 3 chars, max 10 limit)
- [ ] Backend: Sanitization strips HTML tags
- [ ] Backend: Firestore query returns correct substring matches
- [ ] Backend: Results ordered by createdAt descending
- [ ] Backend: CORS headers present
- [ ] Frontend: Debouncing prevents rapid API calls
- [ ] Frontend: Dropdown appears when typing ≥3 chars
- [ ] Frontend: Dropdown shows avatars, count, date
- [ ] Frontend: Keyboard navigation works (↑↓, Enter, Escape)
- [ ] Frontend: Click-outside closes dropdown
- [ ] Integration: Selecting debate navigates with pre-filled panelists
- [ ] Integration: Cache hit shows "View Debate" button
- [ ] Integration: Modifying panelists shows "Generate New Debate"
- [ ] Edge case: Empty Firestore hides dropdown
- [ ] Edge case: Firestore failure doesn't break "Find Panelists"

## Common Issues

**Issue**: Dropdown doesn't appear when typing  
**Fix**: Check browser console for CORS errors; verify ALLOWED_ORIGIN environment variable

**Issue**: Firestore query returns no results  
**Fix**: Verify `topicLowercase` field exists in debate documents; rebuild index if needed

**Issue**: Avatars not loading in dropdown  
**Fix**: Check get-portrait endpoint is running; verify CORS headers

**Issue**: Cache detection always shows "Generate New Debate"  
**Fix**: Verify `isCacheHit` utility comparison logic; check console for panelist data structure

## Performance Monitoring

- Monitor Firestore quota usage in GCP Console (Firestore → Usage tab)
- Target: <500ms p95 response time for autocomplete API
- Alert if daily read quota approaches 50K (free tier limit)
- Consider client-side caching if quota issues arise

## Next Steps

After implementing US6:
1. Collect user feedback on autocomplete UX
2. Monitor Firestore quota usage over 1 week
3. Consider implementing US7 (Panelist Autocomplete) based on US6 success
4. Evaluate need for full-text search (Algolia) if substring matching insufficient
