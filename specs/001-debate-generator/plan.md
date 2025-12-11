# Implementation Plan: AI-Powered Theology/Philosophy Debate Generator

**Branch**: `001-debate-generator` | **Date**: 2025-12-11 | **Spec**: [spec.md](spec.md)
**Input**: Feature specification from `/specs/001-debate-generator/spec.md`

## Summary

Build a web application that generates AI-powered theological/philosophical debates between historical figures. Users enter a topic, receive topic validation and panelist suggestions via Claude API (through GCP Cloud Functions proxy in Go), select 2-5 panelists, then watch the debate stream in real-time as a chat-style conversation with avatars. Completed debates can be exported as PDF.

## Technical Context

**Language/Version**: Go 1.23+ (backend/GCP Functions), JavaScript/React 18+ (frontend)  
**Primary Dependencies**: 
- Frontend: React 18, React Router, Axios (HTTP client), jsPDF (PDF generation), CSS Modules
- Backend: Go 1.23+, Google Cloud Functions SDK, Anthropic Claude API SDK
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
**Project Type**: Web application (frontend + backend)  
**Performance Goals**: 
- Topic validation response: <3s
- First debate response: <5s
- Streaming chunk intervals: <500ms
- UI interaction response: <100ms
- PDF generation: <2s for 5000 words
**Constraints**: 
- No database/persistence (stateless MVP)
- Claude API rate limits (per Anthropic tier)
- GCP Cloud Functions timeout: 60s max per request
- Client-side PDF generation to avoid server overhead
- Mobile-first responsive design (≥375px width)
**Scale/Scope**: 
- MVP: Single-user sessions, no concurrent debate limit
- Expected load: <100 concurrent users initially
- Debate length: ~10-20 exchanges (5000-10000 words typical)
- Frontend: ~15-20 components, 5-8 pages/views
- Backend: 2-3 Cloud Functions

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
│   ├── validate-topic/          # GCP Cloud Function: Topic validation
│   │   ├── main.go
│   │   ├── handler.go
│   │   ├── claude.go            # Claude API client
│   │   ├── validator.go         # Input validation
│   │   └── go.mod
│   ├── suggest-panelists/       # GCP Cloud Function: Panelist suggestions
│   │   ├── main.go
│   │   ├── handler.go
│   │   ├── claude.go
│   │   ├── panelist.go          # Panelist data structures
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
│   │   ├── debateService.js     # Debate generation API calls (SSE)
│   │   └── sanitizer.js         # DOMPurify wrapper for XSS prevention
│   ├── hooks/
│   │   ├── useDebateStream.js   # Custom hook for SSE streaming
│   │   ├── usePanelistSelection.js
│   │   └── useTopicValidation.js
│   ├── pages/
│   │   ├── Home.jsx             # Topic entry page
│   │   ├── PanelistSelection.jsx
│   │   ├── DebateGeneration.jsx
│   │   └── NotFound.jsx
│   ├── utils/
│   │   ├── validation.js        # Client-side input validation
│   │   ├── constants.js         # App constants (max panelists, etc.)
│   │   └── accessibility.js     # A11y utilities
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

# Docker/DevOps files
docker-compose.yml               # Local development orchestration
.dockerignore                    # Docker build exclusions
start-local.sh                   # Quick start script

.github/
└── workflows/
    ├── frontend-ci.yml          # Frontend lint, test, build
    ├── backend-ci.yml           # Backend Go tests
    └── deploy.yml               # Deploy to GCP
```

**Structure Decision**: Web application architecture selected due to separate frontend (React SPA) and backend (GCP Cloud Functions). Frontend handles all UI/UX concerns including streaming display and PDF export. Backend provides three focused functions acting as a secure proxy to Claude API, implementing rate limiting and input validation. No shared state between components - each function is independently deployable.
