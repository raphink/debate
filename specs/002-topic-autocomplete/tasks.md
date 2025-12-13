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

- [ ] T001 Create Firestore composite index for autocomplete (topicLowercase ASC + createdAt DESC) via gcloud command
- [ ] T002 Update backend/shared/firebase/debates.go SaveDebate function to auto-populate topicLowercase field when saving debates

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Backend API endpoint that MUST be complete before frontend autocomplete can function

**⚠️ CRITICAL**: Frontend autocomplete cannot work until this phase is complete

- [ ] T003 Extend backend/functions/list-debates/handler.go to handle optional q query parameter with query validation (≥3 chars), sanitization, and autocomplete mode branching
- [ ] T004 Add AutocompleteDebates query function in backend/functions/list-debates/firestore.go with substring matching (WHERE topicLowercase >= query AND topicLowercase < query+"~"), ordering by createdAt DESC, limiting to 10 results
- [ ] T005 Update frontend/src/services/api.js listDebates function to accept optional query parameter for GET /api/list-debates?q={query}&limit=10

**Checkpoint**: Backend autocomplete API operational - can be tested via curl/Postman with `curl "http://localhost:8084/api/list-debates?q=free"`

---

## Phase 3: User Story 6 - Topic Discovery via History Integration (Priority: P3) 🎯

**Goal**: Enable users to see autocomplete suggestions of previous debates as they type, streamlining topic discovery by combining history browsing with topic entry

**Independent Test**: 
1. Generate 3-5 debates via existing flow (ensure saved to Firestore)
2. Return to home page, type 3+ characters matching existing topics
3. Verify autocomplete dropdown appears with matching debates showing topic, avatars, panelist count, and date
4. Select debate from dropdown → verify navigation to PanelistSelection with pre-filled panelists
5. Make no changes → verify "View Debate" button appears (cache hit)
6. Click "View Debate" → verify redirect to /d/{uuid} without generation
7. Return to home, select same debate, click "Modify Panelists" → change panelists → verify "Generate New Debate" button and full generation flow

### Implementation for User Story 6

#### Backend Enhancements

- [ ] T006 [US6] Add DebateMetadata type to backend/functions/list-debates/types.go for autocomplete response format per contracts/list-debates-autocomplete.json (if not already present)

#### Frontend Utilities

- [ ] T007 [P] [US6] Create cache detection utility in frontend/src/utils/cacheDetection.js with isCacheHit function performing deep comparison (topic exact match + panelist array order-independent)
- [ ] T008 [P] [US6] Create debounce utility hook in frontend/src/hooks/useDebounce.js with 300ms delay for autocomplete queries

#### Frontend Hooks

- [ ] T009 [US6] Create useTopicAutocomplete hook in frontend/src/hooks/useTopicAutocomplete.js managing autocomplete state (suggestions, loading, error) with debounced API calls and cleanup
- [ ] T010 [US6] Update useDebateLoader hook in frontend/src/hooks/useDebateLoader.js to check navigation state for autocomplete-sourced debates and pre-fill panelists

#### Frontend Components

- [ ] T011 [US6] Create TopicAutocompleteDropdown component in frontend/src/components/TopicAutocompleteDropdown/TopicAutocompleteDropdown.jsx displaying suggestions with topic text, panelist avatars, count badge, and generation date
- [ ] T012 [US6] Add TopicAutocompleteDropdown.module.css with dropdown positioning (absolute, below input), hover states, keyboard navigation styles, and loading indicator
- [ ] T013 [US6] Implement keyboard navigation in TopicAutocompleteDropdown (arrow keys, Enter to select, Escape to close) with ARIA accessibility attributes
- [ ] T014 [US6] Update TopicInput component in frontend/src/components/TopicInput/TopicInput.jsx to integrate useTopicAutocomplete hook and render TopicAutocompleteDropdown conditionally
- [ ] T015 [US6] Update Home.jsx in frontend/src/pages/Home.jsx to handle autocomplete selection, navigate to /panelist-selection with state: {source: 'autocomplete', debateId, topic, preFilled: panelists}
- [ ] T016 [US6] Update PanelistSelection.jsx in frontend/src/pages/PanelistSelection.jsx to detect autocomplete source from navigation state, pre-fill panelists from state.preFilled, and run cache detection
- [ ] T017 [US6] Add conditional button rendering in PanelistSelection.jsx showing "View Debate" button when cache hit detected, "Modify Panelists" button when pre-filled, and "Generate New Debate" after modifications
- [ ] T018 [US6] Implement cache hit redirect logic in PanelistSelection.jsx to navigate directly to /d/{uuid} when user clicks "View Debate" without triggering debate generation
- [ ] T019 [US6] Add "Modify Panelists" button handler in PanelistSelection.jsx to unlock panelist chip editing and change main button label to "Generate New Debate"

#### Input Sanitization

- [ ] T020 [US6] Add input sanitization to TopicInput component stripping HTML tags and special characters before passing to autocomplete API (reuse backend sanitization logic pattern)

#### Error Handling & Edge Cases

- [ ] T021 [US6] Add graceful degradation in TopicAutocompleteDropdown to hide dropdown on API failures or empty results without blocking "Find Panelists" button
- [ ] T022 [US6] Add loading state indicator in TopicAutocompleteDropdown shown when API response takes >300ms
- [ ] T023 [US6] Implement click-outside detection in TopicAutocompleteDropdown to close dropdown when user clicks outside
- [ ] T024 [US6] Add dropdown auto-hide logic when user input length drops below 3 characters

**Checkpoint**: At this point, User Story 6 should be fully functional and testable independently

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements affecting multiple components

- [ ] T025 [P] Update README.md with autocomplete feature documentation including setup instructions for Firestore index
- [ ] T026 [P] Add autocomplete feature section to quickstart.md validation checklist
- [ ] T027 Update DEPLOYMENT.md with list-debates function autocomplete enhancement and Firestore index requirements
- [ ] T028 Add error logging for autocomplete API failures in backend/functions/list-debates/handler.go with context about query and Firestore errors
- [ ] T029 [P] Add performance monitoring for autocomplete API response times (consider Cloud Function metrics)
- [ ] T030 Run quickstart.md validation: generate debates, test autocomplete flow, verify cache detection, test graceful degradation

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
  - T001 (Firestore index) is CRITICAL for query performance
- **Foundational (Phase 2)**: Depends on T001, T002 completion - BLOCKS all User Story 6 frontend work
  - Backend API must be operational before frontend integration
- **User Story 6 (Phase 3)**: Depends on Foundational (Phase 2) completion
  - All frontend tasks depend on T005 (API service function)
- **Polish (Phase 4)**: Depends on User Story 6 completion

### User Story 6 Internal Dependencies

**Backend (can run in parallel once T001-T002 done)**:
- T003, T004, T006 can proceed in parallel

**Frontend Utilities (can run in parallel once T005 done)**:
- T007 (cache detection) - no dependencies
- T008 (debounce hook) - no dependencies

**Frontend Hooks (sequential dependencies)**:
- T009 (useTopicAutocomplete) - depends on T005, T008
- T010 (useDebateLoader update) - depends on T007

**Frontend Components (complex dependencies)**:
- T011, T012, T013 (TopicAutocompleteDropdown) - depends on T009
- T014 (TopicInput update) - depends on T009, T011
- T015 (Home.jsx update) - depends on T014
- T016 (PanelistSelection detection) - depends on T007, T010
- T017, T018, T019 (PanelistSelection buttons) - depends on T016

**Cross-cutting**:
- T020 (sanitization) - can run in parallel with T014
- T021-T024 (error handling) - run after corresponding components complete

### Parallel Opportunities

**Setup Phase**:
- T001 and T002 can run in parallel

**Foundational Phase**:
- T003, T004, T006 can all run in parallel once T001-T002 complete
- T005 depends on T003 completion

**User Story 6 - Backend**:
- T006 can run in parallel with T003-T004

**User Story 6 - Frontend Utilities**:
- T007, T008 can run in parallel once T005 complete

**User Story 6 - Components (after T011)**:
- T012, T013 can run in parallel with T011

**Polish Phase**:
- T025, T026, T029 can run in parallel

---

## Parallel Example: User Story 6 Frontend Utilities

```bash
# Launch foundational backend tasks together:
Task: "Extend backend/functions/list-debates/handler.go for autocomplete mode"
Task: "Add AutocompleteDebates query function in backend/functions/list-debates/firestore.go"
Task: "Add DebateMetadata type to backend/functions/list-debates/types.go"

# Launch frontend utilities together (after T005):
Task: "Create cache detection utility in frontend/src/utils/cacheDetection.js"
Task: "Create debounce utility hook in frontend/src/hooks/useDebounce.js"

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
2. Add Frontend Utilities (T007-T008) → Test independently
3. Add TopicAutocompleteDropdown (T011-T013) → Test component in isolation
4. Integrate with TopicInput (T014) → Test autocomplete dropdown appears
5. Add navigation logic (T015) → Test navigation to PanelistSelection
6. Add cache detection (T016-T019) → Test full workflow
7. Polish and deploy (Phase 4)

### Single Developer Strategy

Sequential execution in priority order:
1. Phase 1: Setup (T001-T002) - ~30 minutes
2. Phase 2: Foundational (T003-T005) - ~2-3 hours
3. Phase 3: User Story 6 - ~6-8 hours
   - Backend: T006
   - Utilities: T007-T008
   - Components: T011-T024 (most time-consuming)
4. Phase 4: Polish - ~1-2 hours

**Total Estimated Time**: ~10-14 hours for complete feature

---

## Notes

- [P] tasks = different files, no dependencies - can run in parallel
- [US6] label maps all implementation tasks to User Story 6
- User Story 6 should be independently completable and testable
- Commit after each task or logical group (e.g., complete component with its styles)
- Stop at Phase 2 checkpoint to validate backend before frontend work
- Avoid: same file conflicts (coordinate T014, T020 edits to TopicInput)
- Graceful degradation is CRITICAL - autocomplete failures must never block normal topic entry workflow
- Cache detection must be synchronous (<50ms) - no API calls in comparison logic
- Debouncing is mandatory to prevent Firestore quota exhaustion
- All input must be sanitized before Firestore queries
- Extending list-debates maintains single source of truth and reduces deployment complexity
