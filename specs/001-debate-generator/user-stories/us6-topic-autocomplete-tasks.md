# US6 Topic Autocomplete - Implementation Tasks

## Phase 1: Backend Infrastructure (5 tasks)

### Task 1.1: Create Function Scaffolding
- [X] Create directory: `backend/functions/autocomplete-topics/`
- [X] Create `go.mod` with dependencies
- [X] Create `handler.go` with package and imports
- [X] Create `cmd/main.go` entry point for Cloud Functions

### Task 1.2: Define Types and Structs
- [X] Create `types.go`:
  - `DebateSummary` struct (ID, Topic, Panelists, PanelistCount, CreatedAt) with JSON tags
  - `PanelistSummary` struct (ID, Name, Slug) with JSON tags
  - `AutocompleteResponse` struct with Debates slice
- [X] Test: Ensure JSON marshaling works correctly

### Task 1.3: Implement Query Logic
- [X] Create `handler.go` with `queryDebates()` function:
  - Accept context, query string, limit
  - Normalize query: `q := strings.ToLower(query)`
  - Fetch recent debates from Firestore:
    ```go
    OrderBy("startedAt", firestore.Desc).
    Limit(100)
    ```
  - Filter client-side: `strings.Contains(strings.ToLower(topic), q)`
  - Return first `limit` matches
- [X] Create `transformToSummary()` helper function
- [X] Test: Query with "eth" should match "Ethics of AI" anywhere in topic

### Task 1.4: Implement HTTP Handler
- [X] Create `AutocompleteTopicsHandler()` in `handler.go`:
  - Set CORS headers (allow all origins)
  - Handle OPTIONS preflight
  - Check method is GET, return 405 if not
  - Parse `q` query param, validate >= 3 chars (400 if invalid)
  - Parse `limit` query param, default to 10
  - Initialize Firestore client if needed
  - Call `queryDebates()`
  - Return JSON: `{"debates": [...]}`
  - Handle errors with appropriate status codes
- [X] Test: Manual curl request should return JSON

### Task 1.5: Create Dockerfile and Deployment
- [X] Create `Dockerfile` in `backend/functions/autocomplete-topics/`:
  - Multi-stage build with golang:1.24-alpine
  - Build binary using distroless runtime
  - Expose 8080, CMD to run binary
- [X] Update `docker-compose.yml` to include autocomplete-topics service
- [X] Update `deploy.sh` to deploy autocomplete-topics function

**Checkpoint**: Backend endpoint functional, returns matching debates

---

## Phase 2: Frontend Autocomplete UI (8 tasks)

### Task 2.1: API Service Method
- [ ] Update `frontend/src/services/api.js`:
  - Add `autocompleteTopics(query, limit = 10)` function
  - Validate query length >= 3, return empty if not
  - Build URL with query params: `/api/autocomplete-topics?q=${query}&limit=${limit}`
  - Fetch GET request with JSON headers
  - Return parsed JSON or throw error
- [ ] Test: Call with "ethics" should return debates array

### Task 2.2: Create useTopicAutocomplete Hook
- [ ] Create `frontend/src/hooks/useTopicAutocomplete.js`:
  - Accept `query` and `enabled` props
  - State: `suggestions`, `loading`, `error`
  - useEffect with 300ms debounce timer
  - Only trigger if `enabled && query.length >= 3`
  - Call `api.autocompleteTopics(query)`
  - Update states appropriately
  - Clear suggestions if query < 3 chars
- [ ] Test: Rapid typing should only trigger one API call after 300ms

### Task 2.3: Create TopicAutocompleteDropdown Component - Structure
- [ ] Create `frontend/src/components/TopicInput/TopicAutocompleteDropdown.jsx`:
  - Props: `suggestions`, `onSelect`, `loading`, `visible`, `onClose`
  - State: `selectedIndex` for keyboard navigation
  - Refs: `dropdownRef` for click-outside detection
  - PropTypes validation
- [ ] Create `TopicAutocompleteDropdown.module.css`:
  - `.dropdown`: absolute positioning, shadow, rounded corners
  - `.list`, `.item`, `.topic`, `.badge` styles
  - `.selected` highlight style
  - `.loading` spinner/message style

### Task 2.4: Implement Keyboard Navigation
- [ ] In `TopicAutocompleteDropdown.jsx`, add useEffect for keyboard events:
  - Listen for ArrowUp/Down: update `selectedIndex` (wrap around)
  - Listen for Enter: call `onSelect(suggestions[selectedIndex])`
  - Listen for Escape: call `onClose()`
  - Prevent default on all handled keys
  - Clean up listener on unmount
- [ ] Test: Arrow keys navigate, Enter selects, Escape closes

### Task 2.5: Implement Click-Outside Detection
- [ ] Add useEffect for click-outside:
  - Listen for mousedown events
  - Check if target is outside `dropdownRef.current`
  - Call `onClose()` if outside
  - Clean up listener
- [ ] Test: Clicking outside dropdown closes it

### Task 2.6: Render Dropdown UI
- [ ] Implement render logic:
  - Return null if `!visible` or empty suggestions
  - If loading: show loading indicator
  - Else: map suggestions to list items
  - Each item: topic text + panelist count badge
  - Apply `.selected` class to `selectedIndex`
  - onClick handler calls `onSelect(debate)`
  - onMouseEnter updates `selectedIndex`
- [ ] Test: Suggestions render with topic and badge

### Task 2.7: Style Autocomplete Dropdown
- [ ] Finalize CSS in `TopicAutocompleteDropdown.module.css`:
  - Gradient badge matching brand colors (#667eea → #764ba2)
  - Hover/selected states with smooth transitions
  - Responsive max-height (400px), scrollable
  - Subtle shadow and border
  - Loading state centered
- [ ] Test: Visual QA across desktop/mobile

### Task 2.8: Integrate into Home.jsx
- [ ] Update `frontend/src/pages/Home.jsx`:
  - Import `useTopicAutocomplete` and `TopicAutocompleteDropdown`
  - Add state: `showAutocomplete`, `selectedDebateMetadata`
  - Call hook: `const { suggestions, loading } = useTopicAutocomplete(topic, showAutocomplete)`
  - Update topic input: `onChange` sets topic and `showAutocomplete=true`
  - Add `onFocus` to show autocomplete
  - Render `<TopicAutocompleteDropdown>` below input (relative positioning wrapper)
  - Implement `handleAutocompleteSelect`:
    - Store debate metadata
    - Set topic
    - Close dropdown
    - Navigate to `/panelist-selection` with state: `{debateId, topic, panelists, skipValidation: true}`
- [ ] Test: Type topic → dropdown appears → select → navigates with correct state

**Checkpoint**: Autocomplete fully functional on home page

---

## Phase 3: Cache Detection and Modified Flow (6 tasks)

### Task 3.1: Create Cache Detection Utility
- [ ] Create `frontend/src/utils/cacheDetection.js`:
  - Export `isCacheHit(originalDebate, currentTopic, currentPanelists)` function
  - Return false if `!originalDebate`
  - Compare `originalDebate.topic === currentTopic` (exact match)
  - Compare panelist count
  - Sort panelist IDs from both arrays
  - Compare sorted ID arrays for deep equality
  - Return true only if all checks pass
- [ ] Add unit tests for various scenarios (match, topic diff, panelist diff, order diff)

### Task 3.2: Update PanelistSelection - Read Navigation State
- [ ] Update `frontend/src/pages/PanelistSelection.jsx`:
  - Import `useLocation` from react-router-dom
  - Extract from `location.state`: `debateId`, `topic`, `panelists`, `skipValidation`
  - Add state: `isLocked`, `showCacheIndicator`, `originalDebateData`
  - Store original debate data if provided
- [ ] Test: Navigate from autocomplete → state correctly populated

### Task 3.3: Implement Pre-fill Logic
- [ ] In `PanelistSelection.jsx`, add useEffect:
  - Check if `panelists` provided in location.state
  - If yes: `setSelectedPanelists(panelists)`, `setIsLocked(true)`, `setShowCacheIndicator(true)`
  - Store `originalDebateData = { id: debateId, topic, panelists }`
- [ ] Render locked chips: disable remove button, gray background
- [ ] Show cache indicator banner: "✓ Using cached debate"
- [ ] Test: Pre-filled panelists appear locked with indicator

### Task 3.4: Implement "Modify Panelists" Button
- [ ] Add "Modify Panelists" button when `isLocked === true`
  - Position above "Generate Debate" button
  - Style: orange/warning color to indicate change
  - onClick: `setIsLocked(false)`, `setShowCacheIndicator(false)`
- [ ] When unlocked: enable chip removal/addition
- [ ] Update cache indicator: hide or change to "Generating new debate"
- [ ] Test: Click "Modify" → chips become editable, indicator changes

### Task 3.5: Implement Cache Hit Detection on Generate
- [ ] Update `handleGenerateDebate()` in `PanelistSelection.jsx`:
  - Import `isCacheHit` from utils
  - Before generating, check if `originalDebateData` exists
  - If yes: call `isCacheHit(originalDebateData, topic, selectedPanelists)`
  - If cache hit: navigate to `/d/${originalDebateData.id}` (skip generation)
  - If cache miss: proceed with normal generation flow
  - Log detection result for debugging
- [ ] Test: Unchanged panelists → redirect to cached debate
- [ ] Test: Modified panelists → new debate generated

### Task 3.6: Add Visual Indicators
- [ ] Update CSS in `PanelistSelection.module.css`:
  - `.cacheIndicator`: green background, border, icon
  - `.chip.locked`: gray background, opacity 0.7, cursor not-allowed
  - `.modifyButton`: orange gradient, hover state
  - `.generateButton`: conditional text based on cache state
- [ ] Update button text: "Load Debate" vs "Generate Debate"
- [ ] Test: Visual states clear and distinct

**Checkpoint**: Full cache detection and modify flow functional

---

## Phase 4: Testing and Polish (7 tasks)

### Task 4.1: Backend Integration Test
- [ ] Create `backend/functions/autocomplete-topics/integration_test.go`:
  - Setup test Firestore emulator
  - Insert sample debates with various topics
  - Test case 1: Query "ethics" returns matching debates
  - Test case 2: Results ordered by createdAt DESC
  - Test case 3: Limit parameter respected (request 5, get max 5)
  - Test case 4: Query < 3 chars returns 400 error
- [ ] Run tests: `go test ./...`

### Task 4.2: Frontend Hook Test
- [ ] Create `frontend/src/hooks/useTopicAutocomplete.test.js`:
  - Mock `api.autocompleteTopics`
  - Test debouncing: rapid queries → single API call
  - Test min chars: query < 3 → no API call
  - Test loading states
  - Test error handling
- [ ] Run tests: `npm test`

### Task 4.3: Frontend Component Test
- [ ] Create `frontend/src/components/TopicInput/TopicAutocompleteDropdown.test.jsx`:
  - Test rendering suggestions
  - Test keyboard navigation (ArrowUp/Down, Enter, Escape)
  - Test click selection
  - Test click-outside closes dropdown
  - Test loading state display
- [ ] Run tests and achieve >80% coverage

### Task 4.4: Cache Detection Unit Test
- [ ] Create `frontend/src/utils/cacheDetection.test.js`:
  - Test exact match returns true
  - Test topic mismatch returns false
  - Test panelist count mismatch returns false
  - Test panelist ID mismatch returns false
  - Test order independence (sorted IDs)
- [ ] Run tests

### Task 4.5: E2E Happy Path Test
- [ ] Manual E2E test scenario:
  1. Generate debate: "Ethics of AI in healthcare" with 3 panelists
  2. Return to home page
  3. Type "eth" in topic input
  4. Verify autocomplete dropdown appears
  5. Verify "Ethics of AI in healthcare (3 panelists)" shown
  6. Click suggestion
  7. Verify navigation to PanelistSelection
  8. Verify 3 panelists pre-filled and locked
  9. Click "Generate Debate" without changes
  10. Verify redirect to `/d/{uuid}` with cached debate
  11. Verify no API call to generate-debate
- [ ] Document results

### Task 4.6: E2E Modified Flow Test
- [ ] Manual E2E test scenario:
  1-8. Same as happy path
  9. Click "Modify Panelists"
  10. Remove one panelist, add different one
  11. Click "Generate Debate"
  12. Verify new debate generated (loading, streaming)
  13. Verify new UUID in URL
  14. Verify new debate saved to Firestore
- [ ] Document results

### Task 4.7: Error Handling and Edge Cases
- [ ] Test network error: disconnect, type topic → verify graceful degradation
- [ ] Test Firestore empty: no debates → autocomplete hidden
- [ ] Test slow network: verify loading indicator after 300ms
- [ ] Test concurrent selections: click multiple items rapidly → no race conditions
- [ ] Test back button: navigate back from cache hit → state preserved
- [ ] Document all edge cases and fixes

**Checkpoint**: All tests passing, feature production-ready

---

## Deployment Tasks (3 tasks)

### Deploy 1: Firestore Index
- [ ] Ensure `firestore.indexes.json` committed
- [ ] Run: `gcloud firestore indexes create --database=(default)`
- [ ] Monitor index creation status in Firebase console
- [ ] Wait for "Index created successfully" (may take 5-10 minutes)

### Deploy 2: Backend Function
- [ ] Run: `./deploy.sh` or deploy autocomplete-topics individually
- [ ] Verify function deployed in GCP console
- [ ] Test endpoint with curl: `curl "https://FUNCTION_URL/autocomplete-topics?q=ethics"`
- [ ] Check function logs for any errors
- [ ] Verify CORS headers in response

### Deploy 3: Frontend Build
- [ ] Build frontend: `cd frontend && npm run build`
- [ ] Deploy to hosting: `firebase deploy --only hosting`
- [ ] Test autocomplete in production environment
- [ ] Monitor browser console for errors
- [ ] Test across devices (desktop, mobile, tablet)

---

## Summary

**Total Tasks**: 31
- Phase 1 Backend: 7 tasks
- Phase 2 Frontend UI: 8 tasks
- Phase 3 Cache Detection: 6 tasks
- Phase 4 Testing: 7 tasks
- Deployment: 3 tasks

**Estimated Effort**: 12-16 hours
- Backend: 3-4 hours
- Frontend: 5-6 hours
- Testing: 3-4 hours
- Deployment/QA: 1-2 hours

**Dependencies**:
- US5 (Firestore integration) must be complete
- Existing US1/US2 flows must be functional
- Firestore index creation may take time (plan accordingly)

**Success Criteria**:
- [ ] Autocomplete appears within 500ms of typing 3+ characters
- [ ] Cache hit detection 100% accurate
- [ ] All tests passing (backend + frontend)
- [ ] No console errors in production
- [ ] Graceful degradation if API fails
- [ ] Responsive on all devices
