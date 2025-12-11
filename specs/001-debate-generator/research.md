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

## Summary of Key Decisions

| Area | Decision | Rationale |
|------|----------|-----------|
| Streaming | Server-Sent Events (SSE) | Simpler than WebSockets, native browser support |
| PDF Generation | Client-side jsPDF | Reduces server load, immediate download |
| Backend Architecture | 3 independent Cloud Functions | Better scaling, cost optimization, error isolation |
| Avatar Strategy | Public domain + AI-generated | Legal compliance, consistent styling |
| Security | Multi-layer sanitization (DOMPurify + Go) | Defense in depth against XSS |
| Accessibility | axe-core + manual testing | Meets WCAG 2.1 Level AA requirements |

All technical unknowns resolved. Ready to proceed to Phase 1 (data model and contracts design).
