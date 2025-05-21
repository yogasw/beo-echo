// TestAddMemberHandler_InvalidRequest tests the invalid request case for AddMember handler
func TestAddMemberHandler_InvalidRequest(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)

	// Create a mock service
	mockRepo := new(MockWorkspaceRepository)
	service := NewWorkspaceService(mockRepo)
	handler := NewWorkspaceHandler(service)

	// Create a Gin router
	router := gin.New()
	router.POST("/workspaces/:workspaceID/members", func(c *gin.Context) {
		// Simulate JWT middleware setting userID
		c.Set("userID", "test-user-id")

		// Skip the database check
		c.Set("isSystemOwner", true)

		handler.AddMember(c)
	})

	// Setup request with missing email
	reqBody := map[string]string{
		"role": "member",
	}
	reqJSON, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/workspaces/workspace-123/members", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusBadRequest, resp.Code)

	var respBody map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &respBody)

	assert.False(t, respBody["success"].(bool))
	assert.Contains(t, respBody["message"].(string), "required")

	// Verify mocks
	mockRepo.AssertExpectations(t)
}
