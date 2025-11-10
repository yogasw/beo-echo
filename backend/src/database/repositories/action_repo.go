package repositories

import (
	"beo-echo/backend/src/actions"
	"beo-echo/backend/src/database"
	"context"

	"gorm.io/gorm"
)

// actionRepository implements the actions.ActionRepository interface
type actionRepository struct {
	db *gorm.DB
}

// NewActionRepository creates a new action repository
func NewActionRepository(db *gorm.DB) actions.ActionRepository {
	return &actionRepository{db: db}
}

// CreateAction creates a new action in the database
func (r *actionRepository) CreateAction(ctx context.Context, action *database.Action) error {
	return r.db.WithContext(ctx).Create(action).Error
}

// GetActionByID retrieves an action by its ID with associated filters
func (r *actionRepository) GetActionByID(ctx context.Context, id string) (*database.Action, error) {
	var action database.Action
	err := r.db.WithContext(ctx).Preload("Filters").First(&action, "id = ?", id).Error
	if err != nil {
		return nil, err
	}
	return &action, nil
}

// GetActionsByProjectID retrieves all actions for a project, ordered by priority
func (r *actionRepository) GetActionsByProjectID(ctx context.Context, projectID string) ([]database.Action, error) {
	var actions []database.Action
	err := r.db.WithContext(ctx).
		Preload("Filters").
		Where("project_id = ?", projectID).
		Order("priority ASC, created_at ASC").
		Find(&actions).Error
	return actions, err
}

// UpdateAction updates an existing action
func (r *actionRepository) UpdateAction(ctx context.Context, action *database.Action) error {
	return r.db.WithContext(ctx).Save(action).Error
}

// DeleteAction deletes an action by ID
func (r *actionRepository) DeleteAction(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&database.Action{}, "id = ?", id).Error
}

// CreateActionFilter creates a new action filter
func (r *actionRepository) CreateActionFilter(ctx context.Context, filter *database.ActionFilter) error {
	return r.db.WithContext(ctx).Create(filter).Error
}

// GetActionFilters retrieves all filters for an action
func (r *actionRepository) GetActionFilters(ctx context.Context, actionID string) ([]database.ActionFilter, error) {
	var filters []database.ActionFilter
	err := r.db.WithContext(ctx).Where("action_id = ?", actionID).Find(&filters).Error
	return filters, err
}

// DeleteActionFilters deletes all filters for an action
func (r *actionRepository) DeleteActionFilters(ctx context.Context, actionID string) error {
	return r.db.WithContext(ctx).Delete(&database.ActionFilter{}, "action_id = ?", actionID).Error
}

// GetEnabledActionsByProjectAndPoint retrieves enabled actions for a project and execution point
// Results are ordered by priority (lower number = higher priority) and creation time
func (r *actionRepository) GetEnabledActionsByProjectAndPoint(ctx context.Context, projectID string, executionPoint database.ExecutionPoint) ([]database.Action, error) {
	var actions []database.Action
	err := r.db.WithContext(ctx).
		Preload("Filters").
		Where("project_id = ?", projectID).
		Where("enabled = ?", true).
		Where("execution_point = ?", executionPoint).
		Order("priority ASC, created_at ASC").
		Find(&actions).Error
	return actions, err
}
