# Unified Topic Input - Implementation Tasks

## Phase 1: Component Refactoring (6 tasks)

### T1: Extract PanelistChipSelector Component
**Files**: 
- `frontend/src/components/PanelistChips/PanelistChipSelector.jsx` (new)
- `frontend/src/components/PanelistChips/PanelistChip.jsx` (new)
- `frontend/src/components/PanelistChips/PanelistChipSelector.module.css` (new)

**Implementation**:
- [ ] Create `PanelistChip.jsx` with:
  - Props: `panelist`, `onRemove`, `disabled`
  - Avatar display (40x40)
  - Name label
  - Remove button (X icon)
  - CSS module for styling
- [ ] Create `PanelistChipSelector.jsx` with:
  - Props: `value`, `onChange`, `disabled`, `max`
  - State: selected panelists array
  - Chip list rendering (horizontal scroll if needed)
  - "+" button to add panelists (placeholder for now)
  - Max count validation (default 10)
- [ ] Create CSS module:
  - Flexbox layout with wrap
  - Mobile-responsive (vertical stack on small screens)
  - Disabled state styling
  - Empty state placeholder
- [ ] Test chip add/remove functionality

**Acceptance**:
- Chips render with avatar and name
- Remove button works correctly
- Max limit enforced
- Component reusable across pages

---

### T2: Create UnifiedTopicInput Component
**Files**:
- `frontend/src/components/TopicInput/UnifiedTopicInput.jsx` (new)
- `frontend/src/components/TopicInput/UnifiedTopicInput.module.css` (new)

**Implementation**:
- [ ] Create component scaffold with props:
  - `value`, `onChange`, `onSelectSuggestion`, `onSubmit`
  - `error`, `disabled`, `placeholder`
- [ ] Integrate `useTopicAutocomplete` hook:
  - Debounce 300ms
  - Min 3 chars
  - Return suggestions, loading, error
- [ ] Implement keyboard navigation:
  - State: `selectedIndex` (default -1)
  - ArrowDown: increment index (max: suggestions.length - 1)
  - ArrowUp: decrement index (min: -1)
  - Enter: select if index >= 0, else submit
  - Escape: close dropdown
- [ ] Manage dropdown visibility:
  - State: `showDropdown`
  - Show when: value.length >= 3 && suggestions.length > 0
  - Hide when: blur, escape, selection, < 3 chars
- [ ] Render structure:
  ```jsx
  <div className={styles.container}>
    <input {...inputProps} />
    {showDropdown && <TopicAutocompleteDropdown />}
    {error && <div className={styles.error}>{error}</div>}
  </div>
  ```
- [ ] Add click-outside hook to dismiss dropdown
- [ ] Style with CSS module (relative positioning for dropdown)

**Acceptance**:
- Autocomplete appears on 3+ chars
- Keyboard navigation works
- Selection triggers callback
- Error displays correctly

---

### T3: Refactor Home.jsx State Management
**Files**:
- `frontend/src/pages/Home.jsx`
- `frontend/src/pages/Home.module.css`

**Implementation**:
- [ ] Remove old state:
  - Delete separate autocomplete search state
  - Delete "or" divider elements
- [ ] Add new state:
  ```jsx
  const [topic, setTopic] = useState('');
  const [selectedPanelists, setSelectedPanelists] = useState([]);
  const [validationError, setValidationError] = useState(null);
  const [isValidating, setIsValidating] = useState(false);
  ```
- [ ] Create `handleAutocompleteSelect` function:
  ```jsx
  const handleAutocompleteSelect = (debate) => {
    navigate('/select-panelists', {
      state: {
        debateId: debate.id,
        topic: debate.topic,
        panelists: debate.panelists,
        skipValidation: true
      }
    });
  };
  ```
- [ ] Create `handleFindPanelists` function:
  ```jsx
  const handleFindPanelists = async () => {
    setIsValidating(true);
    setValidationError(null);
    
    try {
      const result = await validateTopic(topic, selectedPanelists);
      
      if (result.isRelevant) {
        navigate('/select-panelists', {
          state: {
            topic: result.topic,
            panelists: selectedPanelists,
            skipValidation: false
          }
        });
      } else {
        setValidationError(result.validationMessage);
      }
    } catch (error) {
      setValidationError('Failed to validate topic. Please try again.');
    } finally {
      setIsValidating(false);
    }
  };
  ```
- [ ] Update render to use new components

**Acceptance**:
- State properly initialized
- Handlers navigate correctly
- Validation errors display

---

### T4: Update Home.jsx Render Method
**Files**:
- `frontend/src/pages/Home.jsx`
- `frontend/src/pages/Home.module.css`

**Implementation**:
- [ ] Replace dual inputs with unified structure:
  ```jsx
  return (
    <div className={styles.container}>
      <h1>Generate a Debate</h1>
      <p className={styles.subtitle}>
        Enter a topic to search previous debates or create a new one
      </p>
      
      <UnifiedTopicInput
        value={topic}
        onChange={setTopic}
        onSelectSuggestion={handleAutocompleteSelect}
        error={validationError}
        disabled={isValidating}
        placeholder="What should they debate?"
      />
      
      <div className={styles.section}>
        <label className={styles.label}>
          Select Panelists <span className={styles.optional}>(Optional)</span>
        </label>
        <PanelistChipSelector
          value={selectedPanelists}
          onChange={setSelectedPanelists}
          disabled={isValidating || topic.length < 3}
          max={10}
        />
      </div>
      
      <button
        onClick={handleFindPanelists}
        disabled={topic.length < 3 || isValidating}
        className={styles.submitButton}
      >
        {isValidating ? 'Validating...' : 'Find Panelists →'}
      </button>
    </div>
  );
  ```
- [ ] Update CSS module:
  - Remove `.autocompleteSection`, `.divider`, `.manualSection`
  - Add `.container` with max-width 600px, centered
  - Add `.section` for panelist chips
  - Add `.submitButton` full-width styling
  - Add `.optional` subtle text styling
- [ ] Test responsive layout (mobile, tablet, desktop)

**Acceptance**:
- Single input displays correctly
- Chip selector below input
- Button at bottom
- Mobile-responsive

---

### T5: Enhance Validation API with Panelist Context
**Files**:
- `frontend/src/services/topicService.js`
- `backend/functions/validate-topic/handler.go`
- `backend/functions/validate-topic/types.go`

**Implementation**:
- [ ] Update `topicService.js`:
  ```javascript
  export const validateTopic = async (topic, panelists = []) => {
    const response = await fetch(`${API_BASE}/validate-topic`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({
        topic,
        panelists: panelists.map(p => p.id || p.slug)
      })
    });
    
    if (!response.ok) throw new Error('Validation failed');
    return response.json();
  };
  ```
- [ ] Update backend `types.go`:
  ```go
  type ValidationRequest struct {
      Topic     string   `json:"topic"`
      Panelists []string `json:"panelists,omitempty"` // NEW
  }
  ```
- [ ] Update `handler.go` to use panelist context:
  - If panelists provided, enhance AI prompt:
    ```
    "Topic: {topic}
     Panelists: {panelist names}
     Is this a good debate topic FOR THESE PANELISTS?"
    ```
  - Suggest names compatible with provided panelists
- [ ] Test validation with and without panelists

**Acceptance**:
- API accepts optional panelists array
- Validation considers panelist context
- Suggestions improved when panelists provided

---

### T6: Add Panelist Mini-Selector Modal
**Files**:
- `frontend/src/components/PanelistChips/PanelistMiniSelector.jsx` (new)
- `frontend/src/components/PanelistChips/PanelistMiniSelector.module.css` (new)

**Implementation**:
- [ ] Create modal component:
  - Props: `isOpen`, `onClose`, `onSelect`, `selectedIds`
  - Fetch famous panelists list (reuse existing data)
  - Search/filter functionality
  - Grid of panelist cards (avatar + name)
  - Select button on each card
  - Close button
- [ ] Integrate with PanelistChipSelector:
  - State: `isModalOpen`
  - "+" button opens modal
  - Modal onSelect adds to chips
  - Modal closes after selection
- [ ] Style modal:
  - Overlay with backdrop
  - Centered card (max-width 800px)
  - Grid layout for panelists (responsive)
  - Mobile-friendly (full screen on small devices)
- [ ] Add keyboard support (Escape to close)

**Acceptance**:
- Modal opens on "+" click
- Panelists selectable
- Chips update on selection
- Modal closeable

---

## Phase 2: Integration & Testing (4 tasks)

### T7: Integration Testing
**Files**:
- `frontend/src/pages/Home.test.jsx` (new/update)
- `frontend/src/components/TopicInput/UnifiedTopicInput.test.jsx` (new)

**Implementation**:
- [ ] Test autocomplete flow:
  - Render Home, type 3+ chars
  - Verify autocomplete appears
  - Select suggestion
  - Verify navigation with correct state
- [ ] Test manual flow without chips:
  - Type topic, click "Find Panelists"
  - Mock validation response
  - Verify navigation
- [ ] Test manual flow with chips:
  - Type topic, add 2 chips
  - Click "Find Panelists"
  - Verify chips passed to validation
- [ ] Test validation error:
  - Mock error response
  - Verify error displays
  - Verify retry works
- [ ] Test keyboard navigation:
  - Arrow keys navigate suggestions
  - Enter selects
  - Escape closes

**Acceptance**:
- All flows tested
- Coverage > 80%
- No regressions

---

### T8: Update Navigation State Handling
**Files**:
- `frontend/src/pages/PanelistSelection.jsx`

**Implementation**:
- [ ] Verify navigation state handling:
  - Read `location.state` for debateId, topic, panelists, skipValidation
  - Pre-fill chips if panelists provided
  - Lock/unlock mode based on skipValidation
- [ ] Test autocomplete selection path:
  - Navigate from Home with autocomplete selection
  - Verify chips pre-filled and locked
  - Verify "Load Debate" button shown if cache hit
- [ ] Test manual entry path:
  - Navigate from Home with validation
  - Verify chips pre-filled (if added on Home)
  - Verify chips editable
- [ ] No code changes needed (already supports both flows)

**Acceptance**:
- Both flows work correctly
- No breaking changes
- State properly consumed

---

### T9: Accessibility Audit
**Files**:
- All new components

**Implementation**:
- [ ] UnifiedTopicInput accessibility:
  - Add ARIA combobox pattern:
    ```jsx
    <input
      role="combobox"
      aria-expanded={showDropdown}
      aria-controls="autocomplete-listbox"
      aria-activedescendant={selectedIndex >= 0 ? `option-${selectedIndex}` : undefined}
    />
    <ul role="listbox" id="autocomplete-listbox">
      <li role="option" id={`option-${index}`} aria-selected={...}>
    ```
  - Error linked with aria-describedby
  - Screen reader announces suggestion count
- [ ] PanelistChipSelector accessibility:
  - Chips have aria-label: "Remove {name}"
  - "+" button has aria-label: "Add panelist"
  - Disabled state has aria-disabled
- [ ] Keyboard navigation:
  - Tab order logical
  - Focus visible
  - Escape closes modal/dropdown
- [ ] Run axe-core audit
- [ ] Test with screen reader (NVDA/JAWS/VoiceOver)

**Acceptance**:
- No critical accessibility issues
- WCAG 2.1 AA compliant
- Screen reader compatible

---

### T10: E2E Testing
**Files**:
- `frontend/cypress/e2e/unified-topic-input.cy.js` (new)

**Implementation**:
- [ ] Test autocomplete selection:
  ```javascript
  cy.visit('/');
  cy.get('[placeholder="What should they debate?"]').type('ethics');
  cy.get('[role="option"]').first().click();
  cy.url().should('include', '/select-panelists');
  cy.get('.panelistChip').should('have.length.greaterThan', 0);
  ```
- [ ] Test manual entry:
  ```javascript
  cy.visit('/');
  cy.get('[placeholder="What should they debate?"]').type('climate change');
  cy.get('button').contains('Find Panelists').click();
  cy.url().should('include', '/select-panelists');
  ```
- [ ] Test with chips:
  ```javascript
  cy.visit('/');
  cy.get('[placeholder="What should they debate?"]').type('ethics');
  cy.get('[aria-label="Add panelist"]').click();
  cy.get('.panelistCard').first().click();
  cy.get('button').contains('Find Panelists').click();
  cy.url().should('include', '/select-panelists');
  ```
- [ ] Test error handling:
  ```javascript
  cy.intercept('POST', '/validate-topic', { isRelevant: false, message: 'Invalid' });
  cy.visit('/');
  cy.get('[placeholder="What should they debate?"]').type('asdf');
  cy.get('button').contains('Find Panelists').click();
  cy.get('.error').should('contain', 'Invalid');
  ```

**Acceptance**:
- All E2E tests pass
- Flows work end-to-end
- No flaky tests

---

## Phase 3: Polish & Deploy (3 tasks)

### T11: Update Documentation
**Files**:
- `README.md`
- `specs/001-debate-generator/improvements/unified-topic-input.md`

**Implementation**:
- [ ] Update README with new flow:
  - Screenshot of unified input
  - Updated user flow diagram
  - Mention optional panelist selection
- [ ] Add migration notes:
  - What changed from previous version
  - API compatibility notes
- [ ] Update component documentation:
  - Props for new components
  - Usage examples

**Acceptance**:
- Documentation accurate
- Screenshots current
- Examples clear

---

### T12: Performance Optimization
**Files**:
- All components

**Implementation**:
- [ ] Memoize expensive components:
  ```jsx
  export default React.memo(PanelistChip, (prev, next) => 
    prev.panelist.id === next.panelist.id && prev.disabled === next.disabled
  );
  ```
- [ ] Debounce autocomplete (300ms)
- [ ] Cancel in-flight requests on unmount:
  ```javascript
  useEffect(() => {
    const controller = new AbortController();
    fetchAutocomplete(query, { signal: controller.signal });
    return () => controller.abort();
  }, [query]);
  ```
- [ ] Lazy load modal component:
  ```javascript
  const PanelistMiniSelector = lazy(() => import('./PanelistMiniSelector'));
  ```
- [ ] Run Lighthouse audit (target: 90+ performance score)

**Acceptance**:
- No unnecessary re-renders
- Requests cancelled properly
- Performance score > 90

---

### T13: Deploy to Production
**Files**:
- Deployment scripts

**Implementation**:
- [ ] Build production bundle:
  ```bash
  cd frontend && npm run build
  ```
- [ ] Verify no console errors
- [ ] Deploy backend (validate-topic enhanced):
  ```bash
  ./deploy.sh
  ```
- [ ] Deploy frontend:
  ```bash
  docker build -t debate-frontend .
  docker-compose up -d
  ```
- [ ] Smoke test on production:
  - Test autocomplete selection
  - Test manual entry
  - Test error handling
- [ ] Monitor for errors (first 24 hours)

**Acceptance**:
- Deployment successful
- No critical errors
- User flows working

---

## Summary

### Task Breakdown
- **Phase 1**: 6 tasks (component refactoring)
- **Phase 2**: 4 tasks (integration & testing)
- **Phase 3**: 3 tasks (polish & deploy)
- **Total**: 13 tasks

### Estimated Effort
- Phase 1: 12-16 hours
- Phase 2: 6-8 hours
- Phase 3: 4-6 hours
- **Total**: 22-30 hours

### Dependencies
```
T1 (Chips) → T4 (Render)
T2 (Unified Input) → T4 (Render)
T1, T2 → T3 (State)
T3 → T4 (Render)
T4 → T5 (API)
T1 → T6 (Modal)
T1-T6 → T7 (Integration)
T7 → T8 (Navigation)
T8 → T9 (A11y)
T9 → T10 (E2E)
T10 → T11 (Docs)
T11 → T12 (Perf)
T12 → T13 (Deploy)
```

### Critical Path
T1 → T2 → T3 → T4 → T5 → T7 → T8 → T9 → T10 → T12 → T13
