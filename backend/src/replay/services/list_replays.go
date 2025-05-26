package services

import (
	"context"
	"fmt"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog"
)

// ListReplays retrieves all replays for a project
func (s *ReplayService) ListReplays(ctx context.Context, projectID string) ([]database.Replay, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("project_id", projectID).
		Msg("listing replays for project")

	// Validate project exists
	_, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("project not found")
		return nil, fmt.Errorf("project not found: %w", err)
	}

	replays, err := s.repo.FindByProjectID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("failed to list replays")
		return nil, fmt.Errorf("failed to list replays: %w", err)
	}

	log.Info().
		Str("project_id", projectID).
		Int("count", len(replays)).
		Msg("successfully listed replays")

	return replays, nil
}
