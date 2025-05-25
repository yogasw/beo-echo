package services

import (
	"context"
	"fmt"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog"
)

// GetReplayLogs retrieves execution logs for replays
func (s *ReplayService) GetReplayLogs(ctx context.Context, projectID string, replayID *string) ([]database.RequestLog, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("project_id", projectID).
		Str("replay_id", func() string {
			if replayID != nil {
				return *replayID
			}
			return "all"
		}()).
		Msg("getting replay execution logs")

	// Validate project exists
	_, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("project not found")
		return nil, fmt.Errorf("project not found: %w", err)
	}

	logs, err := s.repo.FindReplayLogs(ctx, projectID, replayID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("failed to get replay logs")
		return nil, fmt.Errorf("failed to get replay logs: %w", err)
	}

	log.Info().
		Str("project_id", projectID).
		Int("count", len(logs)).
		Msg("successfully retrieved replay logs")

	return logs, nil
}
