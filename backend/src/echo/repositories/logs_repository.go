package repositories

import (
	"beo-echo/backend/src/database"

	"gorm.io/gorm"
)

// LogRepository handles database operations for logs
type LogRepository struct {
	DB *gorm.DB
}

// NewLogRepository creates a new log repository
func NewLogRepository(db *gorm.DB) *LogRepository {
	return &LogRepository{
		DB: db,
	}
}

// GetLogs retrieves logs with pagination
func (r *LogRepository) GetLogs(page, pageSize int, projectID string) ([]database.RequestLog, int64, error) {
	var logs []database.RequestLog
	var total int64

	query := r.DB.Model(&database.RequestLog{})

	// Filter by project if specified
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	// Get total count
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Get paginated results
	offset := (page - 1) * pageSize
	if err := query.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// GetLatestLogs retrieves the most recent logs up to a limit
func (r *LogRepository) GetLatestLogs(limit int, projectID string) ([]database.RequestLog, error) {
	var logs []database.RequestLog

	query := r.DB.Model(&database.RequestLog{})

	// Filter by project if specified
	if projectID != "" {
		query = query.Where("project_id = ?", projectID)
	}

	// Get latest logs
	if err := query.Order("created_at DESC").Limit(limit).Find(&logs).Error; err != nil {
		return nil, err
	}

	return logs, nil
}
