# Feature Specification: Debate History with Reactive Search

**Feature ID**: US7-debate-history  
**Date**: 2025-12-13  
**Status**: Draft  
**Owner**: Product Team

## Overview

Enable users to browse all previously generated debates and quickly find specific conversations using real-time search. This feature provides a centralized view of debate history with instant filtering as users type, making it easy to rediscover and share past theological/philosophical discussions.

## User Stories

### US7.1: Browse Debate History (P1)
**As a** user who has generated debates  
**I want to** see a list of all my previous debates  
**So that** I can revisit interesting conversations and share them with others

**Acceptance Criteria**:
- View displays all debates sorted by most recent first
- Each debate card shows: topic, panelist names (max 3 visible + count), timestamp
- Clicking a debate card navigates to the full debate viewer
- Debates load progressively (pagination) if more than 20 exist
- Empty state message shown when no debates exist yet

**Test Scenario**:
```
GIVEN I have generated 5 debates in the past week
WHEN I navigate to the debate history page
THEN I see all 5 debates sorted newest to oldest
AND each card displays the topic and up to 3 panelist names
AND clicking any card takes me to that debate's viewer page
```

### US7.2: Search Debates (P1)
**As a** user browsing debate history  
**I want to** search debates by topic or panelist name  
**So that** I can quickly find specific conversations

**Acceptance Criteria**:
- Search input filters debates instantly as user types (reactive)
- Search matches topic text and panelist names (case-insensitive)
- Results update within 100ms of typing
- Shows "No results found" when search yields no matches
- Clearing search restores full list

**Test Scenario**:
```
GIVEN I have 10 debates on various topics including "ethics of AI" and "theodicy"
WHEN I type "ai" in the search box
THEN only debates with "AI" in topic or panelist names appear
AND results update instantly without clicking search button
AND typing backspace to clear shows all debates again
```

### US7.3: Navigate to History from Home (P2)
**As a** user on the home page  
**I want to** access my debate history with one click  
**So that** I can quickly revisit past debates

**Acceptance Criteria**:
- "View Debate History" link/button visible on home page
- Clicking navigates to debate history page
- Link is prominent but doesn't distract from primary flow

**Test Scenario**:
```
GIVEN I am on the home page
WHEN I click "View Debate History"
THEN I navigate to /debates page showing my debate history
```

## Functional Requirements

### Core Functionality

**FR-101**: System shall display debates in reverse chronological order (newest first)  
**FR-102**: System shall show topic, panelists (max 3 + count), and timestamp for each debate  
**FR-103**: System shall implement client-side reactive search filtering debate list  
**FR-104**: System shall match search query against topic text and all panelist names (case-insensitive)  
**FR-105**: System shall update search results within 100ms of keystroke  
**FR-106**: System shall navigate to debate viewer when user clicks debate card  
**FR-107**: System shall display empty state message when no debates exist  
**FR-108**: System shall display "No results" message when search yields no matches  

### API Requirements

**FR-109**: Backend shall provide GET endpoint to fetch debates with pagination  
**FR-110**: Backend shall support query parameters: `limit` (default 20), `offset` (default 0)  
**FR-111**: Backend shall return debates with: id, topic, panelists (id, name), startedAt timestamp  
**FR-112**: Backend shall handle Firestore query errors gracefully with 500 response  

### UI/UX Requirements

**FR-113**: Search input shall have placeholder text "Search debates by topic or panelist..."  
**FR-114**: Debate cards shall use card layout with hover effects  
**FR-115**: Timestamp shall display in relative format (e.g., "2 hours ago", "3 days ago")  
**FR-116**: Panelist names shall display as comma-separated list with ellipsis for overflow  
**FR-117**: Page shall be responsive on mobile (≥375px width)  

### Navigation

**FR-118**: Home page shall include navigation link to debate history  
**FR-119**: Debate viewer page shall include navigation link back to debate history  
**FR-120**: Debate history page shall include navigation link back to home  

## Success Criteria

**SC-101**: Users can view all debates in under 2 seconds (for up to 100 debates)  
**SC-102**: Search results update within 100ms of typing  
**SC-103**: Clicking a debate card navigates to viewer within 200ms  
**SC-104**: Page is fully responsive on screens ≥375px width  
**SC-105**: Empty state is clear and includes CTA to create first debate  

## Non-Functional Requirements

### Performance
- Initial page load: <2s for 20 debates
- Search filter execution: <100ms
- Navigation to debate viewer: <200ms

### Accessibility
- Search input has proper ARIA labels
- Debate cards are keyboard navigable (Tab key)
- Focus states clearly visible
- Screen reader announces search result count

### Browser Support
- Chrome 90+, Firefox 88+, Safari 14+, Edge 90+
- Mobile: iOS Safari 14+, Android Chrome 90+

## Out of Scope

- Debate deletion (no user data management in MVP)
- Sorting options (chronological only)
- Filtering by date range or panelist
- Infinite scroll (pagination only)
- Share buttons (debates are shareable via URL already)
- Export multiple debates at once

## Dependencies

- Firestore database with debates stored from previous sessions
- Existing debate viewer page at `/d/:id`
- React Router for navigation

## Assumptions

- Debates are already being saved to Firestore with id, topic, panelists, startedAt
- Firestore security rules allow public read access to debates collection
- Average user has <100 debates (pagination handles growth)
- Search is client-side (sufficient for MVP scale)

## Open Questions

None - design is straightforward and leverages existing architecture.

## Acceptance Scenarios

### Scenario 1: First-time user with no debates
```
GIVEN I am a new user who has never generated a debate
WHEN I navigate to /debates
THEN I see an empty state message: "No debates yet. Create your first debate!"
AND a button to navigate to home page
```

### Scenario 2: Browse and select debate
```
GIVEN I have generated 8 debates
WHEN I navigate to /debates
THEN I see all 8 debates sorted newest first
WHEN I click on the 3rd debate card
THEN I navigate to /d/{id} and see that debate's full conversation
```

### Scenario 3: Search filters results
```
GIVEN I have debates on topics: "AI ethics", "theodicy", "free will"
WHEN I type "theo" in search
THEN only "theodicy" debate appears
WHEN I clear the search
THEN all 3 debates reappear
```

### Scenario 4: No search results
```
GIVEN I have 5 debates with no mention of "quantum"
WHEN I search for "quantum"
THEN I see "No debates found matching 'quantum'"
AND the search input remains focused
```
