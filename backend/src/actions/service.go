package actions

import (
	"beo-echo/backend/src/database"
	"context"
	"encoding/json"
	"errors"
	"regexp"
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
	repo ActionRepository
}

// NewActionService creates a new action service
func NewActionService(repo ActionRepository) *ActionService {
	return &ActionService{repo: repo}
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
		if err := s.validateReplaceTextConfig(action.Config); err != nil {
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
		if err := s.validateReplaceTextConfig(action.Config); err != nil {
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

// ReplaceTextConfig represents the configuration for replace_text action type
type ReplaceTextConfig struct {
	Target      string `json:"target"`       // request_body, response_body, request_header, response_header
	Pattern     string `json:"pattern"`      // String or regex pattern to find
	Replacement string `json:"replacement"`  // Replacement text
	UseRegex    bool   `json:"use_regex"`    // Whether to use regex matching
	HeaderKey   string `json:"header_key"`   // Header key (only for header targets)
}

// validateReplaceTextConfig validates the configuration for replace_text actions
func (s *ActionService) validateReplaceTextConfig(configJSON string) error {
	if configJSON == "" {
		return errors.New("config is required for replace_text action")
	}

	var config ReplaceTextConfig
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return errors.New("invalid config JSON: " + err.Error())
	}

	// Validate target
	validTargets := map[string]bool{
		"request_body":     true,
		"response_body":    true,
		"request_header":   true,
		"response_header":  true,
	}
	if !validTargets[config.Target] {
		return errors.New("invalid target: must be request_body, response_body, request_header, or response_header")
	}

	// Validate header key for header targets
	if (config.Target == "request_header" || config.Target == "response_header") && config.HeaderKey == "" {
		return errors.New("header_key is required for header targets")
	}

	// Validate pattern
	if config.Pattern == "" {
		return errors.New("pattern is required")
	}

	// Validate regex if used
	if config.UseRegex {
		if _, err := regexp.Compile(config.Pattern); err != nil {
			return errors.New("invalid regex pattern: " + err.Error())
		}
	}

	return nil
}

// RequestData represents the data that can be modified by actions
type RequestData struct {
	Method  string
	Path    string
	Headers map[string]string
	Body    string
}

// ResponseData represents the response data that can be modified by actions
type ResponseData struct {
	StatusCode int
	Headers    map[string]string
	Body       string
}

// ExecuteBeforeRequestActions executes all enabled actions that run before the request
func (s *ActionService) ExecuteBeforeRequestActions(ctx context.Context, projectID string, reqData *RequestData) error {
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
			if !s.matchFilters(action.Filters, reqData, nil) {
				log.Debug().Str("action_id", action.ID).Msg("action filters did not match, skipping")
				continue
			}
		}

		// Execute the action based on type
		if err := s.executeAction(ctx, &action, reqData, nil); err != nil {
			log.Error().Err(err).Str("action_id", action.ID).Msg("failed to execute action")
			// Continue with other actions even if one fails
			continue
		}

		log.Debug().Str("action_id", action.ID).Msg("action executed successfully")
	}

	return nil
}

// ExecuteAfterRequestActions executes all enabled actions that run after the request
func (s *ActionService) ExecuteAfterRequestActions(ctx context.Context, projectID string, reqData *RequestData, respData *ResponseData) error {
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
			if !s.matchFilters(action.Filters, reqData, respData) {
				log.Debug().Str("action_id", action.ID).Msg("action filters did not match, skipping")
				continue
			}
		}

		// Execute the action based on type
		if err := s.executeAction(ctx, &action, reqData, respData); err != nil {
			log.Error().Err(err).Str("action_id", action.ID).Msg("failed to execute action")
			// Continue with other actions even if one fails
			continue
		}

		log.Debug().Str("action_id", action.ID).Msg("action executed successfully")
	}

	return nil
}

// executeAction executes a single action
func (s *ActionService) executeAction(ctx context.Context, action *database.Action, reqData *RequestData, respData *ResponseData) error {
	switch action.Type {
	case database.ActionTypeReplaceText:
		return s.executeReplaceTextAction(action, reqData, respData)
	default:
		return errors.New("unsupported action type: " + string(action.Type))
	}
}

// executeReplaceTextAction executes a replace_text action
func (s *ActionService) executeReplaceTextAction(action *database.Action, reqData *RequestData, respData *ResponseData) error {
	var config ReplaceTextConfig
	if err := json.Unmarshal([]byte(action.Config), &config); err != nil {
		return err
	}

	// Perform replacement based on target
	switch config.Target {
	case "request_body":
		if reqData != nil {
			reqData.Body = s.replaceText(reqData.Body, config.Pattern, config.Replacement, config.UseRegex)
		}
	case "response_body":
		if respData != nil {
			respData.Body = s.replaceText(respData.Body, config.Pattern, config.Replacement, config.UseRegex)
		}
	case "request_header":
		if reqData != nil && reqData.Headers != nil {
			if value, ok := reqData.Headers[config.HeaderKey]; ok {
				reqData.Headers[config.HeaderKey] = s.replaceText(value, config.Pattern, config.Replacement, config.UseRegex)
			}
		}
	case "response_header":
		if respData != nil && respData.Headers != nil {
			if value, ok := respData.Headers[config.HeaderKey]; ok {
				respData.Headers[config.HeaderKey] = s.replaceText(value, config.Pattern, config.Replacement, config.UseRegex)
			}
		}
	}

	return nil
}

// replaceText performs the actual text replacement
func (s *ActionService) replaceText(text, pattern, replacement string, useRegex bool) string {
	if useRegex {
		re, err := regexp.Compile(pattern)
		if err != nil {
			return text // Return original on error
		}
		return re.ReplaceAllString(text, replacement)
	}
	return strings.ReplaceAll(text, pattern, replacement)
}

// matchFilters checks if any filter matches the request/response data (OR logic)
func (s *ActionService) matchFilters(filters []database.ActionFilter, reqData *RequestData, respData *ResponseData) bool {
	// OR logic: if any filter matches, return true
	for _, filter := range filters {
		if s.matchFilter(filter, reqData, respData) {
			return true
		}
	}
	return false
}

// matchFilter checks if a single filter matches
func (s *ActionService) matchFilter(filter database.ActionFilter, reqData *RequestData, respData *ResponseData) bool {
	var value string

	// Get the value to compare based on filter type
	switch filter.Type {
	case "method":
		if reqData != nil {
			value = reqData.Method
		}
	case "path":
		if reqData != nil {
			value = reqData.Path
		}
	case "header":
		if reqData != nil && reqData.Headers != nil {
			value = reqData.Headers[filter.Key]
		}
	case "status_code":
		if respData != nil {
			value = string(rune(respData.StatusCode))
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
