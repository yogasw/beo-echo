package workspaces

import (
	"beo-echo/backend/src/database"
	"context"
	"testing"
	"time"

	"github.com/rs/zerolog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

// testWorkspaceRepo implements WorkspaceRepository for testing
type testWorkspaceRepo struct {
	db *gorm.DB
}

func (r *testWorkspaceRepo) GetUserWorkspaces(ctx context.Context, userID string) ([]database.Workspace, error) {
	var workspaces []database.Workspace
	var userWorkspaces []database.UserWorkspace

	err := r.db.Where("user_id = ?", userID).Find(&userWorkspaces).Error
	if err != nil {
		return nil, err
	}

	var workspaceIDs []string
	for _, uw := range userWorkspaces {
		workspaceIDs = append(workspaceIDs, uw.WorkspaceID)
	}

	if len(workspaceIDs) > 0 {
		err = r.db.Where("id IN ?", workspaceIDs).Find(&workspaces).Error
	}

	return workspaces, err
}

func (r *testWorkspaceRepo) GetUserWorkspacesWithRoles(ctx context.Context, userID string) ([]WorkspaceWithRole, error) {
	var results []WorkspaceWithRole

	query := `
		SELECT w.id, w.name, uw.role 
		FROM workspaces w 
		JOIN user_workspaces uw ON w.id = uw.workspace_id 
		WHERE uw.user_id = ?
	`

	err := r.db.Raw(query, userID).Scan(&results).Error
	return results, err
}

func (r *testWorkspaceRepo) CreateWorkspace(ctx context.Context, workspace *database.Workspace, userID string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Create workspace
		if err := tx.Create(workspace).Error; err != nil {
			return err
		}

		// Add user as owner
		userWorkspace := &database.UserWorkspace{
			UserID:      userID,
			WorkspaceID: workspace.ID,
			Role:        "owner",
		}

		return tx.Create(userWorkspace).Error
	})
}

func (r *testWorkspaceRepo) CountUserWorkspaces(ctx context.Context, userID string) (int, error) {
	var count int64
	err := r.db.Model(&database.UserWorkspace{}).Where("user_id = ?", userID).Count(&count).Error
	return int(count), err
}

func (r *testWorkspaceRepo) GetUserByID(ctx context.Context, userID string) (*database.User, error) {
	var user database.User
	err := r.db.First(&user, "id = ?", userID).Error
	return &user, err
}

func (r *testWorkspaceRepo) CheckWorkspaceRole(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error) {
	var userWorkspace database.UserWorkspace
	err := r.db.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userWorkspace).Error
	if err != nil {
		return nil, err
	}
	return &userWorkspace, nil
}

func (r *testWorkspaceRepo) IsUserWorkspaceAdmin(ctx context.Context, userID string, workspaceID string) (bool, error) {
	var count int64
	err := r.db.Model(&database.UserWorkspace{}).
		Where("user_id = ? AND workspace_id = ? AND role IN ?", userID, workspaceID, []string{"owner", "admin"}).
		Count(&count).Error
	return count > 0, err
}

func (r *testWorkspaceRepo) GetUserByEmail(ctx context.Context, email string) (*database.User, error) {
	var user database.User
	err := r.db.Where("email = ?", email).First(&user).Error
	return &user, err
}

func (r *testWorkspaceRepo) AddUserToWorkspace(ctx context.Context, workspaceID string, userID string, role string) error {
	userWorkspace := &database.UserWorkspace{
		UserID:      userID,
		WorkspaceID: workspaceID,
		Role:        role,
	}
	return r.db.Create(userWorkspace).Error
}

func (r *testWorkspaceRepo) GetWorkspaceMembers(ctx context.Context, workspaceID string) ([]WorkspaceMember, error) {
	var members []WorkspaceMember

	query := `
		SELECT u.id, u.name, u.email, uw.role 
		FROM users u 
		JOIN user_workspaces uw ON u.id = uw.user_id 
		WHERE uw.workspace_id = ?
	`

	err := r.db.Raw(query, workspaceID).Scan(&members).Error
	return members, err
}

func (r *testWorkspaceRepo) GetWorkspacesWithAutoInvite(ctx context.Context) ([]database.Workspace, error) {
	var workspaces []database.Workspace
	err := r.db.Where("auto_invite_domains IS NOT NULL AND auto_invite_domains != ''").Find(&workspaces).Error
	return workspaces, err
}

func (r *testWorkspaceRepo) CheckUserWorkspaceMembership(ctx context.Context, userID string, workspaceID string) (*database.UserWorkspace, error) {
	return r.CheckWorkspaceRole(ctx, userID, workspaceID)
}

func (r *testWorkspaceRepo) CreateUserWorkspaceMembership(ctx context.Context, membership *database.UserWorkspace) error {
	return r.db.Create(membership).Error
}

func (r *testWorkspaceRepo) CreateProject(ctx context.Context, project *database.Project) error {
	return r.db.Create(project).Error
}

func (r *testWorkspaceRepo) CreateEndpoint(ctx context.Context, endpoint *database.MockEndpoint) error {
	return r.db.Create(endpoint).Error
}

func (r *testWorkspaceRepo) CreateResponse(ctx context.Context, response *database.MockResponse) error {
	return r.db.Create(response).Error
}

func (r *testWorkspaceRepo) CheckProjectAliasExists(ctx context.Context, alias string) (bool, error) {
	var count int64
	err := r.db.Model(&database.Project{}).Where("alias = ?", alias).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *testWorkspaceRepo) GetAllWorkspaces(ctx context.Context) ([]database.Workspace, error) {
	var workspaces []database.Workspace
	err := r.db.Find(&workspaces).Error
	return workspaces, err
}

func TestCreateDemoWorkspace(t *testing.T) {
	// Setup test database
	database.SetupTestEnvironment(t)

	// Get database connection
	db := database.DB

	// Create workspace service with test repository
	workspaceService := &WorkspaceService{
		repo: &testWorkspaceRepo{db: db},
	}

	// Setup context with logger
	logger := zerolog.New(zerolog.NewTestWriter(t))
	ctx := logger.WithContext(context.Background())

	// Test data
	testUserName := "Test User"
	testUserEmail := "testuser@example.com"
	workspaceName := "Test Demo Workspace"

	// Step 1: Create a test user first
	testUser := &database.User{
		Name:                 testUserName,
		Email:                testUserEmail,
		Password:             "test-password",
		IsActive:             true,
		MaxWorkspaces:        nil, // Use system default
		MaxProjectsWorkspace: nil, // Use system default
	}

	err := db.Create(testUser).Error
	require.NoError(t, err, "Failed to create test user")
	require.NotEmpty(t, testUser.ID, "Test user ID should not be empty")

	t.Logf("Created test user with ID: %s", testUser.ID)

	// Step 2: Test CreateDemoWorkspace
	workspace, project, err := workspaceService.CreateDemoWorkspace(ctx, testUser.ID, testUserName, workspaceName)

	// Assertions
	require.NoError(t, err, "CreateDemoWorkspace should not return error")
	require.NotNil(t, workspace, "Workspace should not be nil")
	require.NotNil(t, project, "Project should not be nil")

	// Verify workspace properties
	assert.Equal(t, workspaceName, workspace.Name, "Workspace name should match")
	assert.NotEmpty(t, workspace.ID, "Workspace ID should not be empty")

	// Verify project properties
	assert.Equal(t, "Demo Project", project.Name, "Project name should be 'Demo Project'")
	assert.Equal(t, workspace.ID, project.WorkspaceID, "Project should belong to the workspace")
	assert.Equal(t, database.ModeMock, project.Mode, "Project mode should be mock")
	assert.Equal(t, "running", project.Status, "Project status should be running")
	assert.NotEmpty(t, project.Alias, "Project alias should not be empty")
	assert.NotEmpty(t, project.ID, "Project ID should not be empty")

	t.Logf("Created workspace ID: %s, project ID: %s, alias: %s", workspace.ID, project.ID, project.Alias)

	// Step 3: Verify user is owner of the workspace
	userWorkspace, err := workspaceService.CheckWorkspaceRole(ctx, testUser.ID, workspace.ID)
	require.NoError(t, err, "Should be able to check workspace role")
	require.NotNil(t, userWorkspace, "User workspace relationship should exist")
	assert.Equal(t, "owner", userWorkspace.Role, "User should be owner of the workspace")
	assert.Equal(t, testUser.ID, userWorkspace.UserID, "User ID should match")
	assert.Equal(t, workspace.ID, userWorkspace.WorkspaceID, "Workspace ID should match")

	t.Logf("Verified user %s is %s of workspace %s", testUser.ID, userWorkspace.Role, workspace.ID)

	// Step 4: Verify endpoints were created
	var endpoints []database.MockEndpoint
	err = db.Where("project_id = ?", project.ID).Find(&endpoints).Error
	require.NoError(t, err, "Should be able to query endpoints")
	assert.Len(t, endpoints, 2, "Should have exactly 2 endpoints (GET and POST)")

	// Find GET and POST endpoints
	var getEndpoint, postEndpoint *database.MockEndpoint
	for i := range endpoints {
		if endpoints[i].Method == "GET" && endpoints[i].Path == "/testing" {
			getEndpoint = &endpoints[i]
		}
		if endpoints[i].Method == "POST" && endpoints[i].Path == "/testing" {
			postEndpoint = &endpoints[i]
		}
	}

	// Verify GET endpoint
	require.NotNil(t, getEndpoint, "GET /testing endpoint should exist")
	assert.Equal(t, "GET", getEndpoint.Method, "Method should be GET")
	assert.Equal(t, "/testing", getEndpoint.Path, "Path should be /testing")
	assert.True(t, getEndpoint.Enabled, "Endpoint should be enabled")
	assert.Equal(t, "random", getEndpoint.ResponseMode, "Response mode should be random")

	// Verify POST endpoint
	require.NotNil(t, postEndpoint, "POST /testing endpoint should exist")
	assert.Equal(t, "POST", postEndpoint.Method, "Method should be POST")
	assert.Equal(t, "/testing", postEndpoint.Path, "Path should be /testing")
	assert.True(t, postEndpoint.Enabled, "Endpoint should be enabled")
	assert.Equal(t, "random", postEndpoint.ResponseMode, "Response mode should be random")

	t.Logf("Verified endpoints - GET ID: %s, POST ID: %s", getEndpoint.ID, postEndpoint.ID)

	// Step 5: Verify responses were created
	var getResponses []database.MockResponse
	err = db.Where("endpoint_id = ?", getEndpoint.ID).Find(&getResponses).Error
	require.NoError(t, err, "Should be able to query GET responses")
	assert.Len(t, getResponses, 1, "GET endpoint should have exactly 1 response")

	var postResponses []database.MockResponse
	err = db.Where("endpoint_id = ?", postEndpoint.ID).Find(&postResponses).Error
	require.NoError(t, err, "Should be able to query POST responses")
	assert.Len(t, postResponses, 2, "POST endpoint should have exactly 2 responses (success and error)")

	// Verify GET response
	getResponse := getResponses[0]
	assert.Equal(t, 200, getResponse.StatusCode, "GET response should have status 200")
	assert.True(t, getResponse.Enabled, "GET response should be enabled")
	assert.Contains(t, getResponse.Body, "Hello from GET /testing!", "GET response should contain expected message")
	assert.Contains(t, getResponse.Body, "users", "GET response should contain users data")
	assert.Contains(t, getResponse.Headers, "Content-Type", "GET response should have content-type header")

	// Verify POST responses
	var successResponse, errorResponse *database.MockResponse
	for i := range postResponses {
		if postResponses[i].StatusCode == 201 {
			successResponse = &postResponses[i]
		}
		if postResponses[i].StatusCode == 400 {
			errorResponse = &postResponses[i]
		}
	}

	// Verify POST success response
	require.NotNil(t, successResponse, "POST success response should exist")
	assert.Equal(t, 201, successResponse.StatusCode, "POST success response should have status 201")
	assert.True(t, successResponse.Enabled, "POST success response should be enabled")
	assert.Contains(t, successResponse.Body, "Resource created successfully!", "POST success response should contain expected message")
	assert.Equal(t, 500, successResponse.DelayMS, "POST success response should have 500ms delay")

	// Verify POST error response
	require.NotNil(t, errorResponse, "POST error response should exist")
	assert.Equal(t, 400, errorResponse.StatusCode, "POST error response should have status 400")
	assert.True(t, errorResponse.Enabled, "POST error response should be enabled")
	assert.Contains(t, errorResponse.Body, "Validation failed", "POST error response should contain expected message")
	assert.Equal(t, 200, errorResponse.DelayMS, "POST error response should have 200ms delay")

	t.Logf("Verified responses - GET: %d, POST success: %d, POST error: %d",
		len(getResponses), 1, 1)

	// Step 6: Verify user can access the workspace (test workspace listing)
	userWorkspaces, err := workspaceService.GetUserWorkspaces(ctx, testUser.ID)
	require.NoError(t, err, "Should be able to get user workspaces")
	assert.Len(t, userWorkspaces, 1, "User should have exactly 1 workspace")
	assert.Equal(t, workspace.ID, userWorkspaces[0].ID, "Workspace ID should match")
	assert.Equal(t, workspaceName, userWorkspaces[0].Name, "Workspace name should match")

	// Step 7: Verify user workspaces with roles
	userWorkspacesWithRoles, err := workspaceService.GetUserWorkspacesWithRoles(ctx, testUser.ID)
	require.NoError(t, err, "Should be able to get user workspaces with roles")
	assert.Len(t, userWorkspacesWithRoles, 1, "User should have exactly 1 workspace with role")
	assert.Equal(t, workspace.ID, userWorkspacesWithRoles[0].ID, "Workspace ID should match")
	assert.Equal(t, workspaceName, userWorkspacesWithRoles[0].Name, "Workspace name should match")
	assert.Equal(t, "owner", userWorkspacesWithRoles[0].Role, "User should be owner")

	t.Logf("✅ All tests passed! Successfully created demo workspace with complete setup")
}

func TestCreateDemoWorkspace_UserLimitExceeded(t *testing.T) {
	// Setup test database
	database.SetupTestEnvironment(t)

	// Get database connection
	db := database.DB

	// Create workspace service with test repository
	workspaceService := &WorkspaceService{
		repo: &testWorkspaceRepo{db: db},
	}

	// Setup context with logger
	logger := zerolog.New(zerolog.NewTestWriter(t))
	ctx := logger.WithContext(context.Background())

	// Create a user with workspace limit of 1
	testUser := &database.User{
		Name:                 "Limited User",
		Email:                "limited@example.com",
		Password:             "test-password",
		IsActive:             true,
		MaxProjectsWorkspace: nil,
		MaxWorkspaces:        intPtr(1), // Limit to 1 workspace
	}

	err := db.Create(testUser).Error
	require.NoError(t, err, "Failed to create test user")

	// Create first workspace (should succeed)
	_, _, err = workspaceService.CreateDemoWorkspace(ctx, testUser.ID, testUser.Name, "First Workspace")
	require.NoError(t, err, "First workspace creation should succeed")

	// Try to create second workspace (should fail due to limit)
	_, _, err = workspaceService.CreateDemoWorkspace(ctx, testUser.ID, testUser.Name, "Second Workspace")
	require.Error(t, err, "Second workspace creation should fail due to limit")
	assert.Contains(t, err.Error(), "workspace limit exceeded", "Error should mention workspace limit")

	t.Logf("✅ Workspace limit test passed - correctly prevented exceeding limit")
}

func TestCreateDemoWorkspace_UniqueAliases(t *testing.T) {
	// Setup test database
	database.SetupTestEnvironment(t)

	// Get database connection
	db := database.DB

	// Create workspace service with test repository
	workspaceService := &WorkspaceService{
		repo: &testWorkspaceRepo{db: db},
	}

	// Setup context with logger
	logger := zerolog.New(zerolog.NewTestWriter(t))
	ctx := logger.WithContext(context.Background())

	// Create test user
	testUser := &database.User{
		Name:                 "Alias Test User",
		Email:                "aliastest@example.com",
		Password:             "test-password",
		IsActive:             true,
		MaxProjectsWorkspace: nil,
		MaxWorkspaces:        nil,
	}

	err := db.Create(testUser).Error
	require.NoError(t, err, "Failed to create test user")

	// Create first workspace
	_, _, err = workspaceService.CreateDemoWorkspace(ctx, testUser.ID, testUser.Name, "Workspace 1")
	require.NoError(t, err, "First workspace creation should succeed")

	// Add small delay to ensure different timestamps
	time.Sleep(1 * time.Millisecond)

	// Create second workspace
	_, _, err = workspaceService.CreateDemoWorkspace(ctx, testUser.ID, testUser.Name, "Workspace 2")
	require.NoError(t, err, "Second workspace creation should succeed")
}

// Helper function to create int pointer
func intPtr(i int) *int {
	return &i
}
