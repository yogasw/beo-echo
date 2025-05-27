package repositories

import (
	"gorm.io/gorm"
)

// ResponseRepository handles response-related database operations
type responseRepository struct {
	db *gorm.DB
}

// NewResponseRepository creates a new response repository
func NewResponseRepository(db *gorm.DB) responseRepository {
	return responseRepository{
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
