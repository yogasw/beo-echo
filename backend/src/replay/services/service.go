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
	FindByID(ctx context.Context, id string) (*database.Replay, error)
	Create(ctx context.Context, replay *database.Replay) error
	Update(ctx context.Context, replay *database.Replay) error
	Delete(ctx context.Context, id string) error

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

// CreateReplayRequest represents the request payload for creating a replay
type CreateReplayRequest struct {
	Name       string            `json:"name"`
	FolderID   *string           `json:"folder_id"`
	Protocol   string            `json:"protocol" binding:"required"`
	Method     string            `json:"method" binding:"required"`
	Url        string            `json:"url" binding:"required"`
	Service    string            `json:"service"`
	MethodName string            `json:"method_name"`
	Headers    map[string]string `json:"headers"`
	Payload    string            `json:"payload"`
	Metadata   map[string]string `json:"metadata"`
	Path       []string          `json:"path"`
}

// ExecuteReplayRequest represents the request payload for executing a replay test
type ExecuteReplayRequest struct {
	Protocol string            `json:"protocol" binding:"required"` // http, https, ws, grpc
	Method   string            `json:"method" binding:"required"`   // HTTP method or operation type
	URL      string            `json:"url" binding:"required"`      // Target URL
	Headers  map[string]string `json:"headers"`                     // Request headers
	Body     string            `json:"body"`                        // Request body/payload
	Query    map[string]string `json:"query"`                       // Query parameters
	Metadata map[string]string `json:"metadata"`                    // Additional protocol-specific metadata
}

// ExecuteReplayResponse represents the response from executing a replay
type ExecuteReplayResponse struct {
	ReplayID        string            `json:"replay_id"`
	StatusCode      int               `json:"status_code"`
	ResponseBody    string            `json:"response_body"`
	ResponseHeaders map[string]string `json:"response_headers"`
	LatencyMS       int               `json:"latency_ms"`
	Error           string            `json:"error,omitempty"`
	LogID           string            `json:"log_id"`
}
