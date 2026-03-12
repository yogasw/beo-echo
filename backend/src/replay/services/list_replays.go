package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// ListReplays retrieves all replays for a project
func (s *ReplayService) ListReplays(ctx context.Context, projectID string) (*ListReplaysResponse, error) {
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

	replays, err := s.repo.FindByProjectID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("failed to list replays")
		return nil, fmt.Errorf("failed to list replays: %w", err)
	}

	folders, err := s.repo.FindFoldersByProjectID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("failed to list folders")
		return nil, fmt.Errorf("failed to list folders: %w", err)
	}

	return &ListReplaysResponse{
		Replays: replays,
		Folders: folders,
	}, nil
}
