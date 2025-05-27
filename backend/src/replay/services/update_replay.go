package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog"
)

// UpdateReplay updates an existing replay configuration
func (s *ReplayService) UpdateReplay(ctx context.Context, replayID string, req UpdateReplayRequest) (*database.Replay, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("replay_id", replayID).
		Msg("updating replay")

	// Get existing replay
	replay, err := s.repo.FindByID(ctx, replayID)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Msg("replay not found")
		return nil, fmt.Errorf("replay not found: %w", err)
	}

	// Update fields if provided
	if req.Name != nil {
		replay.Name = *req.Name
	}

	if req.FolderID != nil {
		replay.FolderID = req.FolderID
	}

	if req.Protocol != nil {
		protocol := database.ReplayProtocol(strings.ToLower(*req.Protocol))
		if protocol != database.ReplayProtocolHTTP {
			log.Error().
				Str("protocol", *req.Protocol).
				Msg("invalid protocol specified")
			return nil, fmt.Errorf("invalid protocol: %s", *req.Protocol)
		}
		replay.Protocol = protocol
	}

	if req.Method != nil {
		replay.Method = strings.ToUpper(*req.Method)
	}

	if req.Url != nil {
		replay.Url = *req.Url
	}

	if req.Headers != nil {
		headersJSON, err := json.Marshal(*req.Headers)
		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to marshal headers")
			return nil, fmt.Errorf("invalid headers format: %w", err)
		}
		replay.Headers = string(headersJSON)
	}

	if req.Payload != nil {
		replay.Payload = *req.Payload
	}

	// Update in database
	err = s.repo.Update(ctx, replay)
	if err != nil {
		log.Error().
			Err(err).
			Str("replay_id", replayID).
			Msg("failed to update replay")
		return nil, fmt.Errorf("failed to update replay: %w", err)
	}

	log.Info().
		Str("replay_id", replayID).
		Str("name", replay.Name).
		Msg("successfully updated replay")

	return replay, nil
}
