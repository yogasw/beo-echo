package services

import (
	"context"
	"fmt"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog"
)

// GetReplay retrieves a replay by ID
func (s *ReplayService) GetReplay(ctx context.Context, replayID string) (*database.Replay, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("replay_id", replayID).
		Msg("getting replay details")

	replay, err := s.repo.FindByID(ctx, replayID)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Msg("replay not found")
		return nil, fmt.Errorf("replay not found: %w", err)
	}

	log.Info().
		Str("replay_id", replayID).
		Str("alias", replay.Alias).
		Msg("successfully retrieved replay")

	return replay, nil
}
