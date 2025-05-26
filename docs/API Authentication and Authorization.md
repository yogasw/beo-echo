# Beo Echo API Authentication and Authorization

This document provides the technical details for the authentication and authorization system in Beo Echo.

## Authentication Flow

### JWT-Based Authentication

Beo Echo uses JSON Web Tokens (JWT) for authentication. The authentication flow is as follows:

1. **User Registration**
   - User provides email, password, and name
   - System creates a new user record with hashed password
   - System creates a default workspace for the user
   - System generates a JWT token and returns it

2. **User Login**
   - User provides email and password
   - System verifies credentials
   - System generates a JWT token and returns it

3. **Protected Resource Access**
   - User includes JWT token in the `Authorization` header
   - System validates the token
   - If valid, the request proceeds; otherwise, returns 401 Unauthorized

### JWT Token Structure

```json
{
  "user_id": "uuid-string",
  "email": "user@example.com",
  "name": "User Name",
  "is_owner": false,
  "exp": 1652375421,
  "iat": 1652289021,
  "nbf": 1652289021
}
```

## Multi-Workspace Architecture

### Data Model

```
User
├── Workspaces (via UserWorkspace)
│   ├── Projects
│   │   ├── Endpoints
│   │   │   └── Responses
│   │   └── ProxyTargets
│   └── Members (other users)
└── Identities (SSO)
```

### Workspace Roles

- **Admin**: Can manage all workspace resources and members
- **Member**: Can use all workspace resources but can't manage members

### System-wide Roles

- **Owner**: Special role with system-wide administrative privileges

## Access Control

### Project Access

- Projects belong to workspaces
- Users can access projects through workspace membership
- System owners can access all projects

### Workspace Access

- Users can be members of multiple workspaces
- Each user-workspace relationship has a role (admin/member)
- Only workspace admins can manage workspace membership

## API Routes Structure

### Authentication Routes

```
POST /api/auth/login
POST /api/auth/register
```

### Workspace Management Routes

```
GET  /api/workspaces
POST /api/workspaces
GET  /api/workspaces/:workspaceID/role
```

### Project Management Routes

```
GET  /api/workspaces/:workspaceID/projects
POST /api/workspaces/:workspaceID/projects
GET  /api/workspaces/:workspaceID/projects/:projectId
PUT  /api/workspaces/:workspaceID/projects/:projectId
DEL  /api/workspaces/:workspaceID/projects/:projectId
```

### Endpoint Management Routes

```
GET  /api/workspaces/:workspaceID/projects/:projectId/endpoints
POST /api/workspaces/:workspaceID/projects/:projectId/endpoints
GET  /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id
PUT  /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id
DEL  /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id
```

### Response Management Routes

```
GET  /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses
POST /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses
GET  /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId
PUT  /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId
DEL  /api/workspaces/:workspaceID/projects/:projectId/endpoints/:id/responses/:responseId
```

### Logs Management Routes

```
GET  /api/workspaces/:workspaceID/projects/:projectId/logs
GET  /api/workspaces/:workspaceID/projects/:projectId/logs/stream
```

## Middleware

### JWTAuthMiddleware

Validates JWT tokens and sets user context for downstream handlers:

```go
c.Set("userID", claims.UserID)
c.Set("isOwner", claims.IsOwner)
```

### WorkspaceProjectAccessMiddleware

Ensures the requested project belongs to the specified workspace and the user has access:

1. Verifies workspace access
2. Confirms project belongs to workspace
3. Sets workspace and project context

## Initialization

On first startup, the system creates:
- Default admin user (email: admin@admin.com, password: admin)
- Demo workspace
- Demo project

## Security Considerations

- Passwords are hashed using bcrypt
- JWT tokens expire after 24 hours
- All project access is isolated by workspace
- System validates proper resource ownership before operations
