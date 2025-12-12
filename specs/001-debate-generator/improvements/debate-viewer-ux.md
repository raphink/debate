# Specification: Debate Viewer UX Improvements

**Feature**: Enhanced Debate Viewer Interface  
**Priority**: P2 (UX Polish)  
**Created**: 2025-12-12

## Problem Statement

The debate viewer page (`/d/{uuid}`) currently displays completed debates using the same interface as live debate generation. This creates UX issues:

1. **Auto-scroll control is unnecessary** - debates are already complete, so there's no new content to scroll to
2. **No navigation to create new debates** - users must manually navigate back to home or use browser back button

## User Requirements

### FR-001: Remove Auto-Scroll Toggle for Cached Debates
**As a** user viewing a cached debate  
**I want** the auto-scroll control hidden  
**So that** I'm not presented with irrelevant UI controls

**Acceptance Criteria:**
- Auto-scroll toggle is NOT visible on `/d/{uuid}` (DebateViewer)
- Auto-scroll toggle IS visible during live generation (DebateGeneration)
- Chat messages are fully scrolled to bottom by default in viewer (no toggle needed)

### FR-002: Add "New Debate" Navigation Button
**As a** user viewing a cached debate  
**I want** a "New Debate" button at the bottom of the page  
**So that** I can easily create a new debate without manual navigation

**Acceptance Criteria:**
- "New Debate" button appears at bottom of debate viewer
- Button navigates to home page (`/`)
- Button uses existing design system styles (gradient, hover effects)
- Button is NOT shown during live generation (DebateGeneration page has its own controls)

## Success Criteria

- SC-001: Viewer mode has simplified UI without streaming controls
- SC-002: "New Debate" button visible and functional in viewer mode
- SC-003: Live generation mode retains full streaming controls
- SC-004: Component reuse maintained (DebateView shared between both modes)

## Out of Scope

- Additional viewer controls (search, jump to message, etc.)
- Debate editing or commenting features
- Social sharing features beyond existing ShareButton

## Technical Constraints

- DebateView component must remain shared between DebateViewer and DebateGeneration
- No breaking changes to existing DebateView props API
- Maintain accessibility standards for navigation

## User Journey

1. User clicks shared debate link `/d/{uuid}`
2. DebateViewer loads and displays complete debate
3. User scrolls through conversation (no auto-scroll toggle visible)
4. User reaches bottom of debate
5. User sees "New Debate" button
6. User clicks button → navigates to home page
7. User creates new debate

## Edge Cases

- Very short debates (1-2 messages): "New Debate" button still visible at bottom
- Mobile viewports: Button remains accessible and properly styled
- Error states (404, network failure): "New Debate" button in error message (already implemented)
