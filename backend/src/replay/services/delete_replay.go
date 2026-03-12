package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// DeleteReplay removes a replay
func (s *ReplayService) DeleteReplay(ctx context.Context, replayID string) error {
	log := zerolog.Ctx(ctx)

	// Verify replay exists
	_, err := s.repo.FindByID(ctx, replayID)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Msg("replay not found")
		return fmt.Errorf("replay not found: %w", err)
	}

	err = s.repo.Delete(ctx, replayID)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Msg("failed to delete replay")
		return fmt.Errorf("failed to delete replay: %w", err)
	}

	return nil
}
