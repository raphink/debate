# AI-Powered Theology/Philosophy Debate Generator

Generate engaging debates between historical theological and philosophical figures on topics of your choice.

## Features

✨ **AI-Powered Debates** - Generate dynamic conversations between historical figures using Claude AI  
🎭 **Real Portraits** - Automatic fetching of portrait images from Wikimedia Commons  
💬 **Live Streaming** - Watch debates unfold in real-time with progressive message streaming  
📱 **Mobile PWA** - Install as a standalone app on iOS and Android devices  
📄 **PDF Export** - Export debates with portraits and chat bubble formatting  
✍️ **Markdown Formatting** - Supports inline formatting (*italic*, **bold**, ***bold italic***) in messages  
🔍 **Topic Autocomplete** - Discover and view previous debates as you type with smart suggestions  
🎨 **Beautiful UI** - Modern, responsive design with gradient effects and animations  
🔒 **Secure** - CORS protection, input sanitization, and rate limiting

## Quick Start

### Prerequisites

- **Docker & Docker Compose** - [Download](https://www.docker.com/get-started) (recommended)
- **Summon** - [Install](https://github.com/cyberark/summon#install) (secret manager for Docker workflow)
- **Summon GCP Plugin** - Must be installed at `/usr/local/lib/summon/gcloud`
- **GCP Secret Manager Access** - Required secrets: `anthropic-api-key`, `gcp-project-id`
- **OR** Manual setup:
  - **Go 1.24+** - [Download](https://golang.org/dl/)
  - **Node.js 18+** - [Download](https://nodejs.org/)
  - **gcloud CLI** - [Install](https://cloud.google.com/sdk/docs/install)
- **Anthropic API Key** - [Get one](https://console.anthropic.com) and store in GCP Secret Manager

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/raphink/debate.git
   ```

2. **Configure secrets in GCP Secret Manager**
   
   Store your secrets in GCP Secret Manager (one-time setup):
   ```bash
   # Set your Anthropic API key
   echo -n "sk-ant-api03-YOUR_KEY_HERE" | gcloud secrets create anthropic-api-key --data-file=-
   
   # Set your GCP project ID
   echo -n "your-project-id" | gcloud secrets create gcp-project-id --data-file=-
   ```
   
   The `secrets.yml` file in the project root maps these secrets to environment variables for local development.

### Running with Docker (Recommended)

The easiest way to run the entire application locally using secrets from GCP Secret Manager:

```bash
# Verify summon and GCP plugin are installed
which summon
ls -l /usr/local/lib/summon/gcloud

# Start all services (backend functions + frontend) with summon
summon -p gcloud docker-compose up --build

# Or run in detached mode
summon -p gcloud docker-compose up -d --build

# Or use the convenience script (validates all prerequisites)
./start-local.sh

# View logs
summon -p gcloud docker-compose logs -f

# Stop all services
summon -p gcloud docker-compose down
```

**Why Summon?** Secrets are stored in GCP Secret Manager and injected securely into Docker containers. This eliminates `.env` files and ensures your local development uses the same secret source as production, preventing configuration drift.

The application will be available at:
- **Frontend**: http://localhost:3000
- **Backend Functions**:
  - Topic Validation: http://localhost:8080
  - Panelist Suggestions: http://localhost:8081
  - Debate Generation: http://localhost:8082

### Running Manually (Without Docker)

**Export secrets from GCP Secret Manager:**

```bash
# Option 1: Use summon to inject secrets into your shell
summon -p gcloud env | grep -E "ANTHROPIC_API_KEY|GCP_PROJECT_ID" > .env.local
source .env.local

# Option 2: Export directly using gcloud CLI
export ANTHROPIC_API_KEY=$(gcloud secrets versions access latest --secret="anthropic-api-key")
export GCP_PROJECT_ID=$(gcloud secrets versions access latest --secret="gcp-project-id")
```

**Install dependencies:**

```bash
# Frontend dependencies
cd frontend
npm install

# Backend dependencies (for each Cloud Function)
cd ../backend/functions/validate-topic && go mod download
cd ../suggest-panelists && go mod download
cd ../generate-debate && go mod download
```

**Start Backend Functions** (3 separate terminals):

```bash
# Terminal 1 - Topic Validation (port 8080)
cd backend/functions/validate-topic
go run main.go

# Terminal 2 - Panelist Suggestions (port 8081)
cd backend/functions/suggest-panelists
PORT=8081 go run main.go

# Terminal 3 - Debate Generation (port 8082)
cd backend/functions/generate-debate
PORT=8082 go run main.go
```

**Start Frontend** (new terminal):

```bash
cd frontend
npm start
# Opens http://localhost:3000
```

### Testing the Application

1. **Topic Validation**: Enter a theological/philosophical topic like "Should Christians defy unjust laws?"
2. **Panelist Selection**: Browse suggested historical figures and select 2-5 panelists
3. **Debate Generation**: Watch the debate stream in real-time with avatars
4. **PDF Export**: Download the completed debate as a formatted PDF

## Project Structure

```
debate/
├── backend/
│   ├── functions/
│   │   ├── validate-topic/      # Topic validation Cloud Function
│   │   ├── suggest-panelists/   # Panelist suggestion Cloud Function
│   │   └── generate-debate/     # Debate generation Cloud Function
│   └── shared/                  # Shared utilities
│       ├── sanitize/            # XSS prevention
│       ├── errors/              # Error handling
│       ├── ratelimit/           # Rate limiting
│       └── auth/                # API key management
├── frontend/
│   ├── public/
│   │   └── avatars/             # Historical figure avatars
│   └── src/
│       ├── components/          # React components
│       ├── pages/               # Page components
│       ├── services/            # API services
│       ├── hooks/               # Custom React hooks
│       └── utils/               # Utilities
└── specs/
    └── 001-debate-generator/    # Feature specification
```

## Development

### PWA Icon Generation

If you need to regenerate the PWA icons from the SVG source:

```bash
cd frontend
./generate-icons.sh
```

Requires either `rsvg-convert` (librsvg) or ImageMagick:
```bash
# macOS
brew install librsvg
# or
brew install imagemagick
```

### Mobile Installation

The app can be installed as a PWA on mobile devices:

**iOS (Safari)**:
1. Open the app in Safari
2. Tap the Share button
3. Select "Add to Home Screen"
4. Confirm installation

**Android (Chrome)**:
1. Open the app in Chrome
2. Tap the menu (three dots)
3. Select "Install app" or "Add to Home Screen"
4. Confirm installation

The installed app runs in standalone mode without browser chrome for a native app experience.

### Running Tests

**Frontend**:
```bash
cd frontend
npm test                # Unit tests
npm run test:a11y      # Accessibility tests
npm run lint           # ESLint
```

**Backend**:
```bash
cd backend/functions/validate-topic
go test ./...          # Unit tests
golangci-lint run      # Linting
```

### Code Quality

- **Linting**: ESLint (frontend), golangci-lint (backend)
- **Formatting**: Prettier (frontend), gofmt (backend)
- **Testing**: Jest + React Testing Library (frontend), Go testing (backend)
- **Accessibility**: axe-core automated testing, WCAG 2.1 Level AA compliance

## Documentation

- **Specification**: [specs/001-debate-generator/spec.md](specs/001-debate-generator/spec.md)
- **Technical Plan**: [specs/001-debate-generator/plan.md](specs/001-debate-generator/plan.md)
- **Data Model**: [specs/001-debate-generator/data-model.md](specs/001-debate-generator/data-model.md)
- **API Contracts**: [specs/001-debate-generator/contracts/](specs/001-debate-generator/contracts/)
- **Developer Guide**: [specs/001-debate-generator/quickstart.md](specs/001-debate-generator/quickstart.md)
- **Tasks**: [specs/001-debate-generator/tasks.md](specs/001-debate-generator/tasks.md)

## Contributing

See [.specify/memory/constitution.md](.specify/memory/constitution.md) for coding principles and quality standards.

## License

Copyright © 2025 Raphael Pinson
