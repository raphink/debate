# AI-Powered Theology/Philosophy Debate Generator

Generate engaging debates between historical theological and philosophical figures on topics of your choice.

## Quick Start

### Prerequisites

- **Docker & Docker Compose** - [Download](https://www.docker.com/get-started) (recommended)
- **OR** Manual setup:
  - **Go 1.23+** - [Download](https://golang.org/dl/)
  - **Node.js 18+** - [Download](https://nodejs.org/)
- **Anthropic API Key** - [Get one](https://console.anthropic.com)

### Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/raphink/debate.git
   cd debate
   git checkout 001-debate-generator
   ```

2. **Configure environment variables**
   ```bash
   cp .env.example .env
   # Edit .env and add your ANTHROPIC_API_KEY
   ```

### Running with Docker (Recommended)

The easiest way to run the entire application locally:

```bash
# Start all services (backend functions + frontend)
docker-compose up --build

# Or run in detached mode
docker-compose up -d --build

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

The application will be available at:
- **Frontend**: http://localhost:3000
- **Backend Functions**:
  - Topic Validation: http://localhost:8080
  - Panelist Suggestions: http://localhost:8081
  - Debate Generation: http://localhost:8082

### Running Manually (Without Docker)

**Install dependencies first:**

**Install dependencies first:**

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
export ANTHROPIC_API_KEY=$(grep ANTHROPIC_API_KEY ../../.env | cut -d '=' -f2)
go run main.go

# Terminal 2 - Panelist Suggestions (port 8081)
cd backend/functions/suggest-panelists
export ANTHROPIC_API_KEY=$(grep ANTHROPIC_API_KEY ../../.env | cut -d '=' -f2)
PORT=8081 go run main.go

# Terminal 3 - Debate Generation (port 8082)
cd backend/functions/generate-debate
export ANTHROPIC_API_KEY=$(grep ANTHROPIC_API_KEY ../../.env | cut -d '=' -f2)
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
