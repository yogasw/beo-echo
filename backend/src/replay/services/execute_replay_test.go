package services

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/database/repositories"
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

	// Initialize database
	err := database.CheckAndHandle()
	require.NoError(t, err, "Failed to initialize database")

	// Cleanup after test
	t.Cleanup(func() {
		utils.CleanupTestFolders()
	})

	// Create test workspace and project
	testUser, testWorkspace, err := database.CreateTestWorkspace(
		"replay_test@example.com",
		"Replay Test User",
		"Test Workspace",
	)
	require.NoError(t, err, "Failed to create test workspace")

	testProject, err := database.CreateTestProject(
		testWorkspace.ID,
		"Test Project",
		"test-project",
	)
	require.NoError(t, err, "Failed to create test project")

	// Cleanup test data
	defer database.CleanupTestWorkspaceAndProject(testUser.ID, testWorkspace.ID, testProject.ID)

	// Setup repository and service
	db := database.DB
	replayRepo := repositories.NewReplayRepository(db)
	replayService := NewReplayService(replayRepo)

	ctx := context.Background()

	t.Run("ExecuteReplay_Success_HTTPBin", func(t *testing.T) {
		// Test data using httpbin.org (a free HTTP testing service)
		req := ExecuteReplayRequest{
			Protocol: "https",
			Method:   "GET",
			URL:      "https://httpbin.org/get",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"User-Agent":   "beo-echo-test",
			},
			Query: map[string]string{
				"test_param": "test_value",
				"replay":     "true",
			},
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.NotEmpty(t, response.ReplayID, "ReplayID should be generated")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")
		assert.Greater(t, response.LatencyMS, 0, "Latency should be greater than 0")
		assert.Empty(t, response.Error, "Error field should be empty on success")

		// Verify response headers are captured
		assert.NotEmpty(t, response.ResponseHeaders, "Response headers should be captured")
		assert.Contains(t, response.ResponseHeaders, "Content-Type", "Should contain Content-Type header")

		t.Logf("Replay ID: %s", response.ReplayID)
		t.Logf("Latency: %d ms", response.LatencyMS)
		t.Logf("Response Headers: %+v", response.ResponseHeaders)
	})

	t.Run("ExecuteReplay_Success_POST_with_Body", func(t *testing.T) {
		// Test POST request with body
		req := ExecuteReplayRequest{
			Protocol: "https",
			Method:   "POST",
			URL:      "https://httpbin.org/post",
			Headers: map[string]string{
				"Content-Type": "application/json",
			},
			Body: `{"test": "data", "number": 123}`,
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")
		assert.Greater(t, response.LatencyMS, 0, "Latency should be greater than 0")

		// For httpbin.org/post, the response should echo back the request data
		assert.Contains(t, response.ResponseBody, "test", "Response should contain the posted data")
		assert.Contains(t, response.ResponseBody, "data", "Response should contain the posted data")
	})

	t.Run("ExecuteReplay_InvalidProject", func(t *testing.T) {
		// Test with non-existent project ID
		req := ExecuteReplayRequest{
			Protocol: "https",
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
		req := ExecuteReplayRequest{
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

	t.Run("ExecuteReplay_HTTPError_NotFound", func(t *testing.T) {
		// Test with URL that returns 404
		req := ExecuteReplayRequest{
			Protocol: "https",
			Method:   "GET",
			URL:      "https://httpbin.org/status/404",
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions - HTTP errors should not cause service errors
		assert.NoError(t, err, "HTTP 404 should not cause service error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 404, response.StatusCode, "Should receive 404 status code")
		assert.NotEmpty(t, response.ReplayID, "ReplayID should still be generated")
		assert.Greater(t, response.LatencyMS, 0, "Latency should be recorded")
	})

	t.Run("ExecuteReplay_NetworkError", func(t *testing.T) {
		// Test with invalid URL that will cause network error
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "GET",
			URL:      "http://non-existent-domain-12345.com",
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions - Network errors should be captured in response
		assert.NoError(t, err, "Network errors should not cause service error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.NotEmpty(t, response.Error, "Error field should contain network error")
		assert.Equal(t, 0, response.StatusCode, "Status code should be 0 for network errors")
		assert.NotEmpty(t, response.ReplayID, "ReplayID should still be generated")
		assert.Greater(t, response.LatencyMS, 0, "Latency should be recorded")
	})

	t.Run("ExecuteReplay_InvalidURL", func(t *testing.T) {
		// Test with malformed URL
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "GET",
			URL:      "not-a-valid-url",
			Query: map[string]string{
				"param": "value",
			},
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.Error(t, err, "Should return error for invalid URL")
		assert.Nil(t, response, "Response should be nil on error")
		assert.Contains(t, err.Error(), "invalid URL format", "Error should indicate invalid URL format")
	})

	t.Run("ExecuteReplay_WithQueryParams", func(t *testing.T) {
		// Test with multiple query parameters
		req := ExecuteReplayRequest{
			Protocol: "https",
			Method:   "GET",
			URL:      "https://httpbin.org/get",
			Headers: map[string]string{
				"User-Agent": "beo-echo-test-agent",
			},
			Query: map[string]string{
				"param1": "value1",
				"param2": "value2",
				"search": "test query",
			},
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK")

		// httpbin.org/get returns the query parameters in the response
		assert.Contains(t, response.ResponseBody, "param1", "Response should contain query parameter")
		assert.Contains(t, response.ResponseBody, "value1", "Response should contain query parameter value")
		assert.Contains(t, response.ResponseBody, "param2", "Response should contain query parameter")
		assert.Contains(t, response.ResponseBody, "test query", "Response should contain encoded query parameter")
	})
}
