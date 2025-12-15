---
title: "Smart Topic Matching - Task List"
description: "Task breakdown for improving autocomplete topic matching (KISS approach)"
created: 2025-12-14
---

# Tasks: Smart Topic Matching (US6 Enhancement)

## Overview

Implement normalized topic matching using KISS approach: tokenize queries/topics (≥3 char words), bag-of-words matching, sort by match count then recency. No stop word lists, no complex scoring.

**Organization**: Sequential implementation with clear testing checkpoints

**Conventions**:
- **[P]**: Can be done in parallel with adjacent tasks
- **Status**: `[ ]` Not started, `[~]` In progress, `[X]` Complete
- **Files**: Absolute paths from repo root
- **Dependencies**: Listed explicitly when task requires previous completion

---

## Phase 1: Normalization Foundation 🔧

### Backend Utilities

- [X] T001 Create `backend/functions/list-debates/normalize.go` with constant `minTokenLength = 3` and function `NormalizeAndTokenize(text string) []string` that: converts to lowercase, replaces hyphens/slashes with spaces, removes punctuation using regex `[^a-z0-9\s]+`, splits on whitespace, filters tokens to keep only len(token) ≥ 3

- [X] T002 Add `CountMatchingTokens(queryTokens, topicTokens []string) int` function to normalize.go that counts how many query tokens appear in topic tokens (bag-of-words), returning the count as weight (partial matches supported)

### Unit Tests

- [X] T003 Create `backend/functions/list-debates/normalize_test.go` with `TestNormalizeAndTokenize` covering: lowercase conversion, hyphen replacement ("climate-change" → ["climate", "change"]), slash replacement ("AI/ML" → ["ai", "ml"]), punctuation removal, token filtering (≥3 chars), edge cases (empty string, all short words like "is it ok", special chars only)

- [X] T004 Add `TestCountMatchingTokens` to normalize_test.go verifying: all query tokens found = full count, partial matches return partial count, bag-of-words matching (order-independent), empty inputs

---

## Phase 2: Autocomplete Integration 🔌

### Modify Existing Code

- [X] T005 [Depends on T001-T002] Update `autocompleteDebates()` in `backend/functions/list-debates/firestore.go` to call `queryTokens := NormalizeAndTokenize(query)` before loop and return empty slice if len(queryTokens) == 0

- [X] T006 Add struct `type matchWithWeight struct { debate DebateSummary; weight int }` inside `autocompleteDebates()` function and change matches slice to `var matches []matchWithWeight`

- [X] T007 Update matching loop in `autocompleteDebates()` to: call `topicTokens := NormalizeAndTokenize(topicText)` for each debate, call `weight := CountMatchingTokens(queryTokens, topicTokens)`, append to matches if weight > 0 with `matches = append(matches, matchWithWeight{debate, weight})`

- [X] T008 Add sorting logic after loop using `sort.Slice(matches, ...)` to sort by weight DESC (primary), then startedAt DESC (secondary tie-breaker)

- [X] T009 Update result extraction to: create `results := make([]DebateSummary, 0, min(10, len(matches)))`, loop through matches extracting just debate field for top 10, return results

### Integration Tests

- [ ] T010 [Depends on T005-T009] Add test cases to `backend/functions/list-debates/firestore_test.go` verifying: "do animals" finds "Should animals have rights?" (short word filtered, match found), "climate-change policy" matches "Climate change and policy reform" (hyphen split, bag-of-words), "AI dangers" finds topics with both "dangers" (or "risks" won't match - no synonyms)

- [ ] T011 Add weight-based sorting test to firestore_test.go: create debates with varying token matches, verify higher weight appears first regardless of recency

- [ ] T012 Add recency tie-breaker test to firestore_test.go: create debates with identical match weights but different startedAt, verify newer appears first

- [ ] T013 Add edge case tests to firestore_test.go covering: query with all short words ("is it ok"), query with punctuation only, empty topic text

---

## Phase 3: Testing & Validation ✅

### Local Testing

- [ ] T014 [P] Run all unit tests locally with `cd backend/functions/list-debates && go test -v` and verify 100% pass rate

- [ ] T015 [P] Test performance locally by running autocomplete queries 100 times and measuring average response time (target: <50ms normalization overhead)

### Manual Testing

- [ ] T016 Start local development environment with `./start-local.sh` and test autocomplete in browser UI with queries: "do animals", "climate-change", "AI/ML risks", "should we ban"

- [ ] T017 Verify autocomplete dropdown shows results in correct order (higher match count first, then recency for ties)

- [ ] T018 Verify edge cases: query with all 1-2 char words shows no results, hyphenated compounds match successfully

---

## Phase 4: Deployment 🚀

### Pre-Deployment

- [ ] T019 Run `go mod tidy` in list-debates directory to verify dependencies (stdlib only, no new deps expected)

- [ ] T020 Run `go build` in list-debates directory to verify no compilation errors

- [ ] T021 Review code changes and ensure no breaking changes to API contract (response format unchanged)

### Staging Deployment

- [ ] T022 Deploy to staging environment using existing deployment script

- [ ] T023 Test autocomplete on staging with real Firestore data: verify \"do animals\" finds animal debates, verify hyphenated queries work

- [ ] T024 Monitor Cloud Function logs on staging for errors or unexpected behavior

- [ ] T025 Measure response time on staging and verify <500ms total (NFR-001 from US6) with <50ms added overhead

### Production Deployment

- [ ] T026 Deploy to production environment using existing deployment script

- [ ] T027 Monitor Cloud Function logs for first 1 hour after deployment

- [ ] T028 Monitor Cloud Function metrics (invocations, errors, latency) in Google Cloud Console

- [ ] T029 Verify Firestore read usage is unchanged (still ~50 reads per autocomplete query)

### Post-Deployment Validation

- [ ] T030 Test production autocomplete with various queries from different browsers/locations

- [ ] T031 Verify user-facing behavior: autocomplete finds more relevant debates with varied phrasing

- [ ] T032 Document any discovered edge cases for future iteration

---

## Task Summary

**Total Tasks**: 32
**Phases**: 4 (Foundation → Integration → Testing → Deployment)
**Parallel Opportunities**: T014-T015 (local testing)

**Critical Path**: T001-T004 → T005-T009 → T010-T013 → T014-T018 → T019-T032

**Estimated Effort**:
- Phase 1: ~2 hours (normalization + tests, simplified KISS approach)
- Phase 2: ~2 hours (integration + tests)
- Phase 3: ~1.5 hours (validation)
- Phase 4: ~2 hours (deployment + monitoring)
- **Total**: ~7.5 hours

**Dependencies**:
- No external dependencies (uses existing Go stdlib)
- No Firestore schema changes
- No frontend changes required
- Backwards compatible with existing autocomplete API
