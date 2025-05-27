package repositories

import (
	"context"
	"errors"

	"gorm.io/gorm"

	"beo-echo/backend/src/database"
)

// replayRepository implements services.replayRepository
type replayRepository struct {
	db *gorm.DB
}

// NewreplayRepository creates a new replay repository
func NewReplayRepository(db *gorm.DB) replayRepository {
	return replayRepository{db: db}
}

// FindByProjectID finds all replays for a specific project
func (r *replayRepository) FindByProjectID(ctx context.Context, projectID string) ([]database.Replay, error) {
	var replays []database.Replay

	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Find(&replays).Error

	if err != nil {
		return nil, err
	}

	return replays, nil
}

// FindByID finds a replay by its ID
func (r *replayRepository) FindByID(ctx context.Context, id string) (*database.Replay, error) {
	var replay database.Replay

	err := r.db.WithContext(ctx).
		Where("id = ?", id).
		First(&replay).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("replay not found")
		}
		return nil, err
	}

	return &replay, nil
}

// Create creates a new replay
func (r *replayRepository) Create(ctx context.Context, replay *database.Replay) error {
	return r.db.WithContext(ctx).Create(replay).Error
}

// Update updates an existing replay
func (r *replayRepository) Update(ctx context.Context, replay *database.Replay) error {
	return r.db.WithContext(ctx).Save(replay).Error
}

// Delete deletes a replay by ID
func (r *replayRepository) Delete(ctx context.Context, id string) error {
	result := r.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&database.Replay{})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return errors.New("replay not found")
	}

	return nil
}

// CreateRequestLog creates a new request log entry
func (r *replayRepository) CreateRequestLog(ctx context.Context, log *database.RequestLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

// FindReplayLogs finds request logs for replay executions
func (r *replayRepository) FindReplayLogs(ctx context.Context, projectID string, replayID *string) ([]database.RequestLog, error) {
	var logs []database.RequestLog

	query := r.db.WithContext(ctx).
		Where("project_id = ? AND source = ?", projectID, database.RequestSourceReplay).
		Order("created_at DESC")

	// If specific replay ID is provided, filter by it
	// We can identify replay logs by matching the path (target_url) with request logs
	// This is a simple approach - in a more sophisticated system,
	// we might add a replay_id field to RequestLog
	if replayID != nil {
		// Find the replay to get its target URL
		var replay database.Replay
		err := r.db.WithContext(ctx).
			Where("id = ?", *replayID).
			First(&replay).Error
		if err != nil {
			return nil, err
		}

		query = query.Where("path = ?", replay.Url)
	}

	err := query.Find(&logs).Error
	if err != nil {
		return nil, err
	}

	return logs, nil
}

// FindProjectByID finds a project by ID for validation
func (r *replayRepository) FindProjectByID(ctx context.Context, projectID string) (*database.Project, error) {
	var project database.Project

	err := r.db.WithContext(ctx).
		Where("id = ?", projectID).
		First(&project).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}

	return &project, nil
}
