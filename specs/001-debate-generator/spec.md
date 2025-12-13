# Feature Specification: AI-Powered Theology/Philosophy Debate Generator

**Feature Branch**: `001-debate-generator`  
**Created**: 2025-12-11  
**Status**: Draft  
**Input**: User description: "web app to generate debates on specific topics"

## Clarifications

### Session 2025-12-13

- Q: How should the debate completion be determined? → A: Token/word count threshold (~5000 words generated)
- Q: Are autocomplete features (main spec US6 & US7) deferred or should they be implemented? → A: Deferred - US6 & US7 are post-MVP, only basic debate history (list-debates) is current scope
- Q: What should the Firestore security model be? → A: Completely locked down - no direct client access. All reads AND writes happen exclusively via backend Cloud Functions
- Q: What should the portrait fallback strategy be? → A: SVG placeholder only - standardize on placeholder-avatar.svg for all missing portraits

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

1. **Given** user is entering a topic, **When** user wants to suggest specific panelists, **Then** user can type names that appear as removable chips (press comma, Tab, or Enter to add, max 5)
2. **Given** user has added panelist chips, **When** user wants to remove one, **Then** user can click the × button on the chip to remove it
3. **Given** topic has been validated as relevant, **When** validation response streams in, **Then** user sees 8-20 historical figures appear progressively with avatar, name, handle (id), tagline, and bio (each panelist sent as complete JSON line as soon as Claude generates it)
4. **Given** user suggested panelist names during topic entry, **When** validation response streams in, **Then** AI MUST include suggested names unless they are fictional/non-existent or completely unrelated to intellectual discourse, inferring positions from their known works even if they never directly addressed the topic
5. **Given** user clicks "Find Panelists", **When** validation begins, **Then** input section is hidden and panelists appear one by one in the main view
6. **Given** panelists are streaming in, **When** user views the page, **Then** loading animation appears at the bottom of the panelist list (not blocking the view)
7. **Given** panelist list is displayed, **When** user reviews the suggestions, **Then** panelists represent diverse time periods across the last 2000 years (roughly 25% ancient/early church 0-500 AD, 25% medieval/reformation 500-1700 AD, 25% modern 1700-1950 AD, 25% contemporary 1950-present)
8. **Given** user views panelist list, **When** user clicks on a panelist card, **Then** panelist is added to selection (maximum 5 total)
9. **Given** user has selected 5 panelists, **When** user attempts to select another, **Then** system prevents selection and displays message "Maximum 5 panelists allowed"
10. **Given** user has selected panelists, **When** user clicks a selected panelist again, **Then** panelist is deselected and removed from selection
11. **Given** streaming completes, **When** user has selected panelists, **Then** "Clear Selection" and "Generate Debate" buttons appear with gradient styling and hover effects
12. **Given** panelist list is displayed, **When** user reviews panelist details, **Then** each panelist shows relevant credentials, historical period, and known position on the topic

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
6. **Given** debate is nearing completion (approximately 5000 words generated), **When** panelists have made their main arguments, **Then** moderator provides a concluding summary that synthesizes the key points and ends the debate
7. **Given** debate is streaming, **When** user clicks on a panelist avatar, **Then** a modal opens displaying the panelist's name, tagline, and full biography
8. **Given** panelist modal is open, **When** user clicks outside modal or presses Escape key, **Then** modal closes and returns focus to debate view
9. **Given** debate is streaming, **When** user toggles auto-scroll, **Then** conversation view automatically follows new messages (when enabled) or remains at current scroll position (when disabled, default)
10. **Given** debate is generating, **When** an error occurs (API timeout, rate limit), **Then** user sees friendly error message with option to retry
11. **Given** Claude uses inline Markdown formatting, **When** debate messages stream in with *italic*, **bold**, or ***bold italic*** text, **Then** formatting is rendered correctly in chat bubbles

---

### User Story 4 - PDF Export (Priority: P2)

User exports completed debate as a formatted PDF document for offline reading, sharing, or archival purposes.

**Why this priority**: Provides additional value and shareability but is not essential for core debate generation functionality. Users can still read debates on-screen without export.

**Independent Test**: Can be tested by generating a complete debate and verifying PDF export produces a well-formatted document with all panelist responses, avatars, and metadata.

**Acceptance Scenarios**:

1. **Given** debate generation has completed, **When** user clicks "Export as PDF", **Then** PDF download begins within 2 seconds
2. **Given** PDF is generated, **When** user opens the PDF, **Then** document includes debate topic, panelist profiles with circular portrait avatars, complete conversation in chat bubble format with portraits, and timestamp
3. **Given** user exports PDF, **When** PDF renders, **Then** text is readable, portrait images are embedded correctly (both Wikimedia URLs and local avatars), chat bubbles match web UI styling, and page breaks don't split individual responses awkwardly
4. **Given** PDF includes panelist portraits, **When** portrait is from Wikimedia Commons (absolute URL), **Then** image is fetched and embedded in PDF
5. **Given** PDF includes panelist portraits, **When** portrait is local avatar (relative path), **Then** image is correctly resolved and embedded in PDF

---

### User Story 5 - Debate Sharing and Caching (Priority: P2)

User shares completed debate via URL that loads from cached storage, allowing debates to be revisited and shared with others without regeneration.

**Why this priority**: Adds significant value for sharing and persistence but is not required for core debate generation. Users can still generate and view debates without sharing capability.

**Independent Test**: Can be tested by completing a debate, copying share URL, opening in new browser/incognito, and verifying debate loads from Firestore with identical content.

**Acceptance Scenarios**:

1. **Given** user starts debate generation, **When** "Generate Debate" is clicked, **Then** system generates unique UUID and updates browser URL to /d/{uuid}
2. **Given** debate generation completes successfully, **When** final message streams in, **Then** system automatically saves complete debate to Firestore with UUID as document ID
3. **Given** debate is saved to Firestore, **When** save operation fails, **Then** user can still view and export debate (Firestore save is non-blocking enhancement)
4. **Given** debate is displayed, **When** user clicks "Share" button, **Then** debate URL (https://domain/d/{uuid}) is copied to clipboard and confirmation appears
5. **Given** user has debate URL, **When** user visits URL in new browser or shares with others, **Then** debate loads from Firestore cache within 2 seconds showing complete conversation
6. **Given** debate loads from Firestore, **When** page renders, **Then** all panelist avatars, message formatting, and metadata display identically to original generation
7. **Given** user visits shared debate URL, **When** debate doesn't exist in Firestore (404), **Then** system shows friendly "Debate not found" message with link to create new debate
8. **Given** debate is loaded from cache, **When** user views page, **Then** PDF export and share functions work identically to freshly generated debates

---
### User Story 6 - Topic Discovery via History Integration (Priority: P3) [DEFERRED POST-MVP]

User discovers previous debate topics directly within the topic input field as an autocomplete dropdown, allowing quick re-use of existing debates with their original panelists or modifications.

**Status**: DEFERRED - Not included in current implementation scope. Basic debate history browsing via list-debates is available instead.

**Why this priority**: Quality-of-life enhancement that streamlines the workflow by combining topic discovery and input in a single interface. Users can quickly access previous debates or use them as starting points for variations.

**Independent Test**: Can be tested by generating several debates with different topics, returning to home page, typing partial topic text, and verifying autocomplete shows matching previous topics. Selecting a topic should pre-fill panelists and optionally load cached debate.

**Acceptance Scenarios**:

1. **Given** user types in topic input field, **When** user has typed ≥3 characters, **Then** system displays dropdown showing up to 10 matching previous topics ordered by recency
2. **Given** topic autocomplete dropdown is displayed, **When** user views a suggestion, **Then** entry shows full topic text and count of panelists (e.g., "3 panelists")
3. **Given** user types in topic input field, **When** input matches previous topics, **Then** matching topics are highlighted/narrowed in real-time as user continues typing
4. **Given** user selects a topic from autocomplete dropdown, **When** topic is selected, **Then** system validates topic (skipping Claude validation) and navigates to panelist selection with original panelists pre-selected
5. **Given** user is on panelist selection with pre-filled panelists from history, **When** user makes no changes to panelist list, **Then** system detects cache hit and loads debate directly from Firestore (bypassing generation)
6. **Given** user is on panelist selection with pre-filled panelists, **When** user clicks "Modify Panelists" button, **Then** system allows editing chips and generates new debate if changes are made
7. **Given** topic autocomplete is loading, **When** API call is in progress, **Then** system shows subtle loading indicator without blocking typing
8. **Given** no previous topics match user input or Firestore fails, **When** user types, **Then** autocomplete dropdown is hidden and user can submit new topic normally (graceful degradation)
9. **Given** user selects historical topic with modified panelists, **When** user proceeds to generate debate, **Then** system generates new debate (no cache hit) and saves as new debate instance
 [DEFERRED POST-MVP]

User receives intelligent panelist chip suggestions based on historical debate data, with normalized name matching to handle duplicate panelists identified differently across debates.

**Status**: DEFERRED - Not included in current implementation scope. Manual panelist name suggestions via chip input remain available

User receives intelligent panelist chip suggestions based on historical debate data, with normalized name matching to handle duplicate panelists identified differently across debates.

**Why this priority**: Quality-of-life enhancement that leverages Firestore data to improve UX, but not essential for core functionality. Users can still manually input panelist suggestions or use Claude's defaults.

**Independent Test**: Can be tested by generating debates with common panelists, then verifying autocomplete dropdown suggests previously used panelists when user starts typing a name.

**Acceptance Scenarios**:

1. **Given** user is on panelist selection page, **When** user starts typing in a chip input field, **Then** system shows autocomplete dropdown with matching panelists from Firestore history
2. **Given** autocomplete dropdown is displayed, **When** multiple panelists with similar names exist (e.g., "Augustine of Hippo", "St. Augustine", "Augustine"), **Then** system deduplicates via normalized matching (lowercase, title removal) and displays the most frequently used variant
3. **Given** autocomplete dropdown is displayed, **When** user views suggestions, **Then** panelists are ranked by frequency (most common first) and limited to top 10 matches
4. **Given** user selects autocomplete suggestion, **When** panelist is added to chips, **Then** system uses the historical panelist data (id, name, slug) from the selected suggestion
5. **Given** no historical panelists match user input, **When** user types, **Then** autocomplete dropdown is hidden and user can create new chip manually (no degradation of existing functionality)
6. **Given** Firestore read fails or is unavailable, **When** user types in chip input, **Then** autocomplete feature gracefully degrades to manual chip creation without errors
7. **Given** autocomplete API is slow (>500ms), **When** user is typing, **Then** system shows subtle loading indicator in dropdown without blocking input

---
### Edge Cases

- What happens when Claude API is unavailable or times out during topic validation?
- What happens when Claude API fails mid-stream during debate generation?
- How does system handle topics in non-English languages?
- What happens when Claude suggests fewer than 5 panelists for an obscure topic?
- What happens when extremely long debate responses that exceed typical length?
- What happens when user closes browser during active debate generation?
- How does system handle rate limiting from Claude API during high usage?
- What happens when GCP function proxy returns malformed JSON?
- How does system handle panelist avatars that fail to load?
- What happens when Wikimedia API returns 403 Forbidden (requires proper User-Agent header)?
- How does frontend distinguish between absolute portrait URLs and relative avatar paths?
- What happens when user installs PWA on mobile and launches offline without service worker?
- How does PWA manifest handle different screen sizes and orientations?
- What happens when Firestore save fails after debate generation completes?
- How does system handle loading debate from Firestore when document doesn't exist?
- What happens when Firestore read operation times out or fails?
- How does system prevent duplicate debate saves for the same UUID?
- What happens when user navigates away during Firestore save operation?
- What happens when Application Default Credentials are not configured for Firestore access?
- How does system handle Firestore authentication failures in local development vs production?

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
- **FR-004d**: System MUST provide separate async portrait service (get-portrait Cloud Function) to fetch panelist portrait URLs from Wikimedia Commons API after panelists stream in, keeping validation fast and non-blocking
- **FR-004e**: Portrait service MUST fetch 300px portrait images suitable for 48x48px circular display with proper User-Agent header "DebateApp/1.0" to avoid 403 Forbidden errors
- **FR-004f**: Portrait service MUST fall back to placeholder-avatar.svg (SVG format standardized) if Wikimedia API fails or returns no suitable image
- **FR-004g**: Portrait URLs MUST be cached in thread-safe in-memory map (sync.RWMutex) to avoid redundant API calls during debate generation
- **FR-004h**: Frontend MUST check if portrait URLs are absolute (http/https prefix) before prepending PUBLIC_URL/avatars/ path to avoid treating Wikimedia URLs as relative paths
- **FR-004i**: All components displaying avatars MUST use placeholder-avatar.svg as the canonical fallback when avatarUrl is empty, null, or fails to load
- **FR-005**: Users MUST be able to select between 2 and 5 panelists from the suggested list
- **FR-006**: System MUST visually distinguish selected vs unselected panelists in the UI
- **FR-007**: System MUST prevent debate generation unless at least 2 panelists are selected
- **FR-008**: System MUST send debate configuration (topic + selected panelists) to Claude API via GCP function proxy
- **FR-009**: System MUST stream debate responses progressively and display them in real-time
- **FR-010**: System MUST parse streaming responses to identify which panelist or moderator is speaking
- **FR-011**: System MUST display each panelist's response in a distinct chat bubble with their avatarwhen word count reaches approximately 5000 words or when arguments are naturally exhauste
- **FR-011a**: System MUST include a neutral moderator who introduces the debate, may intervene between panelist exchanges, and provides a concluding summary at the end
- **FR-011b**: Moderator responses MUST be visually distinguished from panelist responses with unique avatar and styling
- **FR-012**: System MUST show loading/typing indicators while waiting for next response
- **FR-012a**: System MUST provide toggleable auto-scroll control for debate view (disabled by default)
- **FR-012b**: System MUST make panelist avatars clickable to display panelist details in a modal
- **FR-012c**: Modal MUST display panelist name, tagline, and biography with accessible close controls (X button, Escape key, click outside)
- **FR-013**: System MUST handle API errors gracefully with user-friendly error messages
- **FR-014**: System MUST provide retry mechanism for failed API calls
- **FR-015**: System MUST allow PDF export of completed debates with chat bubble formatting and panelist portraits
- **FR-016**: PDF export MUST include topic, panelist profiles with circular portrait avatars, complete conversation in chat bubble format with portraits, and generation timestamp
- **FR-016a**: PDF export MUST embed portrait images from avatarUrl (both absolute Wikimedia URLs and relative local paths) as circular avatars in chat bubbles
- **FR-016b**: PDF export MUST render messages in chat bubble format matching the web UI, with speaker identification and visual distinction
- **FR-017**: System MUST sanitize all Claude API outputs before rendering to prevent XSS attacks (per Constitution Principle V)
- **FR-018**: System MUST rate-limit API requests to prevent abuse (per Constitution Principle V)
- **FR-019**: System MUST validate and sanitize user topic input and suggested panelist names before sending to Claude API (per Constitution Principle V)
- **FR-020**: UI MUST be keyboard-navigable for accessibility (per Constitution Principle III)
- **FR-021**: System MUST maintain minimum 4.5:1 contrast ratio for text (per Constitution Principle III)
- **FR-022**: Backend services MUST restrict CORS to localhost in development and raphink.github.io in production via ALLOWED_ORIGIN environment variable (per Constitution Principle V)
- **FR-023**: Application MUST provide PWA manifest for mobile installation with app name, description, icons, theme colors, and display mode
- **FR-023a**: PWA manifest MUST include multiple icon sizes (192x192, 512x512) for various mobile platforms
- **FR-023b**: Application MUST use standalone display mode to provide app-like experience when installed
- **FR-024**: System MUST render inline Markdown formatting (*italic*, **bold**, ***bold italic***) in debate messages for both web UI and PDF export
- **FR-024a**: Markdown rendering MUST properly escape HTML to prevent XSS attacks while preserving formatting
- **FR-025**: Backend MUST generate a unique UUID (v4) when receiving debate generation request to identify the debate session
- **FR-025a**: UUID MUST be cryptographically random using Go's uuid.New() to ensure uniqueness and unpredictability
- **FR-025b**: Backend MUST include generated UUID in SSE response header: X-Debate-Id
- **FR-025c**: Backend MUST expose X-Debate-Id header to frontend JavaScript via Access-Control-Expose-Headers CORS header
- **FR-026**: Backend MUST save completed debates to Firestore with UUID as document ID for caching and sharing
- **FR-026a**: Firestore document MUST include complete debate data: topic, panelists, messages, status, timestamps, and metadata
- **FR-026b**: Backend MUST write to Firestore automatically after debate generation completes successfully
- **FR-026c**: Firestore save failures MUST NOT prevent user from viewing the debate (graceful degradation, logged only)
- **FR-026d**: Backend MUST use Firebase Admin SDK for all Firestore operations (no client-side Firestore access)
- **FR-026e**: Backend MUST initialize Firestore with GCP_PROJECT_ID environment variable to specify target project
- **FR-027**: System MUST provide shareable URLs in format /d/{uuid} that load debates via backend API
- **FR-027a**: Backend MUST provide GET /api/get-debate?id={uuid} endpoint to retrieve saved debates
- **FR-027b**: Frontend MUST display complete debate with all original formatting, avatars, and metadata when loading from backend
- **FR-027c**: System MUST provide "Share" button that copies debate URL to clipboard with visual confirmation
- **FR-027d**: Loading debate from backend API MUST complete within 2 seconds on stable connection
- **FR-027e**: Frontend MUST update browser URL to /d/{uuid} after receiving debate ID from backend (History API, no page reload)
- **FR-028**: Firestore security rules MUST prevent all direct client access - all reads and writes happen exclusively via backend Cloud Functions (get-debate, list-debates, generate-debate)
- **FR-028a**: Backend Cloud Functions MUST validate debate ID format before querying Firestore
- **FR-028b**: Backend API MUST return 404 for non-existent debates and 500 for Firestore errors
- **FR-028c**: Firestore documents MUST be immutable (no updates or deletes after creation)
- **FR-029**: [DEFERRED POST-MVP] System SHOULD provide topic autocomplete endpoint: GET /api/autocomplete-topics?q={query}&limit=10 (US6 - deferred)
- **FR-029a**: [DEFERRED] Topic autocomplete endpoint SHOULD search historical debate topics from Firestore matching query substring (case-insensitive)
- **FR-029b**: [DEFERRED] Topic autocomplete response SHOULD return debates ordered by creation timestamp descending (newest first)
- **FR-029c**: [DEFERRED] Topic autocomplete response SHOULD include: debate ID, topic text, panelist count, panelist IDs/names, created timestamp
- **FR-029d**: [DEFERRED] Topic autocomplete dropdown SHOULD appear when user types ≥3 characters in topic input field on home page
- **FR-029e**: [DEFERRED] Topic autocomplete SHOULD show up to 10 matching results with full topic text and panelist count (e.g., "3 panelists")
- **FR-029f**: [DEFERRED] Selecting topic from autocomplete SHOULD skip Claude validation and pre-fill panelists on panelist selection page
- **FR-029g**: [DEFERRED] System SHOULD detect cache hit when topic + panelists match exactly and load debate from Firestore without regenerating
- **FR-029h**: [DEFERRED] Panelist selection page SHOULD show "Modify Panelists" button when panelists are pre-filled from history
- **FR-029i**: [DEFERRED] Topic autocomplete SHOULD degrade gracefully if Firestore unavailable (hide dropdown, allow manual topic entry)
- **FR-030**: [DEFERRED POST-MVP] System SHOULD provide panelist autocomplete endpoint: GET /api/autocomplete-panelists?q={query} (US7 - deferred)
- **FR-030a**: [DEFERRED] Autocomplete endpoint SHOULD aggregate panelists from all historical debates in Firestore
- **FR-030b**: [DEFERRED] Autocomplete SHOULD normalize panelist names for deduplication (lowercase, remove titles like "St.", "Dr.", strip punctuation)
- **FR-030c**: [DEFERRED] Autocomplete SHOULD use fuzzy matching to deduplicate similar panelists (e.g., "Augustine of Hippo", "St. Augustine", "Augustine")
- **FR-030d**: [DEFERRED] Autocomplete response SHOULD return panelists ranked by frequency (most common first), limited to top 10 matches
- **FR-030e**: [DEFERRED] Autocomplete SHOULD return canonical panelist data (id, name, slug) from the most frequently used variant
- **FR-030f**: [DEFERRED] Autocomplete dropdown SHOULD appear when user types ≥2 characters in chip input field on panelist selection page
- **FR-030g**: [DEFERRED] Autocomplete feature SHOULD degrade gracefully if Firestore unavailable (hide dropdown, allow manual chip creation)
- **FR-030h**: [DEFERRED] Autocomplete API SHOULD respond within 500ms; if slower, frontend shows loading indicator without blocking input

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

- User authentication and account management (debates are public via UUID URLs)
- User-specific debate history or saved debates (no user accounts)
- Editing or regenerating portions of completed debates
- Deleting or modifying debates after creation (Firestore documents are immutable)
- Debate analytics or usage tracking beyond basic metadata
- Sharing debates via social media integrations (users can copy/paste URLs manually)
- User customization of panelist avatars or bios
- Multi-language support (non-English topics)
- Real-time collaborative debate watching with multiple users
- Audio/video generation of debates
- Direct creation of custom panelists with user-defined bios (AI evaluates suggested names instead)
- User acting as moderator with ability to ask questions during debate (future enhancement)
- Payment or subscription features
- Admin panel or content moderation tools
- Debate search or discovery features
- TTL (Time-To-Live) policies for auto-deleting old debates (relying on Firestore free tier capacity)
