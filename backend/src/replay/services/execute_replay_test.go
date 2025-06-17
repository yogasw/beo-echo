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

	t.Run("ExecuteReplay_Success_HTTPBin", func(t *testing.T) {
		// Test data using httpbin.org (a free HTTP testing service)
		req := ExecuteReplayRequest{
			Protocol: "http",
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
			Protocol: "http",
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
			Protocol: "http",
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
			URL:      "https://non-existent-domain-12345.com",
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
		assert.NoError(t, err, "Network errors should not cause service error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.NotEmpty(t, response.Error, "Error field should contain invalid URL error")
		assert.Equal(t, 0, response.StatusCode, "Status code should be 0 for invalid URL")
		assert.NotEmpty(t, response.ReplayID, "ReplayID should still be generated")
	})

	t.Run("ExecuteReplay_WithQueryParams", func(t *testing.T) {
		// Test with multiple query parameters
		req := ExecuteReplayRequest{
			Protocol: "http",
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

	t.Run("ExecuteReplay_POST_JSON_Complex", func(t *testing.T) {
		// Test POST request with complex JSON body
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "POST",
			URL:      "https://httpbin.org/post",
			Headers: map[string]string{
				"Content-Type": "application/json",
				"Accept":       "application/json",
			},
			Body: `{
				"user": {
					"name": "Test User",
					"email": "test@example.com",
					"age": 30,
					"preferences": {
						"theme": "dark",
						"notifications": true
					},
					"tags": ["test", "api", "json"]
				},
				"metadata": {
					"client": "beo-echo-test",
					"version": "1.0.0"
				}
			}`,
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")

		// Verify response contains the complex JSON data we sent
		assert.Contains(t, response.ResponseBody, "Test User", "Response should contain user name")
		assert.Contains(t, response.ResponseBody, "test@example.com", "Response should contain email")
		assert.Contains(t, response.ResponseBody, "notifications", "Response should contain nested preferences")
		assert.Contains(t, response.ResponseBody, "tags", "Response should contain array data")
		assert.Contains(t, response.ResponseBody, "beo-echo-test", "Response should contain metadata")
	})

	t.Run("ExecuteReplay_POST_Form_Data", func(t *testing.T) {
		// Test POST request with form data
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "POST",
			URL:      "https://httpbin.org/post",
			Headers: map[string]string{
				"Content-Type": "application/x-www-form-urlencoded",
			},
			Body: "username=testuser&password=secret&remember=true&options=option1&options=option2",
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")

		// Verify response contains the form data we sent
		assert.Contains(t, response.ResponseBody, "username", "Response should contain form field")
		assert.Contains(t, response.ResponseBody, "testuser", "Response should contain form value")
		assert.Contains(t, response.ResponseBody, "password", "Response should contain form field")
		assert.Contains(t, response.ResponseBody, "options", "Response should contain multiple values field")
	})

	t.Run("ExecuteReplay_GET_Multiple_Params", func(t *testing.T) {
		// Test GET with multiple parameters including special characters
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "GET",
			URL:      "https://httpbin.org/get",
			Query: map[string]string{
				"id":          "12345",
				"filter":      "status:active",
				"sort":        "name:asc",
				"search":      "test query with spaces",
				"tags":        "golang,api,testing",
				"special":     "!@#$%^&*()",
				"coordinates": "37.7749,-122.4194",
				"limit":       "50",
				"offset":      "0",
				"include":     "details,metadata",
			},
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")

		// Verify response contains the query parameters we sent
		assert.Contains(t, response.ResponseBody, "12345", "Response should contain numeric ID")
		assert.Contains(t, response.ResponseBody, "status:active", "Response should contain filter value")
		assert.Contains(t, response.ResponseBody, "test query with spaces", "Response should contain encoded search query")
		assert.Contains(t, response.ResponseBody, "golang,api,testing", "Response should contain comma-separated tags")
		assert.Contains(t, response.ResponseBody, "limit", "Response should contain pagination params")
		assert.Contains(t, response.ResponseBody, "include", "Response should contain include params")
	})

	t.Run("ExecuteReplay_POST_Multipart_Form", func(t *testing.T) {
		// NOTE: For multipart/form-data with file attachments, we would normally use proper multipart encoding
		// However, for this test, we're using httpbin's post endpoint which can parse a simple JSON body
		// representing what a multipart form would contain. In a real implementation, you would use
		// proper multipart form data construction with boundaries.

		// This test simulates a multipart form request by setting the Content-Type header appropriately
		// and providing a specially formatted body that httpbin will recognize
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "POST",
			URL:      "https://httpbin.org/post",
			Headers: map[string]string{
				"Content-Type": "application/json", // Using JSON to simulate the form data for httpbin
			},
			Body: `{
				"file_simulation": "This is file content that would be uploaded",
				"filename": "test.txt",
				"form_field1": "value1",
				"form_field2": "value2"
			}`,
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")

		// Verify response contains our simulated file and form fields
		assert.Contains(t, response.ResponseBody, "file_simulation", "Response should contain file field")
		assert.Contains(t, response.ResponseBody, "This is file content", "Response should contain file content")
		assert.Contains(t, response.ResponseBody, "form_field1", "Response should contain form field")
		assert.Contains(t, response.ResponseBody, "value1", "Response should contain form value")
	})

	t.Run("ExecuteReplay_Auth_Headers", func(t *testing.T) {
		// Test request with different types of authentication headers
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "GET",
			URL:      "https://httpbin.org/headers",
			Headers: map[string]string{
				"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
				"X-API-Key":     "test-api-key-12345",
				"Cookie":        "session=abc123; user=testuser",
			},
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")

		// Verify response contains the authentication headers we sent
		assert.Contains(t, response.ResponseBody, "Authorization", "Response should contain Authorization header")
		assert.Contains(t, response.ResponseBody, "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9", "Response should contain Bearer token")
		assert.Contains(t, response.ResponseBody, "X-Api-Key", "Response should contain API key header")
		assert.Contains(t, response.ResponseBody, "test-api-key-12345", "Response should contain API key value")
		assert.Contains(t, response.ResponseBody, "Cookie", "Response should contain Cookie header")
		assert.Contains(t, response.ResponseBody, "session=abc123", "Response should contain session cookie")
	})

	t.Run("ExecuteReplay_OAuth2_Bearer_Token", func(t *testing.T) {
		// Test request with OAuth2 Bearer token authentication
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "GET",
			URL:      "https://httpbin.org/bearer",
			Headers: map[string]string{
				"Authorization": "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiaWF0IjoxNTE2MjM5MDIyfQ.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c",
			},
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")

		// httpbin.org/bearer validates the Bearer token and returns information about it
		assert.Contains(t, response.ResponseBody, "authenticated", "Response should indicate authenticated status")
		assert.Contains(t, response.ResponseBody, "true", "Response should show authenticated as true")
		assert.Contains(t, response.ResponseBody, "token", "Response should contain token information")
	})

	t.Run("ExecuteReplay_Basic_Auth", func(t *testing.T) {
		// Test request with Basic Authentication
		// The header value is "Basic " + base64("username:password")
		// For username "user" and password "pass", this is "Basic dXNlcjpwYXNz"
		req := ExecuteReplayRequest{
			Protocol: "http",
			Method:   "GET",
			URL:      "https://httpbin.org/basic-auth/user/pass",
			Headers: map[string]string{
				"Authorization": "Basic dXNlcjpwYXNz",
			},
		}

		// Execute the replay
		response, err := replayService.ExecuteReplay(ctx, testProject.ID, req)

		// Assertions
		assert.NoError(t, err, "ExecuteReplay should not return an error")
		assert.NotNil(t, response, "Response should not be nil")
		assert.Equal(t, 200, response.StatusCode, "Should receive 200 OK from httpbin")
		assert.NotEmpty(t, response.ResponseBody, "Response body should not be empty")

		// httpbin.org/basic-auth/user/pass validates the Basic Auth credentials
		assert.Contains(t, response.ResponseBody, "authenticated", "Response should indicate authenticated status")
		assert.Contains(t, response.ResponseBody, "true", "Response should show authenticated as true")
		assert.Contains(t, response.ResponseBody, "user", "Response should contain username")
	})
}
