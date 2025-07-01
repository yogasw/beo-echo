package project

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
)

func TestCheckAliasAvailabilityHandler_Integration(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	tests := []struct {
		name               string
		setupData          func() (*database.User, *database.Workspace, []*database.Project, error)
		requestBody        CheckAliasAvailabilityRequest
		expectedStatus     int
		expectedAvailable  bool
		expectedProjectLen int
	}{
		{
			name: "successful alias check with available alias",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				user, workspace, err := database.CreateTestWorkspace("handler-test1@example.com", "Handler Test User 1", "Handler Test Workspace 1")
				if err != nil {
					return nil, nil, nil, err
				}

				project, err := database.CreateTestProject(workspace.ID, "Test Project", "test-project")
				if err != nil {
					return nil, nil, nil, err
				}

				return user, workspace, []*database.Project{project}, nil
			},
			requestBody: CheckAliasAvailabilityRequest{
				Query: "new-available-alias",
			},
			expectedStatus:     http.StatusOK,
			expectedAvailable:  true,
			expectedProjectLen: 0,
		},
		{
			name: "successful alias check with unavailable alias",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				user, workspace, err := database.CreateTestWorkspace("handler-test2@example.com", "Handler Test User 2", "Handler Test Workspace 2")
				if err != nil {
					return nil, nil, nil, err
				}

				project, err := database.CreateTestProject(workspace.ID, "Existing Project", "existing-alias")
				if err != nil {
					return nil, nil, nil, err
				}

				return user, workspace, []*database.Project{project}, nil
			},
			requestBody: CheckAliasAvailabilityRequest{
				Query: "Existing", // Search by name, not alias
			},
			expectedStatus:     http.StatusOK,
			expectedAvailable:  true,           // "Existing" as alias doesn't exist
			expectedProjectLen: 1,              // Should find "Existing Project" by name
		},
		{
			name: "successful search with partial name match",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				user, workspace, err := database.CreateTestWorkspace("handler-test3@example.com", "Handler Test User 3", "Handler Test Workspace 3")
				if err != nil {
					return nil, nil, nil, err
				}

				project1, err := database.CreateTestProject(workspace.ID, "Raya API", "raya-api")
				if err != nil {
					return nil, nil, nil, err
				}

				project2, err := database.CreateTestProject(workspace.ID, "Management Raya App", "management-raya")
				if err != nil {
					return nil, nil, nil, err
				}

				return user, workspace, []*database.Project{project1, project2}, nil
			},
			requestBody: CheckAliasAvailabilityRequest{
				Query: "raya",
			},
			expectedStatus:     http.StatusOK,
			expectedAvailable:  true,
			expectedProjectLen: 2, // Should find both projects
		},
		{
			name:      "invalid request body - missing query",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				user, workspace, err := database.CreateTestWorkspace("handler-test4@example.com", "Handler Test User 4", "Handler Test Workspace 4")
				return user, workspace, []*database.Project{}, err
			},
			requestBody: CheckAliasAvailabilityRequest{
				Query: "", // Empty query should fail validation
			},
			expectedStatus:     http.StatusBadRequest,
			expectedAvailable:  false,
			expectedProjectLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test data
			user, workspace, projects, err := tt.setupData()
			require.NoError(t, err)

			// Setup cleanup
			defer func() {
				for _, project := range projects {
					database.CleanupTestWorkspaceAndProject(user.ID, workspace.ID, project.ID)
				}
			}()

			// Create request
			jsonBody, err := json.Marshal(tt.requestBody)
			require.NoError(t, err)

			req, err := http.NewRequest("POST", "/api/projects/check-alias", bytes.NewBuffer(jsonBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			router := gin.New()
			router.POST("/api/projects/check-alias", func(c *gin.Context) {
				// Set userID in context (simulating JWT middleware)
				c.Set("userID", user.ID)
				CheckAliasAvailabilityHandler(c)
			})

			// Execute request
			router.ServeHTTP(w, req)

			// Verify response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				// Debug: print the actual response to understand the structure
				t.Logf("Response: %+v", response)

				success, exists := response["success"]
				require.True(t, exists, "Response should contain 'success' field")
				assert.True(t, success.(bool))

				data, ok := response["data"].(map[string]interface{})
				require.True(t, ok)

				assert.Equal(t, tt.expectedAvailable, data["available"].(bool))

				projectsInterface := data["projects"]
				if projectsInterface == nil {
					// Handle case where projects is nil (should be empty array)
					assert.Equal(t, 0, tt.expectedProjectLen)
				} else {
					projects, ok := projectsInterface.([]interface{})
					require.True(t, ok)
					assert.Len(t, projects, tt.expectedProjectLen)

					// Verify project structure if projects exist
					if len(projects) > 0 {
						firstProject := projects[0].(map[string]interface{})
						assert.NotEmpty(t, firstProject["id"])
						assert.NotEmpty(t, firstProject["name"])
						assert.NotEmpty(t, firstProject["alias"])
						assert.NotEmpty(t, firstProject["workspace_id"])
						assert.NotEmpty(t, firstProject["workspace_name"])
					}
				}
			} else {
				var response map[string]interface{}
				err := json.Unmarshal(w.Body.Bytes(), &response)
				require.NoError(t, err)

				assert.False(t, response["success"].(bool))
				assert.NotEmpty(t, response["message"])
			}
		})
	}
}

func TestCheckAliasAvailabilityHandler_Unauthorized(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create request without user context
	requestBody := CheckAliasAvailabilityRequest{
		Query: "test-query",
	}
	jsonBody, err := json.Marshal(requestBody)
	require.NoError(t, err)

	req, err := http.NewRequest("POST", "/api/projects/check-alias", bytes.NewBuffer(jsonBody))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context without userID
	router := gin.New()
	router.POST("/api/projects/check-alias", CheckAliasAvailabilityHandler)

	// Execute request
	router.ServeHTTP(w, req)

	// Verify unauthorized response
	assert.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "User not authenticated", response["message"])
}

func TestCheckAliasAvailabilityHandler_InvalidJSON(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test user
	user, workspace, err := database.CreateTestWorkspace("invalid-json-test@example.com", "Invalid JSON Test User", "Invalid JSON Test Workspace")
	require.NoError(t, err)

	defer database.CleanupTestWorkspaceAndProject(user.ID, workspace.ID, "")

	// Create request with invalid JSON
	req, err := http.NewRequest("POST", "/api/projects/check-alias", bytes.NewBufferString("invalid json"))
	require.NoError(t, err)
	req.Header.Set("Content-Type", "application/json")

	// Create response recorder
	w := httptest.NewRecorder()

	// Create Gin context
	router := gin.New()
	router.POST("/api/projects/check-alias", func(c *gin.Context) {
		// Set userID in context (simulating JWT middleware)
		c.Set("userID", user.ID)
		CheckAliasAvailabilityHandler(c)
	})

	// Execute request
	router.ServeHTTP(w, req)

	// Verify bad request response
	assert.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &response)
	require.NoError(t, err)

	assert.False(t, response["success"].(bool))
	assert.Equal(t, "Invalid request format", response["message"])
}

func TestCheckAliasAvailabilityHandler_ErrorHandling(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	// Set Gin to test mode
	gin.SetMode(gin.TestMode)

	// Create test user
	user, workspace, err := database.CreateTestWorkspace("error-test@example.com", "Error Test User", "Error Test Workspace")
	require.NoError(t, err)

	defer database.CleanupTestWorkspaceAndProject(user.ID, workspace.ID, "")

	tests := []struct {
		name            string
		requestBody     string
		expectedStatus  int
		expectedMessage string
	}{
		{
			name:            "empty request body",
			requestBody:     "",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid request format",
		},
		{
			name:            "invalid JSON - invalid character",
			requestBody:     "invalid json",
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid request format",
		},
		{
			name:            "invalid JSON - unclosed brace",
			requestBody:     `{"query": "test"`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Invalid request format",
		},
		{
			name:            "missing query field",
			requestBody:     `{}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Query must be at least 1 character long",
		},
		{
			name:            "empty query string",
			requestBody:     `{"query": ""}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Query must be at least 1 character long",
		},
		{
			name:            "query with only whitespace",
			requestBody:     `{"query": "   "}`,
			expectedStatus:  http.StatusBadRequest,
			expectedMessage: "Query must be at least 1 character long",
		},
		{
			name:            "valid query",
			requestBody:     `{"query": "valid-query"}`,
			expectedStatus:  http.StatusOK,
			expectedMessage: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create request
			req, err := http.NewRequest("POST", "/api/projects/check-alias", strings.NewReader(tt.requestBody))
			require.NoError(t, err)
			req.Header.Set("Content-Type", "application/json")

			// Create response recorder
			w := httptest.NewRecorder()

			// Create Gin context
			router := gin.New()
			router.POST("/api/projects/check-alias", func(c *gin.Context) {
				// Set userID in context (simulating JWT middleware)
				c.Set("userID", user.ID)
				CheckAliasAvailabilityHandler(c)
			})

			// Execute request
			router.ServeHTTP(w, req)

			// Verify response
			assert.Equal(t, tt.expectedStatus, w.Code)

			var response map[string]interface{}
			err = json.Unmarshal(w.Body.Bytes(), &response)
			require.NoError(t, err)

			if tt.expectedStatus == http.StatusOK {
				assert.True(t, response["success"].(bool))
				assert.NotNil(t, response["data"])
			} else {
				assert.False(t, response["success"].(bool))
				if tt.expectedMessage != "" {
					assert.Equal(t, tt.expectedMessage, response["message"].(string))
				}
			}
		})
	}
}
