package database

import (
	"testing"
)

// Example demonstrating the new test setup pattern
// This shows how to use the new InitTestWorkspaceWithProject and InitTestWorkspaceOnly functions

func ExampleTestWithWorkspaceAndProject(t *testing.T) {
	// Initialize test workspace with project and defer cleanup
	setup, err := InitTestWorkspaceWithProject(
		"user@example.com",
		"Test User",
		"Test Workspace",
		"Test Project",
		"test-project-alias",
	)
	if err != nil {
		t.Fatalf("Failed to initialize test workspace: %v", err)
	}
	defer setup.Cleanup() // This will clean up everything

	// Use the setup data
	user := setup.User
	workspace := setup.Workspace
	project := setup.Project

	// Your test logic here...
	t.Logf("Created user: %s", user.Email)
	t.Logf("Created workspace: %s", workspace.Name)
	t.Logf("Created project: %s", project.Name)

	// No need for manual cleanup - defer handles it
}

func ExampleTestWithWorkspaceOnly(t *testing.T) {
	// Initialize test workspace only (no project) and defer cleanup
	setup, err := InitTestWorkspaceOnly(
		"user2@example.com",
		"Test User 2",
		"Test Workspace 2",
	)
	if err != nil {
		t.Fatalf("Failed to initialize test workspace: %v", err)
	}
	defer setup.Cleanup() // This will clean up everything

	// Use the setup data
	user := setup.User
	workspace := setup.Workspace

	// Your test logic here...
	t.Logf("Created user: %s", user.Email)
	t.Logf("Created workspace: %s", workspace.Name)

	// If you need a project later, create it manually:
	project, err := CreateTestProject(workspace.ID, "Manual Project", "manual-alias")
	if err != nil {
		t.Fatalf("Failed to create project: %v", err)
	}

	// Note: Manual project will be cleaned up automatically when workspace is cleaned up
	t.Logf("Created project: %s", project.Name)

	// No need for manual cleanup - defer handles it
}
