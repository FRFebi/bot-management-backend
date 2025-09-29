# Bot Management Backend

Golang backend service using Fiber framework for the Bot Management Dashboard.

## Project Structure

```
backend/
├── cmd/
│   └── server/          # Application entry points
├── internal/            # Private application code
│   ├── config/          # Configuration management
│   ├── handlers/        # HTTP handlers
│   ├── middleware/      # Custom middleware
│   ├── models/          # Data models
│   └── services/        # Business logic
├── pkg/                 # Public packages
│   ├── logger/          # Logging utilities
│   └── utils/           # Utility functions
├── docker-compose.dev.yml
├── Dockerfile
├── Makefile
└── .env.example
```

## Quick Start

### Prerequisites
- Go 1.21+
- Docker & Docker Compose (optional)

### Development Setup

1. **Clone and install dependencies:**
```bash
make deps
```

2. **Copy environment variables:**
```bash
cp .env.example .env
```

3. **Run the application:**
```bash
make run
```

The server will start on `http://localhost:8080`

### Docker Development

Run with full development stack (PostgreSQL + Redis):
```bash
make docker-dev
```

## Available Commands

```bash
make help          # Show all available commands
make build         # Build the application
make run           # Run the application
make test          # Run tests
make docker-dev    # Run development environment
make lint          # Run linter
make fmt           # Format code
```

## API Endpoints

### Health Check
- `GET /health` - Returns service status

### API v1
- `GET /api/v1/` - API information

## Configuration

Environment variables (see `.env.example`):
- `PORT` - Server port (default: 8080)
- `ENVIRONMENT` - Environment mode
- `DB_*` - Database configuration
- `JWT_*` - JWT configuration
- `REDIS_*` - Redis configuration

## Status

✅ **Backend Issue #1 Completed:**
- [x] Initialize Go module
- [x] Setup Fiber framework
- [x] Create internal and pkg directory structure
- [x] Setup development environment with Docker
- [x] Configure environment variables
- [x] Setup logging framework
- [x] Create basic health check endpoint

## Next Steps

Ready to proceed with:
- Backend Issue #2: Database Setup and Schema Implementation
- Backend Issue #3: Authentication and Authorization System