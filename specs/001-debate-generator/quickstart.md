# Developer Quickstart: AI-Powered Theology/Philosophy Debate Generator

**Feature**: 001-debate-generator  
**Last Updated**: 2025-12-11  
**Prerequisites**: Go 1.23+, Node.js 18+, GCP account with billing enabled

## Overview

This guide gets you from zero to running the debate generator locally in under 15 minutes. The application consists of:

- **Frontend**: React 18+ SPA (port 3000)
- **Backend**: 3 Go Cloud Functions (local emulation via Functions Framework)
- **External API**: Anthropic Claude API (requires API key)

---

## Initial Setup

### 1. Clone and Navigate

```bash
cd /Users/raphink/src/github.com/raphink/debate
git checkout -b 001-debate-generator
```

### 2. Install Dependencies

**Frontend**:
```bash
cd frontend
npm install
```

**Backend** (each function independently):
```bash
# Topic validation function
cd backend/functions/validate-topic
go mod download

# Panelist suggestion function  
cd ../suggest-panelists
go mod download

# Debate generation function
cd ../generate-debate
go mod download
```

### 3. Configure API Keys

Create `.env` file in project root:

```bash
# .env
ANTHROPIC_API_KEY=sk-ant-api03-...  # Get from https://console.anthropic.com
GCP_PROJECT_ID=your-project-id       # For local emulation metadata
```

**Security Note**: Never commit `.env` to version control. Add to `.gitignore` immediately.

---

## Running Locally

### Backend (Cloud Functions Emulation)

Use Google Cloud Functions Framework to run functions locally:

**Terminal 1 - Validate Topic Function**:
```bash
cd backend/functions/validate-topic
export ANTHROPIC_API_KEY=$(grep ANTHROPIC_API_KEY ../../.env | cut -d '=' -f2)
go run cmd/main.go
# Runs on http://localhost:8080
```

**Terminal 2 - Suggest Panelists Function**:
```bash
cd backend/functions/suggest-panelists
export ANTHROPIC_API_KEY=$(grep ANTHROPIC_API_KEY ../../.env | cut -d '=' -f2)
PORT=8081 go run cmd/main.go
# Runs on http://localhost:8081
```

**Terminal 3 - Generate Debate Function**:
```bash
cd backend/functions/generate-debate
export ANTHROPIC_API_KEY=$(grep ANTHROPIC_API_KEY ../../.env | cut -d '=' -f2)
PORT=8082 go run cmd/main.go
# Runs on http://localhost:8082
```

### Frontend (React Dev Server)

**Terminal 4 - React Application**:
```bash
cd frontend
npm start
# Runs on http://localhost:3000
```

**Environment Configuration**:  
Edit `frontend/.env.development`:
```bash
REACT_APP_VALIDATE_TOPIC_URL=http://localhost:8080
REACT_APP_SUGGEST_PANELISTS_URL=http://localhost:8081
REACT_APP_GENERATE_DEBATE_URL=http://localhost:8082
```

---

## Testing the Flow

### 1. Topic Validation (User Story 1)

**cURL Test**:
```bash
curl -X POST http://localhost:8080 \
  -H "Content-Type: application/json" \
  -d '{"topic": "Should Christians defy authorities when the law is unfair?"}'
```

**Expected Response**:
```json
{
  "isRelevant": true,
  "message": "This topic is suitable for theological/philosophical debate.",
  "suggestedRefinement": ""
}
```

**Browser Test**:
1. Navigate to http://localhost:3000
2. Enter topic: "Should Christians defy authorities when the law is unfair?"
3. Click "Validate Topic"
4. Verify validation result appears within 3 seconds

### 2. Panelist Suggestions (User Story 2)

**cURL Test**:
```bash
curl -X POST http://localhost:8081 \
  -H "Content-Type: application/json" \
  -d '{"topic": "Should Christians defy authorities when the law is unfair?"}'
```

**Expected Response**:
```json
{
  "panelists": [
    {
      "id": "Augustine354",
      "name": "Augustine of Hippo",
      "avatarUrl": "/avatars/augustine.jpg",
      "tagline": "Bishop of Hippo, Doctor of the Church",
      "bio": "4th-century theologian who wrote extensively on the relationship between earthly and divine authority in 'The City of God'."
    },
    {
      "id": "MLKJr",
      "name": "Martin Luther King Jr.",
      "avatarUrl": "/avatars/mlk.jpg",
      "tagline": "Civil Rights Leader, Baptist Minister",
      "bio": "20th-century activist who advocated for civil disobedience against unjust laws in 'Letter from Birmingham Jail'."
    }
    // ... up to 20 panelists
  ]
}
```

**Browser Test**:
1. After topic validation, view suggested panelists
2. Verify each panelist shows avatar, name, handle (id), tagline, bio
3. Select 2-5 panelists by clicking cards
4. Verify selection limit (max 5) is enforced

### 3. Debate Generation (User Story 3)

**cURL Test (SSE Stream)**:
```bash
curl -X POST http://localhost:8082 \
  -H "Content-Type: application/json" \
  -d '{
    "topic": "Should Christians defy authorities when the law is unfair?",
    "panelists": ["Augustine354", "MLKJr"]
  }' \
  --no-buffer
```

**Expected SSE Stream**:
```
data: {"type":"message","panelistId":"Augustine354","content":"I must begin by acknowledging...","done":false}

data: {"type":"message","panelistId":"MLKJr","content":"Brother Augustine raises...","done":false}

data: {"type":"done"}
```

**Browser Test**:
1. After selecting panelists, click "Generate Debate"
2. Verify first response appears within 5 seconds
3. Watch debate stream progressively in chat bubbles
4. Verify each panelist has correct avatar and name
5. Verify loading indicators between turns

### 4. PDF Export (User Story 4)

**Browser Test**:
1. After debate completes, click "Export as PDF"
2. Verify PDF download begins within 2 seconds
3. Open PDF and verify:
   - Debate topic in header
   - Panelist profiles included
   - Complete conversation with avatars
   - Timestamp and formatting

---

## Development Workflow

### Code Organization

**Frontend Components** (follow Single Responsibility Principle):
- `TopicInput/` - Topic entry form (US1)
- `ValidationResult/` - Validation feedback (US1)
- `PanelistGrid/` - Panelist display grid (US2)
- `PanelistCard/` - Individual panelist card (US2)
- `PanelistSelector/` - Selection management (US2)
- `DebateView/` - Chat-style debate display (US3)
- `ChatBubble/` - Individual debate message (US3)
- `PDFExport/` - PDF generation button (US4)

**Backend Functions** (stateless, focused):
- `validate-topic/` - Claude API topic validation
- `suggest-panelists/` - Claude API panelist suggestions
- `generate-debate/` - Claude API streaming debate generation

### Running Tests

**Frontend Tests**:
```bash
cd frontend
npm test                    # Run Jest tests
npm run test:a11y          # Run axe-core accessibility tests
npm run lint               # ESLint
```

**Backend Tests**:
```bash
cd backend/functions/validate-topic
go test ./...              # Unit tests
go test -race ./...        # Race condition detection
golangci-lint run          # Linting
```

### Adding a New Component

1. Create component directory: `frontend/src/components/MyComponent/`
2. Add component file: `MyComponent.jsx`
3. Add test file: `MyComponent.test.jsx`
4. Add styles: `MyComponent.module.css`
5. Export from `index.js`

Example component structure:
```javascript
// MyComponent.jsx
import React from 'react';
import styles from './MyComponent.module.css';

export const MyComponent = ({ prop1, prop2 }) => {
  return (
    <div className={styles.container}>
      {/* Component JSX */}
    </div>
  );
};
```

### Adding a New Cloud Function

1. Create function directory: `backend/functions/my-function/`
2. Initialize Go module: `go mod init github.com/raphink/debate/functions/my-function`
3. Create `main.go` with HTTP handler
4. Add to local test script with unique port

---

## Debugging Tips

### Backend Issues

**Function won't start**:
```bash
# Check Go version (must be 1.23+)
go version

# Verify environment variable is set
echo $ANTHROPIC_API_KEY
```

**API key errors**:
- Verify key starts with `sk-ant-api03-`
- Check Anthropic console for key status
- Ensure key has sufficient credits

**CORS errors**:
- Verify Cloud Functions Framework includes CORS headers
- Check `Access-Control-Allow-Origin: *` in response

### Frontend Issues

**React won't start**:
```bash
# Check Node version (must be 18+)
node --version

# Clear cache and reinstall
rm -rf node_modules package-lock.json
npm install
```

**API calls failing**:
- Verify backend functions are running (check terminals 1-3)
- Check browser Network tab for exact error
- Verify `.env.development` URLs match running ports

**SSE stream not working**:
- Check browser console for EventSource errors
- Verify backend sends `Content-Type: text/event-stream`
- Test with cURL to isolate frontend vs backend issue

---

## Common Tasks

### Update Dependencies

**Frontend**:
```bash
cd frontend
npm update
npm audit fix  # Security patches
```

**Backend**:
```bash
cd backend/functions/validate-topic
go get -u ./...
go mod tidy
```

### Format Code

**Frontend**:
```bash
cd frontend
npm run format  # Prettier
```

**Backend**:
```bash
cd backend
go fmt ./...
```

### Check Accessibility

```bash
cd frontend
npm run test:a11y
# Runs axe-core against all components
# Verifies WCAG 2.1 Level AA compliance
```

---

## Deployment Preview (Not for Local Dev)

When ready to deploy to GCP:

```bash
# Deploy Cloud Functions (from backend/)
gcloud functions deploy validate-topic \
  --gen2 \
  --runtime=go123 \
  --trigger-http \
  --allow-unauthenticated \
  --set-secrets=ANTHROPIC_API_KEY=anthropic-api-key:latest

# Deploy frontend (from frontend/)
npm run build
gcloud app deploy  # Or use Cloud Storage + CDN
```

**Note**: Full deployment instructions will be in `tasks.md` after running `/speckit.tasks`.

---

## Next Steps

1. **Read the spec**: Review `specs/001-debate-generator/spec.md` for full requirements
2. **Check data model**: Review `specs/001-debate-generator/data-model.md` for entity definitions
3. **Review contracts**: Check `specs/001-debate-generator/contracts/` for API schemas
4. **Run tasks**: After `/speckit.tasks` generates task list, follow dependency order

---

## Getting Help

- **Constitution**: See `.specify/memory/constitution.md` for coding principles
- **Spec Questions**: Check `specs/001-debate-generator/spec.md` acceptance scenarios
- **API Reference**: Review OpenAPI contracts in `contracts/` directory
- **Research Decisions**: See `specs/001-debate-generator/research.md` for technical rationale
