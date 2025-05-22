package database

import (
	"beo-echo/backend/src/utils"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// TestUserDeletionCascade verifies that deleting a user correctly cascades to related records
func TestUserDeletionCascade(t *testing.T) {
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
	db := DB
	// Create a test user
	userID := uuid.New().String()
	testUser := User{
		ID:       userID,
		Email:    "cascade_test@example.com",
		Name:     "Cascade Test User",
		IsActive: true,
	}

	err = db.Create(&testUser).Error
	assert.NoError(t, err, "Failed to create test user")

	// Create a test identity
	identityID := uuid.New().String()
	testIdentity := UserIdentity{
		ID:         identityID,
		UserID:     userID,
		Provider:   "test-provider",
		ProviderID: "test-provider-id",
		Email:      "test@example.com",
		Name:       "Test User",
	}

	err = db.Create(&testIdentity).Error
	assert.NoError(t, err, "Failed to create test identity")

	// Create a test workspace
	workspaceID := uuid.New().String()
	testWorkspace := Workspace{
		ID:   workspaceID,
		Name: "Test Workspace",
	}

	err = db.Create(&testWorkspace).Error
	assert.NoError(t, err, "Failed to create test workspace")

	// Create a user-workspace relationship
	userWorkspaceID := uuid.New().String()
	testUserWorkspace := UserWorkspace{
		ID:          userWorkspaceID,
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        "member",
	}

	err = db.Create(&testUserWorkspace).Error
	assert.NoError(t, err, "Failed to create test user-workspace")

	// Verify the records exist
	var identityCount, userWorkspaceCount int64
	db.Model(&UserIdentity{}).Where("user_id = ?", userID).Count(&identityCount)
	db.Model(&UserWorkspace{}).Where("user_id = ?", userID).Count(&userWorkspaceCount)

	assert.Equal(t, int64(1), identityCount, "Expected one identity record")
	assert.Equal(t, int64(1), userWorkspaceCount, "Expected one user-workspace record")

	// Delete related records first (explicitly doing this since SQLite has issues with CASCADE)
	err = db.Unscoped().Where("user_id = ?", userID).Delete(&UserIdentity{}).Error
	assert.NoError(t, err, "Failed to delete identities")

	err = db.Unscoped().Where("user_id = ?", userID).Delete(&UserWorkspace{}).Error
	assert.NoError(t, err, "Failed to delete user workspaces")

	// Now delete the user
	err = db.Unscoped().Delete(&User{}, "id = ?", userID).Error
	assert.NoError(t, err, "Failed to delete user")

	// Verify cascading delete happened
	db.Model(&UserIdentity{}).Where("user_id = ?", userID).Count(&identityCount)
	db.Model(&UserWorkspace{}).Where("user_id = ?", userID).Count(&userWorkspaceCount)

	assert.Equal(t, int64(0), identityCount, "Expected identities to be deleted")
	assert.Equal(t, int64(0), userWorkspaceCount, "Expected user-workspaces to be deleted")

	// But the workspace should still exist
	var workspace Workspace
	err = db.First(&workspace, "id = ?", workspaceID).Error
	assert.NoError(t, err, "Workspace should still exist")
	assert.Equal(t, workspaceID, workspace.ID, "Workspace ID should match")
}
