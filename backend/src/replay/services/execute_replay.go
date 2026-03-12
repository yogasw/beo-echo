package services

import (
	"context"
	"fmt"
	"strings"

	"github.com/rs/zerolog"

	"beo-echo/backend/src/replay/models"
	"beo-echo/backend/src/replay/protocol"
	httpprotocol "beo-echo/backend/src/replay/protocol/http"
)

// ExecuteReplay executes a replay request with the provided configuration
func (s *ReplayService) ExecuteReplay(ctx context.Context, projectID string, req models.ExecuteReplayRequest) (*models.ExecuteReplayResponse, error) {
	log := zerolog.Ctx(ctx)

	// Validate project exists
	_, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("project not found")
		return nil, fmt.Errorf("project not found: %w", err)
	}

	var executor protocol.Executor
	protocolName := strings.ToLower(req.Protocol)

	switch protocolName {
	case "http", "https":
		executor = httpprotocol.NewExecutor()
	default:
		return nil, fmt.Errorf("unsupported protocol: %s (supported: http, https)", req.Protocol)
	}

	resp, err := executor.Execute(ctx, projectID, req)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
