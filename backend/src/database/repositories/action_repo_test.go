package repositories

import (
	"beo-echo/backend/src/database"
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActionRepository_Integration(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	// Create test workspace and project
	setup, err := database.InitTestWorkspaceWithProject(
		"test@example.com",
		"Test User",
		"Test Workspace",
		"Test Project",
		"test-project",
	)
	require.NoError(t, err)
	defer setup.Cleanup()

	// Create repository
	repo := NewActionRepository(database.DB)
	ctx := context.Background()

	t.Run("CreateAction_Success", func(t *testing.T) {
		action := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Test Replace Action",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointAfterRequest,
			Enabled:        true,
			Priority:       0,
			Config:         `{"target":"response_body","pattern":"test","replacement":"production","use_regex":false}`,
		}

		err := repo.CreateAction(ctx, action)
		assert.NoError(t, err)
		assert.NotEmpty(t, action.ID)

		// Cleanup
		defer database.DB.Unscoped().Delete(&database.Action{}, "id = ?", action.ID)
	})

	t.Run("GetActionByID_Success", func(t *testing.T) {
		// Create action
		action := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Test Get Action",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointAfterRequest,
			Enabled:        true,
			Priority:       1,
			Config:         `{"target":"response_body","pattern":"old","replacement":"new","use_regex":false}`,
		}
		err := repo.CreateAction(ctx, action)
		require.NoError(t, err)
		defer database.DB.Unscoped().Delete(&database.Action{}, "id = ?", action.ID)

		// Get action
		retrieved, err := repo.GetActionByID(ctx, action.ID)
		assert.NoError(t, err)
		assert.Equal(t, action.ID, retrieved.ID)
		assert.Equal(t, action.Name, retrieved.Name)
		assert.Equal(t, action.Type, retrieved.Type)
	})

	t.Run("GetActionsByProjectID_Success", func(t *testing.T) {
		// Create multiple actions
		action1 := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Action Priority 0",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointAfterRequest,
			Enabled:        true,
			Priority:       0,
			Config:         `{"target":"response_body","pattern":"test","replacement":"prod","use_regex":false}`,
		}
		err := repo.CreateAction(ctx, action1)
		require.NoError(t, err)
		defer database.DB.Unscoped().Delete(&database.Action{}, "id = ?", action1.ID)

		action2 := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Action Priority 1",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointBeforeRequest,
			Enabled:        true,
			Priority:       1,
			Config:         `{"target":"request_body","pattern":"dev","replacement":"staging","use_regex":false}`,
		}
		err = repo.CreateAction(ctx, action2)
		require.NoError(t, err)
		defer database.DB.Unscoped().Delete(&database.Action{}, "id = ?", action2.ID)

		// Get all actions for project
		actions, err := repo.GetActionsByProjectID(ctx, setup.Project.ID)
		assert.NoError(t, err)
		assert.GreaterOrEqual(t, len(actions), 2)

		// Verify ordering by priority
		if len(actions) >= 2 {
			assert.LessOrEqual(t, actions[0].Priority, actions[1].Priority)
		}
	})

	t.Run("UpdateAction_Success", func(t *testing.T) {
		// Create action
		action := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Original Name",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointAfterRequest,
			Enabled:        true,
			Priority:       0,
			Config:         `{"target":"response_body","pattern":"original","replacement":"updated","use_regex":false}`,
		}
		err := repo.CreateAction(ctx, action)
		require.NoError(t, err)
		defer database.DB.Unscoped().Delete(&database.Action{}, "id = ?", action.ID)

		// Update action
		action.Name = "Updated Name"
		action.Enabled = false
		err = repo.UpdateAction(ctx, action)
		assert.NoError(t, err)

		// Verify update
		retrieved, err := repo.GetActionByID(ctx, action.ID)
		assert.NoError(t, err)
		assert.Equal(t, "Updated Name", retrieved.Name)
		assert.False(t, retrieved.Enabled)
	})

	t.Run("DeleteAction_Success", func(t *testing.T) {
		// Create action
		action := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Action To Delete",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointAfterRequest,
			Enabled:        true,
			Priority:       0,
			Config:         `{"target":"response_body","pattern":"test","replacement":"prod","use_regex":false}`,
		}
		err := repo.CreateAction(ctx, action)
		require.NoError(t, err)

		// Delete action
		err = repo.DeleteAction(ctx, action.ID)
		assert.NoError(t, err)

		// Verify deletion
		_, err = repo.GetActionByID(ctx, action.ID)
		assert.Error(t, err)
	})

	t.Run("GetEnabledActionsByProjectAndPoint_Success", func(t *testing.T) {
		// Create enabled actions at different execution points
		enabledAfter := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Enabled After Request",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointAfterRequest,
			Enabled:        true,
			Priority:       0,
			Config:         `{"target":"response_body","pattern":"test","replacement":"prod","use_regex":false}`,
		}
		err := repo.CreateAction(ctx, enabledAfter)
		require.NoError(t, err)
		defer database.DB.Unscoped().Delete(&database.Action{}, "id = ?", enabledAfter.ID)

		enabledBefore := &database.Action{
			ProjectID:      setup.Project.ID,
			Name:           "Enabled Before Request",
			Type:           database.ActionTypeReplaceText,
			ExecutionPoint: database.ExecutionPointBeforeRequest,
			Enabled:        true,
			Priority:       0,
			Config:         `{"target":"request_body","pattern":"test","replacement":"prod","use_regex":false}`,
		}
		err = repo.CreateAction(ctx, enabledBefore)
		require.NoError(t, err)
		defer database.DB.Unscoped().Delete(&database.Action{}, "id = ?", enabledBefore.ID)

		// Get enabled actions for after_request
		afterActions, err := repo.GetEnabledActionsByProjectAndPoint(ctx, setup.Project.ID, database.ExecutionPointAfterRequest)
		assert.NoError(t, err)

		// Verify we find the after_request action
		foundAfter := false
		foundBefore := false

		for _, action := range afterActions {
			if action.ID == enabledAfter.ID {
				foundAfter = true
			}
			if action.ID == enabledBefore.ID {
				foundBefore = true
			}
			// All returned actions should be enabled and after_request
			assert.Equal(t, database.ExecutionPointAfterRequest, action.ExecutionPoint)
			assert.True(t, action.Enabled)
		}

		// Should find the after_request action
		assert.True(t, foundAfter, "Should find the after_request action")
		// Should NOT find before_request action
		assert.False(t, foundBefore, "Should not find before_request action")
	})
}
