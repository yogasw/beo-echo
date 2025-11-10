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
	if action.Type != database.ActionTypeReplaceText && action.Type != database.ActionTypeRunJavascript {
		return errors.New("unsupported action type")
	}

	// Validate execution point
	if action.ExecutionPoint != database.ExecutionPointBeforeRequest &&
		action.ExecutionPoint != database.ExecutionPointAfterRequest {
		return errors.New("invalid execution point")
	}

	// Validate config based on action type
	if action.Type == database.ActionTypeReplaceText {
		if err := s.modules.ValidateReplaceTextConfig(action.Config); err != nil {
			log.Error().Err(err).Msg("invalid replace_text config")
			return err
		}
	} else if action.Type == database.ActionTypeRunJavascript {
		if err := s.modules.ValidateRunJavascriptConfig(action.Config); err != nil {
			log.Error().Err(err).Msg("invalid run_javascript config")
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
	if action.Config != existing.Config {
		if action.Type == database.ActionTypeReplaceText {
			if err := s.modules.ValidateReplaceTextConfig(action.Config); err != nil {
				log.Error().Err(err).Msg("invalid replace_text config")
				return err
			}
		} else if action.Type == database.ActionTypeRunJavascript {
			if err := s.modules.ValidateRunJavascriptConfig(action.Config); err != nil {
				log.Error().Err(err).Msg("invalid run_javascript config")
				return err
			}
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

// UpdateActionPriority updates the priority of an action and reorders other actions
// Optimized to only update affected actions within the same execution point
// Actions can only be reordered within their execution point group (before_request or after_request)
func (s *ActionService) UpdateActionPriority(ctx context.Context, actionID string, newPriority int) error {
	log := zerolog.Ctx(ctx)

	// Get the action to update
	action, err := s.repo.GetActionByID(ctx, actionID)
	if err != nil {
		log.Error().Err(err).Str("action_id", actionID).Msg("action not found")
		return errors.New("action not found")
	}

	// Get all actions for the same project
	allActions, err := s.repo.GetActionsByProjectID(ctx, action.ProjectID)
	if err != nil {
		log.Error().Err(err).Str("project_id", action.ProjectID).Msg("failed to get project actions")
		return err
	}

	// Filter actions to only include those with the same execution_point
	sameExecutionPointActions := make([]database.Action, 0)
	for _, a := range allActions {
		if a.ExecutionPoint == action.ExecutionPoint {
			sameExecutionPointActions = append(sameExecutionPointActions, a)
		}
	}

	// Validate new priority is within bounds of the same execution point group (1-based)
	if newPriority < 1 || newPriority > len(sameExecutionPointActions) {
		return errors.New("invalid priority: must be between 1 and " + strconv.Itoa(len(sameExecutionPointActions)) +
			" for " + string(action.ExecutionPoint) + " actions")
	}

	oldPriority := action.Priority

	// If priority hasn't changed, nothing to do
	if oldPriority == newPriority {
		log.Debug().
			Str("action_id", actionID).
			Int("priority", newPriority).
			Msg("priority unchanged, skipping update")
		return nil
	}

	// Collect only affected actions that need updating (within same execution point)
	affectedActions := make([]*database.Action, 0)
	affectedCount := 0

	// Reorder priorities - only update affected actions in the same execution point
	if oldPriority < newPriority {
		// Moving down: shift actions up (decrease priority)
		for i := range sameExecutionPointActions {
			if sameExecutionPointActions[i].ID == actionID {
				sameExecutionPointActions[i].Priority = newPriority
				affectedActions = append(affectedActions, &sameExecutionPointActions[i])
				affectedCount++
			} else if sameExecutionPointActions[i].Priority > oldPriority && sameExecutionPointActions[i].Priority <= newPriority {
				sameExecutionPointActions[i].Priority--
				affectedActions = append(affectedActions, &sameExecutionPointActions[i])
				affectedCount++
			}
		}
	} else {
		// Moving up: shift actions down (increase priority)
		for i := range sameExecutionPointActions {
			if sameExecutionPointActions[i].ID == actionID {
				sameExecutionPointActions[i].Priority = newPriority
				affectedActions = append(affectedActions, &sameExecutionPointActions[i])
				affectedCount++
			} else if sameExecutionPointActions[i].Priority >= newPriority && sameExecutionPointActions[i].Priority < oldPriority {
				sameExecutionPointActions[i].Priority++
				affectedActions = append(affectedActions, &sameExecutionPointActions[i])
				affectedCount++
			}
		}
	}

	// Update only affected actions
	for _, a := range affectedActions {
		if err := s.repo.UpdateAction(ctx, a); err != nil {
			log.Error().Err(err).Str("action_id", a.ID).Msg("failed to update action priority")
			return errors.New("failed to update priorities")
		}
	}

	log.Info().
		Str("action_id", actionID).
		Str("execution_point", string(action.ExecutionPoint)).
		Int("old_priority", oldPriority).
		Int("new_priority", newPriority).
		Int("affected_actions", affectedCount).
		Int("total_in_group", len(sameExecutionPointActions)).
		Msg("action priority updated successfully")

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
	case database.ActionTypeRunJavascript:
		return s.modules.ExecuteRunJavascriptAction(action, req, resp)
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
