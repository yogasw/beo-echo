package repositories

import (
	"gorm.io/gorm"
)

// ResponseRepository handles response-related database operations
type responseRepository struct {
	db *gorm.DB
}

// NewResponseRepository creates a new response repository
func NewResponseRepository(db *gorm.DB) *responseRepository {
	return &responseRepository{
		db: db,
	}
}

// ValidateResponseHierarchy checks if a response exists and belongs to the specified project/endpoint
func (r *responseRepository) ValidateResponseHierarchy(projectID string, endpointID string, responseID string) (bool, error) {
	// Check if endpoint exists and belongs to the project
	var count int64
	err := r.db.Model(&struct{}{}).
		Table("mock_endpoints").
		Where("id = ? AND project_id = ?", endpointID, projectID).
		Count(&count).Error

	if err != nil || count == 0 {
		return false, err
	}

	// Check if response exists and belongs to the endpoint
	err = r.db.Model(&struct{}{}).
		Table("mock_responses").
		Where("id = ? AND endpoint_id = ?", responseID, endpointID).
		Count(&count).Error

	if err != nil || count == 0 {
		return false, err
	}

	return true, nil
}

// ReorderResponses updates the priority of responses based on the provided order
// Priority is assigned in descending order: first item gets highest priority
func (r *responseRepository) ReorderResponses(endpointID string, responseOrder []string) error {
	// Start a transaction to ensure atomicity
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Calculate highest priority based on array length
	totalResponses := len(responseOrder)

	// Update priority for each response in the order provided
	// First item in array gets highest priority, last item gets lowest priority
	for i, responseID := range responseOrder {
		priority := totalResponses - i // Descending priority: highest first

		result := tx.Model(&struct{}{}).
			Table("mock_responses").
			Where("id = ? AND endpoint_id = ?", responseID, endpointID).
			Update("priority", priority)

		if result.Error != nil {
			tx.Rollback()
			return result.Error
		}

		if result.RowsAffected == 0 {
			tx.Rollback()
			return gorm.ErrRecordNotFound
		}
	}

	return tx.Commit().Error
}
