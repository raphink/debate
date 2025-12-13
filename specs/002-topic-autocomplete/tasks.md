---
description: "Task list for Topic Discovery via History Integration (US6)"
---

# Tasks: Topic Discovery via History Integration (US6)

**Input**: Design documents from `/specs/002-topic-autocomplete/`
**Prerequisites**: plan.md, spec.md, quickstart.md, contracts/list-debates-autocomplete.json
**Dependencies**: Requires US5 (Debate Caching & Sharing) - Firestore debates collection must be populated

**Tests**: Not explicitly requested in specification - omitted per template guidelines

**Organization**: Tasks organized around single user story (US6) with setup, foundational, and implementation phases

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US6 for all implementation tasks)
- Include exact file paths in descriptions

## Path Conventions

- **Web app structure**: `backend/functions/`, `frontend/src/`
- Backend: Go 1.24, Cloud Functions
- Frontend: React 18+, React Router

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Firestore indexing and data preparation for autocomplete queries

- [ ] T001 Verify Firestore debates collection exists and is accessible (no special index needed - will fetch and filter in code)
## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Backend API endpoint that MUST be complete before frontend autocomplete can function

**⚠️ CRITICAL**: Frontend autocomplete cannot work until this phase is complete

- [X] T002 Extend backend/functions/list-debates/handler.go to handle optional q query parameter with query validation (≥3 chars), sanitization, and autocomplete mode branching
- [X] T003 Add AutocompleteDebates query function in backend/functions/list-debates/firestore.go that fetches recent debates, filters by topic substring (case-insensitive) in code, and returns top 10 matches ordered by createdAt DESC
- [X] T004 Update frontend/src/services/api.js listDebates function to accept optional query parameter for GET /api/list-debates?q={query}&limit=10

**Checkpoint**: Backend autocomplete API operational - can be tested via curl/Postman with `curl "http://localhost:8084/api/list-debates?q=free"`

---

## Phase 3: User Story 6 - Topic Discovery via History Integration (Priority: P3) 🎯

**Goal**: Enable users to see autocomplete suggestions of previous debates as they type, providing quick access to view existing debates without re-entering topics or regenerating content

**Independent Test**: 
1. Generate 3-5 debates via existing flow (ensure saved to Firestore)
2. Return to home page, type 3+ characters matching existing topics
3. Verify autocomplete dropdown appears with matching debates showing topic, avatars, panelist count, and date
4. Select debate from dropdown → verify navigation to /d/{debate.id} with complete debate displayed
5. Navigate back to home → verify can select other autocomplete suggestions or create new debate via "Find Panelists"

### Implementation for User Story 6

#### Backend Enhancements

- [X] T005 [US6] Add DebateMetadata type to backend/functions/list-debates/types.go for autocomplete response format per contracts/list-debates-autocomplete.json (if not already present)

#### Frontend Utilities

- [X] T006 [P] [US6] Create debounce utility hook in frontend/src/hooks/useDebounce.js with 300ms delay for autocomplete queries

#### Frontend Hooks

- [X] T007 [US6] Create useTopicAutocomplete hook in frontend/src/hooks/useTopicAutocomplete.js managing autocomplete state (suggestions, loading, error) with debounced API calls and cleanup

#### Frontend Components

- [X] T008 [US6] Create TopicAutocompleteDropdown component in frontend/src/components/TopicAutocompleteDropdown/TopicAutocompleteDropdown.jsx displaying suggestions with topic text, panelist avatars, count badge, and generation date
- [X] T009 [US6] Add TopicAutocompleteDropdown.module.css with dropdown positioning (absolute, below input), hover states, keyboard navigation styles, and loading indicator
- [X] T010 [US6] Implement keyboard navigation in TopicAutocompleteDropdown (arrow keys, Enter to select, Escape to close) with ARIA accessibility attributes
- [X] T011 [US6] Update TopicInput component in frontend/src/components/TopicInput/TopicInput.jsx to integrate useTopicAutocomplete hook and render TopicAutocompleteDropdown conditionally
- [X] T012 [US6] Update Home.jsx in frontend/src/pages/Home.jsx to handle autocomplete selection and navigate to /d/{debate.id} to view existing debate
- [X] T013 [US6] **CHANGED**: PanelistSelection pre-fill not needed - autocomplete navigates directly to debate viewer. Task marked complete but implementation uses different approach (view existing debate instead of pre-fill for regeneration)

#### Input Sanitization

- [X] T014 [US6] Add input sanitization to TopicInput component stripping HTML tags and special characters before passing to autocomplete API (reuse backend sanitization logic pattern)

#### Error Handling & Edge Cases

- [X] T015 [US6] Add graceful degradation in TopicAutocompleteDropdown to hide dropdown on API failures or empty results without blocking "Find Panelists" button
- [X] T016 [US6] Add loading state indicator in TopicAutocompleteDropdown shown when API response takes >300ms
- [X] T017 [US6] Implement click-outside detection in TopicAutocompleteDropdown to close dropdown when user clicks outside
- [X] T018 [US6] Add dropdown auto-hide logic when user input length drops below 3 characters

**Checkpoint**: At this point, User Story 6 should be fully functional and testable independently

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements affecting multiple components

- [X] T019 [P] Update README.md with autocomplete feature documentation
- [X] T020 [P] Add autocomplete feature section to quickstart.md validation checklist
- [X] T021 Update DEPLOYMENT.md with list-debates function autocomplete enhancement
- [X] T022 Add error logging for autocomplete API failures in backend/functions/list-debates/handler.go with context about query and Firestore errors
- [X] T023 [P] Add performance monitoring for autocomplete API response times (consider Cloud Function metrics)
- [X] T024 Run quickstart.md validation: generate debates, test autocomplete flow, verify direct navigation to debate viewer, test graceful degradation

---

## Implementation Notes

**UX Simplification (2025-12-13)**: During implementation, the navigation flow was simplified from the original specification:
- **Original Design**: Select autocomplete → navigate to PanelistSelection with pre-filled panelists → modify/keep panelists → generate new debate
- **Implemented Design**: Select autocomplete → navigate directly to /d/{debate.id} to view existing debate
- **Rationale**: Simpler UX with fewer clicks; clear separation between viewing existing debates (autocomplete) and creating new ones ("Find Panelists" button)
- **Tasks Affected**: T012 (navigation destination changed), T013 (PanelistSelection pre-fill not needed)
- **Specification Updated**: 2025-12-13 to align with implemented behavior

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - verification only (T001)
- **Foundational (Phase 2)**: Depends on T001 completion - BLOCKS all User Story 6 frontend work
  - Backend API must be operational before frontend integration
- **User Story 6 (Phase 3)**: Depends on Foundational (Phase 2) completion
  - All frontend tasks depend on T004 (API service function)
- **Polish (Phase 4)**: Depends on User Story 6 completion

### User Story 6 Internal Dependencies

**Backend (can run in parallel once T001 done)**:
- T002, T003, T005 can proceed in parallel

**Frontend Utilities (can run in parallel once T004 done)**:
- T006 (debounce hook) - no dependencies

**Frontend Hooks**:
- T007 (useTopicAutocomplete) - depends on T004, T006

**Frontend Components (sequential dependencies)**:
- T008, T009, T010 (TopicAutocompleteDropdown) - depends on T007
- T011 (TopicInput update) - depends on T007, T008
- T012 (Home.jsx update) - depends on T011
- T013 (marked CHANGED - not needed in revised flow)

**Cross-cutting**:
- T014 (sanitization) - can run in parallel with T011
- T015-T018 (error handling) - run after corresponding components complete

### Parallel Opportunities

**Setup Phase**:
- T001 and T002 can run in parallel

**Foundational Phase**:
- T002, T003, T005 can all run in parallel once T001 complete
- T004 depends on T002 completion

**User Story 6 - Backend**:
- T005 can run in parallel with T002-T003

**User Story 6 - Frontend**:
- T006 standalone once T004 complete

**User Story 6 - Components**:
- T009, T010 can run in parallel with T008
- T012, T013 can run in parallel

**Polish Phase**:
- T019, T020, T023 can run in parallel

---

## Parallel Example: User Story 6 Frontend Utilities

```bash
# Launch foundational backend tasks together:
Task: "Extend backend/functions/list-debates/handler.go for autocomplete mode"
Task: "Add AutocompleteDebates query function in backend/functions/list-debates/firestore.go"
Task: "Add DebateMetadata type to backend/functions/list-debates/types.go"

# Launch frontend components together (after T008):
Task: "Add TopicAutocompleteDropdown.module.css with dropdown positioning"
Task: "Implement keyboard navigation in TopicAutocompleteDropdown"

# Launch documentation tasks together:
Task: "Update README.md with autocomplete feature documentation"
Task: "Add autocomplete feature section to quickstart.md"
Task: "Add performance monitoring for autocomplete API response times"
```

---

## Implementation Strategy

### MVP First (Backend API Only)

1. Complete Phase 1: Setup (Firestore indexing)
2. Complete Phase 2: Foundational (Backend API extension)
3. **STOP and VALIDATE**: Test backend API via curl/Postman
   - curl "http://localhost:8084/api/list-debates?q=free"
   - Verify JSON response with debate metadata
4. Proceed to frontend only after backend validation

### Incremental Delivery

1. Complete Setup + Foundational → Backend API ready
2. Add Frontend Utilities (T006) → Test debounce independently
3. Add TopicAutocompleteDropdown (T008-T010) → Test component in isolation
4. Integrate with TopicInput (T011) → Test autocomplete dropdown appears
5. Add navigation logic (T012) → Test direct navigation to debate viewer
6. Polish and deploy (Phase 4)

### Single Developer Strategy

Sequential execution in priority order:
1. Phase 1: Setup (T001) - ~5 minutes
2. Phase 2: Foundational (T002-T004) - ~2-3 hours
3. Phase 3: User Story 6 - ~4-6 hours
   - Backend: T005
   - Utilities: T006
   - Components: T007-T018 (most time-consuming)
4. Phase 4: Polish - ~1-2 hours

**Total Estimated Time**: ~7-11 hours for complete feature

---

## Notes

- [P] tasks = different files, no dependencies - can run in parallel
- [US6] label maps all implementation tasks to User Story 6
- User Story 6 should be independently completable and testable
- Commit after each task or logical group (e.g., complete component with its styles)
- Stop at Phase 2 checkpoint to validate backend before frontend work
- Avoid: same file conflicts (coordinate T011, T014 edits to TopicInput)
- Graceful degradation is CRITICAL - autocomplete failures must never block normal topic entry workflow
- Debouncing is mandatory to prevent Firestore quota exhaustion
- All input must be sanitized before Firestore queries
- Extending list-debates maintains single source of truth and reduces deployment complexity
- Direct navigation to debate viewer simplifies UX - autocomplete for viewing existing debates, "Find Panelists" for creating new ones
- No Firestore index needed - fetch recent debates and filter by substring in code for true full-text search
