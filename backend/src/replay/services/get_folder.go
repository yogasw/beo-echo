package services

import (
	"beo-echo/backend/src/database"
	"context"
	"fmt"
)

// GetFolder retrieves a single ReplayFolder by ID scoped to the project
func (s *ReplayService) GetFolder(ctx context.Context, projectID string, folderID string) (*database.ReplayFolder, error) {
	folder, err := s.repo.FindFolderByID(ctx, projectID, folderID)
	if err != nil {
		return nil, fmt.Errorf("folder not found: %w", err)
	}
	return folder, nil
}
