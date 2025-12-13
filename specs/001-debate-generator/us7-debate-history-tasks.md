# Implementation Tasks: Debate History with Reactive Search

**Feature**: US7-debate-history  
**Date**: 2025-12-13  
**Plan**: [us7-debate-history-plan.md](us7-debate-history-plan.md)

## Task Breakdown

### Phase 1: Backend - list-debates Cloud Function (2-3 hours)

**T201: Setup list-debates function structure** (15 min)
- [X] Create `backend/functions/list-debates/` directory
- [X] Initialize `go.mod` with module name `github.com/raphink/debate/functions/list-debates`
- [X] Add dependencies: `cloud.google.com/go/firestore`, `google.golang.org/api/option`
- [X] Create `main.go`, `handler.go`, `firestore.go`, `types.go`
- [X] Create `cmd/main.go` for local development server
- Files: `backend/functions/list-debates/{main.go,handler.go,firestore.go,types.go,go.mod,cmd/main.go}`

**T202: Implement request/response types** (15 min)
- [X] Define `ListDebatesRequest` struct with Limit, Offset fields
- [X] Define `DebateSummary` struct with ID, Topic, Panelists, StartedAt
- [X] Define `PanelistInfo` struct with ID, Name
- [X] Define `ListDebatesResponse` struct with Debates, Total, HasMore
- [X] Add JSON tags to all structs
- Files: `backend/functions/list-debates/types.go`

**T203: Implement Firestore query logic** (45 min)
- [X] Create `queryDebates()` function accepting Firestore client, limit, offset
- [X] Build query: `Collection("debates").OrderBy("startedAt", firestore.Desc)`
- [X] Apply `Limit(limit)` and `Offset(offset)`
- [X] Execute query and iterate over documents
- [X] Map Firestore documents to `DebateSummary` structs
- [X] Implement separate count query to get total debate count
- [X] Return debates slice, total count, error
- Files: `backend/functions/list-debates/firestore.go`

**T204: Implement HTTP handler with validation** (30 min)
- [X] Parse query parameters `limit` (default 20, max 100) and `offset` (default 0)
- [X] Validate limit is between 1 and 100, return 400 if invalid
- [X] Validate offset is >= 0, return 400 if invalid
- [X] Initialize Firestore client with credentials from environment
- [X] Call `queryDebates()` with validated parameters
- [X] Handle Firestore errors and return 500 with error message
- [X] Calculate `hasMore = (offset + len(debates)) < total`
- [X] Return JSON response with debates, total, hasMore
- Files: `backend/functions/list-debates/handler.go`

**T205: Add CORS support** (15 min)
- [X] Add CORS headers to response: `Access-Control-Allow-Origin`, `Access-Control-Allow-Methods`
- [X] Read allowed origin from environment variable `ALLOWED_ORIGIN`
- [X] Handle preflight OPTIONS requests
- [X] Apply CORS middleware to handler
- Files: `backend/functions/list-debates/handler.go`

**T206: Create Cloud Functions entry point** (10 min)
- [X] Implement `ListDebates()` function with Cloud Functions HTTP signature
- [X] Register handler in `main.go`
- [X] Ensure function name matches GCP deployment requirements
- Files: `backend/functions/list-debates/main.go`

**T207: Create local development server** (10 min)
- [X] Implement HTTP server in `cmd/main.go` on port 8086
- [X] Route requests to handler function
- [X] Add logging for local debugging
- Files: `backend/functions/list-debates/cmd/main.go`

**T208: Create Dockerfile for list-debates** (15 min)
- [X] Create multi-stage Dockerfile (golang:1.24-alpine → distroless)
- [X] Set build context to `./backend/functions/list-debates`
- [X] Copy source and build binary from `./cmd/main.go`
- [X] Expose port 8080
- [X] Set entrypoint to run binary
- Files: `backend/functions/list-debates/Dockerfile`

**T209: Add list-debates to docker-compose** (10 min)
- [X] Add `list-debates` service to `docker-compose.yml`
- [X] Set build context and Dockerfile path
- [X] Map port 8086:8080
- [X] Add environment variables: `GOOGLE_APPLICATION_CREDENTIALS`, `ALLOWED_ORIGIN`
- [X] Mount secrets volume
- Files: `docker-compose.yml`

**T210: Test backend locally** (20 min)
- [ ] Start docker-compose with list-debates service
- [ ] Test `GET /list-debates` returns debates (if any exist)
- [ ] Test `GET /list-debates?limit=5` returns max 5 debates
- [ ] Test `GET /list-debates?limit=200` returns 400 error
- [ ] Test `GET /list-debates?offset=-1` returns 400 error
- [ ] Verify CORS headers are present
- Command: `curl http://localhost:8086/list-debates`

---

### Phase 2: Frontend - DebateHistory Page (3-4 hours)

**T211: Create DebateCard component** (30 min)
- [X] Create `frontend/src/components/DebateCard/DebateCard.jsx`
- [X] Accept props: debate (object), onClick (function)
- [X] Display topic (truncated to 100 chars with ellipsis)
- [X] Display panelists (show first 3 names + "+N more" if applicable)
- [X] Display relative timestamp using `formatRelativeTime()` helper
- [X] Add click handler to call `onClick`
- [X] Make card keyboard accessible (role, tabIndex, onKeyDown)
- Files: `frontend/src/components/DebateCard/DebateCard.jsx`

**T212: Style DebateCard** (20 min)
- [X] Create `frontend/src/components/DebateCard/DebateCard.module.css`
- [X] Card layout: padding, border, border-radius, shadow
- [X] Hover effect: shadow lift, border color change
- [X] Focus state: outline for keyboard navigation
- [X] Topic: font-size 18px, bold, margin-bottom
- [X] Panelists: color gray-600, margin-bottom
- [X] Timestamp: color gray-400, font-size 14px
- [X] Responsive: full width on mobile, fixed width on desktop
- Files: `frontend/src/components/DebateCard/DebateCard.module.css`

**T213: Create SearchInput component** (20 min)
- [X] Create `frontend/src/components/SearchInput/SearchInput.jsx`
- [X] Accept props: value, onChange, placeholder
- [X] Render input with controlled value
- [X] Add search icon (🔍 or SVG)
- [X] Add clear button (×) that appears when value is not empty
- [X] Clear button calls `onChange('')` to reset search
- [X] Add ARIA label for accessibility
- Files: `frontend/src/components/SearchInput/SearchInput.jsx`

**T214: Style SearchInput** (15 min)
- [X] Create `frontend/src/components/SearchInput/SearchInput.module.css`
- [X] Input container: flexbox, border, border-radius, padding
- [X] Search icon: positioned left inside input
- [X] Clear button: positioned right, cursor pointer, hover effect
- [X] Input field: flex-grow, no border, focus outline
- [X] Responsive: full width with padding on mobile
- Files: `frontend/src/components/SearchInput/SearchInput.module.css`

**T215: Update debateService with fetchDebateHistory** (15 min)
- [X] Open `frontend/src/services/debateService.js`
- [X] Add `fetchDebateHistory(limit = 20, offset = 0)` function
- [X] Build URL with query parameters using `REACT_APP_LIST_DEBATES_URL`
- [X] Fetch and parse JSON response
- [X] Handle errors and throw with descriptive message
- Files: `frontend/src/services/debateService.js`

**T216: Create helper functions** (15 min)
- [X] Create `frontend/src/utils/formatters.js`
- [X] Implement `formatPanelists(panelists)`: returns "Name1, Name2, Name3 +2 more"
- [X] Implement `formatRelativeTime(isoTimestamp)`: uses date-fns `formatDistanceToNow`
- [X] Implement `truncate(text, maxLength)`: adds "..." if text exceeds length
- Files: `frontend/src/utils/formatters.js`

**T217: Install date-fns dependency** (5 min)
- [X] Add `date-fns` to package.json
- Command: `cd frontend && npm install date-fns`
- Files: `frontend/package.json`

**T218: Create DebateHistory page** (45 min)
- [X] Create `frontend/src/pages/DebateHistory.jsx`
- [X] Add state: debates, searchQuery, isLoading, error
- [X] Add `useEffect` to fetch debates on mount via `fetchDebateHistory()`
- [X] Compute `filteredDebates` with `useMemo` based on searchQuery
- [X] Filter by topic and panelist names (case-insensitive)
- [X] Render SearchInput with searchQuery state
- [X] Render grid of DebateCards from filteredDebates
- [X] Handle empty state: "No debates yet. Create your first debate!"
- [X] Handle no results: "No debates found matching '{searchQuery}'"
- [X] Handle error state with error message
- [X] Add click handler to navigate to `/d/{id}` using `useNavigate`
- Files: `frontend/src/pages/DebateHistory.jsx`

**T219: Style DebateHistory page** (25 min)
- [X] Create `frontend/src/pages/DebateHistory.module.css`
- [X] Page container: max-width, margin auto, padding
- [X] Header: title "Debate History", back link
- [X] Search section: margin-bottom, centered
- [X] Results info: "Showing X of Y debates" below search
- [X] Cards grid: CSS Grid, 1 column mobile, 2-3 columns desktop, gap
- [X] Empty state: centered, icon, message, button
- [X] Loading spinner: centered
- Files: `frontend/src/pages/DebateHistory.module.css`

**T220: Add route for DebateHistory** (10 min)
- [X] Open `frontend/src/App.jsx`
- [X] Import DebateHistory component
- [X] Add route: `<Route path="/debates" element={<DebateHistory />} />`
- [X] Verify route is accessible
- Files: `frontend/src/App.jsx`

---

### Phase 3: Navigation Updates (1 hour)

**T221: Add "View Debate History" to Home page** (15 min)
- [ ] Open `frontend/src/pages/Home.jsx`
- [ ] Add Link/button to "/debates" below topic input section
- [ ] Text: "View Debate History" or "Browse Past Debates"
- [ ] Style as secondary action (not primary gradient)
- Files: `frontend/src/pages/Home.jsx`

**T222: Style history link on Home** (10 min)
- [ ] Open `frontend/src/pages/Home.module.css`
- [ ] Add class for history link: margin-top, text-align center
- [ ] Button style: border, padding, hover effect, no gradient
- Files: `frontend/src/pages/Home.module.css`

**T223: Add "Back to History" to DebateViewer** (15 min)
- [ ] Open `frontend/src/pages/DebateViewer.jsx`
- [ ] Add Link to "/debates" in header section
- [ ] Text: "← Back to History"
- [ ] Position next to "Export as PDF" button
- Files: `frontend/src/pages/DebateViewer.jsx`

**T224: Style back link on DebateViewer** (10 min)
- [ ] Open `frontend/src/pages/DebateViewer.module.css`
- [ ] Add class for back link: display inline, margin-right, hover underline
- Files: `frontend/src/pages/DebateViewer.module.css`

**T225: Add "Back to Home" to DebateHistory** (10 min)
- [ ] Already included in T218 (header with back link)
- [ ] Verify Link to "/" is present in header
- [ ] Text: "← Back to Home"
- Files: `frontend/src/pages/DebateHistory.jsx`

---

### Phase 4: Testing & Validation (2 hours)

**T226: Test backend endpoint** (20 min)
- [ ] Ensure at least 3 debates exist in Firestore (generate if needed)
- [ ] Test `GET /list-debates` returns all debates
- [ ] Test pagination with `limit=1&offset=0`, `limit=1&offset=1`
- [ ] Test invalid parameters return 400
- [ ] Test CORS headers are correct
- [ ] Verify response structure matches contract
- Command: `curl http://localhost:8086/list-debates?limit=10`

**T227: Test search filtering** (15 min)
- [ ] Navigate to `/debates` in browser
- [ ] Type in search box and verify instant filtering
- [ ] Search by topic keyword, verify results
- [ ] Search by panelist name, verify results
- [ ] Search for non-existent term, verify "No results" message
- [ ] Clear search, verify all debates reappear

**T228: Test navigation flow** (15 min)
- [ ] Start at Home, click "View Debate History", verify navigation to `/debates`
- [ ] Click a debate card, verify navigation to `/d/{id}`
- [ ] Click "Back to History", verify navigation to `/debates`
- [ ] Click "Back to Home", verify navigation to `/`

**T229: Test responsive design** (15 min)
- [ ] Open `/debates` in Chrome DevTools mobile view (375px)
- [ ] Verify search input is full width
- [ ] Verify cards stack in single column
- [ ] Test on tablet width (768px), verify 2-column grid
- [ ] Test on desktop (1200px), verify 3-column grid

**T230: Test empty state** (10 min)
- [ ] Clear all debates from Firestore (or use fresh database)
- [ ] Navigate to `/debates`
- [ ] Verify empty state message appears
- [ ] Click "Create your first debate" button (if added)
- [ ] Verify navigation to Home

**T231: Test accessibility** (15 min)
- [ ] Navigate to `/debates` with keyboard only (Tab key)
- [ ] Verify search input is focusable
- [ ] Verify debate cards are focusable with visible focus ring
- [ ] Press Enter on focused card, verify navigation
- [ ] Run axe DevTools accessibility scan, verify no violations

**T232: Test error handling** (10 min)
- [ ] Stop list-debates service to simulate backend failure
- [ ] Navigate to `/debates`
- [ ] Verify error message appears: "Failed to load debates"
- [ ] Restart service, refresh page, verify debates load

**T233: Performance testing** (10 min)
- [ ] Generate 50 debates in Firestore (or use seed script)
- [ ] Navigate to `/debates` and measure load time (<2s goal)
- [ ] Type in search box, measure filter time (<100ms goal)
- [ ] Click debate card, measure navigation time (<200ms goal)

---

### Phase 5: Deployment (30 min)

**T234: Update environment variables** (5 min)
- [ ] Add `REACT_APP_LIST_DEBATES_URL=http://localhost:8086` to `.env.development`
- [ ] Add `REACT_APP_LIST_DEBATES_URL=https://us-central1-PROJECT_ID.cloudfunctions.net/list-debates` to `.env.production`
- Files: `frontend/.env.development`, `frontend/.env.production`

**T235: Deploy list-debates to GCP** (10 min)
- [ ] Deploy function: `gcloud functions deploy list-debates --runtime go124 --trigger-http --allow-unauthenticated`
- [ ] Set environment variable: `ALLOWED_ORIGIN=https://raphink.github.io`
- [ ] Test deployed endpoint with curl
- Command: `curl https://us-central1-PROJECT_ID.cloudfunctions.net/list-debates`

**T236: Update docker-compose.yml** (5 min)
- [ ] Verify list-debates service configuration
- [ ] Test `docker-compose up` starts all services including list-debates
- Files: `docker-compose.yml`

**T237: Deploy frontend** (5 min)
- [ ] Build frontend: `npm run build`
- [ ] Deploy to GitHub Pages: `npm run deploy`
- [ ] Test production site at `https://raphink.github.io/debate/debates`

**T238: Final smoke test** (5 min)
- [ ] Navigate to production site
- [ ] Verify debate history loads
- [ ] Test search functionality
- [ ] Click debate card, verify viewer loads
- [ ] Test on mobile device

---

## Task Summary

**Total Tasks**: 38  
**Estimated Time**: 8-10 hours

**Breakdown by Phase**:
- Phase 1 (Backend): 10 tasks, 2-3 hours
- Phase 2 (Frontend): 10 tasks, 3-4 hours
- Phase 3 (Navigation): 5 tasks, 1 hour
- Phase 4 (Testing): 8 tasks, 2 hours
- Phase 5 (Deployment): 5 tasks, 30 min

**Parallel Tasks**:
- T211-T214 (Frontend components) can run in parallel with T201-T210 (Backend)
- T221-T225 (Navigation) can run independently after T218 completes

**Dependencies**:
- T215-T220 depend on T201-T210 (need backend endpoint)
- T226-T233 depend on all implementation tasks
- T234-T238 depend on all testing tasks

**Critical Path**:
T201→T202→T203→T204→T205→T206→T207→T208→T209→T210→T215→T218→T220→T226→T234→T235→T237
