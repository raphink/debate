# Implementation Plan: AI-Powered Theology/Philosophy Debate Generator

**Branch**: `001-debate-generator` | **Date**: 2025-12-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-debate-generator/spec.md`

## Clarifications (Session 2025-12-13)

- **Debate Completion**: Determined by token/word count threshold (~5000 words generated). Moderator provides concluding summary when this threshold is reached.
- **Autocomplete Features**: US6 (Topic Discovery) and US7 (Panelist Autocompletion) are deferred post-MVP. Current scope includes only basic debate history browsing (list-debates).
- **Firestore Security**: Completely locked down - no direct client access. All reads AND writes happen exclusively via backend Cloud Functions (get-debate, list-debates, generate-debate).
- **Portrait Fallback**: Standardized on placeholder-avatar.svg (SVG format) for all missing portraits across all components.

## Summary

Build a web application that generates AI-powered theological/philosophical debates between historical figures. Users enter a topic and click "Find Panelists" to trigger validation with an engaging "Looking for Panelists" loading animation. Users can optionally suggest up to 5 panelist names they'd like considered (with PRIORITY weighting - included unless fictional/non-existent). The system validates the topic and streams panelists progressively via Server-Sent Events (SSE) - Claude returns each panelist as a complete JSON line, and the backend emits them immediately as they're detected (through GCP Cloud Functions proxy in Go using official Anthropic SDK) for true progressive loading. Users select 2-5 panelists, then watch the debate stream in real-time as a chat-style conversation with avatars. A neutral moderator introduces the debate, may intervene between panelist exchanges to redirect/clarify/summarize, and MUST provide a concluding summary at the end. Completed debates can be exported as PDF.

**Key Architecture Decision**: The validate-topic endpoint uses character-by-character streaming to detect complete JSON lines as Claude generates them. Response format changed from single JSON object to line-by-line format: `{"type":"rejection","message":"..."}` OR multiple `{"type":"panelist","data":{...}}` lines. This eliminates the validation/panelist race condition and provides true incremental streaming. The backend strips markdown code blocks (```json...```) from responses to handle Claude's formatting variations. Both backend services (validate-topic and generate-debate) use the official Anthropic Go SDK (v1.19.0) for reliable streaming.

**Panelist Portrait Service**: A separate async endpoint (get-portrait) fetches real portrait images from Wikimedia Commons API after panelists stream in. The service runs as an independent Cloud Function (port 8082/8083) to keep validation streaming fast and non-blocking.

**Implementation Details**:
- **Wikimedia API Client** (wikimedia.go): Queries Wikipedia API with action=query, prop=pageimages, pithumbsize=300 to fetch 300px thumbnails suitable for 48x48px circular display
- **User-Agent Requirement**: All Wikimedia API requests MUST include User-Agent header "DebateApp/1.0 (https://github.com/raphink/debate; debate@example.com)" to avoid 403 Forbidden errors (Wikimedia policy)
- **In-Memory Cache** (cache.go): Thread-safe map using sync.RWMutex to cache portrait URLs per session, preventing redundant API calls when same panelists appear in debate generation
- **Fallback Strategy**: Returns empty string on Wikimedia failure, frontend falls back to placeholder-avatar.svg (SVG format standardized)
- **Frontend Integration**: portraitService.js calls get-portrait endpoint asynchronously when panelists arrive via useTopicValidation hook, updates avatarUrl in state progressively
- **Canonical Placeholder**: All avatar display components (PanelistCard, DebateBubble, PanelistModal, PanelistSelector) use placeholder-avatar.svg when avatarUrl is empty, null, or fails to load
- **URL Handling Fix**: All avatar display components (PanelistCard, DebateBubble, PanelistModal, PanelistSelector) check if avatarUrl starts with 'http' or '/' before prepending PUBLIC_URL/avatars/ path, enabling both absolute Wikimedia URLs and relative local paths to work correctly
- **CORS Security**: get-portrait handler uses ALLOWED_ORIGIN environment variable (http://localhost:3000 for dev, https://raphink.github.io for prod) matching validate-topic and generate-debate services
- **Local Development**: cmd/main.go provides standalone HTTP server for testing portrait service locally on port 8083

**User-Suggested Panelists**: Treated as PRIORITY requests - included unless clearly invalid (fictional, non-existent, or completely unrelated to intellectual discourse). Claude infers positions from known works/tradition even if they never directly addressed the specific topic.

**Suggested Names Feature**: Users can optionally propose up to 5 panelist names during topic entry via a chip-based input (type name, press comma, Tab, or Enter to create chip with × remove button). The backend sanitizes these names and includes them in the Claude API prompt. Claude evaluates whether the suggested individuals have known, documented positions on the topic and includes them in the panelist list if appropriate. This gives users more control while maintaining quality through AI validation.

**Progressive Loading UX**: When "Find Panelists" is clicked, the input section hides to focus attention on the streaming panelist results. Panelists appear one by one as Claude generates them, with the loading animation positioned at the bottom of the list (not blocking view). Once streaming completes, the panelist selector sidebar appears with beautifully styled gradient buttons.

**Button Design**: All action buttons feature gradient backgrounds (purple for primary actions, gray for secondary), smooth hover animations with lift effects, shine sweeps across on hover, enhanced shadows, and improved focus states for accessibility.

**Moderator Behavior**: The moderator is responsible for introducing the debate, optionally intervening between panelist exchanges (to ask clarifying questions, highlight contrasts, or summarize progress), and providing a concluding summary that synthesizes the key points when word count reaches approximately 5000 words or when arguments are naturally exhausted.

**Debate Completion**: Backend monitors word count during generation. When approximately 5000 words are generated, Claude is signaled to have the moderator provide concluding summary and end the debate gracefully.

**Future Enhancement**: User-as-moderator functionality is out of scope for MVP but documented as a potential future feature where users could interactively ask questions during the debate.

## Technical Context

**Language/Version**: Go 1.24 (backend/GCP Functions), JavaScript/React 18+ (frontend)  
**Primary Dependencies**: 
- Frontend: React 18, React Router, CSS Modules
- Backend: Go 1.24, Google Cloud Functions SDK, Anthropic Go SDK v1.19.0 (official streaming client)
- DevOps: Docker 24+, Docker Compose v2+ (local development)
**Storage**: N/A (stateless, no persistence for MVP)  
**Testing**: 
- Frontend: Jest, React Testing Library, axe-core (accessibility)
- Backend: Go testing package, httptest
**Local Development**: 
- Docker Compose orchestration for all services
- Hot-reload for frontend development
- Independent containerized Cloud Functions
**Target Platform**: 
- Frontend: Modern browsers (Chrome 90+, Firefox 88+, Safari 14+, Edge 90+)
- Backend: Google Cloud Functions (Gen 2), serverless runtime
- Mobile: PWA installable on iOS (Safari 14+) and Android (Chrome 90+) with standalone display mode
- Storage: Cloud Firestore (NoSQL document database) for debate caching
**Project Type**: Web application (frontend + backend)  
**Performance Goals**: 
- Topic validation response: <3s
- First debate response: <5s
- Streaming chunk intervals: <500ms
- UI interaction response: <100ms
- PDF generation: <2s for 5000 words (with portrait embedding)
- Firestore save: <1s (non-blocking, happens after debate completion)
- Firestore read: <2s (debate loading from shareable URL)
**PDF Export Strategy**:
- Client-side generation using jsPDF library
- Chat bubble format matching web UI styling
- Circular portrait images embedded from avatarUrl (Wikimedia URLs and local avatars)
- Portrait images converted to base64 data URLs for embedding
- CORS proxy for cross-origin Wikimedia images
- Automatic page breaks between messages to avoid splitting bubbles
**Firestore Storage Strategy**:
- UUID v4 generated at debate start using Web Crypto API
- Debate saved to Firestore after generation completes
- Document structure: {id, topic, panelists, messages, status, timestamps, metadata}
- **Security Model**: Complete lockdown - NO direct client access (read or write). All Firestore operations exclusively via backend Cloud Functions (get-debate, list-debates, generate-debate)
- Backend uses Application Default Credentials (ADC) or service account for Firestore access
- Graceful degradation: Firestore save failures don't block viewing/export
- Average document size: ~25 KB (topic + 3 panelists + 15 messages)
- Free tier capacity: ~40,000 debates (1 GB storage)
**Constraints**: 
- No user authentication (debates are public via UUID URLs)
- No database persistence for user sessions (stateless MVP, but debates cached in Firestore)
- Claude API rate limits (per Anthropic tier)
- GCP Cloud Functions timeout: 60s max per request
- Client-side PDF generation to avoid server overhead
- Mobile-first responsive design (≥375px width)
- PWA manifest for mobile installation (no service worker for MVP - online-only)
- Firestore free tier limits: 50K reads/day, 20K writes/day, 10 GB egress/month
**Scale/Scope**: 
- MVP: Single-user sessions, no concurrent debate limit
- Expected load: <100 concurrent users initially
- Debate length: ~10-20 exchanges (~5000 words typical, completion threshold)
- Frontend: ~15-20 components, 5-8 pages/views
- Backend: 5 Cloud Functions (validate-topic, generate-debate, get-portrait, get-debate, list-debates)
- **Deferred Features**: US6 (topic autocomplete), US7 (panelist autocomplete) - post-MVP enhancements

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

### ✅ Principle I: User-Centric Design (UX First)
- **Status**: PASS
- **Evidence**: Spec defines 4 prioritized user stories (P1/P2) with independent test criteria; each story delivers standalone value
- **Validation**: Clear acceptance scenarios with measurable outcomes; user feedback built into streaming UX

### ✅ Principle II: Code Quality & Maintainability  
- **Status**: PASS
- **Evidence**: Single Responsibility Principle enforced through component-based React architecture; Go functions focused on single API operations
- **Validation**: Linting required (ESLint for React, golangci-lint for Go); dependencies explicitly versioned in package.json and go.mod

### ✅ Principle III: Responsive & Accessible UI
- **Status**: PASS
- **Evidence**: FR-019 requires keyboard navigation; FR-020 mandates 4.5:1 contrast ratio; SC-008 requires WCAG 2.1 Level AA compliance
- **Validation**: Mobile-first design (≥375px); axe-core automated testing; semantic HTML with ARIA labels

### ✅ Principle IV: Interactive & Performant Experience
- **Status**: PASS
- **Evidence**: SC-001 (3s validation), SC-003 (5s first response), SC-004 (<500ms streaming), SC-007 (<100ms UI response)
- **Validation**: Loading indicators for all async operations; progressive streaming display; retry mechanisms for failures

### ✅ Principle V: AI Safety & Security (NON-NEGOTIABLE)
- **Status**: PASS
- **Evidence**: FR-016 (XSS prevention via output sanitization), FR-017 (rate limiting), FR-018 (input validation)
- **Validation**: 
  - DOMPurify for sanitizing Claude outputs before render
  - Input validation on both client and GCP function
  - API keys stored in GCP Secret Manager (never in code)
  - HTTPS-only communication with Claude API
  - No PII logging in Cloud Functions

### Constitution Compliance Summary
**All 5 principles: PASS** - No violations. Design aligns with UX-first approach, maintains code quality standards, ensures accessibility, meets performance targets, and implements comprehensive AI security controls.

## Project Structure

### Documentation (this feature)

```text
specs/001-debate-generator/
├── plan.md              # This file (/speckit.plan command output)
├── research.md          # Phase 0 output (/speckit.plan command)
├── data-model.md        # Phase 1 output (/speckit.plan command)
├── quickstart.md        # Phase 1 output (/speckit.plan command)
├── contracts/           # Phase 1 output (/speckit.plan command)
│   ├── validate-topic.json      # Topic validation API contract
│   ├── suggest-panelists.json   # Panelist suggestion API contract
│   └── generate-debate.json     # Debate generation API contract
└── tasks.md             # Phase 2 output (/speckit.tasks command - NOT created by /speckit.plan)
```

### Source Code (repository root)

```text
backend/
├── functions/
│   ├── validate-topic/          # GCP Cloud Function: Topic validation + panelist suggestions
│   │   ├── main.go
│   │   ├── handler.go
│   │   ├── claude.go            # Claude API client with streaming
│   │   ├── validator.go         # Input validation
│   │   ├── types.go
│   │   └── go.mod
│   ├── get-portrait/            # GCP Cloud Function: Async portrait fetching from Wikimedia Commons
│   │   ├── main.go              # Cloud Functions init (functions.HTTP)
│   │   ├── cmd/
│   │   │   └── main.go          # Local dev HTTP server (port 8083)
│   │   ├── handler.go           # HTTP handler with CORS and validation
│   │   ├── wikimedia.go         # Wikimedia Commons API client
│   │   ├── cache.go             # Thread-safe in-memory URL cache
│   │   ├── types.go             # Request/response structs
│   │   ├── Dockerfile           # Multi-stage build from ./cmd
│   │   └── go.mod
│   └── generate-debate/         # GCP Cloud Function: Debate generation with streaming
│       ├── main.go
│       ├── handler.go
│       ├── claude.go
│       ├── stream.go            # Server-Sent Events (SSE) streaming
│       └── go.mod
├── shared/                      # Shared utilities across functions
│   ├── auth/                    # API key management
│   ├── sanitize/                # Input sanitization
│   ├── ratelimit/               # Rate limiting logic
│   └── errors/                  # Error handling utilities
└── tests/
    ├── integration/
    └── unit/

frontend/
├── public/
│   ├── index.html
│   ├── manifest.json            # PWA manifest for mobile installation
│   └── avatars/                 # Historical figure avatars
├── src/
│   ├── components/
│   │   ├── TopicInput/          # US1: Topic entry form
│   │   │   ├── TopicInput.jsx
│   │   │   ├── TopicInput.test.jsx
│   │   │   └── TopicInput.module.css
│   │   ├── ValidationResult/    # US1: Validation feedback
│   │   │   ├── ValidationResult.jsx
│   │   │   └── ValidationResult.module.css
│   │   ├── PanelistGrid/        # US2: Panelist display grid
│   │   │   ├── PanelistGrid.jsx
│   │   │   ├── PanelistCard.jsx
│   │   │   └── PanelistGrid.module.css
│   │   ├── PanelistSelector/    # US2: Selection management
│   │   │   ├── PanelistSelector.jsx
│   │   │   └── PanelistSelector.module.css
│   │   ├── DebateView/          # US3: Chat-style debate display
│   │   │   ├── DebateView.jsx
│   │   │   ├── DebateBubble.jsx
│   │   │   ├── TypingIndicator.jsx
│   │   │   ├── ShareButton.jsx  # US5: Share debate URL
│   │   │   └── DebateView.module.css
│   │   ├── PDFExport/           # US4: PDF generation
│   │   │   ├── PDFExport.jsx
│   │   │   └── pdfGenerator.js
│   │   └── common/
│   │       ├── Button/
│   │       ├── ErrorBoundary/
│   │       ├── LoadingSpinner/
│   │       └── ErrorMessage/
│   ├── services/
│   │   ├── api.js               # Axios HTTP client configuration
│   │   ├── topicService.js      # Topic validation API calls
│   │   ├── panelistService.js   # Panelist suggestion API calls
│   │   ├── portraitService.js   # Async portrait URL fetching from get-portrait
│   │   ├── debateService.js     # Debate generation API calls (SSE)
│   │   ├── firestoreService.js  # Firestore read/write operations
│   │   └── sanitizer.js         # DOMPurify wrapper for XSS prevention
│   ├── hooks/
│   │   ├── useDebateStream.js   # Custom hook for SSE streaming + Firestore save
│   │   ├── usePanelistSelection.js
│   │   ├── useTopicValidation.js
│   │   ├── useDebateLoader.js   # Load debate from Firestore by UUID
│   │   ├── useTopicAutocomplete.js  # US6: Autocomplete topics from history
│   │   └── usePanelistAutocomplete.js  # US7: Autocomplete panelist suggestions
│   ├── pages/
│   │   ├── Home.jsx             # Topic entry with integrated autocomplete (US6)
│   │   ├── PanelistSelection.jsx # Panelist selection with "Modify" button for pre-filled (US6)
│   │   ├── DebateGeneration.jsx # Live debate generation (/d/:uuid)
│   │   ├── DebateViewer.jsx     # US5: Load cached debate from Firestore
│   │   └── NotFound.jsx
│   ├── utils/
│   │   ├── validation.js        # Client-side input validation
│   │   ├── cacheDetection.js    # US6: Detect cache hit (topic + panelist matching)
│   │   ├── constants.js         # App constants (max panelists, etc.)
│   │   ├── uuid.js              # UUID v4 generation using Web Crypto API
│   │   ├── markdown.js          # Markdown parsing for *italic*, **bold**, etc.
│   │   └── accessibility.js     # A11y utilities
│   ├── firebase.js              # Firebase SDK initialization
│   ├── App.jsx
│   ├── App.test.jsx
│   ├── index.jsx
│   └── index.css
├── tests/
│   ├── integration/
│   ├── accessibility/           # axe-core tests
│   └── e2e/
├── package.json
├── .eslintrc.json
├── Dockerfile                   # Production frontend container
├── nginx.conf                   # Nginx config for production
└── README.md

# Configuration files
firestore.rules                  # Firestore security rules (public reads, no client writes)
.firebaserc                      # Firebase project configuration
firebase.json                    # Firebase deployment config

# Docker/DevOps files
docker-compose.yml               # Local development orchestration
.dockerignore                    # Docker build exclusions
start-local.sh                   # Quick start script

.github/
└── workflows/
    ├── frontend-ci.yml          # Frontend lint, test, build
    ├── backend-ci.yml           # Backend Go tests
    └── deploy.yml               # Deploy to GCP + Firestore rules
```

**Structure Decision**: Web application architecture selected due to separate frontend (React SPA) and backend (GCP Cloud Functions). Frontend handles all UI/UX concerns including streaming display, PDF export, and Firestore integration. Backend provides three focused functions acting as a secure proxy to Claude API, implementing rate limiting and input validation. Firestore provides debate caching for shareable URLs without requiring user authentication or complex session management. Each function is independently deployable.

---

## Firestore Integration Details

### Backend-Managed Persistence

**Architecture Decision**: Backend (Cloud Functions) manages all Firestore operations. Frontend never directly accesses Firestore - all reads/writes go through backend API endpoints. This provides:
- Centralized access control
- Prevents client-side spam/abuse
- Complete security lockdown (Firestore rules prevent all direct client access)
- Better security (Application Default Credentials or service account only in backend)
- Consistent error handling

**Configuration Requirements**:
- `GCP_PROJECT_ID` environment variable must be set for Firestore initialization
- Firestore uses the `(default)` database in the specified GCP project
- Firestore security rules deployed to prevent all direct client read/write access

**Backend Endpoints for Firestore**:
- `get-debate?id={uuid}` - Retrieve single debate by UUID
- `list-debates?limit={n}&offset={m}` - Paginated debate history browsing
- `generate-debate` - Save completed debate after generation (automatic)
- **Firestore database must be created before first use**: `gcloud firestore databases create --database="(default)" --location=europe-west1 --project=${GCP_PROJECT_ID}`

**Authentication Setup**:
- **Local Development**: Run `gcloud auth application-default login` to create ADC
- **Docker Compose**: Credentials automatically mounted from `~/.config/gcloud/application_default_credentials.json` to container
- **Production**: Use GCP service account attached to Cloud Run/Cloud Functions (automatic ADC)
- **Service Account Permissions**: Requires `Cloud Datastore User` role (or `roles/datastore.user`) for Firestore access

### UUID Generation
- **Location**: Backend (`generate-debate` Cloud Function)
- **Library**: Google UUID (`github.com/google/uuid`)
- **Format**: UUID v4 (128-bit, cryptographically random)
- **Generation point**: When backend receives debate generation request
- **Delivery**: Returned in HTTP header `X-Debate-Id: {uuid}`
- **CORS Exposure**: Custom header exposed via `Access-Control-Expose-Headers: X-Debate-Id` (required for JavaScript access)
- **URL pattern**: `/d/{uuid}` (e.g., `/d/550e8400-e29b-41d4-a716-446655440000`)
- **Uniqueness**: ~10^36 combinations, collision probability negligible

### Data Flow

#### Generation Flow
```
1. Frontend POST /api/generate-debate {topic, panelists}
   → Backend generates UUID
   → Backend starts SSE stream with header: X-Debate-Id: {uuid}
   
2. Frontend extracts UUID from response headers
   → Updates browser URL to /d/{uuid} (History API, no page reload)
   → Continues receiving SSE messages
   
3. Backend accumulates messages during streaming
   → Builds complete debate object in memory
   → On final "done" event, writes to Firestore: debates/{uuid}
   → Non-blocking write (debate success independent of Firestore)
   
4. If Firestore write fails
   → Backend logs error
   → Debate stream continues normally
   → User can still view/export, just can't share
```

#### Retrieval Flow
```
1. User visits /d/{uuid} (shared URL)
   → Frontend calls GET /api/get-debate?id={uuid}
   → Backend reads Firestore: debates/{uuid}
   
2. If document exists
   → Backend returns JSON with complete debate data
   → Frontend renders DebateView (no streaming needed)
   
3. If document not found
   → Backend returns 404 Not Found
   → Frontend shows "Debate not found" message
   
4. If Firestore error
   → Backend returns 500 Internal Error
   → Frontend shows retry button
```

### Data Model

**Firestore Collection**: `debates`

**Document Structure**:
```javascript
{
  // Document ID is the UUID
  id: "550e8400-e29b-41d4-a716-446655440000",
  
  topic: {
    text: "Should Christians defy authorities when the law is unfair?",
    suggestedNames: ["Martin Luther King Jr."],
    isRelevant: true,
    validationMessage: "This topic is well-suited..."
  },
  
  panelists: [
    {
      id: "Augustine354",
      name: "Augustine of Hippo",
      tagline: "4th-5th century theologian...",
      biography: "Early Christian theologian...",
      avatarUrl: "/avatars/Augustine354-avatar.png",
      position: "Would argue that Christians..."
    }
    // ... 1-4 more panelists
  ],
  
  messages: [
    {
      id: "moderator-0",
      panelistId: "moderator",
      panelistName: "Moderator",
      avatarUrl: "/avatars/moderator-avatar.png",
      text: "Welcome to today's debate...",
      timestamp: "2025-12-12T10:30:00Z",
      sequence: 0,
      isComplete: true
    }
    // ... 10-20 more messages
  ],
  
  status: "complete",
  startedAt: "2025-12-12T10:30:00Z",
  completedAt: "2025-12-12T10:32:30Z",
  
  metadata: {
    createdBy: "anonymous",
    userAgent: request.Header.Get("User-Agent"),
    version: "1.0",
    generatedBy: "backend"
  }
}
```

**Size Estimates**:
- Small debate (2 panelists, 10 messages): ~15 KB
- Average debate (3 panelists, 15 messages): ~20 KB
- Large debate (5 panelists, 25 messages): ~35 KB

### Backend Implementation

**New Cloud Function**: `get-debate`
```
Endpoint: GET /api/get-debate?id={uuid}
Purpose: Retrieve saved debate from Firestore
Response: 200 OK with debate JSON, 404 Not Found, or 500 Error
```

**New Cloud Function**: `autocomplete-topics` (US6)
```
Endpoint: GET /api/autocomplete-topics?q={query}&limit=10
Purpose: Search historical debate topics for autocomplete in topic input field
Response: 200 OK with array of matching debates (id, topic, panelistCount, panelists[], startedAt)
Implementation:
  1. Query Firestore debates collection where topic contains query substring (case-insensitive)
  2. Order by startedAt DESC (newest first)
  3. Return up to 'limit' results (default 10)
  4. Include full debate metadata: id, topic text, panelist IDs/names/slugs, created timestamp
  5. Frontend uses this to pre-fill panelists and detect cache hits
Notes: Index on topic field for performance, consider caching recent queries
```

**New Cloud Function**: `autocomplete-panelists` (US7)
```
Endpoint: GET /api/autocomplete-panelists?q={query}
Purpose: Return autocomplete suggestions based on historical panelist data
Response: 200 OK with array of panelist suggestions ranked by frequency
Implementation:
  1. Aggregate all panelists from Firestore debates
  2. Normalize names (lowercase, strip titles/punctuation)
  3. Fuzzy match query against normalized names
  4. Deduplicate similar panelists (e.g., "Augustine" vs "St. Augustine")
  5. Return top 10 matches sorted by frequency (most common first)
  6. Return canonical panelist data from most frequent variant
Notes: Cache aggregation results for 5 minutes to reduce Firestore reads
```

**Modified Cloud Function**: `generate-debate`
```
Changes:
1. Generate UUID at start
2. Include X-Debate-Id header in SSE response
3. Accumulate messages during streaming
4. Write to Firestore on completion
5. Non-blocking save (don't fail debate if Firestore down)
```

**Firebase Admin SDK** (`backend/shared/firebase/`):
```go
// client.go
package firebase

import (
    "context"
    "cloud.google.com/go/firestore"
    firebase "firebase.google.com/go"
)

var client *firestore.Client

func InitFirestore(ctx context.Context) error {
    app, err := firebase.NewApp(ctx, nil)
    if err != nil {
        return err
    }
    
    client, err = app.Firestore(ctx)
    return err
}

func GetClient() *firestore.Client {
    return client
}
```

**Debate Storage** (`backend/shared/firebase/debates.go`):
```go
type DebateDocument struct {
    ID          string    `firestore:"id"`
    Topic       Topic     `firestore:"topic"`
    Panelists   []Panelist `firestore:"panelists"`
    Messages    []Message `firestore:"messages"`
    Status      string    `firestore:"status"`
    StartedAt   time.Time `firestore:"startedAt"`
    CompletedAt time.Time `firestore:"completedAt"`
    Metadata    Metadata  `firestore:"metadata"`
}

func SaveDebate(ctx context.Context, uuid string, debate DebateDocument) error {
    _, err := GetClient().Collection("debates").Doc(uuid).Set(ctx, debate)
    return err
}

func GetDebate(ctx context.Context, uuid string) (*DebateDocument, error) {
    doc, err := GetClient().Collection("debates").Doc(uuid).Get(ctx)
    if err != nil {
        return nil, err
    }
    
    var debate DebateDocument
    if err := doc.DataTo(&debate); err != nil {
        return nil, err
    }
    return &debate, nil
}
```

### Security Rules

**File**: `firestore.rules`

```javascript
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /debates/{debateId} {
      // NO direct client access - all operations via backend API
      allow read, write: if false;
    }
  }
}
```

**Security Notes**:
- Frontend NEVER accesses Firestore directly
- All reads go through `get-debate` API endpoint
- All writes handled by `generate-debate` backend
- Backend uses Firebase Admin SDK (full access via service account)
- No Firebase SDK in frontend (reduced bundle size)
- No sensitive data stored (historical debates are educational content)
- Client cannot write directly (prevents spam/abuse)
- Cloud Functions use service account for authenticated writes
- UUID obscurity provides practical privacy (128-bit unguessable)

### Frontend Changes

**New API Client Methods** (`src/services/api.js`):
```javascript
// Fetch saved debate by UUID
export const getDebateById = async (uuid) => {
  const response = await fetch(`${API_BASE_URL}/get-debate?id=${uuid}`);
  if (!response.ok) {
    if (response.status === 404) {
      throw new Error('Debate not found');
    }
    throw new Error('Failed to load debate');
  }
  return response.json();
};
```

**New Route** (`src/App.jsx`):
```jsx
<Route path="/d/:uuid" element={<DebateViewer />} />
```

**New Component**: `src/pages/DebateViewer.jsx`
- Extracts UUID from URL params
- Calls `getDebateById(uuid)` API
- Shows loading spinner while fetching
- Renders `DebateView` with loaded debate data
- Displays error message if not found or failed

**Modified Hook**: `src/hooks/useDebateStream.js`
- Extract `X-Debate-Id` header from SSE response
- Store UUID in state
- Call `history.pushState()` to update URL to `/d/{uuid}`
- No Firestore interaction (backend handles it)

**New Component**: Share button in `DebateView`
- Copy current URL to clipboard
- Toast notification: "Link copied!"
- Only shown when UUID is available

### Environment Configuration

**Backend** (Cloud Functions deployment):
- Uses Application Default Credentials (ADC)
- No explicit service account key needed
- Firestore automatically initialized from GCP project

**Frontend** (no Firebase SDK):
- No Firebase configuration needed
- No environment variables for Firestore
- Reduced bundle size

**Firestore Rules Deployment**:
```bash
firebase deploy --only firestore:rules
```

### Cost Optimization

**Free Tier Capacity** (per FIRESTORE_PRICING.md):
- Storage: 1 GB = ~40,000 debates
- Reads: 50K/day = 1.5M/month (via backend API)
- Writes: 20K/day = 600K/month (backend only)
- Egress: 10 GB/month = ~400K views

**Expected Usage**:
- MVP: <100 debates/month, <1000 views/month
- Cost: $0/month (well within free tier)

**Future Optimizations** (if needed):
- TTL policy: Auto-delete debates >90 days old
- Compression: Gzip message text (~60% size reduction)
- Deduplication: Hash topic+panelist IDs to detect duplicates

### Error Handling

**Backend Firestore Write Failures** (during generation):
- Non-blocking: SSE stream continues normally
- Logged to Cloud Logging for monitoring
- Debate still viewable/exportable, just not shareable
- No user-facing error message (graceful degradation)

**Backend Firestore Read Failures** (when loading `/d/{uuid}`):
- Return 500 Internal Server Error status
- Include error message in response body
- Frontend shows friendly error with retry button
- Logged to Cloud Logging for debugging

**Frontend Network Failures**:
- Fetch timeout: 10 seconds
- Show loading spinner during fetch
- Display error message with retry button
- Log to console for debugging

**404 Not Found** (debate doesn't exist):
- Backend returns 404 status
- Frontend shows "Debate not found or expired"
- Provide button to create new debate
- Log UUID for investigation
