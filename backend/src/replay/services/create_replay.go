package services

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog"
)

// CreateReplay creates a new replay configuration
func (s *ReplayService) CreateReplay(ctx context.Context, projectID string, req CreateReplayRequest) (*database.Replay, error) {
	log := zerolog.Ctx(ctx)

	log.Info().Str("project_id", projectID).Str("name", req.Name).Msg("creating replay")
	name := req.Name
	if name == "" {
		name = req.Url
	}

	protocol := database.ReplayProtocol(strings.ToLower(req.Protocol))
	if protocol == "" {
		log.Error().
			Str("protocol", req.Protocol).
			Msg("invalid protocol specified")
		return nil, fmt.Errorf("invalid protocol: %s", req.Protocol)
	}

	// Validate project exists
	_, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("project not found")
		return nil, fmt.Errorf("project not found: %w", err)
	}

	// Convert headers to JSON
	headersJSON, err := json.Marshal(req.Headers)
	if err != nil {
		log.Error().
			Err(err).
			Msg("failed to marshal headers")
		return nil, fmt.Errorf("invalid headers format: %w", err)
	}

	replay := &database.Replay{
		Name:      name,
		ProjectID: projectID,
		FolderID:  req.FolderID,
		Protocol:  database.ReplayProtocol(strings.ToLower(req.Protocol)),
		Method:    strings.ToUpper(req.Method),
		Url:       req.Url,
		Headers:   string(headersJSON),
		Payload:   req.Payload,
	}

	err = s.repo.Create(ctx, replay)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("name", req.Name).
			Msg("failed to create replay")
		return nil, fmt.Errorf("failed to create replay: %w", err)
	}

	log.Info().
		Str("replay_id", replay.ID).
		Str("project_id", projectID).
		Str("name", req.Name).
		Msg("successfully created replay")

	return replay, nil
}
