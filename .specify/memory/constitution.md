<!--
SYNC IMPACT REPORT
==================
Version Change: Initial → 1.0.0
Bump Rationale: MAJOR - Initial constitution establishing core governance principles
Added Principles:
  - I. User-Centric Design (UX First)
  - II. Code Quality & Maintainability
  - III. Responsive & Accessible UI
  - IV. Interactive & Performant Experience
  - V. AI Safety & Security (NON-NEGOTIABLE)
Added Sections:
  - Quality Standards
  - Development Workflow
Templates Requiring Updates:
  ✅ plan-template.md - Verified compatibility with constitution check structure
  ✅ spec-template.md - Verified alignment with user story and requirement format
  ✅ tasks-template.md - Verified compatibility with quality gates and testing discipline
  ⚠️ All command files - Should be reviewed for consistency with new principles
Follow-up TODOs: None
==================
-->

# Debate Project Constitution

## Core Principles

### I. User-Centric Design (UX First)

Every feature MUST prioritize end-user experience and usability over technical convenience. User feedback loops and usability testing are mandatory before feature completion.

**Rules**:
- User scenarios MUST be defined and prioritized (P1, P2, P3) before technical design begins
- Each user story MUST be independently testable and deliver standalone value
- UI/UX decisions MUST be justified with user research, accessibility guidelines, or industry best practices
- Error messages MUST be user-friendly, actionable, and never expose technical implementation details
- Feature acceptance REQUIRES user validation or usability testing completion

**Rationale**: Users judge software by their experience, not by its internal architecture. Poorly designed UX leads to user frustration, support burden, and feature abandonment regardless of code quality.

### II. Code Quality & Maintainability

Code MUST be self-documenting, testable, and maintainable. Technical debt requires explicit justification and a remediation plan.

**Rules**:
- All code MUST pass linting and formatting standards before commit
- Functions MUST have single, clear responsibilities (Single Responsibility Principle)
- Magic numbers and hardcoded values MUST be replaced with named constants
- Complex logic MUST include inline comments explaining the "why" not the "what"
- Code duplication beyond 3 lines REQUIRES refactoring into shared utilities
- Dependencies MUST be explicitly declared with version pinning
- Breaking changes REQUIRE version bumps following semantic versioning (MAJOR.MINOR.PATCH)

**Rationale**: Code is read 10x more than it is written. Maintainability directly impacts development velocity, bug reduction, and onboarding efficiency. Poor code quality creates compounding costs over time.

### III. Responsive & Accessible UI

User interfaces MUST be responsive across devices and accessible to users with disabilities, adhering to WCAG 2.1 Level AA standards.

**Rules**:
- All interactive elements MUST be keyboard-navigable
- Color MUST NOT be the only means of conveying information (color-blind accessibility)
- Text MUST maintain minimum contrast ratios: 4.5:1 for normal text, 3:1 for large text
- UI MUST be responsive and functional on mobile (≥375px width), tablet, and desktop viewports
- Images MUST include alt text; decorative images MUST use empty alt attributes
- Form inputs MUST have associated labels (visible or aria-label)
- Focus indicators MUST be clearly visible for keyboard navigation

**Rationale**: Approximately 15% of the global population has some form of disability. Inaccessible UIs exclude users, create legal liability, and violate ethical design principles. Responsive design ensures usability across the growing diversity of devices.

### IV. Interactive & Performant Experience

User interfaces MUST provide immediate feedback and maintain performance standards to ensure smooth, responsive interactions.

**Rules**:
- User actions MUST provide immediate visual feedback (loading states, button press states)
- Page initial load MUST complete within 3 seconds on 3G connections
- Interactive elements MUST respond within 100ms (perceived as instant)
- Long-running operations (>1s) MUST display progress indicators
- UI animations MUST run at 60fps; non-essential animations MUST be cancelable
- API calls MUST implement timeout and retry logic with user-visible error states
- State changes MUST be reflected in the UI within one event loop cycle

**Rationale**: Users perceive delays over 100ms as sluggish and abandon applications that feel unresponsive. Performance is a feature, not an afterthought. Poor interactivity creates frustration and damages trust.

### V. AI Safety & Security (NON-NEGOTIABLE)

AI-generated content, AI model interactions, and user data handling MUST implement security controls, safety guardrails, and privacy protections. Security vulnerabilities are CRITICAL and block all releases.

**Rules**:
- AI model outputs MUST be sanitized before display to prevent XSS and injection attacks
- User inputs to AI systems MUST be validated and rate-limited to prevent abuse
- AI prompts MUST NOT include unsanitized user data or secrets
- Personal Identifiable Information (PII) MUST NOT be logged or sent to external AI services without explicit consent
- AI-generated code MUST be reviewed for security vulnerabilities before merge
- Authentication tokens and API keys MUST be stored in secure vaults, never in code or logs
- All data transmission MUST use TLS 1.3+ encryption
- User data MUST be encrypted at rest using industry-standard algorithms (AES-256 or equivalent)
- Security dependencies MUST be automatically scanned for CVEs; critical vulnerabilities block deployment
- Input validation MUST be performed on both client and server to prevent injection attacks

**Rationale**: AI systems introduce new attack vectors including prompt injection, data poisoning, and unintended information disclosure. Security breaches destroy user trust, create legal liability, and can cause irreparable reputation damage. Safety is non-negotiable.

## Quality Standards

All code submissions MUST meet the following quality gates before merge approval:

- **Code Review**: At least one peer review approval required
- **Automated Tests**: All existing tests MUST pass; new features MUST include tests
- **Linting**: Zero linting errors; warnings require justification
- **Security Scan**: No critical or high-severity vulnerabilities
- **Accessibility**: UI changes MUST pass automated accessibility checks (axe, pa11y, or equivalent)
- **Performance**: No regressions in page load time or interaction responsiveness

Exceptions to quality gates REQUIRE explicit written justification and technical lead approval.

## Development Workflow

### Feature Development Process

1. **Specification** (`/speckit.specify`): Define user stories with priorities and acceptance criteria
2. **Planning** (`/speckit.plan`): Create technical implementation plan with constitution check
3. **Task Breakdown** (`/speckit.tasks`): Generate dependency-ordered task list organized by user story
4. **Analysis** (`/speckit.analyze`): Validate consistency across spec, plan, and tasks
5. **Implementation** (`/speckit.implement`): Execute tasks in phases with continuous testing
6. **Review & Merge**: Code review, quality gates, and merge to main branch

### Constitution Compliance

- All feature plans MUST include a "Constitution Check" section verifying adherence to principles
- Violations MUST be justified in a "Complexity Tracking" table with simpler alternatives documented
- Features that fundamentally violate principles require constitution amendment approval before proceeding

### Amendment Process

- Constitution amendments require creation of a formal proposal with rationale
- Amendments MUST use semantic versioning: MAJOR (backward-incompatible changes), MINOR (new principles), PATCH (clarifications)
- Amendments MUST trigger updates to all dependent templates and command files
- A "Sync Impact Report" MUST be generated documenting all affected artifacts

## Governance

This constitution is the authoritative source for all development practices in the Debate project. It supersedes conflicting guidance in other documents.

**Enforcement**:
- All pull requests MUST verify constitution compliance during code review
- Automated checks SHOULD enforce linting, security scanning, and accessibility standards
- Non-compliance blocks merge until resolved or formally justified

**Guidance Integration**:
- Runtime development guidance is maintained in `.specify/templates/agent-file-template.md`
- Agent-specific workflows are defined in `.github/agents/speckit.*.agent.md`
- All guidance documents MUST align with this constitution

**Review Cycle**:
- Constitution SHOULD be reviewed quarterly for relevance and effectiveness
- User feedback and incident retrospectives MAY trigger out-of-cycle reviews

**Version**: 1.0.0 | **Ratified**: 2025-12-11 | **Last Amended**: 2025-12-11
