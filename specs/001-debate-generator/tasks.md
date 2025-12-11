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
- [X] T009 [P] Create docker-compose.yml to orchestrate all services (validate-topic:8080, suggest-panelists:8081, generate-debate:8082, frontend:3000)
- [X] T010 [P] Create multi-stage Dockerfiles for each backend Cloud Function (golang:1.23-alpine → distroless)
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

**Checkpoint**: Foundation ready - user story implementation can now begin in parallel

---

## Phase 3: User Story 1 - Topic Entry and Validation (Priority: P1) 🎯 MVP

**Goal**: User enters a debate topic and receives immediate feedback on whether it's suitable for theological/philosophical discussion

**Independent Test**: Enter various topics (theological, philosophical, off-topic) and verify validation responses appear within 3 seconds without needing panelist selection or debate generation

### Backend Implementation for User Story 1

- [ ] T027 [P] [US1] Create Topic validation request/response structs in backend/functions/validate-topic/types.go
- [ ] T028 [P] [US1] Implement input validation in backend/functions/validate-topic/validator.go (10-500 chars, HTML stripping)
- [ ] T029 [US1] Implement Claude API client in backend/functions/validate-topic/claude.go (topic relevance check)
- [ ] T030 [US1] Implement HTTP handler in backend/functions/validate-topic/handler.go (CORS, error handling, sanitization)
- [ ] T031 [US1] Create main entry point in backend/functions/validate-topic/main.go (Cloud Function registration)
- [ ] T032 [US1] Add unit tests for validator in backend/functions/validate-topic/validator_test.go
- [ ] T033 [US1] Add integration tests for Claude API client in backend/functions/validate-topic/claude_test.go

### Frontend Implementation for User Story 1

- [ ] T034 [P] [US1] Create TopicInput component in frontend/src/components/TopicInput/TopicInput.jsx (form, character counter)
- [ ] T035 [P] [US1] Create TopicInput styles in frontend/src/components/TopicInput/TopicInput.module.css (mobile-first responsive)
- [ ] T036 [P] [US1] Create ValidationResult component in frontend/src/components/ValidationResult/ValidationResult.jsx (success/error display)
- [ ] T037 [P] [US1] Create ValidationResult styles in frontend/src/components/ValidationResult/ValidationResult.module.css
- [ ] T038 [US1] Implement topicService in frontend/src/services/topicService.js (API call to validate-topic function)
- [ ] T039 [US1] Create useTopicValidation custom hook in frontend/src/hooks/useTopicValidation.js (state management, API call)
- [ ] T040 [US1] Create Home page in frontend/src/pages/Home.jsx (integrate TopicInput, ValidationResult, navigation)
- [ ] T041 [US1] Add client-side validation utilities in frontend/src/utils/validation.js (length check, sanitization)
- [ ] T042 [US1] Add TopicInput component tests in frontend/src/components/TopicInput/TopicInput.test.jsx (Jest, RTL)
- [ ] T043 [US1] Add ValidationResult component tests in frontend/src/components/ValidationResult/ValidationResult.test.jsx
- [ ] T044 [US1] Add accessibility tests for topic validation flow in frontend/tests/accessibility/topic-validation.test.js (axe-core)

**Checkpoint**: User Story 1 complete - users can validate topics independently

---

## Phase 4: User Story 2 - Panelist Discovery and Selection (Priority: P1)

**Goal**: User browses AI-suggested historical figures with known positions on the topic and selects up to 5 panelists for the debate

**Independent Test**: After validating a topic, view suggested panelist list with complete profiles (name, avatar, tagline, bio), select/deselect panelists, and verify selection limit (max 5) is enforced

### Backend Implementation for User Story 2

- [ ] T045 [P] [US2] Create Panelist struct in backend/functions/suggest-panelists/panelist.go (id, name, tagline, bio, avatarUrl, position)
- [ ] T046 [P] [US2] Implement Claude API client in backend/functions/suggest-panelists/claude.go (panelist suggestion with topic context)
- [ ] T047 [US2] Implement HTTP handler in backend/functions/suggest-panelists/handler.go (parse request, call Claude, sanitize responses)
- [ ] T048 [US2] Create main entry point in backend/functions/suggest-panelists/main.go (Cloud Function registration)
- [ ] T049 [US2] Add validation for Panelist data in backend/functions/suggest-panelists/validator.go (alphanumeric ID, no HTML)
- [ ] T050 [US2] Add unit tests for Panelist validation in backend/functions/suggest-panelists/validator_test.go
- [ ] T045 [US2] Add integration tests for Claude API client in backend/functions/suggest-panelists/claude_test.go

### Frontend Implementation for User Story 2

- [ ] T046 [P] [US2] Create PanelistCard component in frontend/src/components/PanelistGrid/PanelistCard.jsx (avatar, name, tagline, bio)
- [ ] T047 [P] [US2] Create PanelistCard styles in frontend/src/components/PanelistGrid/PanelistCard.module.css (selected state, hover effects)
- [ ] T048 [P] [US2] Create PanelistGrid component in frontend/src/components/PanelistGrid/PanelistGrid.jsx (grid layout, keyboard navigation)
- [ ] T049 [P] [US2] Create PanelistGrid styles in frontend/src/components/PanelistGrid/PanelistGrid.module.css (responsive grid)
- [ ] T050 [P] [US2] Create PanelistSelector component in frontend/src/components/PanelistSelector/PanelistSelector.jsx (selection counter, clear button)
- [ ] T051 [P] [US2] Create PanelistSelector styles in frontend/src/components/PanelistSelector/PanelistSelector.module.css
- [ ] T052 [US2] Implement panelistService in frontend/src/services/panelistService.js (API call to suggest-panelists function)
- [ ] T053 [US2] Create usePanelistSelection custom hook in frontend/src/hooks/usePanelistSelection.js (selection state, 2-5 limit enforcement)
- [ ] T054 [US2] Create PanelistSelection page in frontend/src/pages/PanelistSelection.jsx (integrate Grid, Selector, navigation)
- [ ] T055 [US2] Add avatar images to frontend/public/avatars/ directory (placeholder-avatar.png and sample panelist avatars)
- [ ] T056 [US2] Add PanelistCard component tests in frontend/src/components/PanelistGrid/PanelistCard.test.jsx
- [ ] T057 [US2] Add PanelistGrid component tests in frontend/src/components/PanelistGrid/PanelistGrid.test.jsx
- [ ] T058 [US2] Add usePanelistSelection hook tests in frontend/src/hooks/usePanelistSelection.test.js (selection limit logic)
- [ ] T059 [US2] Add accessibility tests for panelist selection flow in frontend/tests/accessibility/panelist-selection.test.js (keyboard navigation, screen reader labels)

**Checkpoint**: User Story 2 complete - users can browse and select panelists independently

---

## Phase 5: User Story 3 - Live Debate Generation with Streaming Display (Priority: P1) 🎯 MVP

**Goal**: User launches debate generation and watches the conversation unfold in real-time as a chat-style interface with panelist avatars

**Independent Test**: After selecting panelists, launch debate and verify responses stream progressively into chat bubbles with correct avatar attribution, loading indicators, and error handling

### Backend Implementation for User Story 3

- [ ] T060 [P] [US3] Create DebateConfiguration request struct in backend/functions/generate-debate/types.go (topic, panelists array)
- [ ] T061 [P] [US3] Create StreamChunk response structs in backend/functions/generate-debate/types.go (message, error, done events)
- [ ] T062 [P] [US3] Implement input validation in backend/functions/generate-debate/validator.go (2-5 panelists, valid topic)
- [ ] T063 [US3] Implement Claude API streaming client in backend/functions/generate-debate/claude.go (SSE with debate prompt)
- [ ] T064 [US3] Implement SSE stream handler in backend/functions/generate-debate/stream.go (chunk parsing, panelist identification)
- [ ] T065 [US3] Implement HTTP handler in backend/functions/generate-debate/handler.go (SSE headers, flush chunks, error recovery)
- [ ] T066 [US3] Create main entry point in backend/functions/generate-debate/main.go (Cloud Function registration)
- [ ] T067 [US3] Add unit tests for debate prompt construction in backend/functions/generate-debate/claude_test.go
- [ ] T068 [US3] Add integration tests for SSE streaming in backend/functions/generate-debate/stream_test.go

### Frontend Implementation for User Story 3

- [ ] T069 [P] [US3] Create DebateBubble component in frontend/src/components/DebateView/DebateBubble.jsx (chat bubble with avatar)
- [ ] T070 [P] [US3] Create DebateBubble styles in frontend/src/components/DebateView/DebateBubble.module.css (left/right alignment per panelist)
- [ ] T071 [P] [US3] Create TypingIndicator component in frontend/src/components/DebateView/TypingIndicator.jsx (animated dots)
- [ ] T072 [P] [US3] Create TypingIndicator styles in frontend/src/components/DebateView/TypingIndicator.module.css
- [ ] T073 [P] [US3] Create DebateView component in frontend/src/components/DebateView/DebateView.jsx (scrollable message list)
- [ ] T074 [P] [US3] Create DebateView styles in frontend/src/components/DebateView/DebateView.module.css (chat interface styling)
- [ ] T075 [US3] Implement debateService in frontend/src/services/debateService.js (EventSource SSE connection)
- [ ] T076 [US3] Create useDebateStream custom hook in frontend/src/hooks/useDebateStream.js (SSE state, message accumulation, error handling)
- [ ] T077 [US3] Create DebateGeneration page in frontend/src/pages/DebateGeneration.jsx (integrate DebateView, generate button, retry logic)
- [ ] T078 [US3] Add message accumulation logic in useDebateStream hook (append chunks to correct message by panelistId)
- [ ] T079 [US3] Add auto-scroll behavior in DebateView component (scroll to latest message during streaming)
- [ ] T080 [US3] Add DebateBubble component tests in frontend/src/components/DebateView/DebateBubble.test.jsx
- [ ] T081 [US3] Add DebateView component tests in frontend/src/components/DebateView/DebateView.test.jsx
- [ ] T082 [US3] Add useDebateStream hook tests in frontend/src/hooks/useDebateStream.test.js (SSE event handling, error states)
- [ ] T083 [US3] Add accessibility tests for debate streaming flow in frontend/tests/accessibility/debate-generation.test.js (screen reader announcements, focus management)

**Checkpoint**: User Story 3 complete - users can generate and watch debates stream in real-time

---

## Phase 6: User Story 4 - PDF Export (Priority: P2)

**Goal**: User exports completed debate as a formatted PDF document for offline reading, sharing, or archival purposes

**Independent Test**: After generating a complete debate, click export button and verify PDF downloads within 2 seconds with all content (topic, panelists, messages, timestamp) properly formatted

### Frontend Implementation for User Story 4

- [ ] T084 [P] [US4] Create pdfGenerator utility in frontend/src/components/PDFExport/pdfGenerator.js (jsPDF integration, debate formatting)
- [ ] T085 [P] [US4] Create PDFExport component in frontend/src/components/PDFExport/PDFExport.jsx (export button, download trigger)
- [ ] T086 [US4] Add PDF header generation in pdfGenerator.js (debate topic, timestamp, page numbers)
- [ ] T087 [US4] Add panelist profile section in pdfGenerator.js (avatars, names, bios)
- [ ] T088 [US4] Add debate conversation rendering in pdfGenerator.js (chat bubbles with avatars, proper page breaks)
- [ ] T089 [US4] Integrate PDFExport component in DebateGeneration page (show after debate completes)
- [ ] T090 [US4] Add error handling for PDF generation failures in PDFExport component
- [ ] T091 [US4] Add pdfGenerator unit tests in frontend/src/components/PDFExport/pdfGenerator.test.js (content formatting, page breaks)
- [ ] T092 [US4] Add PDFExport component tests in frontend/src/components/PDFExport/PDFExport.test.jsx

**Checkpoint**: User Story 4 complete - users can export debates as PDFs

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Final touches, performance optimization, and comprehensive quality assurance

- [ ] T093 [P] Add loading states to all async operations (topic validation, panelist loading, debate generation)
- [ ] T094 [P] Add error retry mechanisms with exponential backoff for all API calls
- [ ] T095 [P] Optimize React component re-renders with React.memo for chat bubbles
- [ ] T096 [P] Add session storage persistence for topic and panelist selection (prevent loss on refresh)
- [ ] T097 [P] Add analytics/logging for user flow completion rates (optional, privacy-preserving)
- [ ] T098 [P] Create 404 NotFound page in frontend/src/pages/NotFound.jsx
- [ ] T099 [P] Add global CSS styles in frontend/src/index.css (typography, color scheme, responsive breakpoints)
- [ ] T100 [P] Add favicon and meta tags in frontend/public/index.html (SEO, PWA manifest)
- [ ] T101 Perform full accessibility audit with axe-core across all pages (verify WCAG 2.1 Level AA)
- [ ] T102 Perform performance audit with Lighthouse (verify SC-007: <100ms UI response)
- [ ] T103 Test complete user journey end-to-end (topic → validation → panelist selection → debate → PDF export)
- [ ] T104 Add comprehensive error scenarios testing (API failures, network timeouts, invalid responses)
- [ ] T105 Review all user-facing text for clarity and Constitution Principle I compliance
- [ ] T106 Security review of all API key handling and XSS prevention measures
- [ ] T107 Setup GitHub Actions CI workflow in .github/workflows/frontend-ci.yml (lint, test, build)
- [ ] T108 Setup GitHub Actions CI workflow in .github/workflows/backend-ci.yml (Go tests, linting)

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

- **Total Tasks**: 108
- **Setup Phase**: 8 tasks
- **Foundational Phase**: 12 tasks (BLOCKING)
- **User Story 1**: 18 tasks (7 backend + 11 frontend)
- **User Story 2**: 21 tasks (7 backend + 14 frontend)
- **User Story 3**: 24 tasks (9 backend + 15 frontend)
- **User Story 4**: 9 tasks (frontend only)
- **Polish Phase**: 16 tasks

**Parallel Opportunities**: ~60% of tasks can run in parallel after foundational phase completes

**MVP Tasks**: 81 tasks (Phases 1-5, excludes US4 and Polish)
