package workspaces

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"beo-echo/backend/src/database"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockWorkspaceRepository is a mock implementation of WorkspaceRepository
type MockWorkspaceRepository struct {
	mock.Mock
}

// GetUserWorkspaces is a mock implementation
func (m *MockWorkspaceRepository) GetUserWorkspaces(ctx context.Context, userID string) ([]database.Workspace, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]database.Workspace), args.Error(1)
}

// GetUserWorkspacesWithRoles is a mock implementation
func (m *MockWorkspaceRepository) GetUserWorkspacesWithRoles(ctx context.Context, userID string) ([]WorkspaceWithRole, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]WorkspaceWithRole), args.Error(1)
}

// CreateWorkspace is a mock implementation
func (m *MockWorkspaceRepository) CreateWorkspace(ctx context.Context, workspace *database.Workspace, userID string) error {
	args := m.Called(ctx, workspace, userID)
	return args.Error(0)
}

// CheckWorkspaceRole is a mock implementation
func (m *MockWorkspaceRepository) CheckWorkspaceRole(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error) {
	args := m.Called(ctx, userID, workspaceID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.UserWorkspace), args.Error(1)
}

// IsUserWorkspaceAdmin is a mock implementation
func (m *MockWorkspaceRepository) IsUserWorkspaceAdmin(ctx context.Context, userID string, workspaceID string) (bool, error) {
	args := m.Called(ctx, userID, workspaceID)
	return args.Bool(0), args.Error(1)
}

// GetAllWorkspaces is a mock implementation
func (m *MockWorkspaceRepository) GetAllWorkspaces(ctx context.Context) ([]database.Workspace, error) {
	args := m.Called(ctx)
	return args.Get(0).([]database.Workspace), args.Error(1)
}

// GetUserByEmail is a mock implementation
func (m *MockWorkspaceRepository) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*database.User), args.Error(1)
}

// AddUserToWorkspace is a mock implementation
func (m *MockWorkspaceRepository) AddUserToWorkspace(ctx context.Context, workspaceID string, userID string, role string) error {
	args := m.Called(ctx, workspaceID, userID, role)
	return args.Error(0)
}

// GetWorkspaceMembers is a mock implementation
func (m *MockWorkspaceRepository) GetWorkspaceMembers(ctx context.Context, workspaceID string) ([]WorkspaceMember, error) {
	args := m.Called(ctx, workspaceID)
	return args.Get(0).([]WorkspaceMember), args.Error(1)
}

// TestAddMember tests the AddMember service method
func TestAddMember(t *testing.T) {
	// Setup
	mockRepo := new(MockWorkspaceRepository)
	service := NewWorkspaceService(mockRepo)
	
	workspaceID := "workspace-123"
	email := "test@example.com"
	role := "member"
	userID := "user-123"
	
	testUser := &database.User{
		ID:    userID,
		Email: email,
		Name:  "Test User",
	}
	
	// Test cases
	t.Run("Add existing user", func(t *testing.T) {
		// Mock repository responses
		mockRepo.On("GetUserByEmail", mock.Anything, email).Return(testUser, nil).Once()
		mockRepo.On("CheckWorkspaceRole", mock.Anything, userID, workspaceID).Return(nil, errors.New("not found")).Once()
		mockRepo.On("AddUserToWorkspace", mock.Anything, workspaceID, userID, role).Return(nil).Once()
		
		// Call the method
		result, err := service.AddMember(context.Background(), workspaceID, email, role)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, userID, result["user_id"])
		assert.Equal(t, email, result["email"])
		assert.Equal(t, "added", result["status"])
		assert.Equal(t, role, result["role"])

		// Verify mocks
		mockRepo.AssertExpectations(t)
	})
	t.Run("User already a member", func(t *testing.T) {
		existingRole := &database.UserWorkspace{
			UserID:      userID,
			WorkspaceID: workspaceID,
			Role:        "admin",
		}
		
		// Mock repository responses
		mockRepo.On("GetUserByEmail", mock.Anything, email).Return(testUser, nil).Once()
		mockRepo.On("CheckWorkspaceRole", mock.Anything, userID, workspaceID).Return(existingRole, nil).Once()
		
		// Call the method
		result, err := service.AddMember(context.Background(), workspaceID, email, role)
		
		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, userID, result["user_id"])
		assert.Equal(t, email, result["email"])
		assert.Equal(t, "already_member", result["status"])
		assert.Equal(t, "admin", result["role"])
		
		// Verify mocks
		mockRepo.AssertExpectations(t)
	})
	
	t.Run("User not found", func(t *testing.T) {
		// Mock repository responses - user not found
		mockRepo.On("GetUserByEmail", mock.Anything, email).Return((*database.User)(nil), errors.New("not found")).Once()
		
		// Call the method
		result, err := service.AddMember(context.Background(), workspaceID, email, role)
		
		// Assertions
		assert.Error(t, err)
		assert.Nil(t, result)
		
		// Verify mocks
		mockRepo.AssertExpectations(t)
	})
}

// TestGetWorkspaceMembers tests the GetWorkspaceMembers service method
func TestGetWorkspaceMembers(t *testing.T) {
	// Setup
	mockRepo := new(MockWorkspaceRepository)
	service := NewWorkspaceService(mockRepo)

	workspaceID := "workspace-123"

	// Define some test members
	testMembers := []WorkspaceMember{
		{ID: "user-1", Name: "User One", Email: "user1@example.com", Role: "admin"},
		{ID: "user-2", Name: "User Two", Email: "user2@example.com", Role: "member"},
	}

	t.Run("Get workspace members successfully", func(t *testing.T) {
		// Mock repository response
		mockRepo.On("GetWorkspaceMembers", mock.Anything, workspaceID).Return(testMembers, nil).Once()

		// Call the method
		members, err := service.GetWorkspaceMembers(context.Background(), workspaceID)

		// Assertions
		assert.NoError(t, err)
		assert.Equal(t, 2, len(members))
		assert.Equal(t, testMembers, members)

		// Verify mocks
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error getting workspace members", func(t *testing.T) {
		// Mock repository response - error case
		mockRepo.On("GetWorkspaceMembers", mock.Anything, workspaceID).Return([]WorkspaceMember{}, errors.New("database error")).Once()

		// Call the method
		members, err := service.GetWorkspaceMembers(context.Background(), workspaceID)

		// Assertions
		assert.Error(t, err)
		assert.Empty(t, members)

		// Verify mocks
		mockRepo.AssertExpectations(t)
	})
}

// TestAddMemberHandler_ValidRequest tests the valid request case for AddMember handler
func TestAddMemberHandler_ValidRequest(t *testing.T) {
	// Setup Gin in test mode
	gin.SetMode(gin.TestMode)
	
	// Create a mock service
	mockRepo := new(MockWorkspaceRepository)
	service := NewWorkspaceService(mockRepo)
	handler := NewWorkspaceHandler(service)
	
	// Create a Gin router with a middleware that skips the database check
	router := gin.New()
	router.POST("/workspaces/:workspaceID/members", func(c *gin.Context) {
		// Simulate JWT middleware setting userID
		c.Set("userID", "test-user-id")
		
		// Skip the database check for user.IsOwner by directly setting isSystemOwner to true
		c.Set("isSystemOwner", true)
		
		// Call the handler with our mock context
		// Mock GetUserByEmail to return a valid user
		mockRepo.On("GetUserByEmail", mock.Anything, "test@example.com").Return(&database.User{
			ID:    "test-user-id",
			Email: "test@example.com",
			Name:  "Test User",
		}, nil).Once()
		
		// Mock AddUserToWorkspace to return success
		mockRepo.On("AddUserToWorkspace", mock.Anything, "workspace-123", "test-user-id", "member").Return(nil).Once()
		
		handler.AddMember(c)
	})

	// Setup request
	reqBody := map[string]string{
		"email": "test@example.com",
		"role":  "member",
	}
	reqJSON, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/workspaces/workspace-123/members", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)

	// Assertions
	assert.Equal(t, http.StatusOK, resp.Code)
	var respBody map[string]interface{}
	json.Unmarshal(resp.Body.Bytes(), &respBody)
	
	assert.True(t, respBody["success"].(bool))
	assert.Contains(t, respBody["message"], "User added to workspace")
	
	// Verify mocks
	mockRepo.AssertExpectations(t)
}

// TestAddMemberHandler_UserNotFound tests the case when a user is not found
func TestAddMemberHandler_UserNotFound(t *testing.T) {
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
		
		// Mock GetUserByEmail to return not found
		mockRepo.On("GetUserByEmail", mock.Anything, "nonexistent@example.com").Return((*database.User)(nil), errors.New("not found")).Once()
		
		handler.AddMember(c)
	})

	// Setup request
	reqBody := map[string]string{
		"email": "nonexistent@example.com",
		"role":  "member",
	}
	reqJSON, _ := json.Marshal(reqBody)

	// Create request
	req, _ := http.NewRequest(http.MethodPost, "/workspaces/workspace-123/members", bytes.NewReader(reqJSON))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Perform request
	router.ServeHTTP(resp, req)
