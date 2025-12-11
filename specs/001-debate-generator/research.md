# Phase 0 Research: AI-Powered Theology/Philosophy Debate Generator

**Feature**: 001-debate-generator  
**Date**: 2025-12-11  
**Purpose**: Resolve technical unknowns and establish best practices for implementation

## Research Topics

### 1. Claude API Streaming with Server-Sent Events (SSE)

**Decision**: Use Anthropic Claude API's streaming mode with Server-Sent Events (SSE) for real-time debate generation.

**Rationale**:
- Claude API supports streaming responses via SSE, enabling progressive text delivery
- SSE is simpler than WebSockets for unidirectional server-to-client streaming
- Native browser EventSource API provides built-in reconnection logic
- Lower overhead than WebSockets for this use case (no bidirectional communication needed)

**Implementation Details**:
- Backend: Use `anthropic.StreamingMessage` in Go SDK with `http.ResponseWriter` flushing
- Frontend: EventSource API to consume SSE stream from GCP Cloud Function
- Message format: JSON-encoded events with `{panelist_id, text_chunk, done}` structure
- Chunking strategy: Stream by sentence boundaries to maintain readability
- **Claude Response Format**: Claude streams responses using `[ID]: text` pattern where ID is panelist handle or "moderator"
  - Format: `[moderator]: Welcome to the debate\n[Augustine354]: Thank you`
  - Backend parses this format to extract speaker ID and message text
  - **Edge Case Handling**: Single SSE chunk may contain multiple `[ID]:` patterns if Claude sends rapid speaker changes
  - Solution: `findNextPattern()` helper scans for subsequent patterns after current message start
  - Each complete message is sent as separate chunk when new speaker detected
  - Final message in buffer flushed when stream ends or next speaker starts

**Alternatives Considered**:
- WebSockets: Rejected due to unnecessary complexity for one-way streaming
- Long polling: Rejected due to poor user experience and higher latency
- HTTP chunked encoding: Rejected as SSE provides better browser support and automatic retry

---

### 2. React Streaming UI Patterns

**Decision**: Use React state updates with `useEffect` and EventSource to progressively append debate chunks to chat bubbles.

**Rationale**:
- React's declarative model handles incremental UI updates efficiently
- Custom hook `useDebateStream` encapsulates SSE connection lifecycle
- Component re-renders are optimized via React.memo for chat bubbles
- Smooth UX without full page reloads

**Implementation Pattern**:
```javascript
// Custom hook for SSE streaming
const useDebateStream = (config) => {
  const [messages, setMessages] = useState([]);
  const [isStreaming, setIsStreaming] = useState(false);
  
  useEffect(() => {
    const eventSource = new EventSource(apiUrl);
    eventSource.onmessage = (event) => {
      const chunk = JSON.parse(event.data);
      setMessages(prev => appendChunk(prev, chunk));
    };
    return () => eventSource.close();
  }, [config]);
  
  return { messages, isStreaming };
};
```

**Best Practices**:
- Virtualize long debate lists if >100 messages (react-window)
- Use `key` prop with stable IDs for efficient reconciliation
- Debounce rapid state updates to avoid excessive re-renders
- Display typing indicator while waiting for next panelist

---

### 3. Client-Side PDF Generation with jsPDF

**Decision**: Generate PDFs client-side using jsPDF library with HTML content rendering.

**Rationale**:
- Eliminates backend dependency and reduces server costs
- User's browser handles compute/memory load
- Immediate download without waiting for server processing
- Modern browsers handle jsPDF efficiently for documents up to ~10MB

**Implementation Approach**:
```javascript
import jsPDF from 'jspdf';
import html2canvas from 'html2canvas';

// Generate PDF from debate HTML
const exportDebateToPDF = async (debateData) => {
  const pdf = new jsPDF('p', 'mm', 'a4');
  
  // Add debate header (topic, panelists, timestamp)
  pdf.setFontSize(16);
  pdf.text(debateData.topic, 20, 20);
  
  // Render chat bubbles as HTML, convert to canvas, add to PDF
  const debateElement = document.getElementById('debate-content');
  const canvas = await html2canvas(debateElement);
  const imgData = canvas.toDataURL('image/png');
  pdf.addImage(imgData, 'PNG', 10, 30, 190, 0);
  
  pdf.save(`debate-${Date.now()}.pdf`);
};
```

**Constraints**:
- Page break logic needed to avoid splitting individual responses
- Avatar images must be embedded as base64 or external URLs
- Font embedding for consistent cross-platform rendering
- Max recommended debate size: 5000-7000 words (~10-15 pages)

**Alternatives Considered**:
- Server-side PDF generation (Puppeteer/Go PDF libs): Rejected due to added complexity and cost
- Native browser print: Rejected due to poor styling control and user friction
- Third-party PDF services: Rejected to avoid external dependencies and privacy concerns

---

### 4. GCP Cloud Functions Architecture

**Decision**: Deploy 3 independent Cloud Functions (Gen 2) for topic validation, panelist suggestion, and debate generation.

**Rationale**:
- Independent scaling: Debate generation has longer execution time (30-60s) vs validation (<3s)
- Easier testing and deployment of individual functions
- Better cost optimization (pay only for actual compute time per function)
- Simplified error isolation and retry logic

**Function Specifications**:

**validate-topic**:
- Runtime: Go 1.23
- Memory: 256MB
- Timeout: 10s
- Trigger: HTTPS
- Cold start: ~500ms
- Estimated cost: <$0.001 per invocation

**suggest-panelists**:
- Runtime: Go 1.23
- Memory: 512MB (JSON parsing of 20 panelists)
- Timeout: 15s
- Trigger: HTTPS
- Cold start: ~800ms
- Estimated cost: <$0.002 per invocation

**generate-debate**:
- Runtime: Go 1.23
- Memory: 1GB (streaming buffer)
- Timeout: 60s (max Cloud Function Gen 2 limit)
- Trigger: HTTPS with SSE streaming
- Cold start: ~1.2s
- Estimated cost: ~$0.01 per invocation

**Best Practices**:
- Use Secret Manager for Claude API keys (never in code or env vars)
- Implement rate limiting with Cloud Memorystore or function-level counters
- Enable request logging to Cloud Logging for debugging
- Use structured logging (JSON format) for better query/analysis
- Set CORS headers for frontend domain

**Alternatives Considered**:
- Single monolithic function: Rejected due to scaling and timeout concerns
- Cloud Run: Rejected as unnecessary for stateless API proxy (Functions simpler)
- App Engine: Rejected due to higher costs and unnecessary complexity

---

### 5. Historical Figure Avatar Sources

**Decision**: Use combination of public domain images (Wikimedia Commons) and AI-generated avatars as fallback.

**Rationale**:
- Wikimedia Commons has extensive collection of historical figure portraits (public domain)
- AI-generated avatars (e.g., via DALL-E, Stable Diffusion) for figures without available images
- Consistent art style improves visual cohesion
- Local storage in `public/avatars/` for performance

**Implementation Strategy**:
1. Curate 50-100 common historical figures' avatars during development
2. Store as optimized PNGs (256x256px, <50KB each)
3. Naming convention: `{slug}-avatar.png` (e.g., `augustine-avatar.png`)
4. Fallback to placeholder if avatar not found
5. Lazy loading with `loading="lazy"` attribute

**Accessibility Considerations**:
- All avatars have descriptive alt text: "{Name}, {Era} {Title/Occupation}"
- Example: "Augustine of Hippo, 4th-5th century theologian and philosopher"
- High contrast against background for visual clarity

**Alternatives Considered**:
- Real-time image fetching from Wikimedia API: Rejected due to latency and reliability concerns
- Generic placeholder for all: Rejected as reduces engagement and educational value
- User-uploaded avatars: Out of scope for MVP

---

### 6. Input Sanitization and XSS Prevention

**Decision**: Implement multi-layer sanitization using DOMPurify (frontend) and Go html package (backend).

**Rationale**:
- Defense in depth: Sanitize at both entry point (backend) and before rendering (frontend)
- DOMPurify is industry-standard for client-side XSS prevention
- Go's `html.EscapeString` provides server-side protection
- Claude API outputs can potentially include malicious formatting if prompted adversarially

**Sanitization Points**:

**Backend (Go Cloud Functions)**:
```go
import "html"

// Sanitize user input before sending to Claude
func sanitizeInput(input string) string {
    // Strip HTML tags
    return html.EscapeString(input)
}
```

**Frontend (React)**:
```javascript
import DOMPurify from 'dompurify';

// Sanitize Claude API response before render
const sanitizeClaudeOutput = (text) => {
  return DOMPurify.sanitize(text, {
    ALLOWED_TAGS: ['b', 'i', 'em', 'strong', 'p', 'br'],
    ALLOWED_ATTR: []
  });
};
```

**Additional Security Measures**:
- Content Security Policy (CSP) headers: `default-src 'self'; script-src 'self'`
- Rate limiting: Max 10 requests per IP per minute (Cloud Armor or function-level)
- Input validation: Topic length 10-500 characters, panelist selection 2-5 only
- No eval() or dangerouslySetInnerHTML in React code

---

### 7. Accessibility Testing Strategy

**Decision**: Automated testing with axe-core + manual keyboard navigation testing.

**Rationale**:
- axe-core catches 57% of WCAG issues automatically (industry benchmark)
- Remaining 43% require manual testing (keyboard nav, screen reader, color contrast in context)
- Integrate axe-core into Jest tests for continuous validation

**Testing Checklist**:
- [ ] All interactive elements keyboard-accessible (Tab, Enter, Space)
- [ ] Focus indicators visible on all elements (2px outline, high contrast)
- [ ] Color contrast ≥4.5:1 for normal text, ≥3.1 for large text
- [ ] ARIA labels on all form inputs and buttons
- [ ] Screen reader announces state changes (validation results, loading states)
- [ ] Alt text on all avatars with descriptive content
- [ ] Skip navigation link for keyboard users
- [ ] No keyboard traps in modal dialogs

**Tools**:
- axe-core: Automated testing in Jest/React Testing Library
- WAVE browser extension: Manual page audits
- NVDA/VoiceOver: Screen reader testing
- Keyboard only: Test all flows without mouse

---

### 8. Local Development with Docker

**Decision**: Use Docker Compose for local development with containerized Cloud Functions and frontend.

**Rationale**:
- Eliminates "works on my machine" issues by standardizing environment
- Simplifies onboarding - single command to start entire stack
- Matches production environment more closely (Cloud Functions run in containers)
- Easier to test service-to-service communication
- No need to manually manage multiple terminal windows for 3+ services

**Implementation Strategy**:
```yaml
# docker-compose.yml structure
services:
  validate-topic:    # Port 8080
  suggest-panelists: # Port 8081
  generate-debate:   # Port 8082
  frontend:          # Port 3000 (nginx serves built React app)
```

**Development Workflow**:
- Each Go Cloud Function has its own Dockerfile (multi-stage build: golang:1.23-alpine → distroless)
- Frontend Dockerfile builds React app and serves via nginx
- docker-compose.yml orchestrates all services with proper networking
- Secrets managed via summon with GCP Secret Manager (same source as production)
- Quick start script (start-local.sh) for one-command launch with summon

**Alternatives Considered**:
- Manual terminal management: Rejected due to poor developer experience
- Tilt/Skaffold: Rejected as overkill for 4-service local dev
- VS Code devcontainers: Rejected to avoid IDE lock-in
- GCP Functions Framework emulator: Rejected due to limited local testing capabilities

**Benefits**:
- New developers productive in minutes, not hours
- Consistent behavior across macOS, Linux, Windows (via WSL2)
- Easy to add auxiliary services later (Redis, PostgreSQL if needed)
- CI/CD can use same Dockerfiles for deployment
- Unified secret management between local dev and production via summon

**Secret Management with Summon**:
- Production secrets stored in GCP Secret Manager
- Local development uses same secrets via summon + gcloud plugin
- summon-gcloud plugin installed at `/usr/local/lib/summon/gcloud`
- secrets.yml defines secret paths (e.g., `ANTHROPIC_API_KEY: gcp/secrets/anthropic-api-key`)
- Command: `summon -p gcloud docker-compose up`
- Eliminates .env file management and prevents secret leakage
- Ensures parity between local and production secret sources

---

### 9. Production Deployment Strategy

**Decision**: Deploy backend to GCP Cloud Functions (Gen 2) and frontend to GitHub Pages.

**Rationale**:
- **Backend (GCP Cloud Functions)**: Serverless, auto-scaling, pay-per-use, native Secret Manager integration
- **Frontend (GitHub Pages)**: Free static hosting, global CDN, HTTPS by default, zero maintenance
- Separation of concerns: Backend can scale independently of frontend
- Cost-effective: Both services have generous free tiers suitable for MVP traffic

**Deployment Architecture**:
```
┌─────────────────┐
│  GitHub Pages   │  ← Static React app (https://raphink.github.io/debate)
│  (Frontend)     │
└────────┬────────┘
         │ HTTPS API calls
         ↓
┌─────────────────────────────────────────────┐
│  GCP Cloud Functions (europe-west1)         │
│  ┌─────────────────────────────────────┐   │
│  │ validate-topic (Go 1.23)            │   │
│  │ suggest-panelists (Go 1.23)         │   │
│  │ generate-debate (Go 1.23)           │   │
│  └─────────────────────────────────────┘   │
└────────┬────────────────────────────────────┘
         │ Secret access
         ↓
┌─────────────────┐
│ GCP Secret Mgr  │  ← ANTHROPIC_API_KEY
└─────────────────┘
```

**Deployment Process**:
1. **Backend**: Use `gcloud functions deploy` with Secret Manager binding
2. **Frontend**: Build React app with production API URLs → deploy to gh-pages branch
3. **Secrets**: Already stored in GCP Secret Manager (same as local dev)

**Infrastructure Details**:
- Cloud Functions Gen 2 (better performance, more control than Gen 1)
- Runtime: Go 1.23
- Region: europe-west1 (Belgium - lowest latency for EU users)
- Memory: 256MB per function (sufficient for API proxy)
- Timeout: 60s (max for streaming responses)
- Concurrency: 80 requests per instance (Cloud Functions default)
- Min instances: 0 (scale to zero when idle)
- Max instances: 100 (prevent runaway costs)

**Frontend Configuration**:
- GitHub Pages serves from `gh-pages` branch
- React app built with production environment variables
- API URLs point to Cloud Functions endpoints
- Service Worker for offline capabilities (future enhancement)

**Cost Estimation** (Monthly for MVP traffic):
- Cloud Functions: Free tier covers ~2M requests
- Secret Manager: Free tier covers 10K accesses
- GitHub Pages: Free for public repos
- Anthropic API: ~$0.01-0.02 per debate (usage-based)
- **Estimated total**: $0-10/month for low-medium traffic

**Deployment Script**: `deploy.sh` in project root automates entire process

**Monitoring & Observability**:
- Cloud Functions logs → Cloud Logging (automatic)
- Error reporting → Cloud Error Reporting (automatic)
- Metrics → Cloud Monitoring (automatic)
- Custom metrics: Track debate completion rate, API latency

**Alternatives Considered**:
- **Vercel/Netlify**: Good alternatives, but GitHub Pages is simpler and free
- **Cloud Run**: More complex than Functions, unnecessary for stateless API proxy
- **Firebase Hosting**: Considered, but GitHub Pages offers same benefits with less vendor lock-in
- **Self-hosted VPS**: Rejected due to maintenance overhead and scaling complexity

---

## Summary of Key Decisions

| Area | Decision | Rationale |
|------|----------|-----------|
| Streaming | Server-Sent Events (SSE) | Simpler than WebSockets, native browser support |
| PDF Generation | Client-side jsPDF | Reduces server load, immediate download |
| Backend Architecture | 3 independent Cloud Functions | Better scaling, cost optimization, error isolation |
| Avatar Strategy | Public domain + AI-generated | Legal compliance, consistent styling |
| Security | Multi-layer sanitization (DOMPurify + Go) | Defense in depth against XSS |
| Accessibility | axe-core + manual testing | Meets WCAG 2.1 Level AA requirements |
| Local Development | Docker Compose + summon | Standardized environment, production parity |
| Production Deployment | GCP Functions + GitHub Pages | Serverless backend, free frontend, minimal ops |

All technical unknowns resolved. Ready to proceed to Phase 1 (data model and contracts design).
