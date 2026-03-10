package services

import (
	"context"
	"fmt"

	"github.com/rs/zerolog"
)

// DeleteFolder deletes a replay folder and all its contents
func (s *ReplayService) DeleteFolder(ctx context.Context, projectID string, folderID string) error {
	log := zerolog.Ctx(ctx)

	// Validate project exists
	_, err := s.repo.FindProjectByID(ctx, projectID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Msg("project not found")
		return fmt.Errorf("project not found: %w", err)
	}

	err = s.repo.DeleteFolder(ctx, projectID, folderID)
	if err != nil {
		log.Error().
			Err(err).
			Str("project_id", projectID).
			Str("folder_id", folderID).
			Msg("failed to delete replay folder")
		return fmt.Errorf("failed to delete folder: %w", err)
	}

	return nil
}
