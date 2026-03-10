package protocol

import (
	"context"

	"beo-echo/backend/src/replay/models"
)

// Executor defines the interface for protocol-specific replay execution
type Executor interface {
	Execute(ctx context.Context, projectID string, req models.ExecuteReplayRequest) (*models.ExecuteReplayResponse, error)
}
