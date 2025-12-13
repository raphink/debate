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
5. Keep or modify panelists → click "Generate Debate" → verify new debate generated with selected panelists

### Implementation for User Story 6

#### Backend Enhancements

- [ ] T006 [US6] Add DebateMetadata type to backend/functions/list-debates/types.go for autocomplete response format per contracts/list-debates-autocomplete.json (if not already present)

#### Frontend Utilities

- [ ] T007 [P] [US6] Create debounce utility hook in frontend/src/hooks/useDebounce.js with 300ms delay for autocomplete queries

#### Frontend Hooks

- [ ] T008 [US6] Create useTopicAutocomplete hook in frontend/src/hooks/useTopicAutocomplete.js managing autocomplete state (suggestions, loading, error) with debounced API calls and cleanup

#### Frontend Components

- [ ] T009 [US6] Create TopicAutocompleteDropdown component in frontend/src/components/TopicAutocompleteDropdown/TopicAutocompleteDropdown.jsx displaying suggestions with topic text, panelist avatars, count badge, and generation date
- [ ] T010 [US6] Add TopicAutocompleteDropdown.module.css with dropdown positioning (absolute, below input), hover states, keyboard navigation styles, and loading indicator
- [ ] T011 [US6] Implement keyboard navigation in TopicAutocompleteDropdown (arrow keys, Enter to select, Escape to close) with ARIA accessibility attributes
- [ ] T012 [US6] Update TopicInput component in frontend/src/components/TopicInput/TopicInput.jsx to integrate useTopicAutocomplete hook and render TopicAutocompleteDropdown conditionally
- [ ] T013 [US6] Update Home.jsx in frontend/src/pages/Home.jsx to handle autocomplete selection, navigate to /panelist-selection with state: {source: 'autocomplete', topic, preFilled: panelists}
- [ ] T014 [US6] Update PanelistSelection.jsx in frontend/src/pages/PanelistSelection.jsx to detect autocomplete source from navigation state and pre-fill panelists from state.preFilled

#### Input Sanitization

- [ ] T015 [US6] Add input sanitization to TopicInput component stripping HTML tags and special characters before passing to autocomplete API (reuse backend sanitization logic pattern)

#### Error Handling & Edge Cases

- [ ] T016 [US6] Add graceful degradation in TopicAutocompleteDropdown to hide dropdown on API failures or empty results without blocking "Find Panelists" button
- [ ] T017 [US6] Add loading state indicator in TopicAutocompleteDropdown shown when API response takes >300ms
- [ ] T018 [US6] Implement click-outside detection in TopicAutocompleteDropdown to close dropdown when user clicks outside
- [ ] T019 [US6] Add dropdown auto-hide logic when user input length drops below 3 characters

**Checkpoint**: At this point, User Story 6 should be fully functional and testable independently

---

## Phase 4: Polish & Cross-Cutting Concerns

**Purpose**: Improvements affecting multiple components

- [ ] T020 [P] Update README.md with autocomplete feature documentation including setup instructions for Firestore index
- [ ] T021 [P] Add autocomplete feature section to quickstart.md validation checklist
- [ ] T022 Update DEPLOYMENT.md with list-debates function autocomplete enhancement and Firestore index requirements
- [ ] T023 Add error logging for autocomplete API failures in backend/functions/list-debates/handler.go with context about query and Firestore errors
- [ ] T024 [P] Add performance monitoring for autocomplete API response times (consider Cloud Function metrics)
- [ ] T025 Run quickstart.md validation: generate debates, test autocomplete flow, verify panelist pre-fill, test graceful degradation

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
- T007 (debounce hook) - no dependencies

**Frontend Hooks**:
- T008 (useTopicAutocomplete) - depends on T005, T007

**Frontend Components (sequential dependencies)**:
- T009, T010, T011 (TopicAutocompleteDropdown) - depends on T008
- T012 (TopicInput update) - depends on T008, T009
- T013 (Home.jsx update) - depends on T012
- T014 (PanelistSelection pre-fill) - can run in parallel with T013

**Cross-cutting**:
- T015 (sanitization) - can run in parallel with T012
- T016-T019 (error handling) - run after corresponding components complete

### Parallel Opportunities

**Setup Phase**:
- T001 and T002 can run in parallel

**Foundational Phase**:
- T003, T004, T006 can all run in parallel once T001-T002 complete
- T005 depends on T003 completion

**User Story 6 - Backend**:
- T006 can run in parallel with T003-T004

**User Story 6 - Frontend**:
- T007 standalone once T005 complete

**User Story 6 - Components**:
- T010, T011 can run in parallel with T009
- T013, T014 can run in parallel

**Polish Phase**:
- T020, T021, T024 can run in parallel

---

## Parallel Example: User Story 6 Frontend Utilities

```bash
# Launch foundational backend tasks together:
Task: "Extend backend/functions/list-debates/handler.go for autocomplete mode"
Task: "Add AutocompleteDebates query function in backend/functions/list-debates/firestore.go"
Task: "Add DebateMetadata type to backend/functions/list-debates/types.go"

# Launch frontend components together (after T009):
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
2. Add Frontend Utilities (T007) → Test debounce independently
3. Add TopicAutocompleteDropdown (T009-T011) → Test component in isolation
4. Integrate with TopicInput (T012) → Test autocomplete dropdown appears
5. Add navigation logic (T013-T014) → Test panelist pre-fill workflow
6. Polish and deploy (Phase 4)

### Single Developer Strategy

Sequential execution in priority order:
1. Phase 1: Setup (T001-T002) - ~30 minutes
2. Phase 2: Foundational (T003-T005) - ~2-3 hours
3. Phase 3: User Story 6 - ~4-6 hours
   - Backend: T006
   - Utilities: T007
   - Components: T008-T019 (most time-consuming)
4. Phase 4: Polish - ~1-2 hours

**Total Estimated Time**: ~8-12 hours for complete feature

---

## Notes

- [P] tasks = different files, no dependencies - can run in parallel
- [US6] label maps all implementation tasks to User Story 6
- User Story 6 should be independently completable and testable
- Commit after each task or logical group (e.g., complete component with its styles)
- Stop at Phase 2 checkpoint to validate backend before frontend work
- Avoid: same file conflicts (coordinate T012, T015 edits to TopicInput)
- Graceful degradation is CRITICAL - autocomplete failures must never block normal topic entry workflow
- Debouncing is mandatory to prevent Firestore quota exhaustion
- All input must be sanitized before Firestore queries
- Extending list-debates maintains single source of truth and reduces deployment complexity
- No cache hit detection needed - panelist names may vary across LLM responses, always generate new debate
