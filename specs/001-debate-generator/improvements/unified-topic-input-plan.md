# Unified Topic Input - Technical Implementation Plan

## Architecture Overview

### Before: Dual Input Architecture
```
Home.jsx
├─ AutocompleteSearchInput
│   └─ useTopicAutocomplete hook
│       └─ TopicAutocompleteDropdown
└─ ManualTopicInput
    └─ useTopicValidation hook
        └─ ValidationResult

Separate state, separate handlers, "or" divider
```

### After: Unified Input Architecture
```
Home.jsx
└─ UnifiedTopicInput
    ├─ useTopicAutocomplete hook
    │   └─ TopicAutocompleteDropdown (overlay)
    ├─ useTopicValidation hook (on submit)
    ├─ usePanelistSelection hook (optional chips)
    └─ PanelistChipSelector (inline)

Single state, merged handlers, progressive disclosure
```

## Component Structure

### New Component Hierarchy

```jsx
Home.jsx
├─ UnifiedTopicInput (new wrapper)
│   ├─ Input field (controlled)
│   ├─ TopicAutocompleteDropdown (conditional)
│   └─ ValidationError (conditional)
├─ PanelistChipSelector (moved from PanelistSelection)
│   └─ PanelistChip[] (add/remove)
└─ Button: "Find Panelists"
```

### Component Responsibilities

#### `UnifiedTopicInput.jsx` (New)
**Purpose**: Single source of truth for topic entry with autocomplete

**State**:
- `topic`: string (controlled input value)
- `showAutocomplete`: boolean (dropdown visibility)
- `selectedIndex`: number (keyboard navigation)

**Props**:
- `value`: string (controlled by parent)
- `onChange`: (topic: string) => void
- `onSelectSuggestion`: (debate: Debate) => void
- `onSubmit`: (topic: string) => void
- `disabled`: boolean
- `error`: string | null

**Behavior**:
- Renders input field with autocomplete overlay
- Manages keyboard navigation (arrows, enter, escape)
- Handles click-outside to dismiss autocomplete
- Triggers validation error display

**Implementation**:
```jsx
const UnifiedTopicInput = ({ value, onChange, onSelectSuggestion, error, disabled }) => {
  const { suggestions, loading } = useTopicAutocomplete(value);
  const [selectedIndex, setSelectedIndex] = useState(-1);
  const [showDropdown, setShowDropdown] = useState(false);
  
  useEffect(() => {
    setShowDropdown(value.length >= 3 && suggestions.length > 0);
  }, [value, suggestions]);
  
  const handleKeyDown = (e) => {
    if (!showDropdown) return;
    
    switch (e.key) {
      case 'ArrowDown':
        setSelectedIndex(i => Math.min(i + 1, suggestions.length - 1));
        e.preventDefault();
        break;
      case 'ArrowUp':
        setSelectedIndex(i => Math.max(i - 1, -1));
        e.preventDefault();
        break;
      case 'Enter':
        if (selectedIndex >= 0) {
          onSelectSuggestion(suggestions[selectedIndex]);
        }
        e.preventDefault();
        break;
      case 'Escape':
        setShowDropdown(false);
        break;
    }
  };
  
  return (
    <div className={styles.container}>
      <input
        value={value}
        onChange={(e) => onChange(e.target.value)}
        onKeyDown={handleKeyDown}
        placeholder="What should they debate?"
        disabled={disabled}
      />
      {showDropdown && (
        <TopicAutocompleteDropdown
          suggestions={suggestions}
          loading={loading}
          selectedIndex={selectedIndex}
          onSelect={onSelectSuggestion}
        />
      )}
      {error && <div className={styles.error}>{error}</div>}
    </div>
  );
};
```

#### `PanelistChipSelector.jsx` (Moved/Refactored)
**Purpose**: Inline panelist selection before validation

**State**:
- `selectedPanelists`: Panelist[]
- `searchQuery`: string (for adding new chips)

**Props**:
- `value`: Panelist[] (controlled)
- `onChange`: (panelists: Panelist[]) => void
- `disabled`: boolean
- `max`: number (default 10)

**Behavior**:
- Renders chips for selected panelists
- "+" button opens mini-selector or autocomplete
- Remove chips with X button
- Validates max panelist count

#### `Home.jsx` (Refactored)
**Purpose**: Orchestrate unified flow

**State**:
```jsx
const [topic, setTopic] = useState('');
const [selectedPanelists, setSelectedPanelists] = useState([]);
const [validationError, setValidationError] = useState(null);
const [isValidating, setIsValidating] = useState(false);
```

**Handlers**:
```jsx
// Handle autocomplete selection (Option A)
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

// Handle manual submission (Option B)
const handleFindPanelists = async () => {
  setIsValidating(true);
  setValidationError(null);
  
  try {
    // Validate topic (with optional panelists context)
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

**Render**:
```jsx
return (
  <div className={styles.container}>
    <h1>Generate a Debate</h1>
    
    <UnifiedTopicInput
      value={topic}
      onChange={setTopic}
      onSelectSuggestion={handleAutocompleteSelect}
      error={validationError}
      disabled={isValidating}
    />
    
    <PanelistChipSelector
      value={selectedPanelists}
      onChange={setSelectedPanelists}
      disabled={isValidating || topic.length < 3}
      max={10}
    />
    
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

## Data Flow Diagrams

### Flow 1: Autocomplete Selection
```
User types "eth" in UnifiedTopicInput
    ↓
useTopicAutocomplete hook debounces (300ms)
    ↓
API call: GET /autocomplete-topics?q=eth
    ↓
TopicAutocompleteDropdown renders suggestions
    ↓
User clicks "Ethics of AI in healthcare"
    ↓
handleAutocompleteSelect fires
    ↓
navigate('/select-panelists', {
  state: {
    debateId: "uuid-123",
    topic: "Ethics of AI in healthcare",
    panelists: [kant, thomson, foot],
    skipValidation: true
  }
})
    ↓
PanelistSelection.jsx pre-fills chips, locks mode
    ↓
User clicks "Generate Debate" or "Load Debate"
```

### Flow 2: Manual Entry (No Selection)
```
User types "climate change impacts"
    ↓
Autocomplete shows suggestions (ignored by user)
    ↓
User continues typing full topic
    ↓
Autocomplete dismisses (no matches or blur)
    ↓
User optionally adds panelist chips (Einstein, Curie)
    ↓
User clicks "Find Panelists"
    ↓
handleFindPanelists fires
    ↓
API call: POST /validate-topic
  body: { topic, panelists: ["einstein", "curie"] }
    ↓
Validation response: { isRelevant: true, ... }
    ↓
navigate('/select-panelists', {
  state: {
    topic: "climate change impacts",
    panelists: [einstein, curie],
    skipValidation: false
  }
})
    ↓
PanelistSelection.jsx shows chips (editable)
    ↓
User reviews/edits panelists, clicks "Generate Debate"
```

### Flow 3: Validation Error
```
User types "asdfghjkl"
    ↓
User clicks "Find Panelists"
    ↓
API call: POST /validate-topic
    ↓
Validation response: { isRelevant: false, message: "..." }
    ↓
setValidationError(message)
    ↓
UnifiedTopicInput displays error inline
    ↓
User edits topic, tries again
```

## State Management

### Home.jsx State
```typescript
interface HomeState {
  topic: string;                    // Controlled input value
  selectedPanelists: Panelist[];    // Optional pre-selected chips
  validationError: string | null;   // Inline error message
  isValidating: boolean;            // Button loading state
}

interface Panelist {
  id: string;
  name: string;
  slug: string;
  avatarUrl: string;
  tagline: string;
  biography: string;
}
```

### Navigation State (unchanged)
```typescript
interface NavigationState {
  debateId?: string;          // From autocomplete selection
  topic: string;              // Always required
  panelists: Panelist[];      // Empty or pre-filled
  skipValidation: boolean;    // true for autocomplete, false for manual
}
```

## API Integration

### Autocomplete Endpoint (No Changes)
```http
GET /autocomplete-topics?q=eth&limit=10

Response:
{
  "debates": [
    {
      "id": "uuid",
      "topic": "Ethics of AI in healthcare",
      "panelists": [...],
      "panelistCount": 3,
      "createdAt": "2025-01-14T10:15:00Z"
    }
  ]
}
```

### Validation Endpoint (Enhanced)
```http
POST /validate-topic

Body:
{
  "topic": "climate change impacts",
  "panelists": ["einstein", "curie"]  // NEW: optional context
}

Response:
{
  "isRelevant": true,
  "topic": "climate change impacts",
  "validationMessage": "Great topic!",
  "suggestedNames": ["Greta Thunberg", "Al Gore", ...]  // Enhanced by panelists
}
```

## CSS/Styling Changes

### Layout Shift
**Before**:
```css
.home {
  display: flex;
  flex-direction: column;
  gap: 2rem;
}

.autocompleteSection { ... }
.divider { ... }  /* Remove */
.manualSection { ... }
```

**After**:
```css
.home {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  max-width: 600px;
  margin: 0 auto;
}

.unifiedInput {
  position: relative;
  width: 100%;
}

.panelistChips {
  min-height: 60px;
  padding: 0.75rem;
  border: 1px dashed #ccc;
  border-radius: 8px;
}

.submitButton {
  width: 100%;
  padding: 1rem;
  font-size: 1.1rem;
  font-weight: 600;
}
```

## File Changes Summary

### New Files
- `frontend/src/components/TopicInput/UnifiedTopicInput.jsx`
- `frontend/src/components/PanelistChips/PanelistChipSelector.jsx`
- `frontend/src/components/PanelistChips/PanelistChip.jsx`
- `frontend/src/components/PanelistChips/PanelistChipSelector.module.css`

### Modified Files
- `frontend/src/pages/Home.jsx` - Refactor to use unified input
- `frontend/src/pages/Home.module.css` - Update layout
- `frontend/src/hooks/usePanelistSelection.js` - Extract for reuse
- `frontend/src/services/topicService.js` - Add panelist context to validation

### Removed Files
- (None - old components left for reference, can delete later)

### Affected Files (Indirect)
- `frontend/src/pages/PanelistSelection.jsx` - No changes (still works)
- `frontend/src/components/TopicInput/TopicAutocompleteDropdown.jsx` - Reused as-is

## Testing Strategy

### Unit Tests
1. **UnifiedTopicInput**:
   - Renders input with placeholder
   - Shows autocomplete on 3+ chars
   - Hides autocomplete on < 3 chars
   - Handles keyboard navigation
   - Calls onSelectSuggestion when item clicked
   - Displays validation error

2. **PanelistChipSelector**:
   - Renders empty state with "+" button
   - Adds chip on selection
   - Removes chip on X click
   - Enforces max limit
   - Disables when parent disabled

3. **Home.jsx**:
   - Autocomplete selection navigates correctly
   - Manual submission validates topic
   - Validation error displays inline
   - Button disabled when topic < 3 chars

### Integration Tests
1. **Autocomplete flow**:
   - Type query → see suggestions → select → navigate with pre-filled data

2. **Manual flow**:
   - Type topic → add chips → submit → validate → navigate with chips

3. **Error handling**:
   - Invalid topic → show error → edit → retry

### E2E Tests
1. Complete autocomplete selection flow
2. Complete manual entry flow with chips
3. Complete manual entry flow without chips
4. Validation error recovery

## Migration Plan

### Phase 1: Feature Flag (Optional)
```jsx
const USE_UNIFIED_INPUT = process.env.REACT_APP_UNIFIED_INPUT === 'true';

return USE_UNIFIED_INPUT ? <UnifiedHome /> : <LegacyHome />;
```

### Phase 2: Gradual Rollout
1. Deploy unified input behind feature flag
2. Test with 10% of users
3. Monitor error rates, conversion metrics
4. Increase to 50%, then 100%

### Phase 3: Cleanup
1. Remove feature flag
2. Delete legacy components
3. Update documentation

## Performance Considerations

### Optimizations
1. **Debouncing**: 300ms for autocomplete queries
2. **Memoization**: Chip components with React.memo
3. **Lazy loading**: Dropdown renders only when visible
4. **Request cancellation**: Abort in-flight autocomplete on new input

### Bundle Size
- New components: ~8KB (gzipped)
- Removed components: ~6KB (gzipped)
- Net increase: ~2KB

## Accessibility Checklist

- ✓ Autocomplete uses ARIA combobox pattern
- ✓ Screen reader announces suggestion count
- ✓ Keyboard navigation (arrows, enter, escape)
- ✓ Focus management on selection/dismissal
- ✓ Error messages associated with input (aria-describedby)
- ✓ Chip removal announces to screen readers
- ✓ Button disabled states have aria-disabled

## Rollback Plan

If critical issues arise:
1. Toggle feature flag to false
2. Revert to legacy Home.jsx
3. Investigate and fix issues
4. Redeploy with fixes

Rollback triggers:
- Error rate > 5%
- User complaints > 10/day
- Conversion rate drop > 20%
