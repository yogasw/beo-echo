package services

import (
	"context"
	"net/http"
	"time"

	"beo-echo/backend/src/database"
)

// ReplayRepository defines data access requirements for replay operations
type ReplayRepository interface {
	// Replay CRUD operations
	FindByProjectID(ctx context.Context, projectID string) ([]database.Replay, error)
	FindFoldersByProjectID(ctx context.Context, projectID string) ([]database.ReplayFolder, error)
	FindByID(ctx context.Context, id string) (*database.Replay, error)
	Create(ctx context.Context, replay *database.Replay) error
	CreateFolder(ctx context.Context, folder *database.ReplayFolder) error
	Update(ctx context.Context, replay *database.Replay) error
	Delete(ctx context.Context, id string) error
	DeleteFolder(ctx context.Context, projectID string, folderID string) error

	// Replay execution logging
	CreateRequestLog(ctx context.Context, log *database.RequestLog) error
	FindReplayLogs(ctx context.Context, projectID string, replayID *string) ([]database.RequestLog, error)

	// Project validation
	FindProjectByID(ctx context.Context, projectID string) (*database.Project, error)
}

// ReplayService implements replay business operations
type ReplayService struct {
	repo   ReplayRepository
	client *http.Client
}

// NewReplayService creates a new replay service
func NewReplayService(repo ReplayRepository) *ReplayService {
	return &ReplayService{
		repo: repo,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// CreateFolderRequest represents the request payload for creating a replay folder
type CreateFolderRequest struct {
	Name     string  `json:"name" binding:"required"`
	ParentID *string `json:"parent_id"`
}

// ListReplaysResponse represents the response for listing replays
type ListReplaysResponse struct {
	Replays []database.Replay       `json:"replays"`
	Folders []database.ReplayFolder `json:"folders"`
}

// CreateReplayRequest represents the request payload for creating a replay
type CreateReplayRequest struct {
	Name     string            `json:"name"`
	FolderID *string           `json:"folder_id"`
	Protocol string            `json:"protocol" binding:"required"`
	Method   string            `json:"method" binding:"required"`
	Url      string            `json:"url" binding:"required"`
	Headers  map[string]string `json:"headers"`
	Payload  string            `json:"payload"`
	Metadata map[string]any    `json:"metadata"` // Additional protocol-specific metadata
	Config   map[string]any    `json:"config"`   // Optional configuration for specific protocols
}

// UpdateReplayRequest represents the request payload for updating a replay
type UpdateReplayRequest struct {
	Name     *string            `json:"name"`
	FolderID *string            `json:"folder_id"`
	Protocol *string            `json:"protocol"`
	Method   *string            `json:"method"`
	Url      *string            `json:"url"`
	Headers  *map[string]string `json:"headers"`
	Payload  *string            `json:"payload"`
	Metadata *map[string]any   `json:"metadata"` // Additional protocol-specific metadata
	Config   *map[string]any   `json:"config"`   // Optional configuration for specific protocols
}

