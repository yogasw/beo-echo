# Beo Echo - Go Backend

This is the Go implementation of the Mockoon Control Panel backend, refactored from the original Node.js/TypeScript implementation.

## Features

- Complete API implementation using Gin Gonic framework
- Command-line interface using Cobra library
- Database models using GORM (replacing Prisma)
  - Support for both SQLite (default) and PostgreSQL
- File and mock instance repositories
- Traefik configuration generators

## Commands

- `go run main.go` - Run the main server (default)
- `go run main.go server` - Run the full server with all services
- `go run main.go api` - Run only the API server without additional services
- `go run main.go generate` - Generate configuration files

Alternatively, you can use the run script:
- `./run.sh` - Run the main server

## Database Configuration

The service supports two database backends:

1. **SQLite** (default): Used when no `DATABASE_URL` environment variable is set. Data is stored in a SQLite file in the configured data directory.

2. **PostgreSQL**: Used when the `DATABASE_URL` environment variable is set. The URL should follow the PostgreSQL connection string format:
   ```
   DATABASE_URL=postgresql://username:password@localhost:5432/dbname
   ```
- `./run.sh api` - Run only the API server
- `./run.sh generate` - Generate configuration files

## Default URLs

By default, the server runs on port 3600:

- Server URL: http://localhost:3600
- API Base URL: http://localhost:3600/mock/api
- Health Check: http://localhost:3600/mock/api/health

To specify a different port, use the `-p` flag:
```
go run main.go -p 8080
```
or
```
./run.sh -p 8080
```

## Development

```bash
# Build the application
make build

# Run the application
make run

# Development with hot-reload
make dev

# Clean build files
make clean
```

## Directory Structure

```
.
├── cmd/             # Command-line interface commands
├── src/             # Source code
│   ├── git-sync/    # Git synchronization functionality
│   ├── health/      # Health check endpoints
│   ├── lib/         # Shared libraries and constants
│   ├── middlewares/ # HTTP middleware components
│   ├── mocks/       # Mock handling functionality
│   ├── prisma/      # Database models and connection
│   ├── scripts/     # Utility scripts
│   ├── server/      # Server implementation
│   ├── traefik/     # Traefik configuration generators
│   ├── types/       # Type definitions
│   └── utils/       # Utility functions
├── logs/            # Log files
└── uploads/         # Uploaded files
```
