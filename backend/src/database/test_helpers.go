package database

import (
	"beo-echo/backend/src/utils"
	"testing"

	"github.com/google/uuid"
)

// CreateTestWorkspace creates a test workspace and user for testing purposes
func CreateTestWorkspace(userEmail string, userName string, workspaceName string) (*User, *Workspace, error) {
	db := DB

	// Create test user
	userID := uuid.New().String()
	testUser := &User{
		ID:       userID,
		Email:    userEmail,
		Name:     userName,
		IsActive: true,
	}

	err := db.Create(testUser).Error
	if err != nil {
		return nil, nil, err
	}

	// Create test workspace
	workspaceID := uuid.New().String()
	testWorkspace := &Workspace{
		ID:   workspaceID,
		Name: workspaceName,
	}

	err = db.Create(testWorkspace).Error
	if err != nil {
		return nil, nil, err
	}

	// Create user-workspace relationship
	userWorkspaceID := uuid.New().String()
	testUserWorkspace := &UserWorkspace{
		ID:          userWorkspaceID,
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        "admin",
	}

	err = db.Create(testUserWorkspace).Error
	if err != nil {
		return nil, nil, err
	}

	return testUser, testWorkspace, nil
}

// CreateTestProject creates a test project in the given workspace
func CreateTestProject(workspaceID string, projectName string, projectAlias string) (*Project, error) {
	return CreateTestProjectWithConfig(workspaceID, projectName, projectAlias, "")
}

// CreateTestProjectWithConfig creates a test project in the given workspace with advance config
func CreateTestProjectWithConfig(workspaceID string, projectName string, projectAlias string, advanceConfig string) (*Project, error) {
	db := DB

	projectID := uuid.New().String()
	testProject := &Project{
		ID:            projectID,
		Name:          projectName,
		WorkspaceID:   workspaceID,
		Mode:          ModeMock,
		Status:        "running",
		Alias:         projectAlias,
		AdvanceConfig: advanceConfig,
	}

	err := db.Create(testProject).Error
	if err != nil {
		return nil, err
	}

	return testProject, nil
}

// CreateTestEndpoint creates a test endpoint in the given project
func CreateTestEndpoint(projectID string, method string, path string) (*MockEndpoint, error) {
	return CreateTestEndpointWithConfig(projectID, method, path, "")
}

// CreateTestEndpointWithConfig creates a test endpoint in the given project with advance config
func CreateTestEndpointWithConfig(projectID string, method string, path string, advanceConfig string) (*MockEndpoint, error) {
	db := DB

	endpointID := uuid.New().String()
	testEndpoint := &MockEndpoint{
		ID:            endpointID,
		ProjectID:     projectID,
		Method:        method,
		Path:          path,
		Enabled:       true,
		ResponseMode:  "random",
		AdvanceConfig: advanceConfig,
	}

	err := db.Create(testEndpoint).Error
	if err != nil {
		return nil, err
	}

	return testEndpoint, nil
}

// CleanupTestWorkspaceAndProject removes test workspace, project, and user
func CleanupTestWorkspaceAndProject(userID, workspaceID, projectID string) {
	db := DB

	// Delete in order due to foreign key constraints
	if projectID != "" {
		db.Unscoped().Delete(&Project{}, "id = ?", projectID)
	}

	if userID != "" && workspaceID != "" {
		db.Unscoped().Delete(&UserWorkspace{}, "user_id = ? AND workspace_id = ?", userID, workspaceID)
	}

	if workspaceID != "" {
		db.Unscoped().Delete(&Workspace{}, "id = ?", workspaceID)
	}

	if userID != "" {
		db.Unscoped().Delete(&User{}, "id = ?", userID)
	}
}

// CleanupTestData removes test workspace, project, endpoint, and user
func CleanupTestData(userID, workspaceID, projectID, endpointID string) {
	db := DB

	// Delete in order due to foreign key constraints
	if endpointID != "" {
		db.Unscoped().Delete(&MockEndpoint{}, "id = ?", endpointID)
	}

	if projectID != "" {
		db.Unscoped().Delete(&Project{}, "id = ?", projectID)
	}

	if userID != "" && workspaceID != "" {
		db.Unscoped().Delete(&UserWorkspace{}, "user_id = ? AND workspace_id = ?", userID, workspaceID)
	}

	if workspaceID != "" {
		db.Unscoped().Delete(&Workspace{}, "id = ?", workspaceID)
	}

	if userID != "" {
		db.Unscoped().Delete(&User{}, "id = ?", userID)
	}
}

// SetupTestEnvironment sets up the test environment with database initialization
// and cleanup. Should be called at the beginning of each test.
func SetupTestEnvironment(t *testing.T) {
	// Setup test environment
	utils.SetupFolderConfigForTest()

	// Initialize database
	err := CheckAndHandle()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		utils.CleanupTestFolders()
	})
}

// TestWorkspaceSetup represents a complete test workspace setup
type TestWorkspaceSetup struct {
	User      *User
	Workspace *Workspace
	Project   *Project
	cleanup   func()
}

// InitTestWorkspaceWithProject creates a complete test setup with workspace and project
// Returns a setup object with cleanup function that should be deferred
func InitTestWorkspaceWithProject(userEmail, userName, workspaceName, projectName, projectAlias string) (*TestWorkspaceSetup, error) {
	// Initialize database
	err := CheckAndHandle()
	if err != nil {
		return nil, err
	}

	// Create workspace and user
	user, workspace, err := CreateTestWorkspace(userEmail, userName, workspaceName)
	if err != nil {
		return nil, err
	}

	// Create project
	project, err := CreateTestProject(workspace.ID, projectName, projectAlias)
	if err != nil {
		// Cleanup workspace if project creation fails
		CleanupTestWorkspaceAndProject(user.ID, workspace.ID, "")
		return nil, err
	}

	setup := &TestWorkspaceSetup{
		User:      user,
		Workspace: workspace,
		Project:   project,
		cleanup: func() {
			CleanupTestWorkspaceAndProject(user.ID, workspace.ID, project.ID)
		},
	}

	return setup, nil
}

// Cleanup should be deferred immediately after InitTestWorkspaceWithProject
func (s *TestWorkspaceSetup) Cleanup() {
	if s.cleanup != nil {
		s.cleanup()
	}
}

// TestWorkspaceOnlySetup represents a test workspace setup without project
type TestWorkspaceOnlySetup struct {
	User      *User
	Workspace *Workspace
	cleanup   func()
}

// InitTestWorkspaceOnly creates a test setup with only workspace (no project)
// Returns a setup object with cleanup function that should be deferred
func InitTestWorkspaceOnly(userEmail, userName, workspaceName string) (*TestWorkspaceOnlySetup, error) {
	// Initialize database
	err := CheckAndHandle()
	if err != nil {
		return nil, err
	}

	// Create workspace and user
	user, workspace, err := CreateTestWorkspace(userEmail, userName, workspaceName)
	if err != nil {
		return nil, err
	}

	setup := &TestWorkspaceOnlySetup{
		User:      user,
		Workspace: workspace,
		cleanup: func() {
			CleanupTestWorkspaceAndProject(user.ID, workspace.ID, "")
		},
	}

	return setup, nil
}

// Cleanup should be deferred immediately after InitTestWorkspaceOnly
func (s *TestWorkspaceOnlySetup) Cleanup() {
	if s.cleanup != nil {
		s.cleanup()
	}
}
