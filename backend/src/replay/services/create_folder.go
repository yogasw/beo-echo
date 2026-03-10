package services

import (
	"context"

	"beo-echo/backend/src/database"
)

// CreateFolder creates a new replay folder
func (s *ReplayService) CreateFolder(ctx context.Context, projectID string, req CreateFolderRequest) (*database.ReplayFolder, error) {
	// Verify project exists
	if _, err := s.repo.FindProjectByID(ctx, projectID); err != nil {
		return nil, err
	}

	folder := &database.ReplayFolder{
		Name:      req.Name,
		Doc:       req.Doc,
		ProjectID: projectID,
		ParentID:  req.ParentID,
	}

	if err := s.repo.CreateFolder(ctx, folder); err != nil {
		return nil, err
	}

	return folder, nil
}
