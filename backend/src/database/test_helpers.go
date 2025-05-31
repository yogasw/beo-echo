package database

import (
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
	db := DB

	projectID := uuid.New().String()
	testProject := &Project{
		ID:          projectID,
		Name:        projectName,
		WorkspaceID: workspaceID,
		Mode:        ModeMock,
		Status:      "running",
		Alias:       projectAlias,
	}

	err := db.Create(testProject).Error
	if err != nil {
		return nil, err
	}

	return testProject, nil
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
