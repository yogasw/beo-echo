# Beo Echo - Go Backend

This is the Go implementation of the Beo Echo API mocking service backend, providing a powerful and flexible way to create, manage, and serve mock APIs similar to tools like Beeceptor and Mockoon.

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
- `./run.sh api` - Run only the API server
- `./run.sh generate` - Generate configuration files

## Database Configuration

The service supports two database backends:

1. **SQLite** (default): Used when no `DATABASE_URL` environment variable is set. Data is stored in a SQLite file in the configured data directory.

2. **PostgreSQL**: Used when the `DATABASE_URL` environment variable is set. The URL should follow the PostgreSQL connection string format:
   ```
   DATABASE_URL=postgresql://username:password@localhost:5432/dbname
   ```

## Default URLs

By default, the server runs on port 3600:

- Server URL: http://localhost:3600
- API Base URL: http://localhost:3600/api/api
- Health Check: http://localhost:3600/api/api/health
- Mock APIs: http://localhost:3600/{project-name}/{endpoint-path}

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
│   ├── auth/         # Authentication functionality
│   ├── caddy/        # Caddy configuration generators
│   ├── database/     # Database models and connection (GORM)
│   ├── health/       # Health check endpoints
│   ├── lib/          # Shared libraries and constants
│   ├── middlewares/  # HTTP middleware components
│   ├── mocks/        # Mock API handling functionality
│   ├── system-config/ # System configuration management
│   ├── types/        # Type definitions
│   └── utils/        # Utility functions
├── logs/            # Log files
└── uploads/         # Uploaded files
```

## Authentication & Authorization

Beo Echo implements a complete authentication system with JWT tokens:

- **Registration**: Create new user accounts
- **Login**: Authenticate and receive a JWT token
- **Workspace Management**: Users can create and join workspaces
- **Role-Based Access Control**: Different permission levels within workspaces
- **Project Access Control**: Projects belong to workspaces with controlled access

To authenticate API requests:
1. Login via `/api/auth/login` to get a JWT token
2. Include the token in the `Authorization` header: `Bearer {your-token}`

## Multi-User Collaboration

The system supports multi-user collaboration through workspaces:

- **Workspaces**: Collection of related projects
- **Projects**: Collection of endpoints, responses, and proxies 
- **User Roles**: Owner, Admin, Member, or Viewer
- **Shared Access**: Team members can access and modify projects according to their permissions
