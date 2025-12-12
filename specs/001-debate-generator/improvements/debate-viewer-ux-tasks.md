# Tasks: Debate Viewer UX Improvements

**Feature**: Enhanced Debate Viewer Interface  
**Spec**: [debate-viewer-ux.md](debate-viewer-ux.md)  
**Plan**: [debate-viewer-ux-plan.md](debate-viewer-ux-plan.md)  
**Created**: 2025-12-12

## Implementation Tasks

### Phase 1: Component Updates

- [X] T001 [P] Add `mode` prop to DebateView component in frontend/src/components/DebateView/DebateView.jsx (default: 'generation')
- [X] T002 [P] Add PropTypes validation for `mode` prop (oneOf: 'viewer', 'generation')
- [X] T003 Import `useNavigate` hook from react-router-dom in DebateView.jsx
- [X] T004 Wrap auto-scroll toggle section in conditional: `{mode === 'generation' && ...}`
- [X] T005 Add "New Debate" button component below messagesEndRef with conditional: `{mode === 'viewer' && ...}`
- [X] T006 Implement button click handler using navigate('/') for home navigation

### Phase 2: Styling

- [X] T007 [P] Add `.newDebateButtonContainer` class in DebateView.module.css (flex center, padding, border-top)
- [X] T008 [P] Add `.newDebateButton` class with gradient background matching design system
- [X] T009 [P] Add hover/active states for button (transform, shadow effects)
- [X] T010 [P] Add mobile responsive styles for button (@media max-width: 768px)

### Phase 3: Parent Component Integration

- [X] T011 Update DebateViewer.jsx to pass `mode="viewer"` prop to DebateView
- [X] T012 Update DebateGeneration.jsx to pass `mode="generation"` prop to DebateView (optional, verifies explicit mode)

### Phase 4: Testing & QA

- [ ] T013 Manual test: Navigate to `/d/{uuid}` → verify auto-scroll toggle hidden
- [ ] T014 Manual test: Verify "New Debate" button visible at bottom of viewer
- [ ] T015 Manual test: Click "New Debate" button → verify navigates to home page
- [ ] T016 Manual test: Generate new debate → verify auto-scroll toggle still visible
- [ ] T017 Manual test: Verify "New Debate" button NOT visible during live generation
- [ ] T018 Accessibility test: Verify button is keyboard accessible (Tab + Enter)
- [ ] T019 Mobile test: Verify button styling on small screens (iPhone, Android)
- [ ] T020 Test error state: Verify existing error buttons still work in DebateViewer

### Phase 5: Header Consistency (Additional Improvements)

- [X] T021 Add header section to DebateViewer.jsx with topic display and Back button
- [X] T022 Display panelists list in DebateViewer header (matching DebateGeneration format)
- [X] T023 Make panelists clickable in DebateViewer header (button elements with onClick)
- [X] T024 Make panelists clickable in DebateGeneration header (convert spans to buttons)
- [X] T025 Import and integrate PanelistModal in DebateViewer for panelist details
- [X] T026 Import and integrate PanelistModal in DebateGeneration for panelist details
- [X] T027 Add state management for selected panelist modal in both pages
- [X] T028 Update DebateViewer.module.css with header, back button, and panelist styles
- [X] T029 Update DebateGeneration.module.css with clickable panelist button styles
- [X] T030 Add "Topic: " prefix to DebateViewer title (matching DebateGeneration format)
- [X] T031 Transform panelist data in DebateViewer to map backend 'biography' field to frontend 'bio' field
- [X] T032 Add CSS styling for topic label and em tag in DebateViewer

## Execution Plan

**Parallel Opportunities:**
- T001-T006 (component logic) can run parallel to T007-T010 (CSS)
- TX] Auto-scroll toggle hidden in viewer mode
- [X] Auto-scroll toggle visible in generation mode
- [X] "New Debate" button visible in viewer mode
- [X] "New Debate" button hidden in generation mode
- [X] Button navigates to home page
- [X] Button styling matches design system
- [ ] Mobile responsive
- [ ] Keyboard accessible
- [X] DebateViewer has header with topic and panelists (matching DebateGeneration)
- [X] Panelists are clickable in both viewer and generation headers
- [X] Clicking panelist opens modal with full bio/position
- [X] Topic displays with "Topic: " prefix in viewer mode
- [X] Backend panelist data (biography field) properly mapped to frontend (bio field)

## Task Statistics

- **Total Tasks**: 32
- **Component Logic**: 6 tasks (T001-T006)
- **Styling**: 4 tasks (T007-T010)
- **Integration**: 2 tasks (T011-T012)
- **Testing**: 8 tasks (T013-T020)
- **Header Consistency**: 12 tasks (T021-T032)
- **Completed**: 24 tasks
- **Remaining**: 8 tasks (testing onlyn viewer mode
- [ ] Auto-scroll toggle visible in generation mode
- [ ] "New Debate" button visible in viewer mode
- [ ] "New Debate" button hidden in generation mode
- [ ] Button navigates to home page
- [ ] Button styling matches design system
- [ ] Mobile responsive
- [ ] Keyboard accessible

## Task Statistics

- **Total Tasks**: 20
- **Component Logic**: 6 tasks (T001-T006)
- **Styling**: 4 tasks (T007-T010)
- **Integration**: 2 tasks (T011-T012)
- **Testing**: 8 tasks (T013-T020)
