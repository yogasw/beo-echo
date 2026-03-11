package services

import (
	"context"
	"net/http"
	"time"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/database/repositories"
)

// ReplayRepository defines data access requirements for replay operations
type ReplayRepository interface {
	// Replay CRUD operations
	FindByProjectID(ctx context.Context, projectID string) ([]repositories.ReplayListRow, error)
	FindFoldersByProjectID(ctx context.Context, projectID string) ([]repositories.ReplayFolderListRow, error)
	FindByID(ctx context.Context, id string) (*database.Replay, error)
	FindFolderByID(ctx context.Context, projectID string, folderID string) (*database.ReplayFolder, error)
	Create(ctx context.Context, replay *database.Replay) error
	CreateFolder(ctx context.Context, folder *database.ReplayFolder) error
	UpdateFolder(ctx context.Context, folder *database.ReplayFolder) error
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
	Doc      string  `json:"doc"`
	ParentID *string `json:"parent_id"`
}

// UpdateFolderRequest represents the request payload for updating a replay folder
type UpdateFolderRequest struct {
	Name           *string `json:"name"`
	Doc            *string `json:"doc"`
	ParentID       *string `json:"parent_id"`
	UpdateParentID bool    `json:"update_parent_id"` // indicates if ParentID should be updated (even to null)
}

// ListReplaysResponse represents the response for listing replays
type ListReplaysResponse struct {
	Replays []repositories.ReplayListRow       `json:"replays"`
	Folders []repositories.ReplayFolderListRow `json:"folders"`
}

// HeaderItem represents a single header with key, value, and description
type HeaderItem struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description"`
}

// CreateReplayRequest represents the request payload for creating a replay
type CreateReplayRequest struct {
	Name     string         `json:"name"`
	Doc      string         `json:"doc"`
	FolderID *string        `json:"folder_id"`
	ParentID *string        `json:"parent_id"` // For histories/saved responses
	Protocol string         `json:"protocol" binding:"required"`
	Method   string         `json:"method" binding:"required"`
	Url      string         `json:"url" binding:"required"`
	Headers  []HeaderItem   `json:"headers"`
	Payload  string         `json:"payload"`
	Metadata map[string]any `json:"metadata"` // Additional protocol-specific metadata
	Config   map[string]any `json:"config"`   // Optional configuration for specific protocols

	// Response fields for creating histories
	IsResponse     bool    `json:"is_response"`
	ResponseStatus *int    `json:"response_status"`
	ResponseMeta   *string `json:"response_meta"`
	ResponseBody   *string `json:"response_body"`
	LatencyMS      *int    `json:"latency_ms"`
}

// UpdateReplayRequest represents the request payload for updating a replay
type UpdateReplayRequest struct {
	Name           *string         `json:"name"`
	Doc            *string         `json:"doc"`
	FolderID       *string         `json:"folder_id"`
	UpdateFolderID bool            `json:"update_folder_id"`
	Protocol       *string         `json:"protocol"`
	Method         *string         `json:"method"`
	Url            *string         `json:"url"`
	Headers        *[]HeaderItem   `json:"headers"`
	Payload        *string         `json:"payload"`
	Metadata       *map[string]any `json:"metadata"` // Additional protocol-specific metadata
	Config         *map[string]any `json:"config"`   // Optional configuration for specific protocols

	// Response fields for updating histories
	ResponseStatus *int            `json:"response_status"`
	ResponseMeta   *string         `json:"response_meta"`
	ResponseBody   *string         `json:"response_body"`
	LatencyMS      *int            `json:"latency_ms"`
}

