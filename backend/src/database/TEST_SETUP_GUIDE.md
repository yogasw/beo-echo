# Test Setup Pattern Guide

## Overview

This guide explains the improved test setup pattern for creating workspace and project test data with proper initialization and cleanup using defer.

## New Initialization Functions

We have added two new helper functions in `backend/src/database/test_helpers.go`:

### 1. `InitTestWorkspaceWithProject` - Complete Setup
Creates a complete test environment with user, workspace, and project:

```go
setup, err := database.InitTestWorkspaceWithProject(
    "user@example.com",     // User email
    "Test User",            // User name  
    "Test Workspace",       // Workspace name
    "Test Project",         // Project name
    "test-project-alias",   // Project alias
)
if err != nil {
    t.Fatalf("Failed to initialize test workspace: %v", err)
}
defer setup.Cleanup() // Always defer cleanup immediately
```

### 2. `InitTestWorkspaceOnly` - Workspace Only Setup
Creates test environment with user and workspace (no project):

```go
setup, err := database.InitTestWorkspaceOnly(
    "user@example.com",     // User email
    "Test User",            // User name
    "Test Workspace",       // Workspace name
)
if err != nil {
    t.Fatalf("Failed to initialize test workspace: %v", err)
}
defer setup.Cleanup() // Always defer cleanup immediately
```

## Benefits of the New Pattern

### 1. **Automatic Cleanup**
- Single `defer setup.Cleanup()` call handles all cleanup
- No risk of forgetting to clean up test data
- Cleanup happens even if test fails or panics

### 2. **Reduced Boilerplate**
- No need to manually create user, workspace, and project
- No need to manage multiple defer statements
- Simplified test initialization

### 3. **Error Handling**
- If project creation fails, workspace is automatically cleaned up
- Consistent error handling across all test setup

### 4. **Database Initialization**
- Database initialization is handled automatically
- No need to call `database.CheckAndHandle()` manually

## Migration Examples

### Before (Old Pattern)
```go
func TestSomething(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)

    // Initialize database for testing
    err := database.CheckAndHandle()
    if err != nil {
        t.Fatalf("Failed to initialize database: %v", err)
    }

    // Create test data using helper function
    testUser, testWorkspace, err := database.CreateTestWorkspace("test@example.com", "Test User", "Test Workspace")
    if err != nil {
        t.Fatalf("Failed to create test workspace: %v", err)
    }
    defer database.CleanupTestWorkspaceAndProject(testUser.ID, testWorkspace.ID, "")

    testProject, err := database.CreateTestProject(testWorkspace.ID, "Test Project", "test-project-123")
    if err != nil {
        t.Fatalf("Failed to create test project: %v", err)
    }
    defer database.CleanupTestWorkspaceAndProject(testUser.ID, testWorkspace.ID, testProject.ID)

    // Test logic...
}
```

### After (New Pattern)
```go
func TestSomething(t *testing.T) {
    // Setup
    gin.SetMode(gin.TestMode)

    // Initialize test workspace with project and defer cleanup
    setup, err := database.InitTestWorkspaceWithProject(
        "test@example.com", 
        "Test User", 
        "Test Workspace", 
        "Test Project", 
        "test-project-123",
    )
    if err != nil {
        t.Fatalf("Failed to initialize test workspace: %v", err)
    }
    defer setup.Cleanup()

    testProject := setup.Project
    // Use setup.User and setup.Workspace if needed

    // Test logic...
}
```

## Usage Guidelines

### 1. **Always Defer Cleanup Immediately**
```go
setup, err := database.InitTestWorkspaceWithProject(...)
if err != nil {
    t.Fatalf("Failed to initialize: %v", err)
}
defer setup.Cleanup() // ✅ Defer immediately after successful creation
```

### 2. **Access Setup Data**
```go
user := setup.User           // *database.User
workspace := setup.Workspace // *database.Workspace  
project := setup.Project     // *database.Project (only in InitTestWorkspaceWithProject)
```

### 3. **Choose the Right Function**
- Use `InitTestWorkspaceWithProject` when you need a complete setup
- Use `InitTestWorkspaceOnly` when you only need workspace or want to create projects manually
- Use existing `CreateTestProject` for additional projects in the same workspace

### 4. **Manual Project Creation (if needed)**
```go
// If using InitTestWorkspaceOnly and need a project later
setup, err := database.InitTestWorkspaceOnly(...)
defer setup.Cleanup()

// Create additional project manually
project, err := database.CreateTestProject(setup.Workspace.ID, "Extra Project", "extra-alias")
if err != nil {
    t.Fatalf("Failed to create project: %v", err)
}
// No need to defer cleanup for manual project - workspace cleanup handles it
```

## Updated Test Files

The following test files have been updated to use the new pattern:

1. `backend/src/echo/handler/endpoint/update_endpoint_handler_test.go`
2. `backend/src/replay/services/execute_replay_test.go`

## Legacy Functions Still Available

The old helper functions are still available for backward compatibility:
- `CreateTestWorkspace()`
- `CreateTestProject()`
- `CleanupTestWorkspaceAndProject()`

However, new tests should use the improved initialization pattern for better maintainability.
