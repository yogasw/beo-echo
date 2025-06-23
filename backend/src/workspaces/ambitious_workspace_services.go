package workspaces

import (
	"beo-echo/backend/src/database"
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog"
)

// CreateAmbitiousWorkspace creates a complete demo workspace with sample project and endpoints
func (s *WorkspaceService) CreateAmbitiousWorkspace(ctx context.Context, userID string, userName string, workspaceName string) (*database.Workspace, *database.Project, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("user_id", userID).
		Str("user_name", userName).
		Str("workspace_name", workspaceName).
		Msg("creating ambitious workspace with demo project")

	// Step 1: Create the workspace
	workspace := &database.Workspace{
		Name: workspaceName,
	}

	err := s.CreateWorkspace(ctx, workspace, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to create workspace for ambitious workspace")
		return nil, nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Step 2: Create demo project with unix timestamp to avoid duplicates
	unixTime := time.Now().Unix()
	demoAlias := fmt.Sprintf("demo_%d", unixTime)

	demoProject := &database.Project{
		Name:        "Demo Project",
		WorkspaceID: workspace.ID,
		Mode:        database.ModeMock,
		Status:      "running",
		Alias:       demoAlias,
	}

	err = s.repo.CreateProject(ctx, demoProject)
	if err != nil {
		log.Error().
			Err(err).
			Str("workspace_id", workspace.ID).
			Str("demo_alias", demoAlias).
			Msg("failed to create demo project")
		return workspace, nil, fmt.Errorf("failed to create demo project: %w", err)
	}

	log.Info().
		Str("project_id", demoProject.ID).
		Str("project_alias", demoAlias).
		Msg("demo project created successfully")

	// Step 3: Create /testing endpoints (GET and POST)
	err = s.createDemoEndpoints(ctx, demoProject.ID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", demoProject.ID).
			Msg("failed to create demo endpoints")
		return workspace, demoProject, fmt.Errorf("failed to create demo endpoints: %w", err)
	}

	log.Info().
		Str("workspace_id", workspace.ID).
		Str("project_id", demoProject.ID).
		Str("project_alias", demoAlias).
		Msg("ambitious workspace created successfully with demo project and endpoints")

	return workspace, demoProject, nil
}

// createDemoEndpoints creates sample /testing endpoints with demo responses
func (s *WorkspaceService) createDemoEndpoints(ctx context.Context, projectID string) error {
	log := zerolog.Ctx(ctx)

	// Create GET /testing endpoint
	getEndpoint := &database.MockEndpoint{
		ProjectID:    projectID,
		Method:       "GET",
		Path:         "/testing",
		Enabled:      true,
		ResponseMode: "random",
	}

	err := s.repo.CreateEndpoint(ctx, getEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create GET endpoint: %w", err)
	}

	// Create GET response with sample data
	getResponse := &database.MockResponse{
		EndpointID: getEndpoint.ID,
		StatusCode: 200,
		Body:       `{"message": "Hello from GET /testing!", "timestamp": "` + time.Now().Format(time.RFC3339) + `", "data": {"users": [{"id": 1, "name": "John Doe", "email": "john@example.com"}, {"id": 2, "name": "Jane Smith", "email": "jane@example.com"}], "total": 2}}`,
		Headers:    `{"Content-Type": "application/json", "X-Demo-Header": "GET Response"}`,
		Priority:   1,
		DelayMS:    0,
		Stream:     false,
		Note:       "Sample GET response with user data",
		Enabled:    true,
	}

	err = s.repo.CreateResponse(ctx, getResponse)
	if err != nil {
		return fmt.Errorf("failed to create GET response: %w", err)
	}

	// Create POST /testing endpoint
	postEndpoint := &database.MockEndpoint{
		ProjectID:    projectID,
		Method:       "POST",
		Path:         "/testing",
		Enabled:      true,
		ResponseMode: "random",
	}

	err = s.repo.CreateEndpoint(ctx, postEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create POST endpoint: %w", err)
	}

	// Create POST success response
	postSuccessResponse := &database.MockResponse{
		EndpointID: postEndpoint.ID,
		StatusCode: 201,
		Body:       `{"message": "Resource created successfully!", "id": ` + fmt.Sprintf("%d", time.Now().Unix()) + `, "status": "created", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`,
		Headers:    `{"Content-Type": "application/json", "X-Demo-Header": "POST Success Response", "Location": "/testing/` + fmt.Sprintf("%d", time.Now().Unix()) + `"}`,
		Priority:   1,
		DelayMS:    500, // Add small delay to simulate processing
		Stream:     false,
		Note:       "Sample POST success response",
		Enabled:    true,
	}

	err = s.repo.CreateResponse(ctx, postSuccessResponse)
	if err != nil {
		return fmt.Errorf("failed to create POST success response: %w", err)
	}

	// Create POST validation error response
	postErrorResponse := &database.MockResponse{
		EndpointID: postEndpoint.ID,
		StatusCode: 400,
		Body:       `{"error": "Validation failed", "message": "Invalid request data", "details": [{"field": "name", "message": "Name is required"}, {"field": "email", "message": "Invalid email format"}], "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`,
		Headers:    `{"Content-Type": "application/json", "X-Demo-Header": "POST Error Response"}`,
		Priority:   2,
		DelayMS:    200,
		Stream:     false,
		Note:       "Sample POST validation error response",
		Enabled:    true,
	}

	err = s.repo.CreateResponse(ctx, postErrorResponse)
	if err != nil {
		return fmt.Errorf("failed to create POST error response: %w", err)
	}

	log.Info().
		Str("project_id", projectID).
		Str("get_endpoint_id", getEndpoint.ID).
		Str("post_endpoint_id", postEndpoint.ID).
		Msg("demo endpoints created successfully with sample responses")

	return nil
}

// CreateAmbitiousWorkspaceWithVariations creates multiple variations of demo workspaces
// This can be used for testing or demo purposes
func (s *WorkspaceService) CreateAmbitiousWorkspaceWithVariations(ctx context.Context, userID string, userName string) ([]*database.Workspace, []*database.Project, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("user_id", userID).
		Str("user_name", userName).
		Msg("creating multiple ambitious workspaces with demo variations")

	var workspaces []*database.Workspace
	var projects []*database.Project

	// Variation 1: Basic API Demo
	workspace1, project1, err := s.CreateAmbitiousWorkspace(ctx, userID, userName, "Basic API Demo")
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create basic API demo workspace: %w", err)
	}
	workspaces = append(workspaces, workspace1)
	projects = append(projects, project1)

	// Variation 2: E-commerce API Demo
	workspace2, project2, err := s.createEcommerceWorkspace(ctx, userID, userName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create e-commerce demo workspace")
		// Continue with other variations even if one fails
	} else {
		workspaces = append(workspaces, workspace2)
		projects = append(projects, project2)
	}

	// Variation 3: User Management API Demo
	workspace3, project3, err := s.createUserManagementWorkspace(ctx, userID, userName)
	if err != nil {
		log.Error().Err(err).Msg("failed to create user management demo workspace")
		// Continue with other variations even if one fails
	} else {
		workspaces = append(workspaces, workspace3)
		projects = append(projects, project3)
	}

	log.Info().
		Str("user_id", userID).
		Int("workspaces_created", len(workspaces)).
		Int("projects_created", len(projects)).
		Msg("ambitious workspaces with variations created successfully")

	return workspaces, projects, nil
}

// createEcommerceWorkspace creates a workspace with e-commerce API endpoints
func (s *WorkspaceService) createEcommerceWorkspace(ctx context.Context, userID string, userName string) (*database.Workspace, *database.Project, error) {
	// Create workspace
	workspace := &database.Workspace{
		Name: "E-commerce API Demo",
	}

	err := s.CreateWorkspace(ctx, workspace, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create e-commerce workspace: %w", err)
	}

	// Create project
	unixTime := time.Now().Unix()
	ecommerceAlias := fmt.Sprintf("ecommerce_%d", unixTime)

	project := &database.Project{
		Name:        "E-commerce API",
		WorkspaceID: workspace.ID,
		Mode:        database.ModeMock,
		Status:      "running",
		Alias:       ecommerceAlias,
	}

	err = s.repo.CreateProject(ctx, project)
	if err != nil {
		return workspace, nil, fmt.Errorf("failed to create e-commerce project: %w", err)
	}

	// Create e-commerce endpoints
	err = s.createEcommerceEndpoints(ctx, project.ID)
	if err != nil {
		return workspace, project, fmt.Errorf("failed to create e-commerce endpoints: %w", err)
	}

	return workspace, project, nil
}

// createUserManagementWorkspace creates a workspace with user management API endpoints
func (s *WorkspaceService) createUserManagementWorkspace(ctx context.Context, userID string, userName string) (*database.Workspace, *database.Project, error) {
	// Create workspace
	workspace := &database.Workspace{
		Name: "User Management API Demo",
	}

	err := s.CreateWorkspace(ctx, workspace, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create user management workspace: %w", err)
	}

	// Create project
	unixTime := time.Now().Unix()
	userMgmtAlias := fmt.Sprintf("usermgmt_%d", unixTime)

	project := &database.Project{
		Name:        "User Management API",
		WorkspaceID: workspace.ID,
		Mode:        database.ModeMock,
		Status:      "running",
		Alias:       userMgmtAlias,
	}

	err = s.repo.CreateProject(ctx, project)
	if err != nil {
		return workspace, nil, fmt.Errorf("failed to create user management project: %w", err)
	}

	// Create user management endpoints
	err = s.createUserManagementEndpoints(ctx, project.ID)
	if err != nil {
		return workspace, project, fmt.Errorf("failed to create user management endpoints: %w", err)
	}

	return workspace, project, nil
}

// createEcommerceEndpoints creates sample e-commerce API endpoints
func (s *WorkspaceService) createEcommerceEndpoints(ctx context.Context, projectID string) error {
	// GET /products endpoint
	productsEndpoint := &database.MockEndpoint{
		ProjectID:    projectID,
		Method:       "GET",
		Path:         "/products",
		Enabled:      true,
		ResponseMode: "random",
	}

	err := s.repo.CreateEndpoint(ctx, productsEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create products endpoint: %w", err)
	}

	// Products response
	productsResponse := &database.MockResponse{
		EndpointID: productsEndpoint.ID,
		StatusCode: 200,
		Body:       `{"products": [{"id": 1, "name": "Laptop", "price": 999.99, "category": "Electronics", "stock": 50}, {"id": 2, "name": "Coffee Mug", "price": 12.99, "category": "Home", "stock": 100}], "total": 2, "page": 1, "per_page": 10}`,
		Headers:    `{"Content-Type": "application/json", "X-API-Version": "1.0"}`,
		Priority:   1,
		DelayMS:    300,
		Stream:     false,
		Note:       "Sample products listing",
		Enabled:    true,
	}

	err = s.repo.CreateResponse(ctx, productsResponse)
	if err != nil {
		return fmt.Errorf("failed to create products response: %w", err)
	}

	// POST /orders endpoint
	ordersEndpoint := &database.MockEndpoint{
		ProjectID:    projectID,
		Method:       "POST",
		Path:         "/orders",
		Enabled:      true,
		ResponseMode: "random",
	}

	err = s.repo.CreateEndpoint(ctx, ordersEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create orders endpoint: %w", err)
	}

	// Orders response
	ordersResponse := &database.MockResponse{
		EndpointID: ordersEndpoint.ID,
		StatusCode: 201,
		Body:       `{"order_id": "` + fmt.Sprintf("ORD_%d", time.Now().Unix()) + `", "status": "pending", "total": 1012.98, "items": [{"product_id": 1, "quantity": 1, "price": 999.99}, {"product_id": 2, "quantity": 1, "price": 12.99}], "created_at": "` + time.Now().Format(time.RFC3339) + `"}`,
		Headers:    `{"Content-Type": "application/json", "X-API-Version": "1.0"}`,
		Priority:   1,
		DelayMS:    800,
		Stream:     false,
		Note:       "Sample order creation",
		Enabled:    true,
	}

	return s.repo.CreateResponse(ctx, ordersResponse)
}

// createUserManagementEndpoints creates sample user management API endpoints
func (s *WorkspaceService) createUserManagementEndpoints(ctx context.Context, projectID string) error {
	// GET /users endpoint
	usersEndpoint := &database.MockEndpoint{
		ProjectID:    projectID,
		Method:       "GET",
		Path:         "/users",
		Enabled:      true,
		ResponseMode: "random",
	}

	err := s.repo.CreateEndpoint(ctx, usersEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create users endpoint: %w", err)
	}

	// Users response
	usersResponse := &database.MockResponse{
		EndpointID: usersEndpoint.ID,
		StatusCode: 200,
		Body:       `{"users": [{"id": 1, "username": "admin", "email": "admin@example.com", "role": "administrator", "active": true, "created_at": "2024-01-01T00:00:00Z"}, {"id": 2, "username": "user1", "email": "user1@example.com", "role": "user", "active": true, "created_at": "2024-01-15T00:00:00Z"}], "total": 2, "page": 1, "per_page": 10}`,
		Headers:    `{"Content-Type": "application/json", "X-API-Version": "1.0"}`,
		Priority:   1,
		DelayMS:    200,
		Stream:     false,
		Note:       "Sample users listing",
		Enabled:    true,
	}

	err = s.repo.CreateResponse(ctx, usersResponse)
	if err != nil {
		return fmt.Errorf("failed to create users response: %w", err)
	}

	// POST /auth/login endpoint
	loginEndpoint := &database.MockEndpoint{
		ProjectID:    projectID,
		Method:       "POST",
		Path:         "/auth/login",
		Enabled:      true,
		ResponseMode: "random",
	}

	err = s.repo.CreateEndpoint(ctx, loginEndpoint)
	if err != nil {
		return fmt.Errorf("failed to create login endpoint: %w", err)
	}

	// Login success response
	loginResponse := &database.MockResponse{
		EndpointID: loginEndpoint.ID,
		StatusCode: 200,
		Body:       `{"access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.sample.token", "token_type": "Bearer", "expires_in": 3600, "user": {"id": 1, "username": "admin", "email": "admin@example.com", "role": "administrator"}}`,
		Headers:    `{"Content-Type": "application/json", "X-API-Version": "1.0"}`,
		Priority:   1,
		DelayMS:    1000,
		Stream:     false,
		Note:       "Sample login success",
		Enabled:    true,
	}

	return s.repo.CreateResponse(ctx, loginResponse)
}
