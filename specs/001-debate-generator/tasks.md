# Tasks: AI-Powered Theology/Philosophy Debate Generator

**Feature Branch**: `001-debate-generator`  
**Input**: Design documents from `/specs/001-debate-generator/`  
**Prerequisites**: plan.md, spec.md, research.md, data-model.md, contracts/

## Format: `[ID] [P?] [Story?] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (e.g., US1, US2, US3)
- Include exact file paths in descriptions

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Project initialization and basic structure

- [X] T001 Create project directory structure per plan.md (backend/functions/, frontend/src/, shared/)
- [X] T002 Initialize Go modules for each Cloud Function (validate-topic, suggest-panelists, generate-debate)
- [X] T003 [P] Initialize React 18 project with package.json and dependencies (react-router, axios, jspdf, dompurify)
- [X] T004 [P] Create .gitignore for Node modules, Go binaries, and .env files
- [X] T005 [P] Setup ESLint configuration in frontend/.eslintrc.json per plan.md code quality standards
- [X] T006 [P] Setup golangci-lint configuration in backend/.golangci.yml
- [X] T007 Create environment variable template (.env.example) with ANTHROPIC_API_KEY and GCP_PROJECT_ID
- [X] T008 [P] Create README.md with quick start instructions per quickstart.md
- [X] T009 [P] Create docker-compose.yml to orchestrate all services (validate-topic:8080, generate-debate:8081, get-portrait:8082, get-debate:8084, frontend:3000)
- [X] T010 [P] Create multi-stage Dockerfiles for each backend Cloud Function (golang:1.24-alpine → distroless). Functions using shared module (generate-debate, get-debate) must build with context: ./backend
- [X] T011 [P] Create multi-stage Dockerfile for frontend (node:18-alpine → nginx:alpine)
- [X] T012 [P] Create .dockerignore to exclude node_modules, .env, build artifacts from Docker context
- [X] T013 [P] Create nginx.conf for frontend container (SPA routing, gzip, security headers)
- [X] T014 [P] Create start-local.sh script for one-command local development startup

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Core infrastructure that MUST be complete before ANY user story can be implemented

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T015 Create shared sanitization utility in backend/shared/sanitize/sanitizer.go (HTML tag removal, XSS prevention)
- [X] T016 [P] Create shared error handling utilities in backend/shared/errors/errors.go (user-friendly error types)
- [X] T017 [P] Create shared rate limiting utility in backend/shared/ratelimit/ratelimit.go (per Constitution Principle V)
- [X] T018 [P] Create shared API key management in backend/shared/auth/keys.go (GCP Secret Manager integration)
- [X] T019 [P] Setup React Router configuration in frontend/src/App.jsx (routes for Home, PanelistSelection, DebateGeneration)
- [X] T020 [P] Create constants file in frontend/src/utils/constants.js (MAX_PANELISTS=5, MIN_TOPIC_LENGTH=10, etc.)
- [X] T021 [P] Create DOMPurify wrapper in frontend/src/services/sanitizer.js (client-side XSS prevention)
- [X] T022 [P] Create Axios HTTP client configuration in frontend/src/services/api.js (base URLs, timeout settings)
- [X] T023 [P] Create common Button component in frontend/src/components/common/Button/Button.jsx (keyboard accessible)
- [X] T024 [P] Create common LoadingSpinner component in frontend/src/components/common/LoadingSpinner/LoadingSpinner.jsx
- [X] T025 [P] Create common ErrorMessage component in frontend/src/components/common/ErrorMessage/ErrorMessage.jsx
- [X] T026 [P] Create ErrorBoundary component in frontend/src/components/common/ErrorBoundary/ErrorBoundary.jsx
- [X] T027 [P] Configure CORS in all backend services (validate-topic, generate-debate, get-portrait) to use ALLOWED_ORIGIN environment variable (localhost for dev, raphink.github.io for prod)

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Topic Entry and Validation (Priority: P1) 🎯 MVP

**Goal**: User enters a debate topic and receives immediate feedback on whether it's suitable for theological/philosophical discussion

**Independent Test**: Enter various topics (theological, philosophical, off-topic) and verify validation responses appear within 3 seconds without needing panelist selection or debate generation

### Backend Implementation for User Story 1

- [X] T027 [P] [US1] Create Topic validation request/response structs in backend/functions/validate-topic/types.go
- [X] T027a [P] [US1] Add SuggestedNames []string field to TopicValidationRequest struct in types.go
- [X] T028 [P] [US1] Implement input validation in backend/functions/validate-topic/validator.go (10-500 chars, HTML stripping)
- [X] T029 [US1] Implement Claude API client in backend/functions/validate-topic/claude.go (topic relevance check)
- [X] T029a [US1] Update ValidateTopicAndSuggestPanelists in claude.go to accept suggestedNames parameter and include in prompt
- [X] T029b [US1] Migrate validate-topic to Anthropic Go SDK v1.19.0 for reliable streaming (replaces manual HTTP/SSE)
- [X] T029c [US1] Implement streamPanelistResponse in claude.go to parse streaming JSON and emit progressive panelist chunks
- [X] T029d [US1] Fix line detection in streamPanelistResponse to only trigger on newline (not closing brace)
- [X] T029e [US1] Add debug logging in claude.go and handler.go to track suggested names and streaming behavior
- [X] T029f [US1] Add fallback parser in streamPanelistResponse for old single-JSON-object format
- [X] T029g [US1] Strip markdown code blocks (```json...```) from Claude responses to handle formatting variations
- [X] T029h [US1] Add rejection format parser to handle {"type":"rejection","message":"..."} responses
- [X] T030 [US1] Implement HTTP handler in backend/functions/validate-topic/handler.go (CORS, error handling, sanitization)
- [X] T030a [US1] Update handler.go to sanitize and limit suggested names to 5 before passing to Claude client
- [X] T030b [US1] Update handler.go to support SSE streaming (set Content-Type: text/event-stream, flush chunks)
- [X] T031 [US1] Create main entry point in backend/functions/validate-topic/main.go (Cloud Function registration)
- [X] T032 [US1] Add unit tests for validator in backend/functions/validate-topic/validator_test.go
- [X] T033 [US1] Add integration tests for Claude API client in backend/functions/validate-topic/claude_test.go

### Frontend Implementation for User Story 1

- [X] T034 [P] [US1] Create TopicInput component in frontend/src/components/TopicInput/TopicInput.jsx (form, character counter)
- [X] T034a [P] [US1] Add suggested panelist names input field to TopicInput component (comma-separated, max 5, optional)
- [X] T034b [US1] Update TopicInput submit button text from "Validate Topic" to "Find Panelists"
- [X] T034c [US1] Implement chip-based panelist suggestions (type + comma/Enter creates chip with × remove button)
- [X] T034d [US1] Update TopicInput and Button components with gradient styling, hover animations, and shine effects
- [X] T034e [US1] Add Tab key support for creating panelist chips (in addition to comma and Enter)
- [X] T035 [P] [US1] Create TopicInput styles in frontend/src/components/TopicInput/TopicInput.module.css (mobile-first responsive)
- [X] T036 [P] [US1] Create ValidationResult component in frontend/src/components/ValidationResult/ValidationResult.jsx (success/error display)
- [X] T037 [P] [US1] Create ValidationResult styles in frontend/src/components/ValidationResult/ValidationResult.module.css
- [X] T038 [US1] Implement topicService in frontend/src/services/topicService.js (streaming SSE connection to validate-topic function)
- [X] T038a [US1] Update topicService to send suggested panelist names array to validate-topic endpoint
- [X] T038b [US1] Implement streaming chunk parsing in topicService (handle validation, panelist, error, done chunks)
- [X] T039 [US1] Create useTopicValidation custom hook in frontend/src/hooks/useTopicValidation.js (state management, streaming API)
- [X] T039a [US1] Update useTopicValidation hook to accept and pass suggestedNames parameter to topicService
- [X] T039b [US1] Update useTopicValidation hook loading message from "Validating..." to "Looking for Panelists"
- [X] T039c [US1] Add progressive panelist state management in useTopicValidation (append panelists as they stream in)
- [X] T040 [US1] Create Home page in frontend/src/pages/Home.jsx (integrate TopicInput, ValidationResult, navigation)
- [X] T040a [US1] Update Home.jsx to display panelists progressively as they stream in (show loading until first panelist)
- [X] T040b [US1] Update Home.jsx to hide input section when validation starts, show panelists with loading at bottom
- [X] T040c [US1] Update Home.jsx handleSubmit to receive and pass suggestedNames array to validate function
- [X] T040d [US1] Update Home.jsx to show PanelistSelector sidebar as soon as first panelist appears (not after streaming completes)
- [X] T041 [US1] Add client-side validation utilities in frontend/src/utils/validation.js (length check, sanitization)
- [X] T042 [US1] Add TopicInput component tests in frontend/src/components/TopicInput/TopicInput.test.jsx (Jest, RTL)
- [X] T043 [US1] Add ValidationResult component tests in frontend/src/components/ValidationResult/ValidationResult.test.jsx
- [X] T044 [US1] Add accessibility tests for topic validation flow in frontend/tests/accessibility/topic-validation.test.js (axe-core)
- [ ] T044a [US1] Enhance LoadingSpinner component with more engaging animation in frontend/src/components/common/LoadingSpinner/LoadingSpinner.jsx
- [X] T044b [US1] Fix Button component to use CSS Modules properly (import styles object, use styles.button classes)

**Checkpoint**: User Story 1 complete - users can validate topics independently

---

## Phase 4: User Story 2 - Panelist Discovery and Selection (Priority: P1)

**Goal**: User browses AI-suggested historical figures with known positions on the topic and selects up to 5 panelists for the debate

**Independent Test**: After validating a topic, view suggested panelist list with complete profiles (name, avatar, tagline, bio), select/deselect panelists, and verify selection limit (max 5) is enforced

**Architecture Note**: Panelists are returned by the validate-topic endpoint (see T029: ValidateTopicAndSuggestPanelists in claude.go) to reduce AI token costs and improve UX responsiveness. There is no separate suggest-panelists backend function. The frontend receives panelist data in the TopicValidationResponse.suggestedPanelists field.

### Frontend Implementation for User Story 2

- [X] T045 [P] [US2] Create PanelistCard component in frontend/src/components/PanelistGrid/PanelistCard.jsx (avatar, name, tagline, bio)
- [X] T046 [P] [US2] Create PanelistCard styles in frontend/src/components/PanelistGrid/PanelistCard.module.css (selected state, hover effects)
- [X] T047 [P] [US2] Create PanelistGrid component in frontend/src/components/PanelistGrid/PanelistGrid.jsx (grid layout, keyboard navigation)
- [X] T048 [P] [US2] Create PanelistGrid styles in frontend/src/components/PanelistGrid/PanelistGrid.module.css (responsive grid)
- [X] T049 [P] [US2] Create PanelistSelector component in frontend/src/components/PanelistSelector/PanelistSelector.jsx (selection counter, clear button)
- [X] T050 [P] [US2] Create PanelistSelector styles in frontend/src/components/PanelistSelector/PanelistSelector.module.css
- [X] T051 [US2] Create usePanelistSelection custom hook in frontend/src/hooks/usePanelistSelection.js (selection state, 2-5 limit enforcement)
- [X] T052 [US2] Create PanelistSelection page in frontend/src/pages/PanelistSelection.jsx (integrate Grid, Selector, navigation)
- [X] T053 [US2] Add avatar images to frontend/public/avatars/ directory (placeholder-avatar.png and sample panelist avatars)
- [ ] T054 [US2] Add PanelistCard component tests in frontend/src/components/PanelistGrid/PanelistCard.test.jsx
- [ ] T055 [US2] Add PanelistGrid component tests in frontend/src/components/PanelistGrid/PanelistGrid.test.jsx
- [ ] T056 [US2] Add usePanelistSelection hook tests in frontend/src/hooks/usePanelistSelection.test.js (selection limit logic)
- [ ] T057 [US2] Add accessibility tests for panelist selection flow in frontend/tests/accessibility/panelist-selection.test.js (keyboard navigation, screen reader labels)

**Checkpoint**: User Story 2 complete - users can browse and select panelists independently

---

## Phase 4b: Panelist Portrait Enhancement (Async)

**Goal**: Progressively enhance panelist avatars with real portraits from Wikimedia Commons

**Independent Test**: After panelists stream in with placeholders, verify portraits load asynchronously and update the UI, or fall back to placeholders gracefully

### Backend Implementation for Portrait Service

- [X] T054a [P] [US2] Create portrait request/response structs in backend/functions/get-portrait/types.go (panelistId, panelistName)
- [X] T054b [P] [US2] Implement Wikimedia Commons API client in backend/functions/get-portrait/wikimedia.go (fetch 300px thumbnails with proper User-Agent)
- [X] T054c [P] [US2] Implement in-memory cache in backend/functions/get-portrait/cache.go (thread-safe map for portrait URLs)
- [X] T054d [US2] Implement HTTP handler in backend/functions/get-portrait/handler.go (validate input, fetch/cache portrait, return URL)
- [X] T054e [US2] Create main entry point in backend/functions/get-portrait/main.go (Cloud Function registration)
- [X] T054f [US2] Create local dev binary in backend/functions/get-portrait/cmd/main.go (HTTP server for local testing)
- [ ] T054g [US2] Add unit tests for Wikimedia API client in backend/functions/get-portrait/wikimedia_test.go
- [ ] T054h [US2] Add cache tests in backend/functions/get-portrait/cache_test.go (concurrent access, TTL)

### Frontend Implementation for Portrait Service

- [X] T054i [P] [US2] Create portraitService in frontend/src/services/portraitService.js (async fetch portrait URLs)
- [X] T054j [US2] Update useTopicValidation hook to fetch portraits when panelists arrive and update avatarUrl state
- [ ] T054k [US2] Add loading shimmer effect to avatars while portraits are being fetched
- [ ] T054l [US2] Ensure portraits are cached in React state to avoid redundant fetches during debate generation
- [X] T054m [US2] Fix PanelistCard.jsx avatar URL handling to check for absolute URLs (http/https prefix) before prepending PUBLIC_URL/avatars/ path, matching pattern from DebateBubble, PanelistModal, and PanelistSelector components

**Checkpoint**: Panelist avatars progressively enhanced with real portraits

---

## Phase 5: User Story 3 - Live Debate Generation with Streaming Display (Priority: P1) 🎯 MVP

**Goal**: User launches debate generation and watches the conversation unfold in real-time as a chat-style interface with panelist avatars

**Independent Test**: After selecting panelists, launch debate and verify responses stream progressively into chat bubbles with correct avatar attribution, loading indicators, and error handling

### Backend Implementation for User Story 3

- [X] T058 [P] [US3] Create DebateConfiguration request struct in backend/functions/generate-debate/types.go (topic, panelists array)
- [X] T059 [P] [US3] Create StreamChunk response structs in backend/functions/generate-debate/types.go (message, error, done events)
- [X] T060 [P] [US3] Implement input validation in backend/functions/generate-debate/validator.go (2-5 panelists, valid topic)
- [X] T061 [US3] Implement Claude API streaming client in backend/functions/generate-debate/claude.go (SSE with debate prompt)
- [X] T061a [US3] Update debate prompt in claude.go to ensure moderator provides concluding summary at end of debate
- [X] T061b [US3] Migrate generate-debate to Anthropic Go SDK v1.19.0 for reliable streaming (replaces manual HTTP/SSE)
- [X] T062 [US3] Implement streaming proxy in claude.go streamResponse function (character-by-character forwarding with pattern buffering)
- [X] T062a [US3] Fix UTF-8 handling in streamResponse (use runes not bytes, WriteRune not WriteByte)
- [X] T063 [US3] Implement HTTP handler in backend/functions/generate-debate/handler.go (SSE headers, flush chunks, error recovery)
- [X] T064 [US3] Create main entry point in backend/functions/generate-debate/main.go (Cloud Function registration)
- [ ] T065 [US3] Add unit tests for debate prompt construction in backend/functions/generate-debate/claude_test.go
- [ ] T066 [US3] Add integration tests for SSE streaming in backend/functions/generate-debate/stream_test.go

### Frontend Implementation for User Story 3

- [X] T067 [P] [US3] Create DebateBubble component in frontend/src/components/DebateView/DebateBubble.jsx (chat bubble with avatar)
- [X] T068 [P] [US3] Create DebateBubble styles in frontend/src/components/DebateView/DebateBubble.module.css (left/right alignment per panelist)
- [X] T069 [P] [US3] Create TypingIndicator component in frontend/src/components/DebateView/TypingIndicator.jsx (animated dots)
- [X] T070 [P] [US3] Create TypingIndicator styles in frontend/src/components/DebateView/TypingIndicator.module.css
- [X] T071 [P] [US3] Create DebateView component in frontend/src/components/DebateView/DebateView.jsx (scrollable message list)
- [X] T072 [P] [US3] Create DebateView styles in frontend/src/components/DebateView/DebateView.module.css (chat interface styling)
- [X] T073 [US3] Implement debateService in frontend/src/services/debateService.js (fetch API SSE connection)
- [X] T074 [US3] Create useDebateStream custom hook in frontend/src/hooks/useDebateStream.js (SSE state, message accumulation, error handling)
- [X] T074a [US3] Simplify useDebateStream after backend streaming redesign (speaker detection handled by backend)
- [X] T075 [US3] Create DebateGeneration page in frontend/src/pages/DebateGeneration.jsx (integrate DebateView, generate button, retry logic)
- [X] T076 [US3] Add message accumulation logic in useDebateStream hook (append chunks to correct message by panelistId)
- [X] T077 [US3] Add auto-scroll toggle control in DebateView component (checkbox, disabled by default, conditional scrollIntoView)
- [X] T077a [P] [US3] Create PanelistModal component in frontend/src/components/DebateView/PanelistModal.jsx (display name, tagline, bio, close controls)
- [X] T077b [P] [US3] Create PanelistModal styles in frontend/src/components/DebateView/PanelistModal.module.css (overlay, centered modal, accessible focus trap)
- [X] T077c [US3] Add clickable avatar handler in DebateBubble component (onClick opens modal with panelist data)
- [ ] T078 [US3] Add DebateBubble component tests in frontend/src/components/DebateView/DebateBubble.test.jsx
- [ ] T079 [US3] Add DebateView component tests in frontend/src/components/DebateView/DebateView.test.jsx
- [ ] T080 [US3] Add useDebateStream hook tests in frontend/src/hooks/useDebateStream.test.js (SSE event handling, error states)
- [ ] T081 [US3] Add accessibility tests for debate streaming flow in frontend/tests/accessibility/debate-generation.test.js (screen reader announcements, focus management)

**Checkpoint**: User Story 3 complete - users can generate and watch debates stream in real-time

---

## Phase 6: User Story 4 - PDF Export (Priority: P2)

**Goal**: User exports completed debate as a formatted PDF document for offline reading, sharing, or archival purposes

**Independent Test**: After generating a complete debate, click export button and verify PDF downloads within 2 seconds with all content (topic, panelists, messages, timestamp) properly formatted

### Frontend Implementation for User Story 4

- [X] T082 [P] [US4] Create pdfGenerator utility in frontend/src/components/PDFExport/pdfGenerator.js (jsPDF integration, debate formatting)
- [X] T083 [P] [US4] Create PDFExport component in frontend/src/components/PDFExport/PDFExport.jsx (export button, download trigger)
- [X] T084 [US4] Add PDF header generation in pdfGenerator.js (debate topic, timestamp, page numbers)
- [X] T085 [US4] Add panelist profile section in pdfGenerator.js (circular portrait avatars, names, bios)
- [X] T086 [US4] Add debate conversation rendering in pdfGenerator.js (chat bubbles with circular portrait avatars, proper page breaks)
- [X] T086a [US4] Implement image loading utility in pdfGenerator.js (fetch and convert portrait URLs to base64 data URLs for PDF embedding)
- [X] T086b [US4] Add CORS proxy support for Wikimedia portrait URLs in PDF generation (handle cross-origin image loading)
- [X] T086c [US4] Implement circular avatar cropping in PDF using jsPDF ellipse clipping path
- [X] T087 [US4] Integrate PDFExport component in DebateGeneration page (show after debate completes)
- [X] T088 [US4] Add error handling for PDF generation failures in PDFExport component
- [ ] T089 [US4] Add pdfGenerator unit tests in frontend/src/components/PDFExport/pdfGenerator.test.js (content formatting, page breaks)
- [ ] T090 [US4] Add PDFExport component tests in frontend/src/components/PDFExport/PDFExport.test.jsx

**Checkpoint**: User Story 4 complete - users can export debates as PDFs

---

## Phase 6.5: User Story 5 - Debate Sharing and Caching (Priority: P2)

**Goal**: User shares completed debate via URL that loads from backend-cached storage (Firestore), allowing debates to be revisited and shared without regeneration

**Independent Test**: Generate a debate, verify URL updates to /d/{uuid}, copy URL, open in new browser/incognito, verify debate loads from backend with identical content

### Backend: Firestore Integration

- [ ] T107 [P] [US5] Add Firebase Admin SDK dependency to backend/go.mod (cloud.google.com/go/firestore, firebase.google.com/go)
- [ ] T108 [P] [US5] Create Firestore client in backend/shared/firebase/client.go (initialize with Application Default Credentials, requires GCP_PROJECT_ID env var)
- [ ] T108a [P] [US5] Update quickstart.md with ADC setup: `gcloud auth application-default login` for local dev, service account key for Docker
- [ ] T109 [P] [US5] Create debate storage service in backend/shared/firebase/debates.go (SaveDebate, GetDebate with DebateDocument struct)
- [ ] T110 [P] [US5] Add UUID generation to backend/functions/generate-debate/handler.go using github.com/google/uuid
- [ ] T111 [P] [US5] Modify generate-debate to include X-Debate-Id header in SSE response with generated UUID
- [ ] T111a [P] [US5] Add Access-Control-Expose-Headers: X-Debate-Id to CORS configuration in generate-debate handler
- [ ] T112 [P] [US5] Modify generate-debate to accumulate messages during streaming and save to Firestore on completion (non-blocking)
- [ ] T113 [P] [US5] Add error handling for Firestore save failures (log error, don't fail debate stream)

### Backend: Get Debate Function

- [ ] T114 [P] [US5] Create new Cloud Function backend/functions/get-debate/ (HTTP GET handler)
- [ ] T115 [P] [US5] Implement get-debate handler in main.go (parse UUID from query param, validate format, query Firestore)
- [ ] T116 [P] [US5] Add error responses for get-debate (404 Not Found, 400 Bad Request, 500 Internal Error)
- [ ] T117 [P] [US5] Add CORS headers to get-debate response for cross-origin requests
- [ ] T118 [P] [US5] Create Dockerfile for get-debate function (multi-stage build, minimal runtime)
- [ ] T119 [P] [US5] Add get-debate deployment configuration to deploy.sh script

### Firestore Security

- [ ] T120 [P] [US5] Create Firestore database: `gcloud firestore databases create --database="(default)" --location=europe-west1`
- [ ] T121 [P] [US5] Create firestore.rules with deny all direct client access (read/write: false)
- [ ] T122 [P] [US5] Create .firebaserc with Firebase project ID configuration
- [ ] T123 [P] [US5] Create firebase.json with Firestore rules deployment configuration
- [ ] T124 [P] [US5] Update DEPLOYMENT.md with Firebase project setup and security rules deployment instructions

### Frontend: API Integration

- [ ] T124 [P] [US5] Add getDebateById method to frontend/src/services/api.js (GET /api/get-debate?id={uuid})
- [ ] T125 [P] [US5] Update useDebateStream hook to extract X-Debate-Id header from SSE response
- [ ] T126 [P] [US5] Update useDebateStream to update browser URL to /d/{uuid} using History API (pushState, no page reload)

### Frontend: Debate Viewer Page

- [ ] T127 [P] [US5] Create useDebateLoader hook in frontend/src/hooks/useDebateLoader.js (fetch debate from backend by UUID parameter)
- [ ] T128 [P] [US5] Create DebateViewer page in frontend/src/pages/DebateViewer.jsx (load and display cached debate, handle loading/error states)
- [ ] T129 [P] [US5] Add /d/:uuid route in App.jsx routing to DebateViewer component
- [ ] T130 [P] [US5] Add "Debate not found" error state in DebateViewer for 404s with link to create new debate
- [ ] T131 [P] [US5] Add retry button for 500 errors in DebateViewer

### Frontend: Share Functionality

- [ ] T132 [P] [US5] Create ShareButton component in frontend/src/components/DebateView/ShareButton.jsx (copy current URL to clipboard)
- [ ] T133 [P] [US5] Add ShareButton to DebateView component with success/failure toast notifications
- [ ] T134 [P] [US5] Style ShareButton with gradient and hover effects matching app design system
- [ ] T135 [P] [US5] Show ShareButton only when debate ID is available (hide during initial generation before UUID received)

### Testing

- [ ] T136 [US5] Backend test: Generate debate → verify Firestore document created → verify document structure matches DebateDocument
- [ ] T137 [US5] Backend test: Call get-debate with valid UUID → verify JSON response matches saved debate
- [ ] T138 [US5] Backend test: Call get-debate with invalid UUID → verify 400 Bad Request response
- [ ] T139 [US5] Backend test: Call get-debate with non-existent UUID → verify 404 Not Found response
- [ ] T140 [US5] Frontend test: Generate debate → verify URL updates to /d/{uuid} → verify share button appears
- [ ] T141 [US5] End-to-end test: Generate debate → copy share URL → open in new browser/incognito → verify identical content loads from backend
- [ ] T142 [US5] Test Firestore save failure handling (graceful degradation, debate still viewable/exportable, just not shareable)
- [ ] T143 [US5] Test share button copy-to-clipboard functionality with success notification

**Checkpoint**: User Story 5 complete - users can share debates via URLs and backend serves cached debates

---

## Phase 11: Recent Debates Discovery (User Story 6)

**Goal**: Display subtle list of recent debates on home page for discovery and quick access

**Independent Test**: Generate multiple debates, return to home page, verify recent debates list displays with topic and avatars

### Backend: List Debates Endpoint

- [ ] T144 [P] [US6] Create list-debates function in backend/shared/firebase/debates.go (query Firestore, order by createdAt desc, limit param)
- [ ] T145 [P] [US6] Create new Cloud Function backend/functions/list-debates/ (HTTP GET handler)
- [ ] T146 [P] [US6] Implement list-debates handler (parse limit from query param, default 10, max 20)
- [ ] T147 [P] [US6] Return minimal debate data: {id, topic, panelistAvatars[], createdAt}
- [ ] T148 [P] [US6] Add CORS headers to list-debates response
- [ ] T149 [P] [US6] Create Dockerfile for list-debates function (multi-stage build, minimal runtime)
- [ ] T150 [P] [US6] Add list-debates deployment to deploy.sh script

### Frontend: Recent Debates Component

- [ ] T151 [P] [US6] Add listRecentDebates method to frontend/src/services/api.js (GET /api/list-debates?limit=10)
- [ ] T152 [P] [US6] Create useRecentDebates hook in frontend/src/hooks/useRecentDebates.js (fetch recent debates, loading/error states)
- [ ] T153 [P] [US6] Create RecentDebates component in frontend/src/components/RecentDebates/RecentDebates.jsx (subtle list UI)
- [ ] T154 [P] [US6] Create DebateListItem component (topic truncated to 60 chars, circular panelist avatars, onClick navigate)
- [ ] T155 [P] [US6] Add RecentDebates to Home page below topic input
- [ ] T156 [P] [US6] Style RecentDebates as subtle suggestion UI (muted colors, not prominent)
- [ ] T157 [P] [US6] Hide RecentDebates gracefully if empty or API fails (no error shown)

### Testing

- [ ] T158 [US6] Backend test: Call list-debates → verify returns debates ordered by createdAt descending
- [ ] T159 [US6] Backend test: Call list-debates with limit=5 → verify returns max 5 debates
- [ ] T160 [US6] Frontend test: Load home page → verify recent debates list renders if debates exist
- [ ] T161 [US6] Frontend test: Click recent debate item → verify navigates to /d/{uuid}
- [ ] T162 [US6] Frontend test: Topic truncation works correctly (60 char limit with ellipsis)

**Checkpoint**: User Story 6 complete - users can discover and access recent debates from home page

### User Story 7 - Panelist Autocomplete (P3)

**Backend - Autocomplete API Function** (Depends: Firestore integration US5)
- [ ] T163 [US7] Backend: Create autocomplete-panelists Cloud Function scaffolding
- [ ] T164 [US7] Backend: Implement Firestore query to aggregate all panelists from debates collection
- [ ] T165 [US7] Backend: Implement name normalization utility (lowercase, strip titles/punctuation)
- [ ] T166 [US7] Backend: Implement fuzzy matching algorithm for panelist name deduplication
- [ ] T167 [US7] Backend: Implement frequency counting for panelist occurrences across debates
- [ ] T168 [US7] Backend: Implement query matching logic (case-insensitive substring/prefix match)
- [ ] T169 [US7] Backend: Implement response ranking (most frequent first, limit to top 10)
- [ ] T170 [US7] Backend: Add 5-minute in-memory cache for aggregated panelist data
- [ ] T171 [US7] Backend: Create autocomplete-panelists contract JSON schema
- [ ] T172 [US7] Backend: Deploy autocomplete-panelists to Cloud Functions (go124 runtime)

**Frontend - Autocomplete Component** (Depends: T163-T172)
- [ ] T173 [US7] Frontend: Create usePanelistAutocomplete.js hook with debouncing (300ms)
- [ ] T174 [US7] Frontend: Create autocomplete API service in topicService.js (GET /api/autocomplete-panelists?q={query})
- [ ] T175 [US7] Frontend: Update chip input component to support autocomplete dropdown
- [ ] T176 [US7] Frontend: Implement autocomplete dropdown UI (Material-UI Autocomplete or custom)
- [ ] T177 [US7] Frontend: Add loading indicator for slow autocomplete responses (>500ms)
- [ ] T178 [US7] Frontend: Implement graceful degradation when autocomplete API fails
- [ ] T179 [US7] Frontend: Add keyboard navigation for autocomplete dropdown (↑↓ arrows, Enter, Escape)
- [ ] T180 [US7] Frontend: Update PanelistSelection.jsx to integrate autocomplete component
- [ ] T181 [US7] Frontend: Add analytics tracking for autocomplete usage (selected vs manual entry)

**Testing User Story 7** (Depends: T173-T181)
- [ ] T182 [US7] Backend test: Verify name normalization ("St. Augustine" → "augustine")
- [ ] T183 [US7] Backend test: Verify fuzzy matching deduplicates similar names
- [ ] T184 [US7] Backend test: Verify frequency ranking returns most common panelists first
- [ ] T185 [US7] Backend test: Verify autocomplete returns max 10 results
- [ ] T186 [US7] Backend test: Verify cache reduces Firestore reads (5-minute TTL)
- [ ] T187 [US7] Frontend test: Type "aug" → verify autocomplete suggests "Augustine of Hippo"
- [ ] T188 [US7] Frontend test: Select autocomplete suggestion → verify chip created with correct data
- [ ] T189 [US7] Frontend test: Autocomplete API fails → verify manual chip creation still works
- [ ] T190 [US7] Frontend test: Type query with no matches → verify dropdown hides gracefully
- [ ] T191 [US7] E2E test: Generate 3 debates with "Augustine" → verify autocomplete suggests him first

**Checkpoint**: User Story 7 complete - users receive intelligent panelist suggestions from historical data

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final touches, performance optimization, and comprehensive quality assurance

- [ ] T132 [P] Add loading states to all async operations (topic validation, panelist loading, debate generation)
- [ ] T133 [P] Add error retry mechanisms with exponential backoff for all API calls
- [ ] T134 [P] Optimize React component re-renders with React.memo for chat bubbles
- [ ] T135 [P] Add session storage persistence for topic and panelist selection (prevent loss on refresh)
- [ ] T136 [P] Add analytics/logging for user flow completion rates (optional, privacy-preserving)
- [ ] T137 [P] Create 404 NotFound page in frontend/src/pages/NotFound.jsx
- [ ] T138 [P] Add global CSS styles in frontend/src/index.css (typography, color scheme, responsive breakpoints)
- [X] T098 [P] Add PWA manifest and meta tags in frontend/public/ (mobile installation support)
- [X] T098a [P] Create manifest.json with app name, description, icons, theme colors, and standalone display mode
- [X] T098b [P] Create app icons in multiple sizes (192x192, 512x512) for iOS and Android
- [X] T098c [P] Add manifest link and apple-touch-icon meta tags to index.html
- [X] T098d [P] Create markdown utility in frontend/src/utils/markdown.js for parsing inline formatting
- [X] T098e [P] Update DebateBubble component to render Markdown formatting with dangerouslySetInnerHTML
- [X] T098f [P] Update PDF generator to render Markdown formatting with appropriate font styles (bold, italic, bold-italic)
- [ ] T139 Perform full accessibility audit with axe-core across all pages (verify WCAG 2.1 Level AA)
- [ ] T140 Perform performance audit with Lighthouse (verify SC-007: <100ms UI response)
- [ ] T141 Test complete user journey end-to-end (topic → validation → panelist selection → debate → PDF export → share)
- [ ] T142 Add comprehensive error scenarios testing (API failures, network timeouts, invalid responses)
- [ ] T143 Review all user-facing text for clarity and Constitution Principle I compliance
- [ ] T144 Security review of all API key handling and XSS prevention measures
- [ ] T145 Setup GitHub Actions CI workflow in .github/workflows/frontend-ci.yml (lint, test, build)
- [ ] T146 Setup GitHub Actions CI workflow in .github/workflows/backend-ci.yml (Go tests, linting)

---

## Dependencies

### Story Completion Order

User Stories 1, 2, and 3 are mostly independent and can be developed in parallel after Phase 2 (Foundational) is complete:

```
Phase 1 (Setup)
    ↓
Phase 2 (Foundational) ← BLOCKING: Must complete before any user story
    ↓
    ├─→ Phase 3 (US1: Topic Validation) ← Independent
    ├─→ Phase 4 (US2: Panelist Selection) ← Depends on US1 for navigation flow
    └─→ Phase 5 (US3: Debate Generation) ← Depends on US1 + US2 for inputs
            ↓
        Phase 6 (US4: PDF Export) ← Depends on US3 for debate data
            ↓
        Phase 7 (Polish)
```

### Parallel Execution Opportunities

**After Phase 2 completes, these tasks can run in parallel:**

**User Story 1 Backend** (T021-T027):
- All backend tasks for validate-topic function can run independently

**User Story 1 Frontend** (T028-T038):
- All frontend components and tests can run in parallel with backend

**User Story 2 Backend** (T039-T045):
- Can start after Phase 2, runs parallel to US1

**User Story 2 Frontend** (T046-T059):
- Can start after Phase 2, runs parallel to US1

**User Story 3**: Must wait for US1 and US2 navigation flow, but backend (T060-T068) and frontend (T069-T083) can be parallel

**User Story 4**: Depends only on US3 completion, all tasks (T084-T092) can run in parallel

**Phase 7**: Many polish tasks (T093-T108) can run in parallel

---

## Implementation Strategy

### MVP Scope (Minimum Viable Product)

The MVP should include **User Stories 1, 2, and 3 only** (Phases 1-5):

- Topic validation
- Panelist selection
- Live debate streaming

This provides complete core functionality. User Story 4 (PDF export) can be added in a subsequent release.

### Incremental Delivery Plan

1. **Sprint 1** (Weeks 1-2): Phase 1 (Setup) + Phase 2 (Foundational)
2. **Sprint 2** (Weeks 3-4): Phase 3 (US1: Topic Validation) - Deliver working topic validation
3. **Sprint 3** (Weeks 5-6): Phase 4 (US2: Panelist Selection) - Deliver working panelist selection
4. **Sprint 4** (Weeks 7-9): Phase 5 (US3: Debate Streaming) - Deliver working debate generation (MVP COMPLETE)
5. **Sprint 5** (Week 10): Phase 6 (US4: PDF Export) - Add export functionality
6. **Sprint 6** (Week 11): Phase 7 (Polish) - Final QA and optimization

### Testing Discipline

- Unit tests written alongside implementation (same sprint)
- Integration tests run after each user story phase completes
- Accessibility tests run at end of each user story implementation
- Full E2E test suite run before phase 7
- No user story marked "complete" until all tests pass

### Quality Gates

Each phase must pass these gates before proceeding:

- **Phase 1**: Project builds successfully, linting passes
- **Phase 2**: All shared utilities have unit tests, security review passes
- **Phase 3**: US1 acceptance scenarios all pass, accessibility audit passes for topic flow
- **Phase 4**: US2 acceptance scenarios all pass, panelist selection limit enforced
- **Phase 5**: US3 acceptance scenarios all pass, streaming performance meets SC-003/SC-004
- **Phase 6**: PDF export meets SC-006 (2s for 5000 words)
- **Phase 7**: All constitution principles verified, Lighthouse score >90

---

## Task Statistics

- **Total Tasks**: 154 (includes portrait service, enhanced PDF export, PWA support, Markdown rendering, and backend-managed Firestore storage)
- **Setup Phase**: 14 tasks (T001-T014)
- **Foundational Phase**: 13 tasks (T015-T027, BLOCKING - includes CORS configuration)
- **User Story 1**: 18 tasks (T027-T044: 7 backend + 11 frontend)
- **User Story 2**: 22 tasks (T045-T057 + T054a-T054m: 13 frontend + 9 portrait service)
  - Original panelist selection: 13 tasks (T045-T057: frontend only, panelists come from validate-topic)
  - Portrait enhancement: 9 tasks (T054a-T054m: 5 backend + 4 frontend, async Wikimedia integration)
- **User Story 3**: 24 tasks (T058-T081: 9 backend + 15 frontend)
- **User Story 4**: 12 tasks (T082-T090 + T086a-T086c: frontend only, includes enhanced PDF with portraits)
  - Original PDF export: 9 tasks (T082-T090)
  - Portrait embedding: 3 tasks (T086a-T086c: image loading, CORS, circular cropping)
- **User Story 5**: 37 tasks (T107-T143: backend-managed Firestore integration, UUID-based shareable URLs)
  - Backend Firestore integration: 7 tasks (T107-T113: Firebase Admin SDK, UUID generation, persistence)
  - Backend get-debate function: 6 tasks (T114-T119: new HTTP endpoint for debate retrieval)
  - Firestore security: 4 tasks (T120-T123: rules, configuration, deployment docs)
  - Frontend API integration: 3 tasks (T124-T126: backend API client, URL updates)
  - Frontend debate viewer: 5 tasks (T127-T131: loading, error handling, routing)
  - Frontend share functionality: 4 tasks (T132-T135: copy to clipboard, notifications)
  - Testing: 8 tasks (T136-T143: backend + frontend + end-to-end)
- **Polish Phase**: 22 tasks (T144-T158 + T098a-T098f: includes PWA manifest and Markdown rendering)
  - Original polish: 16 tasks (renumbered T144-T158 from T132-T146)
  - PWA support: 3 tasks (T098a-T098c: manifest, icons, meta tags)
  - Markdown rendering: 3 tasks (T098d-T098f: utility, web UI, PDF export)

**Parallel Opportunities**: ~70% of tasks can run in parallel after foundational phase completes

**MVP Tasks**: 130 tasks (Phases 1-6.5 including portrait service, enhanced PDF, PWA, Markdown, and backend-managed Firestore storage, excludes remaining Polish)

**Architecture Notes**: 
- Portrait service (get-portrait) runs as independent Cloud Function with async frontend integration
- Firestore operations managed entirely by backend (generate-debate saves, get-debate retrieves) using Firebase Admin SDK
- Frontend has NO direct Firestore access - all operations via backend API endpoints for security and control
- All backend services use ALLOWED_ORIGIN environment variable for CORS security
- Frontend avatar components check for absolute URLs before prepending local path prefix
- PDF export uses async image loading with CORS-enabled fetch, converts portraits to base64 data URLs for embedding
- Firestore security rules deny all direct client access (read/write: false), enforcing API-only pattern
- Chat bubble format in PDF matches web UI with circular portrait avatars and rounded rectangles
- PWA manifest enables mobile installation with standalone display mode (iOS Safari 14+, Android Chrome 90+)
- Icon generation automated via generate-icons.sh script using librsvg or ImageMagick
- Inline Markdown formatting (*italic*, **bold**, ***bold italic***) rendered in web UI via dangerouslySetInnerHTML (HTML-escaped)
- PDF export renders Markdown with jsPDF font styles (normal, bold, italic, bolditalic)
- Firestore integration: UUID v4 generated at debate start using Web Crypto API, complete debate saved after generation
- Shareable URLs: /d/{uuid} pattern loads debates from Firestore cache, public reads with no client writes
- Graceful degradation: Firestore save failures don't block viewing/exporting debates
