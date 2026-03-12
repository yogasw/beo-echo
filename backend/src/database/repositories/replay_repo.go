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
func NewReplayRepository(db *gorm.DB) *replayRepository {
	return &replayRepository{db: db}
}

// ReplayListRow is a minimal projection used by FindByProjectID.
// GORM maps field names to column names automatically — add fields here when you
// need more data in the list view without touching the query string.
type ReplayListRow struct {
	ID         string  `json:"id"`
	Name       string  `json:"name"`
	ProjectID  string  `json:"project_id"`
	FolderID   *string `json:"folder_id"`
	ParentID   *string `json:"parent_id"`
	IsResponse bool    `json:"is_response"`
	Method     string  `json:"method"`
	CreatedAt  string  `json:"created_at"`
	UpdatedAt  string  `json:"updated_at"`
}

// FindByProjectID finds all replays for a specific project.
// Returns only the fields defined in ReplayListRow — no payload, headers, body, etc.
// Full data is loaded on demand via FindByID.
func (r *replayRepository) FindByProjectID(ctx context.Context, projectID string) ([]ReplayListRow, error) {
	var rows []ReplayListRow

	err := r.db.WithContext(ctx).
		Model(&database.Replay{}).
		Where("project_id = ?", projectID).
		Order("created_at DESC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
}

// ReplayFolderListRow is a minimal projection used by FindFoldersByProjectID.
// Add fields here when you need more folder data in the list view — GORM auto-selects
// matching columns without changing the query string.
type ReplayFolderListRow struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	ProjectID string  `json:"project_id"`
	ParentID  *string `json:"parent_id"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt string  `json:"updated_at"`
}

// FindFoldersByProjectID finds all replay folders for a specific project.
// Returns only the fields defined in ReplayFolderListRow — doc and nested relations are excluded.
// Full folder data (including doc) is loaded on demand via a dedicated get endpoint.
func (r *replayRepository) FindFoldersByProjectID(ctx context.Context, projectID string) ([]ReplayFolderListRow, error) {
	var rows []ReplayFolderListRow

	err := r.db.WithContext(ctx).
		Model(&database.ReplayFolder{}).
		Where("project_id = ?", projectID).
		Order("name ASC").
		Scan(&rows).Error

	if err != nil {
		return nil, err
	}

	return rows, nil
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

// FindFolderByID finds a single replay folder by ID scoped to a project.
// Returns the full database.ReplayFolder model (for updates, not list views).
func (r *replayRepository) FindFolderByID(ctx context.Context, projectID string, folderID string) (*database.ReplayFolder, error) {
	var folder database.ReplayFolder

	err := r.db.WithContext(ctx).
		Where("id = ? AND project_id = ?", folderID, projectID).
		First(&folder).Error

	if err != nil {
		return nil, err
	}

	return &folder, nil
}

// Create creates a new replay
func (r *replayRepository) Create(ctx context.Context, replay *database.Replay) error {
	return r.db.WithContext(ctx).Create(replay).Error
}

// CreateFolder creates a new replay folder
func (r *replayRepository) CreateFolder(ctx context.Context, folder *database.ReplayFolder) error {
	return r.db.WithContext(ctx).Create(folder).Error
}

// UpdateFolder updates an existing replay folder
func (r *replayRepository) UpdateFolder(ctx context.Context, folder *database.ReplayFolder) error {
	return r.db.WithContext(ctx).Save(folder).Error
}

// DeleteFolder deletes a folder, all its subfolders, and all replays inside it
func (r *replayRepository) DeleteFolder(ctx context.Context, projectID string, folderID string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Helper function for recursive deletion
		var deleteFolderRecursive func(fID string) error
		deleteFolderRecursive = func(fID string) error {
			// Find all subfolders of the current folder
			var subfolders []database.ReplayFolder
			if err := tx.Where("parent_id = ? AND project_id = ?", fID, projectID).Find(&subfolders).Error; err != nil {
				return err
			}

			// Recursively delete subfolders
			for _, sub := range subfolders {
				if err := deleteFolderRecursive(sub.ID); err != nil {
					return err
				}
			}

			// Delete all replays in this folder
			if err := tx.Where("folder_id = ? AND project_id = ?", fID, projectID).Delete(&database.Replay{}).Error; err != nil {
				return err
			}

			// Delete the folder itself
			if err := tx.Where("id = ? AND project_id = ?", fID, projectID).Delete(&database.ReplayFolder{}).Error; err != nil {
				return err
			}

			return nil
		}

		// Verify the target folder actually exists and belongs to the project
		var targetFolder database.ReplayFolder
		if err := tx.Where("id = ? AND project_id = ?", folderID, projectID).First(&targetFolder).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return errors.New("folder not found")
			}
			return err
		}

		// Start recursive deletion
		return deleteFolderRecursive(folderID)
	})
}

// Update updates an existing replay
func (r *replayRepository) Update(ctx context.Context, replay *database.Replay) error {
	return r.db.WithContext(ctx).Save(replay).Error
}

// Delete deletes a replay by ID, and any children (histories)
func (r *replayRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Delete any children checking this replay as a parent first
		if err := tx.Where("parent_id = ?", id).Delete(&database.Replay{}).Error; err != nil {
			return err
		}

		result := tx.Where("id = ?", id).Delete(&database.Replay{})

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("replay not found")
		}

		return nil
	})
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
