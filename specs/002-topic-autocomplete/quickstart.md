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
│ ↓ fetchDebateHistory(query) │  
└──────────┬──────────────────┘
           │ GET /api/list-debates?q=...
           ↓
┌───────────────────────────────┐
│ list-debates Function         │
│ ├─ Sanitize query             │
│ ├─ Fetch recent debates (100) │
│ ├─ Filter by substring in Go  │
│ │   strings.Contains(lower)   │
│ │   ORDER BY createdAt DESC   │
│ │   LIMIT 10 matches          │
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
│ navigate(`/d/${debate.id}`)│
│                          │
│ Direct to DebateViewer   │
└────────┬─────────────────┘
         │
         ↓
┌─────────────────────────────┐
│  DebateViewer.jsx           │
│  ├─ Load debate by ID       │
│  ├─ Display all messages    │
│  └─ Show panelist info      │
└─────────────────────────────┘
```

## Implementation Steps

### Step 1: Verify Firestore Access (2 min)

**No special Firestore index needed** - Per R001 decision in plan.md, the implementation fetches recent debates (~100) and filters by substring in code. This approach:
- Avoids complex Firestore range queries
- Provides true substring matching (not just prefix)
- Eliminates index creation and maintenance

**Verification**:
```bash
# Check that debates collection exists and has data
gcloud firestore collections describe debates --project=${GCP_PROJECT_ID}
```

### Step 2: Backend Implementation (Already Complete)

**No new Cloud Function needed** - Per R001 decision, autocomplete is integrated into the existing `list-debates` function via the `?q=` query parameter.

**Verify existing implementation** in `backend/functions/list-debates/`:
- `handler.go`: Checks for `q` parameter and calls `AutocompleteDebates()`
- `firestore.go`: Contains `AutocompleteDebates()` function that fetches recent debates and filters in code

### Step 3: Frontend Implementation (1.5 hours)

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
  navigate(`/d/${debate.id}`);
};
```

**No changes needed to PanelistSelection.jsx** - Autocomplete navigates directly to debate viewer.

### Step 6: Test Integration (15 min)

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
- Select debate from dropdown
- Verify navigation to /d/{debate.id} and complete debate displays

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
- [ ] Backend: Substring matching works correctly (case-insensitive)
- [ ] Backend: Results ordered by createdAt descending
- [ ] Backend: CORS headers present
- [ ] Frontend: Debouncing prevents rapid API calls
- [ ] Frontend: Dropdown appears when typing ≥3 chars
- [ ] Frontend: Dropdown shows avatars, count, date
- [ ] Frontend: Keyboard navigation works (↑↓, Enter, Escape)
- [ ] Frontend: Click-outside closes dropdown
- [ ] Integration: Selecting debate navigates to /d/{debate.id}
- [ ] Integration: Debate viewer displays complete existing debate
- [ ] Edge case: Empty Firestore hides dropdown
- [ ] Edge case: Firestore failure doesn't break "Find Panelists"

## Common Issues

**Issue**: Dropdown doesn't appear when typing  
**Fix**: Check browser console for CORS errors; verify ALLOWED_ORIGIN environment variable

**Issue**: Autocomplete returns no results  
**Fix**: Verify debates exist in Firestore; check substring filtering logic in backend/functions/list-debates/firestore.go

**Issue**: Avatars not loading in dropdown  
**Fix**: Check get-portrait endpoint is running; verify CORS headers

**Issue**: Navigation doesn't work after selecting suggestion  
**Fix**: Verify DebateViewer route is configured at /d/:uuid in App.jsx; check debate.id is valid UUID

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
