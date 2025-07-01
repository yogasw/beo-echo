package repositories

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/services"
)

func TestProjectSearchRepository_CheckAliasAvailability_Integration(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	// Create test data
	user, workspace, err := database.CreateTestWorkspace("alias-test@example.com", "Alias Test User", "Alias Test Workspace")
	require.NoError(t, err)

	project, err := database.CreateTestProject(workspace.ID, "Test Project", "existing-alias")
	require.NoError(t, err)

	// Setup cleanup
	defer database.CleanupTestWorkspaceAndProject(user.ID, workspace.ID, project.ID)

	// Create repository
	repo := NewProjectSearchRepository(database.DB)

	tests := []struct {
		name           string
		alias          string
		expectedResult bool
	}{
		{
			name:           "existing alias should be unavailable",
			alias:          "existing-alias",
			expectedResult: false,
		},
		{
			name:           "non-existing alias should be available",
			alias:          "new-alias",
			expectedResult: true,
		},
		{
			name:           "case sensitive check",
			alias:          "EXISTING-ALIAS",
			expectedResult: true, // Should be available since it's case sensitive
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := repo.CheckAliasAvailability(context.Background(), tt.alias)
			require.NoError(t, err)
			assert.Equal(t, tt.expectedResult, result)
		})
	}
}

func TestProjectSearchRepository_SearchProjectsByNameInUserWorkspaces_Integration(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	// Create test data with multiple workspaces and users
	user1, workspace1, err := database.CreateTestWorkspace("search-test1@example.com", "Search Test User 1", "Search Workspace 1")
	require.NoError(t, err)

	user2, workspace2, err := database.CreateTestWorkspace("search-test2@example.com", "Search Test User 2", "Search Workspace 2")
	require.NoError(t, err)

	// Create projects in workspace1 (user1 is member)
	project1, err := database.CreateTestProject(workspace1.ID, "Raya API Service", "raya-api")
	require.NoError(t, err)

	project2, err := database.CreateTestProject(workspace1.ID, "Management Raya App", "management-raya")
	require.NoError(t, err)

	// Create project in workspace2 (user1 is NOT member)
	project3, err := database.CreateTestProject(workspace2.ID, "Raya Mobile App", "raya-mobile")
	require.NoError(t, err)

	// Setup cleanup
	defer func() {
		database.CleanupTestWorkspaceAndProject(user1.ID, workspace1.ID, project1.ID)
		database.CleanupTestWorkspaceAndProject(user1.ID, workspace1.ID, project2.ID)
		database.CleanupTestWorkspaceAndProject(user2.ID, workspace2.ID, project3.ID)
	}()

	// Create repository
	repo := NewProjectSearchRepository(database.DB)

	tests := []struct {
		name           string
		userID         string
		searchQuery    string
		expectedCount  int
		expectedNames  []string
	}{
		{
			name:          "user1 searches for 'raya' - should only find projects in their workspace",
			userID:        user1.ID,
			searchQuery:   "raya",
			expectedCount: 2, // "Raya API Service" and "Management Raya App"
			expectedNames: []string{"Raya API Service", "Management Raya App"},
		},
		{
			name:          "user2 searches for 'raya' - should only find projects in their workspace",
			userID:        user2.ID,
			searchQuery:   "raya",
			expectedCount: 1, // "Raya Mobile App"
			expectedNames: []string{"Raya Mobile App"},
		},
		{
			name:          "user1 searches for 'api' - should find matching project",
			userID:        user1.ID,
			searchQuery:   "api",
			expectedCount: 1, // "Raya API Service"
			expectedNames: []string{"Raya API Service"},
		},
		{
			name:          "user1 searches for non-existing term",
			userID:        user1.ID,
			searchQuery:   "nonexistent",
			expectedCount: 0,
			expectedNames: []string{},
		},
		{
			name:          "case insensitive search",
			userID:        user1.ID,
			searchQuery:   "RAYA",
			expectedCount: 2, // Should find same results as "raya"
			expectedNames: []string{"Raya API Service", "Management Raya App"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			results, err := repo.SearchProjectsByNameInUserWorkspaces(context.Background(), tt.userID, tt.searchQuery)
			require.NoError(t, err)
			assert.Len(t, results, tt.expectedCount)

			// Verify project names
			actualNames := make([]string, len(results))
			for i, result := range results {
				actualNames[i] = result.Name
			}

			for _, expectedName := range tt.expectedNames {
				assert.Contains(t, actualNames, expectedName)
			}

			// Verify workspace information is included
			for _, result := range results {
				assert.NotEmpty(t, result.WorkspaceID)
				assert.NotEmpty(t, result.WorkspaceName)
				assert.NotEmpty(t, result.ID)
				assert.NotEmpty(t, result.Alias)
			}
		})
	}
}

func TestProjectSearchService_CheckAliasAndSearchProjects_Integration(t *testing.T) {
	// Setup test environment
	database.SetupTestEnvironment(t)

	tests := []struct {
		name                string
		setupData           func() (*database.User, *database.Workspace, []*database.Project, error)
		query               string
		expectedAvailable   bool
		expectedProjectsLen int
		expectedProjectName string
	}{
		{
			name: "alias available and projects found",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				// Create test workspace and user
				user, workspace, err := database.CreateTestWorkspace("test@example.com", "Test User", "Test Workspace")
				if err != nil {
					return nil, nil, nil, err
				}

				// Create test projects
				project1, err := database.CreateTestProject(workspace.ID, "Raya Project", "raya-project")
				if err != nil {
					return nil, nil, nil, err
				}

				project2, err := database.CreateTestProject(workspace.ID, "Another Raya App", "another-raya")
				if err != nil {
					return nil, nil, nil, err
				}

				return user, workspace, []*database.Project{project1, project2}, nil
			},
			query:               "test-alias", // This alias doesn't exist, so it should be available
			expectedAvailable:   true,
			expectedProjectsLen: 0, // No projects contain "test-alias" in their name
		},
		{
			name: "alias unavailable and projects found by name",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				// Create test workspace and user
				user, workspace, err := database.CreateTestWorkspace("test2@example.com", "Test User 2", "Test Workspace 2")
				if err != nil {
					return nil, nil, nil, err
				}

				// Create test projects
				project1, err := database.CreateTestProject(workspace.ID, "Raya Project", "raya-project")
				if err != nil {
					return nil, nil, nil, err
				}

				project2, err := database.CreateTestProject(workspace.ID, "Another Project", "another-project")
				if err != nil {
					return nil, nil, nil, err
				}

				return user, workspace, []*database.Project{project1, project2}, nil
			},
			query:               "Raya", // Search by name, not alias
			expectedAvailable:   true,    // "Raya" as alias doesn't exist
			expectedProjectsLen: 1, // One project contains "raya-project" in name
			expectedProjectName: "Raya Project",
		},
		{
			name: "search by partial name match",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				// Create test workspace and user
				user, workspace, err := database.CreateTestWorkspace("test3@example.com", "Test User 3", "Test Workspace 3")
				if err != nil {
					return nil, nil, nil, err
				}

				// Create test projects with "raya" in their names
				project1, err := database.CreateTestProject(workspace.ID, "Raya API", "raya-api")
				if err != nil {
					return nil, nil, nil, err
				}

				project2, err := database.CreateTestProject(workspace.ID, "Management Raya App", "management-raya")
				if err != nil {
					return nil, nil, nil, err
				}

				project3, err := database.CreateTestProject(workspace.ID, "Other Project", "other-project")
				if err != nil {
					return nil, nil, nil, err
				}

				return user, workspace, []*database.Project{project1, project2, project3}, nil
			},
			query:               "raya", // Should find both "Raya API" and "Management Raya App"
			expectedAvailable:   true,    // "raya" alias doesn't exist
			expectedProjectsLen: 2,       // Two projects contain "raya" in their name
		},
		{
			name: "user not in workspace - no projects returned",
			setupData: func() (*database.User, *database.Workspace, []*database.Project, error) {
				// Create workspace and user, but don't join them
				user1, _, err := database.CreateTestWorkspace("test4@example.com", "Test User 4", "Test Workspace 4")
				if err != nil {
					return nil, nil, nil, err
				}

				// Create another workspace with different user
				_, workspace2, err := database.CreateTestWorkspace("test5@example.com", "Test User 5", "Test Workspace 5")
				if err != nil {
					return nil, nil, nil, err
				}

				// Create project in workspace2 (user1 is not a member)
				project, err := database.CreateTestProject(workspace2.ID, "Raya Project", "raya-in-other-workspace")
				if err != nil {
					return nil, nil, nil, err
				}

				return user1, workspace2, []*database.Project{project}, nil
			},
			query:               "raya", // user1 should not see projects from workspace2
			expectedAvailable:   true,
			expectedProjectsLen: 0, // No projects should be returned
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

			// Create repository and service
			repo := NewProjectSearchRepository(database.DB)
			service := services.NewProjectSearchService(repo)

			// Execute test
			result, err := service.CheckAliasAndSearchProjects(context.Background(), user.ID, tt.query)

			// Verify results
			require.NoError(t, err)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expectedAvailable, result.Available)
			assert.Len(t, result.Projects, tt.expectedProjectsLen)

			// Additional verification for specific test cases
			if tt.expectedProjectName != "" && len(result.Projects) > 0 {
				found := false
				for _, project := range result.Projects {
					if project.Name == tt.expectedProjectName {
						found = true
						assert.Equal(t, workspace.ID, project.WorkspaceID)
						assert.Equal(t, workspace.Name, project.WorkspaceName)
						break
					}
				}
				assert.True(t, found, "Expected project name not found in results")
			}
		})
	}
}
