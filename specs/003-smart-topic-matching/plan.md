# Implementation Plan: Smart Topic Matching

**Branch**: `003-smart-topic-matching` | **Date**: 2025-12-14 | **Spec**: [spec.md](spec.md)

## Summary

Improve topic autocomplete matching to find semantically similar topics even when exact wording differs. Uses simplified KISS approach: tokenize queries/topics (≥3 char words), match query tokens against topics (bag-of-words with partial match support), sort by match count then recency. No stop word lists, no complex scoring algorithms.

## Technical Context

**Language/Version**: Go 1.24.0  
**Primary Dependencies**: 
- `cloud.google.com/go/firestore` (existing)
- Go stdlib only for new code (strings, regexp, sort)

**Storage**: Firestore (existing `debates` collection, no schema changes)  
**Testing**: Go testing (`go test`)  
**Target Platform**: Google Cloud Functions (existing `list-debates` function)  
**Project Type**: Web application (backend enhancement only, no frontend changes)  
**Performance Goals**: <50ms added latency for normalization/matching (total <500ms)  
**Constraints**: 
- Backwards compatible API (no breaking changes)
- Must work with existing 50-debate fetch limit
- Firestore read quota unchanged

**Scale/Scope**: Enhancement to existing function, ~200 lines of new code

## Constitution Check

✅ **I. User-Centric Design (UX First)**
- Improves user experience by finding more relevant debates
- No UI changes required (transparent backend improvement)
- Error handling maintains graceful degradation

✅ **II. Code Quality & Maintainability**
- KISS principle applied (rejected complex scoring)
- No external dependencies
- Clear unit test strategy
- Single responsibility functions

✅ **III. Responsive & Accessible UI**
- N/A - Backend only, no UI changes

✅ **IV. Interactive & Performant Experience**
- <50ms overhead target (NFR-001)
- In-memory operations only (no additional network calls)

✅ **V. AI Safety & Security (NON-NEGOTIABLE)**
- Uses existing sanitization (sanitize.StripHTML)
- No new attack vectors introduced
- No PII or sensitive data logged

**Result**: ✅ All applicable principles satisfied

## Project Structure

### Documentation (this feature)

```text
specs/003-smart-topic-matching/
├── spec.md              # Feature specification (completed with clarifications)
├── plan.md              # This file
└── tasks.md             # Task breakdown
```

### Source Code

```text
backend/functions/list-debates/
├── firestore.go         # MODIFY: Update autocompleteDebates()
├── handler.go           # No changes
├── types.go             # No changes
├── normalize.go         # NEW: Normalization and matching utilities
├── normalize_test.go    # NEW: Unit tests for normalization
└── go.mod               # No changes (stdlib only)
```

**Structure Decision**: Extends existing `list-debates` Cloud Function. New code isolated in `normalize.go` with comprehensive unit tests. No changes to API contract, request/response formats, or frontend integration.

## Implementation Strategy

### Phase 1: Normalization Foundation

Create `normalize.go` with:

```go
const minTokenLength = 3

// NormalizeAndTokenize converts text to array of significant tokens
// Steps: lowercase → replace hyphens/slashes → remove punctuation → 
//        split on whitespace → filter ≥3 chars
func NormalizeAndTokenize(text string) []string

// CountMatchingTokens returns count of query tokens found in topic tokens
// Returns weight based on number of matches (0 if no matches, partial matches supported)
func CountMatchingTokens(queryTokens, topicTokens []string) int
```

**Processing Pipeline**:
1. Convert to lowercase
2. Replace `-` and `/` with spaces
3. Remove punctuation: `[^a-z0-9\s]`
4. Split on whitespace
5. Filter tokens: keep only len(token) ≥ 3

### Phase 2: Autocomplete Integration

Update `firestore.go::autocompleteDebates()`:

**Before**:
```go
if strings.Contains(strings.ToLower(topicText), queryLower) {
    matches = append(matches, debate)
}
```

**After**:
```go
queryTokens := NormalizeAndTokenize(query)
if len(queryTokens) == 0 {
    return []DebateSummary{}, nil // Empty query
}

type matchWithWeight struct {
    debate DebateSummary
    weight int
}

var matches []matchWithWeight

for each debate {
    topicTokens := NormalizeAndTokenize(topicText)
    weight := CountMatchingTokens(queryTokens, topicTokens)
    if weight > 0 {
        matches = append(matches, matchWithWeight{debate, weight})
    }
}

// Sort by weight DESC, then startedAt DESC
sort.Slice(matches, func(i, j int) bool {
    if matches[i].weight != matches[j].weight {
        return matches[i].weight > matches[j].weight
    }
    return matches[i].debate.StartedAt.After(matches[j].debate.StartedAt)
})

// Extract top 10
results := make([]DebateSummary, 0, min(10, len(matches)))
for i := 0; i < len(matches) && i < 10; i++ {
    results = append(results, matches[i].debate)
}
```

### Phase 3: Testing & Deployment

**Unit Tests** (`normalize_test.go`):
- Lowercase conversion
- Hyphen/slash replacement
- Punctuation removal
- Token filtering (≥3 chars)
- Edge cases (empty, all short words)
- Matching token count accuracy

**Integration Tests** (extend `firestore_test.go`):
- "do animals" matches "Should animals have rights?"
- "climate-change" matches "Climate change effects"
- Weight-based sorting
- Recency as tie-breaker

**Deployment**:
- No schema changes
- No frontend changes
- Zero downtime deployment
- Monitor latency metrics

## Complexity Tracking

No constitution violations. KISS principle successfully applied.

## Success Criteria

1. ✅ Unit tests pass with >80% coverage
2. ✅ "do animals" finds "Should animals have rights?"
3. ✅ Hyphenated queries work ("climate-change")
4. ✅ Response time <500ms (≤50ms added overhead)
5. ✅ No API contract changes
6. ✅ Graceful degradation on edge cases
