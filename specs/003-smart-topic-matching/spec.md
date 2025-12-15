# Feature Specification: Smart Topic Matching for Autocomplete (US6 Enhancement)

## Clarifications

### Session 2025-12-14

- Q: Should we use complex relevance scoring with weights for position, length ratio, and consecutive matches? → A: No - use simple matching token count as weight, sort by count DESC then recency DESC (KISS principle)
- Q: Should we maintain stop word lists (question starters, fillers) for normalization? → A: No - use length-based filtering (≥3 characters) instead of maintaining stop word lists
- Q: Should matching require consecutive words (substring) or support any order? → A: Bag-of-words matching with partial match support - query tokens can match in any order, weight based on count of matching tokens (more flexible)
- Q: What minimum character length for tokens? → A: Always use ≥3 char threshold (simpler code, catches important short words like "ban", "war", "AI")
- Q: How to handle hyphenated/compound words (e.g., "climate-change", "AI/ML")? → A: Replace hyphens and slashes with spaces before tokenizing (treats compounds as separate words for better matching)

## Overview

**Feature ID**: 003-smart-topic-matching  
**Priority**: P2 (Quality improvement for existing feature)  
**Status**: Planning  
**Dependencies**: Requires 002-topic-autocomplete (US6) - extends existing autocomplete functionality

## Summary

Improve topic autocomplete matching to find semantically similar topics even when the exact wording differs. Users searching for "do animals" should see results for "should animals" since they represent the same core topic. This enhancement makes topic discovery more effective by normalizing and comparing topics based on their semantic content rather than exact text matching.

## Context

Currently, the autocomplete feature (US6) uses exact substring matching:
- User types "do animals"
- System searches for topics containing "do animals" 
- Topics like "Should animals have rights?" are NOT found
- This creates poor user experience when topics exist with similar meaning but different phrasing

**Current Implementation**: `backend/functions/list-debates/firestore.go` line 125-189
- Uses `strings.Contains(strings.ToLower(topicText), queryLower)` for matching
- Simple case-insensitive substring search
- No normalization or semantic comparison

## User Story

**As a** user creating a new debate,  
**I want** autocomplete to find topics regardless of their exact phrasing (e.g., "do animals" matches "should animals"),  
**So that** I can discover existing debates on the same topic without needing to guess the exact wording used.

## Acceptance Scenarios

1. **Given** user types "do animals", **When** autocomplete searches, **Then** system returns debates with topics like "Should animals have rights?" or "Do animals feel emotions?"

2. **Given** user types "should we ban", **When** autocomplete searches, **Then** system returns debates with topics like "Can we ban plastics?" or "Should governments ban social media?"

3. **Given** user types "AI dangers", **When** autocomplete searches, **Then** system returns debates with topics like "Is artificial intelligence dangerous?" or "Are there risks in AI?"

4. **Given** user types with common filler words (is, are, the, should, can, do, does, will), **When** autocomplete normalizes query, **Then** system matches topics focusing on content keywords rather than filler words

5. **Given** user types with punctuation or special characters, **When** autocomplete normalizes query, **Then** system strips punctuation and matches based on words only

6. **Given** multiple topics match after normalization, **When** autocomplete returns results, **Then** system orders by relevance (best match first) then by recency (startedAt)

7. **Given** autocomplete finds no matches after normalization, **When** query fails, **Then** system falls back gracefully showing no results (existing behavior)

## Functional Requirements

### Topic Normalization (Simplified KISS Approach)

- **FR-001**: System MUST normalize both user query and stored topics before comparison
- **FR-002**: System MUST convert all text to lowercase
- **FR-003**: System MUST replace hyphens (-) and slashes (/) with spaces to split compound words (e.g., "climate-change" → "climate change", "AI/ML" → "AI ML")
- **FR-004**: System MUST remove all punctuation (?, !, ., ,, ;, :, ', ", etc.)
- **FR-005**: System MUST split text into tokens (words) by whitespace
- **FR-006**: System MUST filter tokens to keep only words with ≥3 characters
- **FR-007**: System MUST preserve filtered tokens as array for matching (no specific word order required)

### Matching Algorithm (Bag-of-Words)

- **FR-008**: System MUST use bag-of-words matching: query tokens can match topic tokens in any order (order-independent)
- **FR-009**: System MUST calculate match weight as: count of query tokens that appear in topic tokens
- **FR-010**: System MUST skip debates where weight is 0 (no matching tokens found)

### Ranking & Results

- **FR-011**: System MUST order results by match weight (DESC), then by recency/startedAt (DESC) for ties
- **FR-012**: System MUST limit results to 10 debates max (existing constraint)
- **FR-013**: System MUST fetch last 50 debates for filtering (existing constraint)

### Backwards Compatibility

- **FR-014**: System MUST maintain existing API contract (no breaking changes to response format)
- **FR-015**: System MUST maintain existing query parameter validation (≥3 chars, ≤500 chars)
- **FR-016**: System MUST maintain existing error handling and graceful degradation
- **FR-017**: System MUST maintain existing CORS and security behavior

## Non-Functional Requirements

- **NFR-001**: Normalization and matching SHOULD add <50ms to autocomplete response time
- **NFR-002**: System SHOULD log normalized query tokens for debugging (not raw user input)
- **NFR-003**: Code MUST have clear unit tests for normalization and matching logic
- **NFR-004**: Token length threshold (≥3 chars) SHOULD be defined as named constant for maintainability

## Technical Approach

### Implementation Strategy (Simplified KISS)

1. **Create normalization utility** in `backend/functions/list-debates/` (e.g., `normalize.go`)
   - `NormalizeAndTokenize(text string) []string` - returns array of significant tokens (≥3 chars)
   - Process: lowercase → replace hyphens/slashes with spaces → remove punctuation → split by whitespace → filter by length
   
2. **Update `autocompleteDebates` function** in `firestore.go`
   - Normalize query once: `queryTokens := NormalizeAndTokenize(query)`
   - For each debate, normalize topic and count matching tokens
   - Sort by match count (DESC), then recency (DESC)
   
3. **Add matching helper function**
   - `CountMatchingTokens(queryTokens, topicTokens []string) int` - returns count of query tokens found in topic
   - Returns 0 if not all query tokens are present (failed match)

### Example Transformations

| Original Query | Tokens (≥3 chars) | Matches Topic | Topic Tokens | Match Weight |
|---------------|-------------------|---------------|--------------|--------------|
| "do animals" | ["animals"] | "Should animals have rights?" | ["should", "animals", "have", "rights"] | 1 |
| "should we ban" | ["should", "ban"] | "Can we ban plastics?" | ["can", "ban", "plastics"] | 1 |
| "climate-change policy" | ["climate", "change", "policy"] | "Climate change and policy reform" | ["climate", "change", "and", "policy", "reform"] | 3 |
| "AI/ML risks" | ["risks"] | "Are there risks in AI?" | ["are", "there", "risks"] | 1 |

## Edge Cases & Error Handling

- **EC-001**: Query normalizes to empty array (all words <3 chars): Return empty results gracefully
- **EC-002**: Query has only 1-2 character words (e.g., "is AI ok?"): Keep only "AI" token, attempt match
- **EC-003**: Topic text is empty or malformed: Skip that debate, continue processing
- **EC-004**: All 50 debates match after normalization: Return top 10 by weight+recency
- **EC-005**: Multiple debates have identical match weights: Sort by recency (startedAt DESC)
- **EC-006**: Hyphenated topic like "self-driving cars": Becomes ["self", "driving", "cars"] for matching

## Out of Scope

- Advanced NLP or machine learning for semantic similarity
- Multi-language support (English only for now)
- Stemming or lemmatization (e.g., "running" → "run")
- Synonym matching (e.g., "cars" = "automobiles")
- Fuzzy matching for typos (e.g., "animls" → "animals")
- Full-text search indexing (Firestore limitation)
- Caching of normalized topics in Firestore (compute on-the-fly)
- Complex relevance scoring with position/length/consecutive match weights (rejected in favor of KISS)
- Stop word lists maintenance (rejected in favor of length-based filtering)

## Testing Strategy

### Unit Tests Required

1. `normalize_test.go`:
   - Test lowercase conversion
   - Test hyphen/slash replacement ("climate-change", "AI/ML")
   - Test punctuation removal
   - Test tokenization by whitespace
   - Test length filtering (≥3 chars)
   - Test edge cases (empty strings, all short words, special chars only)

2. `firestore_test.go` (extend existing):
   - Test bag-of-words matching (all query tokens must appear in topic)
   - Test match weight calculation (count of matching tokens)
   - Test result ordering (weight DESC, then startedAt DESC)
   - Test partial matches rejected (not all query tokens present)

### Integration Tests

- Test real autocomplete queries through handler
- Verify "do animals" finds "Should animals have rights?"
- Verify "climate-change policy" matches topics with both "climate" and "policy"
- Verify API response time stays within NFR-001 (<50ms overhead)

## Success Metrics

- Users find relevant debates with varied query phrasing (improved recall)
- Autocomplete returns results for "do animals" when "Should animals have rights?" exists
- Hyphenated/compound terms match successfully ("climate-change" finds "climate change")
- No degradation in autocomplete response time (<500ms total)
- Zero breaking changes to API contract or frontend behavior
