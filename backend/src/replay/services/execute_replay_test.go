package services

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/database/repositories"
	"beo-echo/backend/src/replay/models"
	"beo-echo/backend/src/utils"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestExecuteReplay tests the ExecuteReplay function with real HTTP calls
func TestExecuteReplay(t *testing.T) {
	// Setup test environment
	utils.SetupFolderConfigForTest()

	// Cleanup after test
	t.Cleanup(func() {
		utils.CleanupTestFolders()
	})

	// Initialize test workspace with project and defer cleanup
	setup, err := database.InitTestWorkspaceWithProject(
		"replay_test@example.com",
		"Replay Test User",
		"Test Workspace",
		"Test Project",
		"test-project",
	)
	require.NoError(t, err, "Failed to initialize test workspace")
	defer setup.Cleanup()

	testProject := setup.Project

	// Setup repository and service
	db := database.DB
	replayRepo := repositories.NewReplayRepository(db)
	replayService := NewReplayService(replayRepo)

	ctx := context.Background()

	t.Run("ExecuteReplay_InvalidProject", func(t *testing.T) {
		// Test with non-existent project ID
		req := models.ExecuteReplayRequest{
			Protocol: "http",
			Method:   "GET",
			URL:      "https://httpbin.org/get",
		}

		// Execute the replay with invalid project ID
		response, err := replayService.ExecuteReplay(ctx, "invalid-project-id", req)

		// Assertions
		assert.Error(t, err, "Should return error for invalid project ID")
		assert.Nil(t, response, "Response should be nil on error")
		assert.Contains(t, err.Error(), "project not found", "Error should indicate project not found")
	})

	t.Run("ExecuteReplay_InvalidProtocol", func(t *testing.T) {
		// Test with unsupported protocol
		req := models.ExecuteReplayRequest{
			Protocol: "ftp",
			Method:   "GET",
			URL:      "ftp://example.com/file",
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.Error(t, err, "Should return error for unsupported protocol")
		assert.Nil(t, response, "Response should be nil on error")
		assert.Contains(t, err.Error(), "unsupported protocol", "Error should indicate unsupported protocol")
	})
}

