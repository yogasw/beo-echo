package repositories

import (
	"fmt"

	"beo-echo/backend/src/database"

	"gorm.io/gorm"
)

// RuleRepository handles CRUD operations for rules
type ruleRepository struct {
	db *gorm.DB
}

// NewRuleRepository creates a new rule repository that implements the required interface
func NewRuleRepository(db *gorm.DB) ruleRepository {
	return ruleRepository{
		db: db,
	}
}

// FindRulesByResponseID gets all rules for a response
func (r *ruleRepository) FindRulesByResponseID(responseID string) ([]database.MockRule, error) {
	var rules []database.MockRule
	result := r.db.Where("response_id = ?", responseID).Find(&rules)
	if result.Error != nil {
		return nil, result.Error
	}
	return rules, nil
}

// FindRuleByID gets a rule by ID
func (r *ruleRepository) FindRuleByID(ruleID string) (*database.MockRule, error) {
	var rule database.MockRule
	result := r.db.Where("id = ?", ruleID).First(&rule)
	if result.Error != nil {
		return nil, result.Error
	}
	return &rule, nil
}

// CreateRule creates a new rule
func (r *ruleRepository) CreateRule(rule *database.MockRule) error {
	// Check if response exists
	var response database.MockResponse
	if result := r.db.Where("id = ?", rule.ResponseID).First(&response); result.Error != nil {
		return fmt.Errorf("response not found: %w", result.Error)
	}

	// Validate rule type
	switch rule.Type {
	case "header", "query", "body":
		// Valid types
	default:
		return fmt.Errorf("invalid rule type: %s, must be header, query, or body", rule.Type)
	}

	// Validate operator
	switch rule.Operator {
	case "equals", "contains", "regex":
		// Valid operators
	default:
		return fmt.Errorf("invalid operator: %s, must be equals, contains, or regex", rule.Operator)
	}

	// Create rule
	return r.db.Create(rule).Error
}

// UpdateRule updates an existing rule
func (r *ruleRepository) UpdateRule(rule *database.MockRule) error {
	return r.db.Save(rule).Error
}

// DeleteRule deletes a rule by ID
func (r *ruleRepository) DeleteRule(ruleID string) error {
	return r.db.Delete(&database.MockRule{}, "id = ?", ruleID).Error
}

// DeleteRulesByResponseID deletes all rules for a response
func (r *ruleRepository) DeleteRulesByResponseID(responseID string) error {
	return r.db.Delete(&database.MockRule{}, "response_id = ?", responseID).Error
}
