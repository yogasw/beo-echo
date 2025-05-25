package services

import (
	"context"
	"encoding/json"
	"fmt"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog"
)

// CreateReplay creates a new replay configuration
func (s *ReplayService) CreateReplay(ctx context.Context, projectID string, req CreateReplayRequest) (*database.Replay, error) {
	log := zerolog.Ctx(ctx)

	log.Info().
		Str("project_id", projectID).
		Str("alias", req.Alias).
		Str("protocol", req.Protocol).
		Str("method", req.Method).
		Str("target_url", req.TargetURL).
		Msg("creating new replay")

	// Validate project exists
	_, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("project not found")
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Convert headers and metadata to JSON
	headersJSON, err := json.Marshal(req.Headers)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to marshal headers")
		return nil, fmt.Errorf("invalid headers format: %w", err)
	}

	metadataJSON, err := json.Marshal(req.Metadata)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to marshal metadata")
		return nil, fmt.Errorf("invalid metadata format: %w", err)
	}

	replay := &database.Replay{
		Alias:      req.Alias,
		ProjectID:  projectID,
		FolderID:   req.FolderID,
		Protocol:   req.Protocol,
		Method:     req.Method,
		TargetURL:  req.TargetURL,
		Service:    req.Service,
		MethodName: req.MethodName,
		Headers:    string(headersJSON),
		Payload:    req.Payload,
		Metadata:   string(metadataJSON),
		IsMutation: req.IsMutation,
		Path:       req.Path,
	}

	err = s.repo.Create(ctx, replay)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("alias", req.Alias).
			Msg("failed to create replay")
		return nil, fmt.Errorf("failed to create replay: %w", err)
	}

	log.Info().
		Str("replay_id", replay.ID).
		Str("project_id", projectID).
		Str("alias", req.Alias).
		Msg("successfully created replay")

	return replay, nil
}
