# Data Model: AI-Powered Theology/Philosophy Debate Generator

**Feature**: 001-debate-generator  
**Date**: 2025-12-11  
**Purpose**: Define core entities, validation rules, and relationships

## Entity Definitions

### Topic

Represents a user-submitted debate subject for theological/philosophical discussion.

**Attributes**:
- `text` (string, required): The debate topic as entered by user
  - Min length: 10 characters
  - Max length: 500 characters
  - Validation: No HTML tags, trim whitespace
- `isRelevant` (boolean, required): Whether topic is suitable for theology/philosophy debate
  - Set by Claude API validation response
- `validationMessage` (string, optional): Feedback message from validation
  - Max length: 200 characters
  - Example: "This topic is well-suited for theological debate"
- `submittedAt` (timestamp, required): When user submitted the topic
  - ISO 8601 format: `2025-12-11T10:30:00Z`

**Validation Rules**:
- Topic text MUST NOT contain HTML tags or script content
- Topic text MUST be trimmed of leading/trailing whitespace
- Topic MUST be between 10-500 characters after sanitization
- Topic MUST be validated before proceeding to panelist selection

**State Transitions**:
```
[Entered] → [Validating] → [Valid|Invalid]
```

---

### Panelist

Represents a historical figure suggested by Claude API with known position on the debate topic.

**Attributes**:
- `id` (string, required): Unique identifier and social media-style handle
  - Pattern: Alphanumeric only (letters and numbers, no spaces or special characters)
  - Min length: 3 characters
  - Max length: 20 characters
  - Used for avatar mapping, selection tracking, and display in chat
  - Example: "Augustine354", "MLKJr", "Aquinas1225"
- `name` (string, required): Full name of historical figure
  - Max length: 100 characters
  - Example: "Augustine of Hippo"
- `tagline` (string, required): Brief descriptor of panelist
  - Max length: 150 characters
  - Example: "4th-5th century theologian and philosopher"
- `biography` (string, required): Detailed background and viewpoint
  - Max length: 500 characters
  - Includes: Era, key works, relevant positions on topic
  - Example: "Early Christian theologian known for 'Confessions' and 'City of God'. Advocated for natural law theory and believed civil disobedience justified when laws contradict divine law."
- `avatarUrl` (string, required): Path to panelist's avatar image
  - Format: `/avatars/{id}-avatar.png`
  - Example: `/avatars/Augustine354-avatar.png`
  - Fallback: `/avatars/placeholder-avatar.png`
- `position` (string, optional): Panelist's stance on the specific topic
  - Max length: 200 characters
  - Example: "Would argue that Christians have a duty to resist unjust laws"

**Validation Rules**:
- ID MUST be alphanumeric only (regex: `^[a-zA-Z0-9]{3,20}$`)
- Name, tagline, and biography MUST NOT contain HTML tags
- Avatar URL MUST point to existing file or use fallback
- All text fields MUST be sanitized before rendering

**Relationships**:
- Belongs to Topic (suggested for specific topic)
- Can be selected in DebateConfiguration (0 or 1 times per session)

---

### PanelistSelection

Represents the user's choice of panelists for the debate (client-side state).

**Attributes**:
- `selectedPanelists` (array of Panelist, required): Currently selected panelists
  - Min length: 2 (enforced when generating debate)
  - Max length: 5 (enforced during selection)
  - Order preserved: First selected appears first in debate
- `availablePanelists` (array of Panelist, required): All suggested panelists from API
  - Max length: 20
  - Filtered from Claude API response

**Validation Rules**:
- User MUST select at least 2 panelists before generating debate
- User CANNOT select more than 5 panelists
- User CANNOT select the same panelist twice
- Selection state MUST persist during navigation (session storage)

**State Transitions**:
```
[Empty] → [Selecting (1-4)] → [Ready (2-5)] → [Generating]
```

---

### DebateConfiguration

Represents the complete setup for debate generation (sent to backend).

**Attributes**:
- `topic` (Topic, required): The validated debate topic
  - MUST have `isRelevant = true`
- `panelists` (array of Panelist, required): Selected panelists for debate
  - Min length: 2
  - Max length: 5
  - Includes full panelist objects with id, name, biography, position
- `conversationStyle` (string, optional): Future enhancement placeholder
  - Default: "balanced" (for MVP, not configurable)
  - Potential values: "balanced", "adversarial", "socratic"

**Validation Rules**:
- Topic MUST be validated and relevant
- Panelist array MUST contain 2-5 unique panelists
- All panelists MUST have complete profiles (name, bio, position)

**JSON Representation** (sent to generate-debate function):
```json
{
  "topic": {
    "text": "Should Christians defy authorities when the law is unfair?",
    "isRelevant": true
  },
  "panelists": [
    {
      "id": "Augustine354",
      "name": "Augustine of Hippo",
      "biography": "Early Christian theologian...",
      "position": "Would argue that Christians have a duty to resist unjust laws"
    },
    {
      "id": "MLKJr",
      "name": "Martin Luther King Jr.",
      "biography": "20th century civil rights leader...",
      "position": "Believed in nonviolent civil disobedience against unjust laws"
    }
  ]
}
```

---

### DebateMessage

Represents a single contribution from a panelist during debate generation.

**Attributes**:
- `id` (string, required): Unique message identifier
  - Format: `{panelistId}-{sequenceNumber}`
  - Example: "Augustine354-1"
- `panelistId` (string, required): ID of speaking panelist (used for display and correlation)
  - References Panelist.id from configuration
- `panelistName` (string, required): Display name for UI
  - Copy of Panelist.name for quick access
- `avatarUrl` (string, required): Avatar to display with message
  - Copy of Panelist.avatarUrl for quick access
- `text` (string, required): The panelist's contribution
  - Max length: 2000 characters per message
  - Accumulated from streaming chunks
- `timestamp` (timestamp, required): When message was generated
  - ISO 8601 format
- `sequence` (integer, required): Position in conversation
  - Zero-indexed: 0, 1, 2, 3...
  - Used for ordering and display
- `isComplete` (boolean, required): Whether message has finished streaming
  - `false` while chunks are arriving
  - `true` when final chunk received

**Validation Rules**:
- Panelist ID MUST match one of the selected panelists
- Text MUST be sanitized before rendering (XSS prevention)
- Sequence MUST be monotonically increasing
- Text MUST NOT be empty when `isComplete = true`

**State Transitions**:
```
[Streaming] → [Complete] → [Rendered]
```

---

### DebateSession

Represents the complete generated debate (client-side aggregation).

**Attributes**:
- `id` (string, required): Unique session identifier
  - Format: `debate-{timestamp}`
  - Example: "debate-1702301400000"
- `topic` (Topic, required): The debate topic
- `panelists` (array of Panelist, required): Participating panelists (2-5)
- `messages` (array of DebateMessage, required): All debate contributions
  - Ordered by sequence
  - Includes completed and in-progress messages
- `status` (enum, required): Current state of debate
  - Values: "generating" | "complete" | "error"
- `startedAt` (timestamp, required): When debate generation began
- `completedAt` (timestamp, optional): When debate generation finished
  - Null while status = "generating"
- `error` (object, optional): Error details if generation failed
  - `message` (string): User-friendly error message
  - `retryable` (boolean): Whether user can retry
  - `errorCode` (string): Error identifier for debugging

**Validation Rules**:
- Session MUST have at least 2 messages before marking complete
- All messages MUST have `isComplete = true` before session status = "complete"
- Completed session MUST have `completedAt` timestamp

**State Transitions**:
```
[Initializing] → [Generating] → [Complete|Error]
                      ↓
                  [Retrying] ← (if error.retryable = true)
```

---

### StreamChunk

Represents a single Server-Sent Event (SSE) chunk during debate streaming (backend → frontend).

**Attributes**:
- `event` (string, required): SSE event type
  - Values: "message" | "error" | "done"
- `data` (object, required): Event payload
  - For "message" events:
    - `panelistId` (string): Speaking panelist
    - `textChunk` (string): Incremental text content
    - `sequence` (integer): Message sequence number
  - For "error" events:
    - `message` (string): Error description
    - `retryable` (boolean): Can user retry
    - `code` (string): Error code
  - For "done" events:
    - `totalMessages` (integer): Final message count

**JSON Representation** (SSE format):
```
event: message
data: {"panelistId":"Augustine354","textChunk":"I believe that ","sequence":0}

event: message
data: {"panelistId":"Augustine354","textChunk":"divine law supersedes human law.","sequence":0}

event: done
data: {"totalMessages":12}
```

**Validation Rules**:
- Event type MUST be one of: "message", "error", "done"
- Text chunks MUST be sanitized before rendering
- Sequence numbers MUST be consistent within a message

---

## Relationships

```
Topic (1) ←→ (0..20) Panelist [suggested_for]
  ↓
  |
PanelistSelection (1) ←→ (2..5) Panelist [selected_from]
  ↓
  |
DebateConfiguration (1)
  ↓
  |
DebateSession (1) ←→ (2..∞) DebateMessage [contains]
  |
  └→ Panelist (2..5) [participated_by]
```

**Cardinality Notes**:
- One Topic can have 0-20 suggested Panelists
- One PanelistSelection contains 2-5 Panelists
- One DebateSession contains multiple DebateMessages (typically 10-20)
- Each DebateMessage references exactly one Panelist

---

## Data Flow

### 1. Topic Validation Flow
```
User Input → Topic (text) 
  → Backend: validate-topic function
  → Claude API validation
  → Topic (text, isRelevant, validationMessage)
  → Frontend: Display validation result
```

### 2. Panelist Suggestion Flow
```
Topic (validated) 
  → Backend: suggest-panelists function
  → Claude API with topic context
  → Array of Panelist objects (0-20)
  → Frontend: Display PanelistGrid
```

### 3. Debate Generation Flow
```
DebateConfiguration (topic + panelists)
  → Backend: generate-debate function
  → Claude API streaming
  → StreamChunk events (SSE)
  → Frontend: Accumulate into DebateMessage objects
  → DebateSession (messages array)
  → UI: Render DebateView with chat bubbles
```

---

## Validation Summary

| Entity | Critical Validations |
|--------|---------------------|
| Topic | Length 10-500 chars, no HTML, sanitized |
| Panelist | No HTML in text fields, valid avatar URL |
| PanelistSelection | 2-5 unique selections |
| DebateConfiguration | Valid topic + 2-5 panelists |
| DebateMessage | Sanitized text, valid panelist ID |
| DebateSession | ≥2 messages, all complete before done |
| StreamChunk | Valid event type, sanitized chunks |

All entities follow Constitution Principle V (AI Safety) by implementing sanitization and validation at every layer.
