# Implementation Plan: Topic Discovery via History Integration (US6)

**Branch**: `002-topic-autocomplete` | **Date**: 2025-12-13 | **Spec**: [spec.md](spec.md)  
**Input**: Feature specification from `/specs/002-topic-autocomplete/spec.md`  
**Dependencies**: Requires US5 (Debate Caching & Sharing) - Firestore debates collection must be populated

## Clarifications (Session 2025-12-13)

- **Concurrent Writes**: Allow duplicates - multiple debates on same topic are acceptable (last write wins, Firestore default behavior)
- **Button Labeling**: "View Debate" for cache hit, "Modify Panelists" alternate button, "Generate New Debate" after modifications
- **Duplicate Differentiation**: Show panelist avatars + generation date in dropdown to distinguish debates with identical topics
- **Input Sanitization**: All autocomplete query input must be sanitized (strip HTML tags, special characters) before Firestore query
- **Keyboard Navigation**: Yes - arrow keys, Enter to select, Escape to close (standard accessibility practice per WCAG 2.1)
- **CRITICAL UX**: ONE single input field for topics (not two separate flows). Autocomplete suggestions appear as users type but "Find Panelists" button always remains available. Users can select from dropdown OR ignore autocomplete and proceed with normal Claude validation.

## Summary

Enhance the existing topic input field (Home.jsx) with autocomplete suggestions from Firestore debate history. As users type ≥3 characters, a dropdown displays up to 10 matching previous debates ordered by recency, showing topic text, panelist avatars, panelist count, and generation date. Users can either:
- **Option A**: Select from dropdown → navigate to PanelistSelection with pre-filled panelists → detect cache hit → show "View Debate" button
- **Option B**: Ignore autocomplete, click "Find Panelists" → proceed with normal Claude validation (US1 flow)

This is a **gracefully-degrading enhancement** that never blocks the existing workflow. Autocomplete failures or empty results hide the dropdown and allow normal topic entry.

**Key Architecture Decision**: Extend existing `list-debates` Cloud Function with optional `q` query parameter for autocomplete filtering, rather than creating separate function. Cache detection happens client-side via deep comparison utility (topic text + panelist array, order-independent). When cache hit detected, system bypasses debate generation and redirects directly to /d/{uuid}. When user modifies panelists, system generates entirely new debate with new UUID.

## Technical Context

**Language/Version**: Go 1.24 (backend), JavaScript/React 18+ (frontend)  
**Primary Dependencies**: 
- Backend: Existing backend/shared/firebase module (Firestore client already initialized)
- Frontend: Existing React Router navigation state, lodash.isequal for deep comparison (or custom utility)
**Storage**: Cloud Firestore (debates collection already exists from US5)  
**Testing**: 
- Backend: Go testing package, Firestore emulator for integration tests
- Frontend: Jest, React Testing Library, MSW for API mocking
**Local Development**: 
- Use existing docker-compose.yml setup
- Add autocomplete-topics function on port 8085
- Firestore emulator already configured from US5
**Target Platform**: 
- Same as US1-US5 (modern browsers, PWA-enabled)
- No new platform requirements
**Performance Goals**: 
- Autocomplete API response: <500ms (p95)
- Debounce delay: 300ms
- Cache detection: <50ms (synchronous, client-side)
- Dropdown render: <100ms
**Constraints**: 
- Firestore free tier: 50K reads/day, 20K writes/day (autocomplete adds read load)
- Maximum 10 results per query (prevent excessive data transfer)
- Debouncing mandatory to prevent API spam
- No pagination for autocomplete results
**Scale/Scope**: 
- New Components: TopicAutocompleteDropdown.jsx, cacheDetection.js utility
- Modified Components: TopicInput.jsx, Home.jsx, PanelistSelection.jsx
- Modified Backend: list-debates Cloud Function (add `q` query parameter)
- New Firestore Query: Substring search on debates.topic field (case-insensitive)
- Estimated Lines of Code: ~500 frontend, ~50 backend (extending existing function)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ Principle I: User-Centric Design (UX First)
- **Status**: PASS
- **Evidence**: Optional enhancement that improves workflow without disrupting existing flow; graceful degradation on failures
- **Validation**: 10 acceptance scenarios covering selection, cache hit, modifications, and edge cases

### ✅ Principle II: Code Quality & Maintainability  
- **Status**: PASS
- **Evidence**: Reuses existing Firestore client (backend/shared/firebase); follows existing Cloud Function pattern (autocomplete-topics mirrors list-debates structure)
- **Validation**: Sanitization utility reuses backend/shared/sanitize; debounce hook follows existing pattern from useTopicValidation

### ✅ Principle III: Simplicity Over Complexity  
- **Status**: PASS with CLARIFICATION NEEDED
- **Evidence**: Simple substring matching (no complex fuzzy search); client-side cache detection (no server-side state)
- **Question**: Should we implement Firestore composite index for debates.topic + createdAt ordering? (Required for case-insensitive substring + sorting)
- **Decision Needed**: Confirm index creation via gcloud command is acceptable for deployment

### ✅ Principle IV: Iterative Development  
- **Status**: PASS
- **Evidence**: Feature builds incrementally on US5 infrastructure; autocomplete is P3 priority (quality-of-life enhancement, not blocker)
- **Validation**: Can be disabled via feature flag if Firestore quota issues arise

### ✅ Principle V: Pragmatic Constraints  
- **Status**: PASS
- **Evidence**: Firestore free tier supports ~1,666 autocomplete queries/day (assuming 30 reads per query for substring search); debouncing reduces call frequency
- **Risk Mitigation**: Client-side debouncing (300ms), max 10 results limit, graceful degradation on quota exhaustion

**GATE STATUS**: ✅ PASS - Proceed to Phase 0 Research

---

## Phase 0: Research & Decisions

### R001: Reuse Existing list-debates Function
**Decision**: Extend existing `list-debates` Cloud Function with optional `q` query parameter for autocomplete
**Rationale**: list-debates already fetches debates from Firestore and handles pagination. Adding a filter parameter is simpler than creating separate autocomplete-topics function. Reduces deployment complexity and maintains single source of truth.
**Implementation**:
```go
// In backend/functions/list-debates/handler.go
query := r.URL.Query().Get("q")
if query != "" {
    // Autocomplete mode: filter by topic substring
    queryLower := strings.ToLower(sanitize.StripHTML(query))
    
    docs, err := client.Collection("debates").
        Where("topicLowercase", ">=", queryLower).
        Where("topicLowercase", "<", queryLower+"~").
        OrderBy("topicLowercase", firestore.Asc).
        OrderBy("createdAt", firestore.Desc).
        Limit(10).
        Documents(ctx).GetAll()
} else {
    // List mode: existing pagination logic
    // ...existing code...
}
```
**Trade-offs**: Requires storing `topicLowercase` field; only supports prefix matching (not substring); requires composite index
**Benefits**: No new Cloud Function deployment, reuses existing CORS/auth configuration, simpler frontend API surface

### R002: Firestore Composite Index
**Decision**: Create composite index for `topicLowercase` (ASC) + `createdAt` (DESC)
**Command**:
```bash
gcloud firestore indexes composite create \
  --collection-group=debates \
  --field-config=field-path=topicLowercase,order=ascending \
  --field-config=field-path=createdAt,order=descending
```
**Alternative**: Use Firestore auto-indexing (slower, requires first query to trigger index creation)

### R003: Avatar Fetching Strategy for Dropdown
**Decision**: Fetch avatars asynchronously after dropdown renders
**Rationale**: Dropdown must appear quickly (<100ms); portraits can lazy-load progressively
**Implementation**:
- AutocompleteDropdown receives panelist data without avatarUrls initially
- useEffect hook fetches portraits via get-portrait endpoint when dropdown opens
- Display placeholder-avatar.svg during fetch, update progressively

**Alternative Considered**: Pre-fetch all debate avatars during autocomplete query (rejected: slower API response, unnecessary for debates user won't select)

### R004: Cache Detection Deep Comparison
**Decision**: Custom JavaScript utility using lodash.isEqual or manual implementation
**Rationale**: Must compare topic text (exact string match) + panelist array (order-independent, id + name fields only)
**Implementation**:
```javascript
export function isCacheHit(originalDebate, currentState) {
  if (originalDebate.topic !== currentState.topic) return false;
  
  // Extract comparable panelist fields
  const normalize = (p) => ({ id: p.id, name: p.name });
  const origPanelists = originalDebate.panelists.map(normalize).sort((a,b) => a.id.localeCompare(b.id));
  const currPanelists = currentState.selectedPanelists.map(normalize).sort((a,b) => a.id.localeCompare(b.id));
  
  return isEqual(origPanelists, currPanelists);
}
```

### R005: Navigation State Management
**Decision**: Use React Router `navigate(path, { state: {...} })` to pass debate metadata
**Rationale**: Avoid prop drilling; allows PanelistSelection to access original debate ID and panelists from history
**Implementation**:
```javascript
// Home.jsx - on autocomplete selection
navigate('/panelist-selection', {
  state: {
    source: 'autocomplete',
    debateId: selectedDebate.id,
    topic: selectedDebate.topic,
    preFilled: selectedDebate.panelists,
  }
});

// PanelistSelection.jsx
const { state } = useLocation();
const isFromAutocomplete = state?.source === 'autocomplete';
```

### R006: Dropdown Component Library
**Decision**: Custom dropdown component (no external library)
**Rationale**: Existing Button and LoadingSpinner components follow custom styling; maintaining consistency without adding Headless UI or Material-UI dependency
**Implementation**: Create TopicAutocompleteDropdown.jsx with absolute positioning, keyboard event handlers, click-outside detection

---

## Phase 1: Data Model & Contracts

### Data Model Updates

#### Firestore Debates Collection (Extended)
**New Field**:
```typescript
{
  // Existing fields from US5
  id: string;
  topic: string;
  panelists: Panelist[];
  messages: Message[];
  status: 'completed' | 'generating';
  createdAt: Timestamp;
  updatedAt: Timestamp;
  
  // NEW for US6
  topicLowercase: string; // Auto-generated on save: topic.toLowerCase()
}
```

**Migration**: Update backend/shared/firebase/debates.go `SaveDebate` function to auto-populate `topicLowercase`

#### Frontend Navigation State
```typescript
interface AutocompleteNavigationState {
  source: 'autocomplete' | 'manual';
  debateId?: string;           // Original debate UUID (if from autocomplete)
  topic: string;
  preFilled?: Panelist[];      // Pre-selected panelists (if from autocomplete)
}
```

### API Contracts

#### Contract: list-debates (Extended)

**Location**: `/specs/002-topic-autocomplete/contracts/list-debates-autocomplete.json`

```json
{
  "endpoint": "/api/list-debates",
  "method": "GET",
  "description": "List debates with optional topic substring filtering for autocomplete. When `q` parameter provided, returns max 10 matching debates ordered by recency. Without `q`, returns paginated debate history (existing US5 behavior).",
  
  "queryParameters": {
    "q": {
      "type": "string",
      "required": false,
      "minLength": 3,
      "maxLength": 500,
      "description": "[NEW] Optional autocomplete filter (topic substring, case-insensitive). When provided, overrides pagination and returns max 10 matching debates."
    },
    "limit": {
      "type": "integer",
      "required": false,
      "default": 10,
      "min": 1,
      "max": 50,
      "description": "[EXISTING] Maximum results for pagination mode (ignored when q is provided)"
    },
    "cursor": {
      "type": "string",
      "required": false,
      "description": "[EXISTING] Pagination cursor (ignored when q is provided)"
    }
  },
  
  "responses": {
    "200": {
      "description": "Successful autocomplete results",
      "schema": {
        "type": "object",
        "properties": {
          "debates": {
            "type": "array",
            "items": {
              "type": "object",
              "properties": {
                "id": { "type": "string", "format": "uuid" },
                "topic": { "type": "string" },
                "panelistCount": { "type": "integer", "min": 2, "max": 5 },
                "panelists": {
                  "type": "array",
                  "items": {
                    "type": "object",
                    "properties": {
                      "id": { "type": "string" },
                      "name": { "type": "string" },
                      "slug": { "type": "string" },
                      "avatarUrl": { "type": "string", "nullable": true }
                    }
                  }
                },
                "createdAt": { "type": "string", "format": "date-time" }
              }
            }
          }
        }
      },
      "example": {
        "debates": [
          {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "topic": "Is free will compatible with divine omniscience?",
            "panelistCount": 3,
            "panelists": [
              { "id": "augustine", "name": "Augustine of Hippo", "slug": "augustine-of-hippo", "avatarUrl": "https://upload.wikimedia.org/..." },
              { "id": "aquinas", "name": "Thomas Aquinas", "slug": "thomas-aquinas", "avatarUrl": null },
              { "id": "descartes", "name": "René Descartes", "slug": "rene-descartes", "avatarUrl": "https://upload.wikimedia.org/..." }
            ],
            "createdAt": "2025-12-10T14:23:00Z"
          }
        ]
      }
    },
    "400": {
      "description": "Invalid query parameters",
      "schema": {
        "type": "object",
        "properties": {
          "error": { "type": "string" }
        }
      },
      "example": { "error": "Query must be at least 3 characters" }
    },
    "500": {
      "description": "Server error (Firestore failure)",
      "schema": {
        "type": "object",
        "properties": {
          "error": { "type": "string" }
        }
      },
      "example": { "error": "Failed to query debates" }
    }
  },
  
  "headers": {
    "request": {
      "Origin": {
        "required": false,
        "description": "CORS preflight"
      }
    },
    "response": {
      "Access-Control-Allow-Origin": {
        "value": "${ALLOWED_ORIGIN}",
        "description": "CORS header matching frontend domain"
      },
      "Content-Type": {
        "value": "application/json"
      }
    }
  }
}
```

---

## Phase 2: Architecture & Design

### Backend Architecture

#### Modified Cloud Function: list-debates

**File Structure** (No changes, existing files modified):
```
backend/functions/list-debates/
├── Dockerfile               # [EXISTING] No changes
├── go.mod                   # [EXISTING] No changes
├── handler.go               # [MODIFIED] Add q parameter handling
├── main.go                  # [EXISTING] No changes
├── firestore.go             # [MODIFIED] Add autocomplete query function
├── types.go                 # [EXISTING] No changes
└── cmd/
    └── main.go              # [EXISTING] No changes (port 8084)
```

**Handler Logic** (handler.go - modifications):
```go
package listdebates

import (
    "context"
    "encoding/json"
    "log"
    "net/http"
    "strconv"
    "strings"
    
    "cloud.google.com/go/firestore"
    "github.com/raphink/debate/backend/shared/firebase"
    "github.com/raphink/debate/backend/shared/sanitize"
)

type AutocompleteResponse struct {
    Debates []DebateMetadata `json:"debates"`
}

type DebateMetadata struct {
    ID            string      `json:"id"`
    Topic         string      `json:"topic"`
    PanelistCount int         `json:"panelistCount"`
    Panelists     []Panelist  `json:"panelists"`
    CreatedAt     string      `json:"createdAt"`
}

func HandleListDebates(w http.ResponseWriter, r *http.Request) {
    // [EXISTING] CORS headers already present
    
    // [NEW] Check for autocomplete mode
    query := r.URL.Query().Get("q")
    if query != "" {
        // Autocomplete mode
        if len(query) < 3 {
            http.Error(w, `{"error":"Query must be at least 3 characters"}`, http.StatusBadRequest)
            return
        }
        
        // Sanitize and normalize
        query = sanitize.StripHTML(query)
        queryLower := strings.ToLower(query)
        
        // Call autocomplete query function
        ctx := context.Background()
        debates, err := AutocompleteDebates(ctx, queryLower, 10)
        if err != nil {
            log.Printf("Autocomplete error: %v", err)
            http.Error(w, `{"error":"Failed to query debates"}`, http.StatusInternalServerError)
            return
        }
        
        w.Header().Set("Content-Type", "application/json")
        json.NewEncoder(w).Encode(map[string]interface{}{"debates": debates})
        return
    }
    
    // [EXISTING] Pagination mode - existing list-debates logic
    // ...existing code for cursor-based pagination...
}
```

**New Function in firestore.go**:
```go
// AutocompleteDebates queries debates by topic substring (prefix matching)
func AutocompleteDebates(ctx context.Context, queryLower string, limit int) ([]DebateMetadata, error) {
    client := firebase.GetFirestoreClient(ctx)
    
    docs, err := client.Collection("debates").
        Where("topicLowercase", ">=", queryLower).
        Where("topicLowercase", "<", queryLower+"~").
        OrderBy("topicLowercase", firestore.Asc).
        OrderBy("createdAt", firestore.Desc).
        Limit(limit).
        Documents(ctx).GetAll()
    
    if err != nil {
        return nil, err
    }
    
    // Transform to DebateMetadata format
    debates := make([]DebateMetadata, 0, len(docs))
    for _, doc := range docs {
        var debate firebase.DebateDocument
        if err := doc.DataTo(&debate); err != nil {
            continue
        }
        
        debates = append(debates, DebateMetadata{
            ID:            debate.ID,
            Topic:         debate.Topic,
            PanelistCount: len(debate.Panelists),
            Panelists:     debate.Panelists,
            CreatedAt:     debate.CreatedAt.Format(time.RFC3339),
        })
    }
    
    return debates, nil
}
```

### Frontend Architecture

#### Component Structure

```
frontend/src/
├── components/
│   ├── TopicInput/
│   │   ├── TopicInput.jsx              # MODIFIED: Integrate autocomplete dropdown
│   │   ├── TopicInput.module.css       # MODIFIED: Add dropdown positioning styles
│   │   └── TopicAutocompleteDropdown.jsx  # NEW: Dropdown component
│   │       └── TopicAutocompleteDropdown.module.css
├── hooks/
│   ├── useTopicValidation.js           # EXISTING: No changes
│   └── useTopicAutocomplete.js         # NEW: Autocomplete state + debouncing
├── services/
│   ├── api.js                          # MODIFIED: Add autocompleteTopics method
│   └── topicService.js                 # EXISTING: No changes
├── utils/
│   └── cacheDetection.js               # NEW: Deep comparison utility
└── pages/
    ├── Home.jsx                        # MODIFIED: Handle autocomplete selection
    └── PanelistSelection.jsx           # MODIFIED: Pre-fill panelists, cache detection
```

#### New Hook: useTopicAutocomplete

**Location**: `frontend/src/hooks/useTopicAutocomplete.js`

```javascript
import { useState, useEffect, useRef } from 'react';
import { autocompleteTopics } from '../services/api';

export const useTopicAutocomplete = (query, minLength = 3) => {
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);
  const debounceTimer = useRef(null);

  useEffect(() => {
    // Clear previous timer
    if (debounceTimer.current) {
      clearTimeout(debounceTimer.current);
    }

    // Reset if query too short
    if (query.length < minLength) {
      setSuggestions([]);
      setLoading(false);
      return;
    }

    // Debounce: wait 300ms after user stops typing
    setLoading(true);
    debounceTimer.current = setTimeout(async () => {
      try {
        const data = await autocompleteTopics(query, 10);
        setSuggestions(data.debates || []);
        setError(null);
      } catch (err) {
        console.error('Autocomplete error:', err);
        setError(err.message);
        setSuggestions([]);
      } finally {
        setLoading(false);
      }
    }, 300);

    return () => {
      if (debounceTimer.current) {
        clearTimeout(debounceTimer.current);
      }
    };
  }, [query, minLength]);

  return { suggestions, loading, error };
};
```

#### New Component: TopicAutocompleteDropdown

**Location**: `frontend/src/components/TopicInput/TopicAutocompleteDropdown.jsx`

**Features**:
- Absolute positioning below input field
- Keyboard navigation (↑↓ arrows, Enter to select, Escape to close)
- Click-outside detection to close
- Display: topic text, panelist avatars (circular 24px), count badge, date
- Hover/focus states for accessibility
- Loading state (shimmer effect)

**Props**:
```typescript
interface TopicAutocompleteDropdownProps {
  suggestions: DebateMetadata[];
  loading: boolean;
  visible: boolean;
  onSelect: (debate: DebateMetadata) => void;
  onClose: () => void;
}
```

#### Modified Component: TopicInput

**Changes**:
- Integrate useTopicAutocomplete hook
- Render TopicAutocompleteDropdown when suggestions available
- Pass onSelect handler to close dropdown and trigger navigation via Home.jsx callback
- Manage dropdown visibility state (show when typing ≥3 chars, hide on selection or click-outside)

#### Modified Page: Home.jsx

**Changes**:
- Add handler for autocomplete selection: `handleAutocompleteSelect(debate)`
- Navigate to PanelistSelection with state: `{ source: 'autocomplete', debateId, topic, preFilled }`
- Existing "Find Panelists" flow remains unchanged (manual Claude validation)

#### Modified Page: PanelistSelection.jsx

**Changes**:
- Check `location.state.source === 'autocomplete'`
- If autocomplete: pre-select panelists from `location.state.preFilled`
- Run cache detection: `isCacheHit(location.state.debateId, currentPanelists)`
- If cache hit: show "View Debate" button (redirects to /d/{uuid})
- If cache hit: show "Modify Panelists" button (unlocks chips, changes to "Generate New Debate")
- If cache miss (panelists modified): normal debate generation flow

---

## Phase 3: Implementation Details

### Backend Implementation Checklist

**No New Files to Create**

**Files to Modify**:
1. `backend/functions/list-debates/handler.go` - Add `q` parameter handling, autocomplete mode
2. `backend/functions/list-debates/firestore.go` - Add `AutocompleteDebates` function
3. `backend/functions/list-debates/types.go` - Reuse existing DebateMetadata type (if needed)
4. `backend/shared/firebase/debates.go` - Add `topicLowercase` field to SaveDebate function
5. `docker-compose.yml` - No changes needed (list-debates already on port 8084)
6. `deploy.sh` - No changes needed (list-debates already deployed)

### Frontend Implementation Checklist

**Files to Create**:
1. `frontend/src/hooks/useTopicAutocomplete.js` - Debouncing, API calls, state management
2. `frontend/src/components/TopicInput/TopicAutocompleteDropdown.jsx` - Dropdown UI component
3. `frontend/src/components/TopicInput/TopicAutocompleteDropdown.module.css` - Dropdown styles
4. `frontend/src/utils/cacheDetection.js` - Deep comparison utility

**Files to Modify**:
1. `frontend/src/services/api.js` - Modify `listDebates()` to accept optional `query` parameter
2. `frontend/src/components/TopicInput/TopicInput.jsx` - Integrate autocomplete hook + dropdown
3. `frontend/src/components/TopicInput/TopicInput.module.css` - Add dropdown positioning styles
4. `frontend/src/pages/Home.jsx` - Handle autocomplete selection, navigation with state
5. `frontend/src/pages/PanelistSelection.jsx` - Pre-fill panelists, cache detection, button logic

### Deployment Updates

**Firestore Index**:
```bash
gcloud firestore indexes composite create \
  --collection-group=debates \
  --field-config=field-path=topicLowercase,order=ascending \
  --field-config=field-path=createdAt,order=descending \
  --project=${GCP_PROJECT_ID}
```

**Cloud Function Deployment** (deploy.sh - no changes needed):
```bash
# list-debates already deployed from US5
# No new deployment required, just redeploy existing function:
gcloud functions deploy list-debates \
  --gen2 \
  --runtime=go124 \
  --region=us-central1 \
  --source=./backend/functions/list-debates \
  --entry-point=HandleListDebates \
  --trigger-http \
  --allow-unauthenticated \
  --set-env-vars ALLOWED_ORIGIN=https://raphink.github.io,GCP_PROJECT_ID=${GCP_PROJECT_ID}
```

---

## Phase 4: Testing Strategy

### Backend Tests

**Unit Tests** (`handler_test.go`):
- Query parameter validation (min 3 chars, max 10 limit)
- Sanitization correctness (HTML stripping)
- CORS headers present

**Integration Tests** (Firestore emulator):
- Substring matching correctness (prefix matching)
- Ordering by createdAt descending
- Limit enforcement (max 10 results)
- Empty results handling
- Firestore failure graceful degradation

### Frontend Tests

**Unit Tests** (Jest):
- `useTopicAutocomplete` debouncing (300ms delay)
- `cacheDetection.js` deep comparison (order-independent, exact topic match)
- TopicAutocompleteDropdown keyboard navigation

**Component Tests** (React Testing Library):
- TopicInput shows dropdown when typing ≥3 chars
- Dropdown closes on selection
- Dropdown closes on Escape key
- "Find Panelists" button always enabled

**Integration Tests** (MSW):
- Mock autocomplete API response
- Verify navigation state passed to PanelistSelection
- Cache hit detection triggers "View Debate" button
- Panelist modification changes button to "Generate New Debate"

### End-to-End Tests

1. Generate debate → return home → type topic → verify autocomplete appears
2. Select topic from dropdown → verify panelists pre-filled
3. Pre-filled panelists unchanged → verify "View Debate" button → click → verify redirect to /d/{uuid}
4. Click "Modify Panelists" → change list → verify "Generate New Debate" → verify new debate generated
5. Firestore failure → verify dropdown hidden, "Find Panelists" still works

---

## Phase 5: Rollout Plan

### Step 1: Backend Modifications
- Create Firestore composite index
- Modify list-debates handler.go to support `q` parameter
- Add AutocompleteDebates function to firestore.go
- Test autocomplete mode locally
- Verify CORS headers and response format

### Step 2: Frontend Components
- Implement useTopicAutocomplete hook
- Create TopicAutocompleteDropdown component
- Integrate into TopicInput (no breaking changes)

### Step 3: Navigation & Cache Detection
- Update Home.jsx to handle autocomplete selection
- Update PanelistSelection.jsx with pre-fill logic
- Implement cache detection utility

### Step 4: Testing & Validation
- Run backend integration tests (Firestore emulator)
- Run frontend component tests
- Perform manual E2E testing
- Verify graceful degradation (disable Firestore, verify app still works)

### Step 5: Deployment
- Deploy frontend to GitHub Pages
- Redeploy list-debates to GCP (with autocomplete modifications)
- Monitor Firestore quota usage
- Collect user feedback

---

## Success Criteria

**Functional**:
- ✅ Autocomplete appears when typing ≥3 characters
- ✅ Dropdown shows max 10 results ordered by recency
- ✅ Selecting debate navigates to PanelistSelection with pre-filled panelists
- ✅ Cache hit shows "View Debate" button (redirects to /d/{uuid})
- ✅ "Modify Panelists" unlocks chips and generates new debate
- ✅ "Find Panelists" button always works (graceful degradation)

**Non-Functional**:
- ✅ Autocomplete API response <500ms (p95)
- ✅ Debouncing prevents excessive API calls (max 1 per 300ms)
- ✅ Keyboard navigation works (↑↓, Enter, Escape)
- ✅ Dropdown accessible (ARIA labels, focus management)
- ✅ Firestore quota stays within free tier limits

**Edge Cases Handled**:
- ✅ Empty Firestore (no debates): dropdown hidden
- ✅ Firestore failure: graceful degradation to manual entry
- ✅ Network timeout: dropdown closes, "Find Panelists" available
- ✅ Duplicate topics: differentiated by avatars + date

---

## Open Questions / Risks

**Q1**: Firestore substring search performance degrades with large collections (>10K debates). Should we implement full-text search (Algolia, Elasticsearch)?  
**Mitigation**: Monitor query performance; Firestore composite index should handle 10K+ debates efficiently for prefix matching. Defer full-text search to future enhancement.

**Q2**: Storing `topicLowercase` duplicates data. Should we use Firestore Functions to auto-populate on write?  
**Decision**: Manual population in SaveDebate function is simpler; Firestore Functions add deployment complexity.

**Q3**: Should we cache autocomplete results client-side to reduce Firestore reads?  
**Decision**: Defer to observability phase; debouncing already reduces calls significantly. Client-side cache adds complexity (invalidation, stale data).

---

## Next Steps

1. **Phase 0**: Research complete (no additional decisions needed)
2. **Phase 1**: Create contract JSON in `/specs/002-topic-autocomplete/contracts/`
3. **Phase 2**: Generate tasks.md with detailed implementation checklist
4. **Phase 3**: Begin backend implementation (autocomplete-topics function)
5. **Phase 4**: Implement frontend components (useTopicAutocomplete, TopicAutocompleteDropdown)
6. **Phase 5**: Integration testing and deployment
