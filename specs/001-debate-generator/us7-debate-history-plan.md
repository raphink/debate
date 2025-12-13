# Implementation Plan: Debate History with Reactive Search

**Feature**: US7-debate-history  
**Date**: 2025-12-13  
**Spec**: [us7-debate-history-spec.md](us7-debate-history-spec.md)

## Summary

Build a debate history view that displays all previously generated debates in reverse chronological order with instant client-side search filtering. Users can search by topic or panelist name with results updating reactively (<100ms) as they type. Clicking a debate card navigates to the full debate viewer. The backend provides a paginated GET endpoint to fetch debate metadata from Firestore, while the frontend implements reactive search filtering using React state.

**Key Architecture Decisions**:
- **Client-side search**: For MVP scale (<100 debates typical), client-side filtering is faster and simpler than backend search
- **Firestore query**: Order by `startedAt DESC`, paginate with `limit` and `offset` query parameters
- **Debate card format**: Show topic, up to 3 panelist names, relative timestamp, truncated topic if >100 chars
- **Navigation**: Add links from Home → History, History → Viewer, Viewer → History
- **Responsive design**: Mobile-first layout with card grid that adapts to screen size

## Technical Context

**Language/Version**: Go 1.24 (backend), React 18 (frontend)  
**Primary Dependencies**: 
- Backend: Google Cloud Functions SDK, Cloud Firestore SDK
- Frontend: React 18, React Router, date-fns (for relative timestamps)
**API Endpoint**: `GET /list-debates?limit=20&offset=0`  
**Firestore Collection**: `debates` (existing)  
**Performance Goals**:
- Initial load: <2s for 20 debates
- Search filtering: <100ms
- Navigation: <200ms

## Architecture

### Backend: list-debates Cloud Function

**Purpose**: Fetch paginated debate metadata from Firestore

**Endpoint**: `GET /list-debates`

**Query Parameters**:
- `limit` (integer, optional, default=20, max=100): Number of debates to return
- `offset` (integer, optional, default=0): Number of debates to skip

**Response Format** (200 OK):
```json
{
  "debates": [
    {
      "id": "550e8400-e29b-41d4-a716-446655440000",
      "topic": "Should Christians defy authorities when the law is unfair?",
      "panelists": [
        { "id": "Augustine354", "name": "Augustine of Hippo" },
        { "id": "MLKJr", "name": "Martin Luther King Jr." },
        { "id": "Aquinas1225", "name": "Thomas Aquinas" }
      ],
      "startedAt": "2025-12-13T10:30:00Z"
    }
  ],
  "total": 42,
  "hasMore": true
}
```

**Error Responses**:
- 400 Bad Request: Invalid limit/offset parameters
- 500 Internal Server Error: Firestore query failure

**Implementation Details**:
- Use Firestore `Collection("debates").OrderBy("startedAt", firestore.Desc)`
- Apply `Limit(limit)` and `Offset(offset)` for pagination
- Query only needed fields: id, topic, panelists (id, name), startedAt
- Return `total` count from separate count query
- Set `hasMore = (offset + limit) < total`

### Frontend: DebateHistory Component

**Purpose**: Display searchable list of debates

**Component Structure**:
```
DebateHistory (page)
├── SearchInput (search box with debounce)
├── DebateList (grid of cards)
│   └── DebateCard (individual debate)
└── EmptyState (when no debates or no results)
```

**State Management**:
```javascript
const [debates, setDebates] = useState([]);      // All fetched debates
const [searchQuery, setSearchQuery] = useState(''); // Current search text
const [isLoading, setIsLoading] = useState(true);
const [error, setError] = useState(null);

// Filtered debates computed from debates + searchQuery
const filteredDebates = useMemo(() => {
  if (!searchQuery) return debates;
  
  const query = searchQuery.toLowerCase();
  return debates.filter(debate => 
    debate.topic.toLowerCase().includes(query) ||
    debate.panelists.some(p => p.name.toLowerCase().includes(query))
  );
}, [debates, searchQuery]);
```

**Search Implementation**:
- Use controlled input with `onChange` updating `searchQuery` state
- No debounce needed (filtering is instant with <100 debates)
- Display result count: "Showing X of Y debates"
- Clear button appears when search is active

**Card Design**:
```jsx
<div className={styles.debateCard} onClick={() => navigate(`/d/${debate.id}`))}>
  <h3 className={styles.topic}>{truncate(debate.topic, 100)}</h3>
  <div className={styles.panelists}>
    {formatPanelists(debate.panelists)} {/* "Name1, Name2, Name3 +2" */}
  </div>
  <div className={styles.timestamp}>
    {formatRelativeTime(debate.startedAt)} {/* "2 hours ago" */}
  </div>
</div>
```

**Helper Functions**:
- `formatPanelists(panelists)`: Show first 3 names + "+N more" if applicable
- `formatRelativeTime(timestamp)`: Convert ISO timestamp to relative format using `date-fns`
- `truncate(text, maxLength)`: Truncate with ellipsis if too long

### Navigation Updates

**Home Page** (`/`):
- Add "View Debate History" button below topic input section
- Button style: Secondary action (not primary gradient)

**Debate Viewer** (`/d/:id`):
- Add "← Back to History" link in header
- Link appears next to "Export as PDF" button

**Debate History** (`/debates`):
- Add "← Back to Home" link in header
- Link aligned to left of page title

### API Service Layer

**File**: `frontend/src/services/debateService.js` (update existing)

```javascript
export const fetchDebateHistory = async (limit = 20, offset = 0) => {
  const url = `${API_BASE_URL}/list-debates?limit=${limit}&offset=${offset}`;
  const response = await fetch(url);
  
  if (!response.ok) {
    throw new Error(`Failed to fetch debates: ${response.statusText}`);
  }
  
  return response.json();
};
```

## File Structure

### Backend
```
backend/functions/list-debates/
├── main.go              # Cloud Functions entry point
├── handler.go           # HTTP handler with CORS
├── firestore.go         # Firestore query logic
├── types.go             # Request/response structs
├── Dockerfile           # Multi-stage build
├── go.mod
└── cmd/
    └── main.go          # Local dev server (port 8086)
```

### Frontend
```
frontend/src/
├── pages/
│   ├── DebateHistory.jsx          # Main history page
│   └── DebateHistory.module.css   # Page styles
├── components/
│   ├── DebateCard/
│   │   ├── DebateCard.jsx         # Individual debate card
│   │   └── DebateCard.module.css
│   └── SearchInput/
│       ├── SearchInput.jsx        # Reusable search input
│       └── SearchInput.module.css
└── services/
    └── debateService.js           # Update with fetchDebateHistory
```

## API Contract

**File**: `specs/001-debate-generator/contracts/list-debates.json`

```json
{
  "endpoint": "/list-debates",
  "method": "GET",
  "description": "Fetch paginated list of debates from Firestore",
  "queryParameters": {
    "limit": {
      "type": "integer",
      "required": false,
      "default": 20,
      "min": 1,
      "max": 100,
      "description": "Number of debates to return"
    },
    "offset": {
      "type": "integer",
      "required": false,
      "default": 0,
      "min": 0,
      "description": "Number of debates to skip (pagination)"
    }
  },
  "responses": {
    "200": {
      "description": "Success - debates fetched",
      "schema": {
        "debates": {
          "type": "array",
          "items": {
            "id": "string (UUID)",
            "topic": "string",
            "panelists": "array of {id: string, name: string}",
            "startedAt": "string (ISO 8601)"
          }
        },
        "total": "integer (total debate count in collection)",
        "hasMore": "boolean (true if more results exist)"
      },
      "example": {
        "debates": [
          {
            "id": "550e8400-e29b-41d4-a716-446655440000",
            "topic": "Should Christians defy authorities when the law is unfair?",
            "panelists": [
              { "id": "Augustine354", "name": "Augustine of Hippo" },
              { "id": "MLKJr", "name": "Martin Luther King Jr." }
            ],
            "startedAt": "2025-12-13T10:30:00Z"
          }
        ],
        "total": 42,
        "hasMore": true
      }
    },
    "400": {
      "description": "Bad Request - invalid parameters",
      "schema": {
        "error": "string"
      },
      "example": {
        "error": "Invalid limit: must be between 1 and 100"
      }
    },
    "500": {
      "description": "Internal Server Error - Firestore failure",
      "schema": {
        "error": "string"
      },
      "example": {
        "error": "Failed to query debates from Firestore"
      }
    }
  },
  "testCases": [
    {
      "name": "Fetch first page of debates",
      "request": "GET /list-debates?limit=10",
      "expectedStatus": 200,
      "expectedResponse": {
        "debates": "array of 10 debates (or fewer if less exist)",
        "total": "integer",
        "hasMore": "boolean"
      }
    },
    {
      "name": "Fetch second page with offset",
      "request": "GET /list-debates?limit=10&offset=10",
      "expectedStatus": 200,
      "expectedResponse": {
        "debates": "array of next 10 debates",
        "hasMore": "boolean based on total"
      }
    },
    {
      "name": "Invalid limit exceeds max",
      "request": "GET /list-debates?limit=200",
      "expectedStatus": 400,
      "expectedResponse": {
        "error": "Invalid limit: must be between 1 and 100"
      }
    },
    {
      "name": "Negative offset",
      "request": "GET /list-debates?offset=-5",
      "expectedStatus": 400,
      "expectedResponse": {
        "error": "Invalid offset: must be >= 0"
      }
    }
  ]
}
```

## Docker Compose Integration

**Update**: `docker-compose.yml`

```yaml
list-debates:
  build:
    context: ./backend/functions/list-debates
    dockerfile: Dockerfile
  ports:
    - "8086:8080"
  environment:
    - GOOGLE_APPLICATION_CREDENTIALS=/app/secrets/service-account-key.json
    - ALLOWED_ORIGIN=http://localhost:3000
  volumes:
    - ./secrets:/app/secrets:ro
```

**Frontend Environment Variable**:
```env
REACT_APP_LIST_DEBATES_URL=http://localhost:8086
```

## Testing Strategy

### Backend Tests
- Unit tests for query parameter validation
- Integration tests for Firestore queries (using emulator)
- Test pagination: offset=0, offset=10, offset > total
- Test empty collection response

### Frontend Tests
- Render DebateHistory with mock data
- Search filtering updates results correctly
- Empty state shown when no debates
- "No results" shown when search yields nothing
- Navigation to debate viewer on card click

### Manual Testing
1. Create 3 debates with different topics
2. Navigate to `/debates` - verify all 3 appear
3. Search for keyword from one topic - verify filtering
4. Click debate card - verify navigation to viewer
5. Test on mobile viewport (375px width)

## Deployment

1. Deploy `list-debates` Cloud Function to GCP
2. Update frontend environment variable with production URL
3. Deploy frontend to GitHub Pages
4. Verify CORS settings allow frontend origin

## Future Enhancements

- Infinite scroll instead of pagination
- Filter by date range
- Sort options (oldest first, by panelist count)
- Full-text search using Firestore text queries or Algolia
- Debate deletion with confirmation dialog

## Estimated Effort

- Backend (list-debates function): **2-3 hours**
- Frontend (DebateHistory page + components): **3-4 hours**
- Navigation updates: **1 hour**
- Testing: **2 hours**
- **Total: 8-10 hours**
