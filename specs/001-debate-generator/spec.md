# Feature Specification: AI-Powered Theology/Philosophy Debate Generator

**Feature Branch**: `001-debate-generator`  
**Created**: 2025-12-11  
**Status**: Draft  
**Input**: User description: "web app to generate debates on specific topics"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Topic Entry and Validation (Priority: P1)

User enters a debate topic and receives immediate feedback on whether it's suitable for theological/philosophical discussion, creating confidence in topic selection.

**Why this priority**: This is the entry point to the entire application. Without a valid topic, no debate can be generated. It provides immediate value by validating user input and preventing wasted time on irrelevant topics.

**Independent Test**: Can be fully tested by entering various topics (theological, philosophical, off-topic) and verifying validation responses without needing panelist selection or debate generation.

**Acceptance Scenarios**:

1. **Given** user is on the home page, **When** user enters "Should Christians defy authorities when the law is unfair?" and clicks "Find Panelists", **Then** system shows "Looking for Panelists" loading animation and validates topic as relevant for theology/philosophy debate
2. **Given** user has entered a topic, **When** validation completes, **Then** system displays validation result with clear messaging within 3 seconds
3. **Given** user is validating a topic, **When** waiting for response, **Then** system displays engaging loading animation with "Looking for Panelists" message
4. **Given** user enters an off-topic query like "Best pizza toppings", **When** validation runs, **Then** system politely indicates topic is not suitable for theological/philosophical debate and suggests refinement

---

### User Story 2 - Panelist Discovery and Selection (Priority: P1)

User browses AI-suggested historical figures with known positions on the topic and selects up to 5 panelists for the debate, creating a customized panel.

**Why this priority**: Panelist selection is essential for debate generation and provides core user value. The curated list of relevant historical thinkers is a key differentiator and educational feature.

**Independent Test**: Can be tested by validating a topic, viewing the suggested panelist list with complete profiles (name, avatar, tagline, bio), and verifying selection mechanics work correctly up to the 5-panelist limit.

**Acceptance Scenarios**:

1. **Given** user is entering a topic, **When** user wants to suggest specific panelists, **Then** user can enter up to 5 panelist names (comma-separated) in an optional input field
2. **Given** topic has been validated as relevant, **When** validation response streams in, **Then** user sees 8-20 historical figures appear progressively with avatar, name, handle (id), tagline, and bio (each panelist sent as complete JSON line as soon as Claude generates it)
3. **Given** user suggested panelist names during topic entry, **When** validation response streams in, **Then** AI MUST include suggested names unless they are fictional/non-existent or completely unrelated to intellectual discourse, inferring positions from their known works even if they never directly addressed the topic
4. **Given** panelist list is displayed, **When** user reviews the suggestions, **Then** panelists represent diverse time periods across the last 2000 years (roughly 25% ancient/early church 0-500 AD, 25% medieval/reformation 500-1700 AD, 25% modern 1700-1950 AD, 25% contemporary 1950-present)
5. **Given** user views panelist list, **When** user clicks on a panelist card, **Then** panelist is added to selection (maximum 5 total)
6. **Given** user has selected 5 panelists, **When** user attempts to select another, **Then** system prevents selection and displays message "Maximum 5 panelists allowed"
7. **Given** user has selected panelists, **When** user clicks a selected panelist again, **Then** panelist is deselected and removed from selection
8. **Given** panelist list is displayed, **When** user reviews panelist details, **Then** each panelist shows relevant credentials, historical period, and known position on the topic

---

### User Story 3 - Live Debate Generation with Streaming Display (Priority: P1)

User launches debate generation and watches the conversation unfold in real-time as a chat-style interface with panelist avatars, creating an engaging experience.

**Why this priority**: This is the core value proposition - generating the actual debate. Without this, the application provides no meaningful output. Streaming display is critical for user engagement and perceived performance.

**Independent Test**: Can be tested by selecting panelists, launching debate, and verifying that responses stream progressively into chat-style bubbles with correct avatar attribution and formatting.

**Acceptance Scenarios**:

1. **Given** user has selected 2-5 panelists, **When** user clicks "Generate Debate", **Then** debate generation begins with a neutral moderator introducing the topic and panelists, and first response appears within 5 seconds
2. **Given** debate is generating, **When** AI produces responses, **Then** each response appears progressively in a chat bubble with the speaking panelist's avatar and name, or the moderator's avatar for moderation
3. **Given** debate is streaming, **When** a panelist's turn begins, **Then** their chat bubble appears with loading indicator before text streams in
4. **Given** debate is in progress, **When** user views the conversation, **Then** different panelists' responses and moderator interventions are visually distinguishable by avatar and styling
5. **Given** debate includes moderator, **When** moderator intervenes, **Then** moderator may redirect conversation, ask clarifying questions, highlight contrasts, or summarize progress between panelist exchanges
6. **Given** debate is nearing completion, **When** panelists have made their main arguments, **Then** moderator provides a concluding summary that synthesizes the key points and ends the debate
7. **Given** debate is streaming, **When** user clicks on a panelist avatar, **Then** a modal opens displaying the panelist's name, tagline, and full biography
8. **Given** panelist modal is open, **When** user clicks outside modal or presses Escape key, **Then** modal closes and returns focus to debate view
9. **Given** debate is streaming, **When** user toggles auto-scroll, **Then** conversation view automatically follows new messages (when enabled) or remains at current scroll position (when disabled, default)
10. **Given** debate is generating, **When** an error occurs (API timeout, rate limit), **Then** user sees friendly error message with option to retry

---

### User Story 4 - PDF Export (Priority: P2)

User exports completed debate as a formatted PDF document for offline reading, sharing, or archival purposes.

**Why this priority**: Provides additional value and shareability but is not essential for core debate generation functionality. Users can still read debates on-screen without export.

**Independent Test**: Can be tested by generating a complete debate and verifying PDF export produces a well-formatted document with all panelist responses, avatars, and metadata.

**Acceptance Scenarios**:

1. **Given** debate generation has completed, **When** user clicks "Export as PDF", **Then** PDF download begins within 2 seconds
2. **Given** PDF is generated, **When** user opens the PDF, **Then** document includes debate topic, panelist profiles, complete conversation with avatars, and timestamp
3. **Given** user exports PDF, **When** PDF renders, **Then** text is readable, avatars display correctly, and page breaks don't split individual responses awkwardly

---

### Edge Cases

- What happens when Claude API is unavailable or times out during topic validation?
- What happens when Claude API fails mid-stream during debate generation?
- How does system handle topics in non-English languages?
- What happens when Claude suggests fewer than 5 panelists for an obscure topic?
- How does system handle extremely long debate responses that exceed typical length?
- What happens when user closes browser during active debate generation?
- How does system handle rate limiting from Claude API during high usage?
- What happens when GCP function proxy returns malformed JSON?
- How does system handle panelist avatars that fail to load?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST accept text input for debate topics with minimum 10 characters
- **FR-002**: System MUST validate topic relevance for theology/philosophy debates via Claude API through GCP function proxy
- **FR-003**: System MUST allow users to optionally suggest up to 5 panelist names during topic entry (comma-separated input)
- **FR-003a**: System MUST send user-suggested panelist names to Claude API with PRIORITY weighting during topic validation
- **FR-003b**: Claude API MUST include user-suggested panelists unless they are fictional/non-existent or completely unrelated to theology/philosophy/intellectual discourse
- **FR-003c**: Claude API MUST infer positions for user-suggested panelists based on their known works and tradition, even if they never directly addressed the specific topic
- **FR-004**: System MUST request and display 8-20 panelist suggestions with complete profiles (name, avatar URL, handle (id), tagline, biography) using line-by-line JSON streaming via Server-Sent Events (SSE)
- **FR-004a**: Panelist suggestions MUST represent diverse historical periods across 2000 years (approximately 25% ancient/early church 0-500 AD, 25% medieval/reformation 500-1700 AD, 25% modern 1700-1950 AD, 25% contemporary 1950-present)
- **FR-004b**: System MUST stream panelists incrementally - each panelist emitted as a complete JSON line as soon as Claude generates it (character-by-character parsing to detect complete lines)
- **FR-004c**: If topic is not relevant, Claude returns rejection JSON instead of panelists, eliminating validation/panelist race condition
- **FR-005**: Users MUST be able to select between 2 and 5 panelists from the suggested list
- **FR-006**: System MUST visually distinguish selected vs unselected panelists in the UI
- **FR-007**: System MUST prevent debate generation unless at least 2 panelists are selected
- **FR-008**: System MUST send debate configuration (topic + selected panelists) to Claude API via GCP function proxy
- **FR-009**: System MUST stream debate responses progressively and display them in real-time
- **FR-010**: System MUST parse streaming responses to identify which panelist or moderator is speaking
- **FR-011**: System MUST display each panelist's response in a distinct chat bubble with their avatar
- **FR-011a**: System MUST include a neutral moderator who introduces the debate, may intervene between panelist exchanges, and provides a concluding summary at the end
- **FR-011b**: Moderator responses MUST be visually distinguished from panelist responses with unique avatar and styling
- **FR-012**: System MUST show loading/typing indicators while waiting for next response
- **FR-012a**: System MUST provide toggleable auto-scroll control for debate view (disabled by default)
- **FR-012b**: System MUST make panelist avatars clickable to display panelist details in a modal
- **FR-012c**: Modal MUST display panelist name, tagline, and biography with accessible close controls (X button, Escape key, click outside)
- **FR-013**: System MUST handle API errors gracefully with user-friendly error messages
- **FR-014**: System MUST provide retry mechanism for failed API calls
- **FR-015**: System MUST allow PDF export of completed debates
- **FR-016**: PDF export MUST include topic, panelist profiles, complete conversation, and generation timestamp
- **FR-017**: System MUST sanitize all Claude API outputs before rendering to prevent XSS attacks (per Constitution Principle V)
- **FR-018**: System MUST rate-limit API requests to prevent abuse (per Constitution Principle V)
- **FR-019**: System MUST validate and sanitize user topic input and suggested panelist names before sending to Claude API (per Constitution Principle V)
- **FR-020**: UI MUST be keyboard-navigable for accessibility (per Constitution Principle III)
- **FR-021**: System MUST maintain minimum 4.5:1 contrast ratio for text (per Constitution Principle III)

### Key Entities

- **Topic**: User-submitted debate subject; includes validation status and relevance indicator
- **Panelist**: Historical figure with position on topic; attributes include unique handle/identifier (alphanumeric only), name, avatar URL, tagline (brief descriptor), biography (credentials and viewpoint)
- **Moderator**: Neutral facilitator with ID "moderator"; introduces topic and panelists, may intervene between exchanges to redirect/clarify/summarize, MUST provide concluding summary at end of debate
- **Debate Configuration**: Combination of validated topic and selected panelists (2-5); represents user's debate setup
- **Debate Response**: Individual contribution from a panelist during debate generation; includes panelist identifier, response text, timestamp, and position in conversation
- **Debate Session**: Complete generated debate; contains topic, panelist list, ordered responses, generation timestamp, and completion status

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Users can enter a topic and receive validation feedback within 3 seconds on a stable internet connection
- **SC-002**: System successfully suggests relevant panelists for 90% of theology/philosophy topics
- **SC-003**: First debate response appears within 5 seconds of clicking "Generate Debate"
- **SC-004**: Debate responses stream progressively with no more than 500ms gaps between chunks
- **SC-005**: Users can complete the entire flow (topic entry → panelist selection → debate generation) within 2 minutes
- **SC-006**: PDF exports generate and download within 2 seconds for debates up to 5000 words
- **SC-007**: UI remains responsive (interactions respond within 100ms) during debate streaming
- **SC-008**: All interactive elements are keyboard-accessible and meet WCAG 2.1 Level AA standards
- **SC-009**: Error recovery mechanisms allow users to retry failed operations without losing their debate configuration
- **SC-010**: 95% of users successfully generate at least one complete debate on first attempt

## Assumptions *(documentation)*

- Claude API has sufficient context window to handle debate topic + 5 panelist biographies + conversation history
- GCP function proxy provides reliable authentication and request forwarding to Claude API
- Claude API can structure streaming responses in a parseable format that identifies speaking panelist
- Historical figure avatars are available via public domain sources or generated illustrations
- Users have modern browsers with JavaScript enabled and stable internet connection
- Debate topics will primarily be in English (multi-language support is out of scope for MVP)
- PDF generation can be handled client-side or via lightweight server-side library
- Claude API rate limits are sufficient for expected user load without requiring user authentication
- Historical panelist suggestions from Claude will be factually accurate and relevant to the topic

## Out of Scope *(explicit boundaries)*

- User authentication and account management
- Saving/persisting debates to database for later retrieval
- Sharing debates via social media or public links
- User customization of panelist avatars or bios
- Multi-language support (non-English topics)
- Real-time collaborative debate watching with multiple users
- Editing or regenerating portions of completed debates
- Audio/video generation of debates
- Direct creation of custom panelists with user-defined bios (AI evaluates suggested names instead)
- User acting as moderator with ability to ask questions during debate (future enhancement)
- Payment or subscription features
- Admin panel or content moderation tools
