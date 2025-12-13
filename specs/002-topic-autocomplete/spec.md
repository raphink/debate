# Feature Specification: Topic Discovery via History Integration (US6)

## Clarifications (Session 2025-12-13)

- Q: How should the system handle concurrent writes to Firestore when multiple users generate debates with the same topic at the same time? → A: Allow duplicates - multiple debates on the same topic are acceptable (last write wins, Firestore default)
- Q: When a user selects a topic from autocomplete and lands on PanelistSelection with pre-filled panelists, what should the "Generate Debate" button say to clarify the cache behavior? → A: "View Debate" for cached (no changes) with alternate button to change the panelists (in which case it will lead to a new debate as if generated from scratch)
- Q: When autocomplete returns multiple debates with identical topics, how should they be differentiated in the dropdown? → A: Show panelist avatars + generation date. Two debates on the same topic can be distinguished because they have different panelists.
- Q: Should the autocomplete query sanitize user input to prevent injection attacks, or can it safely pass raw input to Firestore? → A: All input should be sanitized
- Q: Should the autocomplete dropdown support keyboard navigation (arrow keys, Enter to select, Escape to close)? → A: Yes - standard accessibility practice
- **CRITICAL UX CLARIFICATION**: There is ONE single input field for topics, not two separate flows. As users type, autocomplete suggestions appear (if matches exist). Users can either: (A) select an existing debate from dropdown → navigate to PanelistSelection with pre-filled panelists, OR (B) ignore autocomplete and click "Find Panelists" button to trigger Claude validation (normal US1 flow with optional suggested panelist names). Autocomplete is optional enhancement, not a separate mode.

## Overview

**Feature ID**: 002-topic-autocomplete  
**Priority**: P3 (Quality-of-life enhancement)  
**Status**: In Development (previously deferred from MVP)  
**Dependencies**: Requires US5 (Debate Caching & Sharing) - Firestore debates collection must be populated

## Summary

Users see autocomplete suggestions of previous debates as they type in the single topic input field. This streamlines topic discovery by combining history browsing with topic entry in a unified interface. Users can select from suggestions to reuse debates or ignore autocomplete and proceed with normal Claude validation.

## User Story

**As a** user creating a new debate,  
**I want** to see autocomplete suggestions of previous topics as I type in the topic input field,  
**So that** I can optionally select from history to reuse debates or discover similar topics, while still being able to proceed with manual topic entry at any time.

## Acceptance Scenarios

1. **Given** user types in topic input field, **When** user has typed ≥3 characters, **Then** system displays dropdown showing up to 10 matching previous topics ordered by recency

2. **Given** topic autocomplete dropdown is displayed, **When** user views a suggestion, **Then** entry shows topic text, panelist avatars (circular thumbnails), panelist count, and generation date (e.g., "3 panelists — Dec 10, 2025")

3. **Given** user types in topic input field, **When** input matches previous topics, **Then** matching topics are highlighted/narrowed in real-time as user continues typing

4. **Given** user selects a topic from autocomplete dropdown, **When** topic is selected, **Then** system navigates directly to panelist selection with original panelists pre-selected (skipping Claude validation)

5. **Given** user types topic and autocomplete appears, **When** user ignores dropdown and clicks "Find Panelists" button, **Then** system proceeds with normal Claude validation flow (US1) with optional suggested panelist names

6. **Given** user is on panelist selection with pre-filled panelists from history, **When** user makes no changes to panelist list, **Then** system detects cache hit and shows "View Debate" button (bypasses generation)

6. **Given** user is on panelist selection with pre-filled panelists from history, **When** user makes no changes to panelist list, **Then** system detects cache hit and shows "View Debate" button (bypasses generation)

7. **Given** user is on panelist selection with pre-filled panelists, **When** user clicks "Modify Panelists" button, **Then** system allows editing chips, changes main button to "Generate New Debate", and generates new debate from scratch if user proceeds

8. **Given** topic autocomplete is loading, **When** API call is in progress, **Then** system shows subtle loading indicator without blocking typing or "Find Panelists" button

9. **Given** no previous topics match user input or Firestore fails, **When** user types, **Then** autocomplete dropdown is hidden and user can click "Find Panelists" normally (graceful degradation)

10. **Given** user selects historical topic with modified panelists, **When** user proceeds to generate debate, **Then** system generates new debate (no cache hit) and saves as new debate instance

## Functional Requirements

### Backend API

- **FR-001**: System MUST provide autocomplete-topics Cloud Function endpoint: GET /api/autocomplete-topics?q={query}&limit=10
- **FR-002**: System MUST sanitize query parameter (strip HTML tags, special characters) before querying Firestore
- **FR-003**: System MUST query Firestore debates collection by topic substring (case-insensitive)
- **FR-004**: System MUST order results by createdAt timestamp descending (most recent first)
- **FR-005**: System MUST limit results to max 10 debates per query
- **FR-006**: System MUST return full debate metadata: {id, topic, panelists: [{id, name, slug, avatarUrl}], panelistCount, createdAt}
- **FR-007**: System MUST add CORS headers to allow cross-origin requests from frontend domain
- **FR-008**: System MUST handle Firestore read failures gracefully (return empty array, log error)

### Frontend UX

- **FR-009**: System MUST trigger autocomplete search when user types ≥3 characters in topic input
- **FR-010**: System MUST debounce autocomplete requests with 300ms delay
- **FR-011**: System MUST display dropdown below topic input field with matching topics
- **FR-012**: System MUST show topic text, panelist avatars, panelist count badge, and generation date for each suggestion (e.g., "3 panelists — Dec 10, 2025")
- **FR-013**: System MUST narrow suggestions in real-time as user continues typing
- **FR-014**: System MUST show subtle loading indicator when API response takes >300ms
- **FR-015**: System MUST hide dropdown if API fails or returns empty results (user can still click "Find Panelists" normally)
- **FR-016**: System MUST NOT disable or block "Find Panelists" button when autocomplete is loading or unavailable

### Panelist Pre-filling

- **FR-017**: System MUST skip Claude validation ONLY when user selects topic from autocomplete dropdown (clicking "Find Panelists" always triggers Claude validation)
- **FR-018**: System MUST navigate to PanelistSelection page with panelists pre-selected from historical debate when user selects from dropdown
- **FR-019**: System MUST pass original debate ID via navigation state for cache detection
- **FR-020**: System MUST display "View Debate" button when cache hit detected (topic + panelists unchanged)
- **FR-021**: System MUST display alternate "Modify Panelists" button when panelists are pre-filled from history
- **FR-022**: System MUST change main button to "Generate New Debate" when user clicks "Modify Panelists" and unlock panelist chips for editing

### Cache Detection

- **FR-023**: System MUST implement deep comparison utility to detect exact match: topic text + panelist array (id, name order-independent)
- **FR-024**: System MUST show "View Debate" button when cache hit detected (clicking redirects to /d/{uuid})
- **FR-025**: System MUST skip debate generation if user clicks "View Debate" button
- **FR-026**: System MUST show visual indicator "Saved debate from [date]" when cache hit detected
- **FR-027**: System MUST generate new debate from scratch (assign new UUID, full Claude streaming) if user modifies panelist list and clicks "Generate New Debate"
- **FR-028**: System MUST show visual indicator "Creating new debate" when user proceeds after modifying panelists

## Non-Functional Requirements

- **NFR-001**: Autocomplete API response time SHOULD be <500ms for typical queries
- **NFR-002**: Debouncing MUST prevent excessive API calls (max 1 request per 300ms)
- **NFR-003**: Dropdown UI MUST be keyboard-accessible (arrow keys, Enter to select, Escape to close)
- **NFR-004**: System MUST handle Firestore quota limits gracefully (read throttling, error messages)
- **NFR-005**: Cache detection comparison MUST complete in <50ms (synchronous, client-side)

## Edge Cases & Error Handling

- **EC-001**: Empty Firestore (no debates): Hide dropdown, allow normal topic entry
- **EC-002**: Firestore read failure: Log error, hide dropdown, continue with normal flow
- **EC-003**: Network timeout (>5s): Cancel autocomplete request, hide dropdown
- **EC-004**: User types <3 characters: Do not trigger autocomplete
- **EC-005**: User clears input field after dropdown displayed: Hide dropdown
- **EC-006**: User clicks outside dropdown: Close dropdown without selection
- **EC-007**: Duplicate topics in Firestore: Show all instances differentiated by panelist avatars and generation dates (multiple debates on same topic are acceptable)
- **EC-008**: User modifies pre-filled panelists then restores original: Treat as cache hit (re-run comparison)

## Out of Scope

- Topic search/filtering beyond simple substring matching
- Advanced fuzzy matching or typo correction for topics
- Topic categorization or tagging
- Autocomplete for panelist names (covered by US7, separate feature)
- User-specific debate history filtering
- Pagination for >10 autocomplete results

## Data Model

### AutocompleteTopicsRequest
```typescript
{
  q: string;        // Query string (min 3 chars)
  limit?: number;   // Max results (default 10, max 10)
}
```

### AutocompleteTopicsResponse
```typescript
{
  debates: Array<{
    id: string;              // UUID
    topic: string;           // Original topic text
    panelistCount: number;   // Number of panelists
    panelists: Array<{       // Full panelist data for pre-filling
      id: string;
      name: string;
      slug: string;
    }>;
    createdAt: Timestamp;    // ISO 8601
  }>;
}
```

### Cache Detection Comparison
```typescript
{
  topic: string;                 // Exact text match
  panelists: Array<{             // Order-independent deep comparison
    id: string;
    name: string;
  }>;
}
```

## Success Metrics

- Autocomplete adoption rate: % of debates created via autocomplete vs manual entry
- Cache hit rate: % of autocomplete selections that result in cache hit (no re-generation)
- API performance: p95 response time <500ms
- User satisfaction: qualitative feedback on workflow improvement

## Testing Strategy

### Unit Tests
- AutocompleteTopics Cloud Function: query parsing, limit enforcement
- Cache detection utility: deep comparison, order independence
- Debounce hook: timing, cancellation

### Integration Tests
- Firestore query correctness: substring matching, case-insensitivity, ordering
- CORS headers validation
- Error handling: empty results, Firestore failures

### End-to-End Tests
1. Generate debate → return home → type 3+ chars → verify autocomplete appears
2. Select topic from dropdown → verify panelists pre-filled on PanelistSelection page
3. Pre-filled panelists unchanged → verify redirect to /d/{uuid} (cache hit)
4. Click "Modify Panelists" → change list → verify new debate generated
5. Firestore disabled → verify graceful degradation to manual topic entry

## Dependencies & Prerequisites

- US5 (Debate Caching & Sharing) fully implemented
- Firestore debates collection populated with sample debates
- list-debates Cloud Function operational (reuses Firestore query logic)
- Navigation state management in React Router functional

## Open Questions

None at this time. Proceed to clarification phase.
