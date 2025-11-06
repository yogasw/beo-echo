package actions

import (
	"beo-echo/backend/src/actions/modules"
	"beo-echo/backend/src/database"
	"context"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/rs/zerolog"
)

// ActionRepository defines data access requirements for action operations
type ActionRepository interface {
	// Action CRUD operations
	CreateAction(ctx context.Context, action *database.Action) error
	GetActionByID(ctx context.Context, id string) (*database.Action, error)
	GetActionsByProjectID(ctx context.Context, projectID string) ([]database.Action, error)
	UpdateAction(ctx context.Context, action *database.Action) error
	DeleteAction(ctx context.Context, id string) error

	// Filter operations
	CreateActionFilter(ctx context.Context, filter *database.ActionFilter) error
	GetActionFilters(ctx context.Context, actionID string) ([]database.ActionFilter, error)
	DeleteActionFilters(ctx context.Context, actionID string) error

	// Execution queries
	GetEnabledActionsByProjectAndPoint(ctx context.Context, projectID string, executionPoint database.ExecutionPoint) ([]database.Action, error)
}

// ActionService implements action business operations
type ActionService struct {
	repo    ActionRepository
	modules *modules.ModulesAction
}

// NewActionService creates a new action service
func NewActionService(repo ActionRepository, modules *modules.ModulesAction) *ActionService {
	return &ActionService{repo: repo, modules: modules}
}

// CreateAction creates a new action for a project
func (s *ActionService) CreateAction(ctx context.Context, action *database.Action) error {
	log := zerolog.Ctx(ctx)

	// Validate required fields
	if action.ProjectID == "" {
		return errors.New("project_id is required")
	}
	if action.Name == "" {
		return errors.New("action name is required")
	}
	if action.Type == "" {
		return errors.New("action type is required")
	}
	if action.ExecutionPoint == "" {
		action.ExecutionPoint = database.ExecutionPointAfterRequest // Default
	}

	// Validate action type
	if action.Type != database.ActionTypeReplaceText {
		return errors.New("unsupported action type")
	}

	// Validate execution point
	if action.ExecutionPoint != database.ExecutionPointBeforeRequest &&
		action.ExecutionPoint != database.ExecutionPointAfterRequest {
		return errors.New("invalid execution point")
	}

	// Validate config for replace_text action
	if action.Type == database.ActionTypeReplaceText {
		if err := s.modules.ValidateReplaceTextConfig(action.Config); err != nil {
			log.Error().Err(err).Msg("invalid replace_text config")
			return err
		}
	}

	// Create the action
	if err := s.repo.CreateAction(ctx, action); err != nil {
		log.Error().Err(err).Str("project_id", action.ProjectID).Msg("failed to create action")
		return err
	}

	log.Info().Str("action_id", action.ID).Str("project_id", action.ProjectID).Msg("action created successfully")
	return nil
}

// GetAction retrieves an action by ID
func (s *ActionService) GetAction(ctx context.Context, id string) (*database.Action, error) {
	return s.repo.GetActionByID(ctx, id)
}

// GetActionsByProject retrieves all actions for a project
func (s *ActionService) GetActionsByProject(ctx context.Context, projectID string) ([]database.Action, error) {
	return s.repo.GetActionsByProjectID(ctx, projectID)
}

// UpdateAction updates an existing action
func (s *ActionService) UpdateAction(ctx context.Context, action *database.Action) error {
	log := zerolog.Ctx(ctx)

	// Get existing action
	existing, err := s.repo.GetActionByID(ctx, action.ID)
	if err != nil {
		return err
	}

	// Validate type hasn't changed
	if action.Type != existing.Type {
		return errors.New("cannot change action type")
	}

	// Validate config if it changed
	if action.Type == database.ActionTypeReplaceText && action.Config != existing.Config {
		if err := s.modules.ValidateReplaceTextConfig(action.Config); err != nil {
			log.Error().Err(err).Msg("invalid replace_text config")
			return err
		}
	}

	// Validate execution point
	if action.ExecutionPoint != database.ExecutionPointBeforeRequest &&
		action.ExecutionPoint != database.ExecutionPointAfterRequest {
		return errors.New("invalid execution point")
	}

	// Update the action
	if err := s.repo.UpdateAction(ctx, action); err != nil {
		log.Error().Err(err).Str("action_id", action.ID).Msg("failed to update action")
		return err
	}

	log.Info().Str("action_id", action.ID).Msg("action updated successfully")
	return nil
}

// DeleteAction deletes an action
func (s *ActionService) DeleteAction(ctx context.Context, id string) error {
	log := zerolog.Ctx(ctx)

	if err := s.repo.DeleteAction(ctx, id); err != nil {
		log.Error().Err(err).Str("action_id", id).Msg("failed to delete action")
		return err
	}

	log.Info().Str("action_id", id).Msg("action deleted successfully")
	return nil
}

// ExecuteBeforeRequestActions executes all enabled actions that run before the request
func (s *ActionService) ExecuteBeforeRequestActions(ctx context.Context, projectID string, req *http.Request) error {
	log := zerolog.Ctx(ctx)

	// Get enabled actions for this project and execution point
	actions, err := s.repo.GetEnabledActionsByProjectAndPoint(ctx, projectID, database.ExecutionPointBeforeRequest)
	if err != nil {
		log.Error().Err(err).Msg("failed to get before_request actions")
		return err
	}

	// Execute actions in priority order (already sorted by repository)
	for _, action := range actions {
		// Check if filters match
		if len(action.Filters) > 0 {
			if !s.matchFilters(action.Filters, req, nil) {
				log.Debug().Str("action_id", action.ID).Msg("action filters did not match, skipping")
				continue
			}
		}

		// Execute the action based on type
		if err := s.executeAction(ctx, &action, req, nil); err != nil {
			log.Error().Err(err).Str("action_id", action.ID).Msg("failed to execute action")
			// Continue with other actions even if one fails
			continue
		}
	}

	return nil
}

// ExecuteAfterRequestActions executes all enabled actions that run after the request
func (s *ActionService) ExecuteAfterRequestActions(ctx context.Context, projectID string, req *http.Request, resp *http.Response) error {
	log := zerolog.Ctx(ctx)

	// Get enabled actions for this project and execution point
	actions, err := s.repo.GetEnabledActionsByProjectAndPoint(ctx, projectID, database.ExecutionPointAfterRequest)
	if err != nil {
		log.Error().Err(err).Msg("failed to get after_request actions")
		return err
	}

	// Execute actions in priority order (already sorted by repository)
	for _, action := range actions {
		// Check if filters match
		if len(action.Filters) > 0 {
			if !s.matchFilters(action.Filters, req, resp) {
				log.Debug().Str("action_id", action.ID).Msg("action filters did not match, skipping")
				continue
			}
		}

		// Execute the action based on type
		if err := s.executeAction(ctx, &action, req, resp); err != nil {
			log.Error().Err(err).Str("action_id", action.ID).Msg("failed to execute action")
			// Continue with other actions even if one fails
			continue
		}
	}

	return nil
}

// executeAction executes a single action
func (s *ActionService) executeAction(ctx context.Context, action *database.Action, req *http.Request, resp *http.Response) error {
	switch action.Type {
	case database.ActionTypeReplaceText:
		return s.modules.ExecuteReplaceTextAction(action, req, resp)
	default:
		return errors.New("unsupported action type: " + string(action.Type))
	}
}

// matchFilters checks if any filter matches the request/response data (OR logic)
func (s *ActionService) matchFilters(filters []database.ActionFilter, req *http.Request, resp *http.Response) bool {
	// OR logic: if any filter matches, return true
	for _, filter := range filters {
		if s.matchFilter(filter, req, resp) {
			return true
		}
	}
	return false
}

// matchFilter checks if a single filter matches
func (s *ActionService) matchFilter(filter database.ActionFilter, req *http.Request, resp *http.Response) bool {
	var value string

	// Get the value to compare based on filter type
	switch filter.Type {
	case "method":
		if req != nil {
			value = req.Method
		}
	case "path":
		if req != nil {
			value = req.URL.Path
		}
	case "header":
		if req != nil && req.Header != nil {
			value = req.Header.Get(filter.Key)
		}
	case "status_code":
		if resp != nil {
			value = strconv.Itoa(resp.StatusCode)
		}
	default:
		return false
	}

	// Compare based on operator
	switch filter.Operator {
	case "equals":
		return value == filter.Value
	case "contains":
		return strings.Contains(value, filter.Value)
	case "starts_with":
		return strings.HasPrefix(value, filter.Value)
	case "ends_with":
		return strings.HasSuffix(value, filter.Value)
	case "regex":
		re, err := regexp.Compile(filter.Value)
		if err != nil {
			return false
		}
		return re.MatchString(value)
	default:
		return false
	}
}
