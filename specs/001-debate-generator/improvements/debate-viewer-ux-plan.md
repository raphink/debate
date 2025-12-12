# Technical Plan: Debate Viewer UX Improvements

**Feature**: Enhanced Debate Viewer Interface  
**Spec**: [debate-viewer-ux.md](debate-viewer-ux.md)  
**Created**: 2025-12-12

## Architecture Overview

### Component Hierarchy

```
DebateViewer (viewer mode)          DebateGeneration (live mode)
       ↓                                     ↓
   DebateView (shared component with mode prop)
       ↓
   ├── DebateBubble (messages)
   ├── TypingIndicator (live only)
   ├── AutoScroll toggle (live only)
   └── New Debate button (viewer only)
```

### Design Approach

**Mode Detection Strategy:**
- Add new optional prop to DebateView: `mode: 'viewer' | 'generation'`
- Default to `'generation'` for backward compatibility
- DebateViewer passes `mode="viewer"`
- DebateGeneration passes `mode="generation"` (or omits, uses default)

**Conditional Rendering:**
- Auto-scroll toggle: `{mode === 'generation' && <AutoScrollToggle />}`
- New Debate button: `{mode === 'viewer' && <NewDebateButton />}`

## File Changes

### 1. Frontend Components

**File: `frontend/src/components/DebateView/DebateView.jsx`**
- Add `mode` prop with PropTypes validation
- Add conditional rendering for auto-scroll toggle
- Add "New Debate" button component at bottom (after messagesEndRef)
- Import `useNavigate` from react-router-dom
- Style button to match existing gradient design system

**File: `frontend/src/components/DebateView/DebateView.module.css`**
- Add `.newDebateButton` class with gradient styling
- Add `.newDebateButtonContainer` for bottom spacing/alignment
- Ensure mobile responsiveness

**File: `frontend/src/pages/DebateViewer.jsx`**
- Pass `mode="viewer"` prop to DebateView

**File: `frontend/src/pages/DebateGeneration.jsx`**
- Pass `mode="generation"` prop to DebateView (optional, uses default)

## Implementation Details

### DebateView Mode Prop

```jsx
const DebateView = ({ 
  messages, 
  panelists, 
  isStreaming, 
  currentPanelistId, 
  debateId, 
  isComplete,
  mode = 'generation' // default for backward compat
}) => {
  const navigate = useNavigate();
  
  // ... existing code ...
  
  return (
    <div className={styles.container}>
      {/* Existing header, messages, etc. */}
      
      {/* Auto-scroll toggle - only in generation mode */}
      {mode === 'generation' && (
        <div className={styles.autoScrollToggle}>
          {/* existing toggle code */}
        </div>
      )}
      
      {/* Messages */}
      <div ref={messagesEndRef} />
      
      {/* New Debate button - only in viewer mode */}
      {mode === 'viewer' && (
        <div className={styles.newDebateButtonContainer}>
          <button 
            className={styles.newDebateButton}
            onClick={() => navigate('/')}
          >
            Create New Debate
          </button>
        </div>
      )}
    </div>
  );
};

DebateView.propTypes = {
  // ... existing props ...
  mode: PropTypes.oneOf(['viewer', 'generation'])
};
```

### CSS Styling

```css
.newDebateButtonContainer {
  display: flex;
  justify-content: center;
  padding: 2rem 1rem;
  margin-top: 2rem;
  border-top: 1px solid rgba(255, 255, 255, 0.1);
}

.newDebateButton {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  font-size: 1rem;
  font-weight: 600;
  padding: 0.75rem 2rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.3s ease;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
}

.newDebateButton:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 12px rgba(102, 126, 234, 0.4);
}

.newDebateButton:active {
  transform: translateY(0);
}

@media (max-width: 768px) {
  .newDebateButton {
    font-size: 0.9rem;
    padding: 0.6rem 1.5rem;
  }
}
```

## Testing Strategy

### Manual Testing
1. Navigate to `/d/{uuid}` → verify auto-scroll toggle hidden
2. Verify "New Debate" button visible at bottom
3. Click "New Debate" → verify navigates to `/`
4. Generate new debate → verify auto-scroll toggle visible
5. Verify "New Debate" button NOT visible during generation
6. Mobile: Verify button styling and responsiveness

### Accessibility
- Button has proper semantic HTML (`<button>`)
- Button text is descriptive ("Create New Debate")
- Keyboard navigation works (Tab + Enter)
- Focus states visible

## Dependencies

- No new NPM packages required
- Uses existing `useNavigate` from react-router-dom
- Uses existing CSS gradient patterns

## Rollout Plan

1. Update DebateView component with mode prop
2. Update DebateViewer to pass `mode="viewer"`
3. Update CSS with new button styles
4. Manual QA testing
5. Deploy to production

## Backward Compatibility

- Default `mode="generation"` ensures existing code works unchanged
- No breaking changes to DebateView API
- Existing consumers (DebateGeneration) continue to work without modification
