package workspaces

import (
	"beo-echo/backend/src/database"
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

// CreateDemoWorkspace creates a complete demo workspace with sample project and endpoints
func (s *WorkspaceService) CreateDemoWorkspace(ctx context.Context, userID string, userName string, workspaceName string) (*database.Workspace, *database.Project, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("user_id", userID).
		Str("user_name", userName).
		Str("workspace_name", workspaceName).
		Msg("creating demo workspace with demo project")

	// Step 1: Create the workspace
	workspace := &database.Workspace{
		Name: workspaceName,
	}

	err := s.CreateWorkspace(ctx, workspace, userID)
	if err != nil {
		log.Error().
			Err(err).
			Str("user_id", userID).
			Msg("failed to create workspace for demo workspace")
		return nil, nil, fmt.Errorf("failed to create workspace: %w", err)
	}

	// Step 2: Create demo project with uuid to avoid duplicates
	demoAlias := uuid.New().String()

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
		Msg("demo workspace created successfully with demo project and endpoints")

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
		Body:       `{"message": "Resource created successfully!", "id": ` + fmt.Sprintf("%d", time.Now().UnixNano()) + `, "status": "created", "timestamp": "` + time.Now().Format(time.RFC3339) + `"}`,
		Headers:    `{"Content-Type": "application/json", "X-Demo-Header": "POST Success Response", "Location": "/testing/` + fmt.Sprintf("%d", time.Now().UnixNano()) + `"}`,
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
