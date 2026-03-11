package services

import (
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/database/repositories"
	"context"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

// UpdateFolder updates an existing replay folder
func (s *ReplayService) UpdateFolder(ctx context.Context, projectID string, folderID string, req UpdateFolderRequest) (*database.ReplayFolder, error) {
	// Verify project exists
	if _, err := s.repo.FindProjectByID(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Fetch target folder directly by ID from DB
	targetFolder, err := s.repo.FindFolderByID(ctx, projectID, folderID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("folder not found")
		}
		return nil, fmt.Errorf("failed to fetch folder: %w", err)
	}

	// If a new ParentID is requested, validate it
	if req.ParentID != nil {
		if *req.ParentID == folderID {
			return nil, fmt.Errorf("invalid parent folder")
		}

		// Fetch all folders for circular-reference check (lightweight projection)
		folders, err := s.repo.FindFoldersByProjectID(ctx, projectID)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch folders: %w", err)
		}

		// Build lookup map
		folderMap := make(map[string]repositories.ReplayFolderListRow, len(folders))
		for _, f := range folders {
			folderMap[f.ID] = f
		}

		// Verify parent exists
		if _, exists := folderMap[*req.ParentID]; !exists {
			return nil, fmt.Errorf("parent folder not found")
		}

		// Check for circular reference
		currentWalkID := *req.ParentID
		for currentWalkID != "" {
			f, ok := folderMap[currentWalkID]
			if !ok || f.ParentID == nil || *f.ParentID == "" {
				break
			}
			if *f.ParentID == folderID {
				return nil, fmt.Errorf("circular folder reference detected")
			}
			currentWalkID = *f.ParentID
		}
	}

	// Apply updates
	if req.Name != nil {
		targetFolder.Name = *req.Name
	}
	if req.Doc != nil {
		targetFolder.Doc = *req.Doc
	}

	// Update ParentID explicitly (can be set to nil)
	if req.UpdateParentID {
		targetFolder.ParentID = req.ParentID
	}

	if err := s.repo.UpdateFolder(ctx, targetFolder); err != nil {
		return nil, fmt.Errorf("failed to save folder: %w", err)
	}

	return targetFolder, nil
}
