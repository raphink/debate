# User Story 6: Topic Autocomplete and History Integration

## Overview

Enable users to discover and reuse previous debates through an integrated autocomplete dropdown on the home page topic input field. When users type a topic, they see suggestions from historical debates, can select one to pre-fill panelists, and either load the cached debate directly or modify panelists to generate a new variation.

## Business Value

**Problem**: Users cannot easily discover previous debates they've generated or reuse successful topic/panelist combinations. This leads to:
- Redundant debate generation (same topic/panelists generated multiple times)
- Inefficient UX (manual re-entry of known good combinations)
- Lack of continuity between debate sessions
- Wasted Claude API credits for duplicate content

**Solution**: Integrated autocomplete that:
- Surfaces relevant historical debates as users type
- Enables one-click reuse of previous topic/panelist combinations
- Detects cache hits to skip redundant generation
- Allows modifications while maintaining discovery

**Impact**: 
- Reduced API costs from duplicate generation prevention
- Improved UX through streamlined topic/panelist selection
- Better content discovery and reuse
- Faster workflow for iterative debate variations

## User Journey

### Happy Path: Reuse Cached Debate

1. User visits home page and begins typing topic: "ethics of AI"
2. After 3 characters ("eth"), autocomplete dropdown appears showing:
   - "Ethics of AI in healthcare" (3 panelists)
   - "Ethical considerations for autonomous vehicles" (4 panelists)
3. User selects first option
4. System validates topic (no Claude call needed) and navigates to PanelistSelection
5. PanelistSelection page shows pre-filled panelist chips: Aristotle, Kant, Peter Singer
6. User clicks "Generate Debate" without modifications
7. System detects cache hit (exact topic + panelist match)
8. System immediately redirects to `/d/{uuid}` showing cached debate
9. No API calls to Claude, instant load from Firestore

### Alternative Path: Modify and Generate New

1-5. Same as happy path
6. User notices pre-filled panelists and clicks "Modify Panelists" button
7. Chips become editable, user removes Peter Singer, adds John Stuart Mill
8. User clicks "Generate Debate"
9. System detects no cache hit (panelist list changed)
10. System generates new debate with modified panelists
11. New debate saved to Firestore with new UUID
12. User sees streaming generation as normal

### Edge Cases

**No Historical Matches**:
- User types topic with no previous matches
- Autocomplete dropdown remains hidden
- User proceeds normally with manual topic submission and Claude validation

**Firestore Unavailable**:
- Autocomplete API fails
- Dropdown hidden gracefully
- User workflow unaffected (degrades to manual input)

**Slow Network**:
- User types quickly
- Loading indicator appears after 300ms delay
- Debouncing prevents excessive API calls
- User can continue typing without blocking

## Requirements

### Functional Requirements

**FR-US6-1: Topic Autocomplete Endpoint**
- MUST provide HTTP GET `/api/autocomplete-topics` endpoint
- MUST accept query parameters: `q` (query string) and `limit` (max results, default 10)
- MUST perform case-insensitive substring search on topic text in Firestore
- MUST return results ordered by creation timestamp DESC (newest first)
- MUST return debate metadata: `{id, topic, panelists: [{id, name, slug}], panelistCount, createdAt}`
- MUST respond with CORS headers for frontend access
- MUST handle errors gracefully with appropriate HTTP status codes

**FR-US6-2: Frontend Autocomplete Integration**
- MUST show autocomplete dropdown when user types ≥3 characters in topic input
- MUST debounce API calls by 300ms to reduce server load
- MUST display up to 10 matching topics with panelist count badge (e.g., "3 panelists")
- MUST support keyboard navigation (arrow keys, enter to select)
- MUST support mouse click selection
- MUST hide dropdown when clicking outside or pressing Escape
- MUST show loading indicator if API response >300ms
- MUST degrade gracefully if API fails (hide dropdown, allow manual input)

**FR-US6-3: Panelist Pre-fill and Navigation**
- MUST skip Claude topic validation when topic selected from autocomplete
- MUST navigate to PanelistSelection with pre-filled panelists from selected debate
- MUST pass debate ID and panelists via navigation state (not URL params)
- MUST display pre-filled panelists as locked chips initially
- MUST provide "Modify Panelists" button to enable editing

**FR-US6-4: Cache Hit Detection**
- MUST compare topic text (exact match, case-sensitive) and panelist array (deep equality by ID/slug)
- MUST detect cache hit when topic + panelists unchanged from historical debate
- MUST redirect directly to `/d/{uuid}` when cache hit detected
- MUST bypass generate-debate API call on cache hit
- MUST show visual indicator: "Loading cached debate..."

**FR-US6-5: Modified Panelist Flow**
- MUST allow panelist chip editing when "Modify Panelists" clicked
- MUST detect cache miss when panelist list modified (add/remove/reorder)
- MUST generate new debate with new UUID when cache miss
- MUST show visual indicator: "Generating new debate..."

### Non-Functional Requirements

**NFR-US6-1: Performance**
- Autocomplete API response time MUST be <200ms for 90th percentile
- Firestore query MUST use index on `topic` field for efficient search
- Frontend debouncing MUST prevent excessive API calls (<1 call per 300ms)

**NFR-US6-2: Scalability**
- Backend MUST handle 100 concurrent autocomplete requests
- Firestore query MUST limit results to prevent excessive data transfer

**NFR-US6-3: UX**
- Autocomplete dropdown MUST NOT block typing or interfere with manual input
- Loading indicators MUST appear only after 300ms to avoid flicker on fast networks
- Keyboard navigation MUST feel native and responsive

**NFR-US6-4: Reliability**
- System MUST function normally if autocomplete fails (graceful degradation)
- Cache hit detection MUST be deterministic and accurate (no false positives)

## Acceptance Criteria

### Backend

- [ ] `autocomplete-topics` Cloud Function exists in `backend/functions/autocomplete-topics/`
- [ ] Function accepts GET requests with `q` and optional `limit` query params
- [ ] Function queries Firestore debates collection using case-insensitive topic substring match
- [ ] Results ordered by `createdAt` DESC, limited to `limit` (default 10)
- [ ] Response includes: `{debates: [{id, topic, panelists, panelistCount, createdAt}]}`
- [ ] CORS headers included in response
- [ ] Error responses use appropriate HTTP status codes (400, 500, etc.)
- [ ] Function deployed via `deploy.sh`

### Frontend - Autocomplete UI

- [ ] `useTopicAutocomplete` hook created in `frontend/src/hooks/`
- [ ] Hook implements 300ms debounce, minimum 3 characters to trigger
- [ ] Hook manages loading/error states
- [ ] `TopicAutocompleteDropdown` component created
- [ ] Dropdown shows topic text + panelist count badge
- [ ] Dropdown supports keyboard navigation (Up/Down/Enter/Escape)
- [ ] Dropdown closes on click outside or Escape
- [ ] Loading indicator shown only after 300ms delay
- [ ] Graceful degradation if API fails (dropdown hidden)

### Frontend - Integration

- [ ] Home.jsx integrates TopicAutocompleteDropdown with topic input
- [ ] Selecting topic from dropdown skips Claude validation
- [ ] Navigation to PanelistSelection includes debate ID + panelists in state
- [ ] PanelistSelection accepts and displays pre-filled panelists
- [ ] "Modify Panelists" button unlocks chips for editing
- [ ] Cache hit detection utility created in `frontend/src/utils/cacheDetection.js`
- [ ] Deep equality check on topic + panelists array
- [ ] Cache hit redirects to `/d/{uuid}` immediately
- [ ] Cache miss triggers normal debate generation flow
- [ ] Visual indicators differentiate "Loading cached debate" vs "Generating new debate"

### Testing

- [ ] Backend: Query `q="ethics"` returns debates with "ethics" in topic
- [ ] Backend: Results ordered by `createdAt` DESC
- [ ] Backend: `limit=5` returns max 5 results
- [ ] Frontend: Typing 3+ chars shows dropdown
- [ ] Frontend: Selecting topic navigates to PanelistSelection with pre-filled panelists
- [ ] Frontend: Unchanged panelists trigger cache hit, redirect to `/d/{uuid}`
- [ ] Frontend: Modified panelists trigger new generation
- [ ] E2E: Generate debate → return home → autocomplete suggests it → select → panelists pre-filled → generate → cache hit

## Technical Design Notes

### Backend Architecture

**Firestore Query Strategy**:
```go
// Case-insensitive substring search requires:
// 1. Lowercase normalization of stored topics
// 2. Index on lowercase_topic field
// 3. Range query: >= query && < query + 'z'

collection.Where("topic_lowercase", ">=", strings.ToLower(query)).
           Where("topic_lowercase", "<", strings.ToLower(query)+"z").
           OrderBy("createdAt", firestore.Desc).
           Limit(limit)
```

**Response Schema**:
```json
{
  "debates": [
    {
      "id": "uuid-string",
      "topic": "Ethics of AI in healthcare",
      "panelists": [
        {"id": "aristotle", "name": "Aristotle", "slug": "aristotle"},
        {"id": "kant", "name": "Immanuel Kant", "slug": "immanuel-kant"}
      ],
      "panelistCount": 2,
      "createdAt": "2025-01-15T10:30:00Z"
    }
  ]
}
```

### Frontend Architecture

**Hook: useTopicAutocomplete**
```javascript
export const useTopicAutocomplete = (query, enabled = true) => {
  const [suggestions, setSuggestions] = useState([]);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState(null);

  // Debounce 300ms, min 3 chars
  useEffect(() => {
    if (!enabled || query.length < 3) {
      setSuggestions([]);
      return;
    }

    const timer = setTimeout(async () => {
      setLoading(true);
      try {
        const data = await api.autocompleteTopics(query);
        setSuggestions(data.debates || []);
        setError(null);
      } catch (err) {
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

**Component: TopicAutocompleteDropdown**
- Renders below topic input field
- Absolutely positioned dropdown
- Shows max 10 items
- Each item: topic text (primary) + badge with panelist count (secondary)
- Keyboard navigation state (selectedIndex)
- Click handlers for selection
- Escape/outside click to close

**Cache Detection Algorithm**:
```javascript
export const isCacheHit = (selectedDebate, currentTopic, currentPanelists) => {
  // Topic must match exactly (case-sensitive for UX clarity)
  if (selectedDebate.topic !== currentTopic) return false;

  // Panelists must match by ID/slug and count
  if (selectedDebate.panelists.length !== currentPanelists.length) return false;

  // Deep equality check on panelist IDs
  const selectedIds = selectedDebate.panelists.map(p => p.id).sort();
  const currentIds = currentPanelists.map(p => p.id).sort();

  return selectedIds.every((id, i) => id === currentIds[i]);
};
```

## Success Metrics

**Quantitative**:
- >80% of repeat users select from autocomplete vs manual input
- 50% reduction in duplicate debate generation (same topic + panelists)
- <200ms autocomplete API response time (p90)
- <1s time to load cached debate on cache hit

**Qualitative**:
- Users report improved workflow for iterative debate generation
- Reduced friction in topic/panelist selection process
- Positive feedback on autocomplete responsiveness and accuracy

## Dependencies

- **US5**: Firestore integration must be complete for historical debate storage
- **US1**: Topic validation flow for integration point
- **US2**: Panelist selection page for pre-fill integration

## Future Enhancements

- **Smart Ranking**: Weight autocomplete results by frequency of access, not just recency
- **Fuzzy Matching**: Support typo tolerance in topic search (Levenshtein distance)
- **Category Filters**: Allow filtering autocomplete by debate category/tags
- **Infinite Scroll**: Load more than 10 results on demand
- **Analytics**: Track autocomplete usage, selection rates, cache hit rates
