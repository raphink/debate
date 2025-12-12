# Unified Topic Input - UI Consolidation

## Overview

Merge the autocomplete search and manual topic entry into a single unified input field that supports both workflows seamlessly. This eliminates UI duplication and creates a more intuitive experience where autocomplete suggestions appear naturally as users type.

## User Story

**As a** user creating a debate  
**I want** a single input field that suggests previous topics while I type  
**So that** I can either reuse an existing debate or create a new one without switching between different input modes

## Current State

### Home Page Layout
```
┌────────────────────────────────────────────┐
│  Search Previous Debates                   │
│  ┌──────────────────────────────────────┐  │
│  │ [Autocomplete Input]                 │  │
│  │   ↓ Dropdown with suggestions        │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  ────────────── or ──────────────          │
│                                            │
│  Create New Debate                         │
│  ┌──────────────────────────────────────┐  │
│  │ [Manual Topic Input]                 │  │
│  └──────────────────────────────────────┘  │
│  [Submit/Validate Button]                  │
└────────────────────────────────────────────┘
```

### Issues
1. **Redundant inputs**: Two text fields for the same purpose (entering a topic)
2. **Cognitive overhead**: Users must choose between "search" and "create" modes
3. **Wasted space**: "or" divider and dual input layout
4. **Unclear affordance**: When to use which input is ambiguous

## Proposed State

### Unified Home Page Layout
```
┌────────────────────────────────────────────┐
│  Enter Debate Topic                        │
│  ┌──────────────────────────────────────┐  │
│  │ What should they debate?             │  │
│  │   ↓ Autocomplete suggestions (if any)│  │
│  └──────────────────────────────────────┘  │
│                                            │
│  Select Panelists (Optional)               │
│  ┌──────────────────────────────────────┐  │
│  │ [Chip: Einstein] [Chip: Curie] [+]   │  │
│  └──────────────────────────────────────┘  │
│                                            │
│  [Find Panelists →]                        │
└────────────────────────────────────────────┘
```

### Benefits
1. **Unified experience**: Single mental model for topic entry
2. **Progressive disclosure**: Autocomplete appears only when relevant
3. **Streamlined UI**: More space for helpful content
4. **Clear call-to-action**: "Find Panelists" button clearly indicates next step
5. **Optional optimization**: Users can pre-select panelists to skip validation

## Detailed Behavior

### Typing Flow

```
State 1: Empty Input
┌────────────────────────┐
│ What should they       │
│ debate?                │
└────────────────────────┘

User types: "eth"

State 2: Autocomplete Active (3+ chars)
┌────────────────────────┐
│ eth|                   │
├────────────────────────┤
│ ▼ Suggestions          │
│ ○ Ethics of AI...      │
│ ○ Ethical implications │
└────────────────────────┘

Option A: User selects suggestion
→ Navigate to /select-panelists
  with pre-filled topic + panelists

Option B: User continues typing
User types: "ics of genetic engineering"

State 3: Manual Entry (no match or ignored)
┌────────────────────────┐
│ ethics of genetic      │
│ engineering|           │
└────────────────────────┘
(Autocomplete dismissed, no matches)

User clicks "Find Panelists"
→ Validate topic
→ Navigate to /select-panelists
  with validated topic (empty panelists)
```

### Autocomplete Dismissal Rules

1. **User selects suggestion**: Dropdown closes, navigate away
2. **No matches found**: Dropdown disappears automatically
3. **Click outside**: Dropdown closes
4. **Escape key**: Dropdown closes
5. **Input loses focus**: Dropdown closes
6. **Continue typing past matches**: Dropdown stays (updates results)

### Early Panelist Selection

```
User types topic: "Climate change"
Autocomplete shows, user ignores
User clicks [+] button below input

┌────────────────────────────────────┐
│ Climate change impacts             │
├────────────────────────────────────┤
│ [Einstein] [Curie] [Hawking] [+]   │
└────────────────────────────────────┘

User clicks "Find Panelists"
→ Validate topic with panelists
→ Navigate to /select-panelists
  with topic + pre-selected panelists
  (allows editing before generation)
```

## Technical Changes

### Component Changes

#### Remove
- Separate autocomplete search input component
- "or" divider element
- Dual input state management

#### Modify
- `Home.jsx`: Single controlled input with autocomplete hook
- `TopicInput` component: Merge validation and autocomplete features
- `useTopicValidation` hook: Integrate with autocomplete state

#### Add
- Inline panelist chip selector on Home page
- `usePanelistSelection` hook on Home (reuse from PanelistSelection)
- "Find Panelists" button with topic validation

### Data Flow

```
Home.jsx State:
├─ topic: string (controlled input value)
├─ selectedPanelists: Panelist[] (optional chips)
├─ autocomplete: { suggestions, loading, error }
└─ validation: { isValid, error } (on button click)

User Actions:
├─ Type → Update topic → Trigger autocomplete
├─ Select suggestion → Navigate with {debateId, topic, panelists}
├─ Add panelist chip → Update selectedPanelists
└─ Click "Find Panelists" → Validate → Navigate with {topic, panelists?}
```

## Acceptance Criteria

### AC1: Unified Input Behavior
- ✓ Single input field accepts topic text
- ✓ Autocomplete dropdown appears on 3+ characters
- ✓ Suggestions update as user types
- ✓ No separate "search" vs "create" mode

### AC2: Autocomplete Selection
- ✓ Clicking suggestion navigates to panelist selection
- ✓ Pre-fills topic and panelists from selected debate
- ✓ Skips topic validation (already validated)

### AC3: Manual Entry Flow
- ✓ User can ignore autocomplete and continue typing
- ✓ Autocomplete dismisses when no matches or user clicks outside
- ✓ "Find Panelists" button triggers topic validation
- ✓ Validation errors display inline (same as current)

### AC4: Early Panelist Selection
- ✓ Optional chip selector appears below topic input
- ✓ Users can add/remove panelists before validation
- ✓ Pre-selected panelists pass to validation endpoint
- ✓ Navigate to PanelistSelection with chips pre-filled

### AC5: Keyboard Navigation
- ✓ Arrow keys navigate autocomplete suggestions
- ✓ Enter on suggestion selects it
- ✓ Enter in empty dropdown validates topic
- ✓ Escape closes autocomplete dropdown

### AC6: Responsive Design
- ✓ Single input scales properly on mobile
- ✓ Autocomplete dropdown adapts to screen size
- ✓ Chip selector wraps on narrow screens

## Edge Cases

1. **Empty topic submission**: Disable "Find Panelists" until topic.length >= 3
2. **Autocomplete during validation**: Disable input during API call
3. **Network failure**: Show error, allow retry
4. **Stale suggestions**: Clear dropdown when input < 3 chars
5. **Panelists without topic**: Disable chips until valid topic entered
6. **Rapid typing**: Debounce autocomplete (300ms) to reduce API calls

## Migration Notes

### Removed Code
- `frontend/src/components/TopicInput/TopicAutocompleteInput.jsx` (if separate)
- Autocomplete-specific sections in `Home.jsx`
- "or" divider styling

### Modified Files
- `frontend/src/pages/Home.jsx`: Consolidate inputs
- `frontend/src/components/TopicInput/TopicInput.jsx`: Add autocomplete
- `frontend/src/pages/Home.module.css`: Update layout

### Backward Compatibility
- No API changes (autocomplete endpoint unchanged)
- Navigation state structure unchanged
- PanelistSelection page works with both flows

## Design Considerations

### Visual Hierarchy
```
Priority 1: Topic input (primary action)
Priority 2: Autocomplete dropdown (contextual)
Priority 3: Panelist chips (optional optimization)
Priority 4: "Find Panelists" button (call-to-action)
```

### Accessibility
- Autocomplete uses ARIA combobox pattern
- Screen reader announces suggestion count
- Keyboard navigation fully supported
- Focus management on selection/dismissal

### Performance
- Debounce autocomplete queries (300ms)
- Cancel in-flight requests on new input
- Limit dropdown to 10 suggestions
- Memoize chip components to prevent re-renders

## Future Enhancements

1. **Recent topics**: Show user's 5 most recent topics below input (no autocomplete)
2. **Topic suggestions**: AI-powered topic ideas based on current events
3. **Panelist recommendations**: Suggest panelists based on topic keywords
4. **Keyboard shortcuts**: Cmd/Ctrl+K to focus topic input
5. **Voice input**: Speech-to-text for topic entry
