# Multi-User and Workspace Support in Beo Echo

This document outlines the multi-user authentication and workspace management features implemented in the Beo Echo API mocking service.

## Overview

Beo Echo now supports:
- User authentication with JWT
- Multi-workspace architecture
- Role-based access control
- Project isolation by workspace

## Authentication Flow

### 1. User Registration

Register a new user account with a personal workspace.

**Request:**
```bash
curl -X POST "http://localhost:3600/mock/api/auth/register" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "email": "john@example.com",
    "password": "secure_password"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Registration successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-id-uuid",
    "email": "john@example.com",
    "name": "John Doe",
    "is_owner": false
  }
}
```

### 2. User Login

Login with email and password to obtain a JWT token.

**Request:**
```bash
curl -X POST "http://localhost:3600/mock/api/auth/login" \
  -H "Content-Type: application/json" \
  -d '{
    "email": "admin@admin.com",
    "password": "admin"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": "user-id-uuid",
    "email": "admin@admin.com",
    "name": "Admin",
    "is_owner": true
  }
}
```

## Workspace Management

### 1. List User's Workspaces

Retrieve all workspaces accessible to the authenticated user.

**Request:**
```bash
curl -X GET "http://localhost:3600/mock/api/workspaces" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "workspace-uuid-1",
      "name": "Demo Workspace",
      "created_at": "2025-05-12T10:15:30Z",
      "updated_at": "2025-05-12T10:15:30Z"
    },
    {
      "id": "workspace-uuid-2",
      "name": "Personal Workspace",
      "created_at": "2025-05-12T11:20:45Z",
      "updated_at": "2025-05-12T11:20:45Z"
    }
  ]
}
```

### 2. Create a New Workspace

Create a new workspace for the authenticated user.

**Request:**
```bash
curl -X POST "http://localhost:3600/mock/api/workspaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "New Project Team"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Workspace created successfully",
  "data": {
    "id": "workspace-uuid-3",
    "name": "New Project Team",
    "created_at": "2025-05-12T14:30:20Z",
    "updated_at": "2025-05-12T14:30:20Z"
  }
}
```

### 3. List Projects in a Workspace

Retrieve all projects within a specific workspace.

**Request:**
```bash
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "project-uuid-1",
      "name": "Demo Project",
      "workspace_id": "workspace-uuid-1",
      "mode": "mock",
      "status": "running",
      "alias": "demo",
      "url": "http://demo.localhost:3600",
      "created_at": "2025-05-12T10:15:30Z",
      "updated_at": "2025-05-12T10:15:30Z"
    }
  ]
}
```

### 4. Check User Role in a Workspace

Check the role of a user in a specific workspace.

**Request:**
```bash
# Check your own role
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/role" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Admins can check other users' roles
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/role?user_id=other-user-uuid" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": {
    "user_id": "user-uuid",
    "workspace_id": "workspace-uuid-1",
    "role": "admin"
  }
}
```

## Projects API (with Workspace Restrictions)

The existing projects API has been updated to respect workspace permissions:

### List Projects (Filtered by Workspace Access)

**Request:**
```bash
curl -X GET "http://localhost:3600/mock/api/projects" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Response:**
```json
{
  "success": true,
  "data": [
    {
      "id": "project-uuid-1",
      "name": "Demo Project",
      "workspace_id": "workspace-uuid-1",
      "mode": "mock",
      "status": "running",
      "alias": "demo",
      "url": "http://demo.localhost:3600",
      "created_at": "2025-05-12T10:15:30Z",
      "updated_at": "2025-05-12T10:15:30Z"
    }
  ]
}
```

### Create Project In a Specific Workspace

**Request:**
```bash
curl -X POST "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "New API Mock",
    "mode": "mock",
    "alias": "new-api",
    "documentation": "My new API mock"
  }'
```

**Response:**
```json
{
  "success": true,
  "message": "Project created successfully",
  "data": {
    "id": "project-uuid-2",
    "name": "New API Mock",
    "workspace_id": "workspace-uuid-1",
    "mode": "mock",
    "status": "running",
    "alias": "new-api",
    "url": "http://new-api.localhost:3600",
    "documentation": "My new API mock",
    "created_at": "2025-05-12T15:30:20Z",
    "updated_at": "2025-05-12T15:30:20Z"
  }
}
```

### Access Project Resources (Hierarchical Routes)

All project resources are now accessed through a hierarchical structure that enforces workspace-project relationships:

```bash
# Get project details using the hierarchical structure
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects/project-uuid-1" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Manage endpoints
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects/project-uuid-1/endpoints" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Create a new endpoint in a workspace project
curl -X POST "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects/project-uuid-1/endpoints" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "method": "GET",
    "path": "/api/users",
    "enabled": true,
    "response_mode": "static",
    "documentation": "Returns list of users"
  }'

# View project logs
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects/project-uuid-1/logs" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
  
# Stream real-time logs from a project
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects/project-uuid-1/logs/stream" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
  
# Manage endpoint responses
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects/project-uuid-1/endpoints/endpoint-uuid-1/responses" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
  
# Create a response for an endpoint
curl -X POST "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/projects/project-uuid-1/endpoints/endpoint-uuid-1/responses" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "status_code": 200,
    "body": "{\"users\": [{\"id\": 1, \"name\": \"John\"}]}",
    "headers": "{\"Content-Type\": \"application/json\"}",
    "delay_ms": 0,
    "priority": 1,
    "enabled": true
  }'
```

Each request is automatically validated to ensure:
1. The user has access to the specified workspace
2. The project belongs to the specified workspace
3. The user has appropriate permissions

## Workspace Management

### List User Workspaces

Retrieve all workspaces accessible to the authenticated user:

```bash
curl -X GET "http://localhost:3600/mock/api/workspaces" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

### Create a New Workspace

Create a new workspace and automatically assign the authenticated user as an admin:

```bash
curl -X POST "http://localhost:3600/mock/api/workspaces" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "My Team Workspace"
  }'
```

### Check User Role in Workspace

Check the authenticated user's role in a specific workspace:

```bash
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/role" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

Admin users can also check other users' roles:

```bash
curl -X GET "http://localhost:3600/mock/api/workspaces/workspace-uuid-1/role?user_id=other-user-uuid" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

## Legacy Route Support

For backward compatibility, the following routes are still supported:

```bash
# List all projects (filtered by workspace access)
curl -X GET "http://localhost:3600/mock/api/projects" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"

# Create a project (without specifying workspace)
curl -X POST "http://localhost:3600/mock/api/projects" \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -d '{
    "name": "Legacy Project Creation",
    "workspace_id": "workspace-uuid-1",
    "mode": "mock",
    "alias": "legacy",
    "documentation": "Created with legacy API"
  }'
```

However, it's recommended to use the new hierarchical routes for better enforcement of workspace-project relationships.

## Access Control Rules

1. **System Owners (`is_owner: true`):**
   - Can access all workspaces and projects
   - Can manage system-wide settings
   - Default admin user is a system owner

2. **Workspace Admins (`role: "admin"`):**
   - Can manage projects within their workspace
   - Can check user roles within their workspace
   - User who creates a workspace is automatically its admin

3. **Workspace Members (`role: "member"`):**
   - Can view projects within their workspace
   - Limited management capabilities (implementation pending)

## Default Setup

On first startup, the system creates:
- Default admin user: `admin@admin.com` with password `admin`
- Demo workspace with the admin user as workspace admin
- Demo project within the demo workspace

## JWT Token Structure

**Payload:**
```json
{
  "user_id": "user-uuid",
  "email": "user@example.com",
  "name": "User Name",
  "is_owner": false,
  "exp": 1715443200,
  "iat": 1715356800,
  "nbf": 1715356800
}
```

## Future Enhancements

Future implementations may include:
1. User management API (invite users to workspaces)
2. Role management within workspaces
3. SSO integration with popular providers
4. API key management per workspace
5. Usage statistics and limits per workspace

## Security Considerations

1. JWT tokens expire after 24 hours
2. Passwords are hashed with bcrypt
3. Authorization is checked on all protected routes
4. System owner flag is protected and not exposed for modification
