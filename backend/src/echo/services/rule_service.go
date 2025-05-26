package services

import (
	"fmt"
	"log"

	"beo-echo/backend/src/database"
)

type ruleRepository interface {
	FindRulesByResponseID(responseID string) ([]database.MockRule, error)
	FindRuleByID(ruleID string) (*database.MockRule, error)
	CreateRule(rule *database.MockRule) error
	UpdateRule(rule *database.MockRule) error
	DeleteRule(ruleID string) error
	DeleteRulesByResponseID(responseID string) error
}

type responseRepository interface {
	ValidateResponseHierarchy(projectID string, endpointID string, responseID string) (bool, error)
}

// RuleService handles business logic for rule management
type RuleService struct {
	RuleRepo     ruleRepository
	ResponseRepo responseRepository
}

// NewRuleService creates a new rule service
func NewRuleService(ruleRepo ruleRepository, responseRepo responseRepository) *RuleService {
	return &RuleService{
		RuleRepo:     ruleRepo,
		ResponseRepo: responseRepo,
	}
}

// GetRules retrieves all rules for a response
func (s *RuleService) GetRules(responseID string) ([]database.MockRule, error) {
	// Get rules directly without additional validation since we validate in the handler
	return s.RuleRepo.FindRulesByResponseID(responseID)
}

// GetRule gets a specific rule by ID and validates it belongs to the specified response
func (s *RuleService) GetRule(ruleID string, responseID string) (*database.MockRule, error) {
	rule, err := s.RuleRepo.FindRuleByID(ruleID)
	if err != nil {
		return nil, fmt.Errorf("rule not found: %w", err)
	}

	// Verify rule belongs to the specified response
	if rule.ResponseID != responseID {
		return nil, fmt.Errorf("rule does not belong to the specified response")
	}

	return rule, nil
}

// ValidateResponseHierarchy checks if a response exists and belongs to the specified project/endpoint
func (s *RuleService) ValidateResponseHierarchy(projectID string, endpointID string, responseID string) (bool, error) {
	return s.ResponseRepo.ValidateResponseHierarchy(projectID, endpointID, responseID)
}

// ValidateRuleHierarchy checks if a rule exists and belongs to the specified response/endpoint/project
func (s *RuleService) ValidateRuleHierarchy(projectID string, endpointID string, responseID string, ruleID string) (bool, error) {
	// First validate response hierarchy
	isValid, err := s.ValidateResponseHierarchy(projectID, endpointID, responseID)
	if err != nil || !isValid {
		return false, fmt.Errorf("invalid response hierarchy: %w", err)
	}

	// Then validate rule belongs to the response
	rule, err := s.RuleRepo.FindRuleByID(ruleID)
	if err != nil {
		return false, fmt.Errorf("rule not found: %w", err)
	}

	return rule.ResponseID == responseID, nil
}

// CreateRule creates a new rule for a response
func (s *RuleService) CreateRule(rule *database.MockRule) (*database.MockRule, error) {
	// Validate rule data
	if rule.Type == "" {
		return nil, fmt.Errorf("rule type is required")
	}

	if rule.Key == "" {
		return nil, fmt.Errorf("rule key is required")
	}

	if rule.Operator == "" {
		return nil, fmt.Errorf("rule operator is required")
	}

	// Value can be empty (checking for absence of a header/query param)

	// Create rule
	err := s.RuleRepo.CreateRule(rule)
	if err != nil {
		log.Printf("Error creating rule: %v", err)
		return nil, fmt.Errorf("error creating rule: %w", err)
	}

	// Get created rule
	createdRule, err := s.RuleRepo.FindRuleByID(rule.ID)
	if err != nil {
		log.Printf("Error retrieving created rule: %v", err)
		return nil, fmt.Errorf("rule created but could not be retrieved: %w", err)
	}

	return createdRule, nil
}

// UpdateRule updates an existing rule
func (s *RuleService) UpdateRule(ruleID string, updates *database.MockRule) (*database.MockRule, error) {
	// Check if rule exists
	existingRule, err := s.RuleRepo.FindRuleByID(ruleID)
	if err != nil {
		return nil, fmt.Errorf("rule not found: %w", err)
	}

	// Apply updates
	if updates.Type != "" {
		existingRule.Type = updates.Type
	}

	if updates.Key != "" {
		existingRule.Key = updates.Key
	}

	if updates.Operator != "" {
		existingRule.Operator = updates.Operator
	}

	// Value can be updated to empty string intentionally
	existingRule.Value = updates.Value

	// Save updates
	err = s.RuleRepo.UpdateRule(existingRule)
	if err != nil {
		log.Printf("Error updating rule: %v", err)
		return nil, fmt.Errorf("error updating rule: %w", err)
	}

	return existingRule, nil
}

// DeleteRule deletes a rule
func (s *RuleService) DeleteRule(ruleID string) error {
	// Check if rule exists
	_, err := s.RuleRepo.FindRuleByID(ruleID)
	if err != nil {
		return fmt.Errorf("rule not found: %w", err)
	}

	// Delete rule
	err = s.RuleRepo.DeleteRule(ruleID)
	if err != nil {
		log.Printf("Error deleting rule: %v", err)
		return fmt.Errorf("error deleting rule: %w", err)
	}

	return nil
}

// DeleteRulesByResponse deletes all rules for a response
func (s *RuleService) DeleteRulesByResponse(responseID string) error {
	err := s.RuleRepo.DeleteRulesByResponseID(responseID)
	if err != nil {
		log.Printf("Error deleting rules for response %s: %v", responseID, err)
		return fmt.Errorf("error deleting rules: %w", err)
	}

	return nil
}
