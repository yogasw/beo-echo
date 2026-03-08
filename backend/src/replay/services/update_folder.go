package services

import (
	"beo-echo/backend/src/database"
	"context"
	"fmt"
)

// UpdateFolder updates an existing replay folder
func (s *ReplayService) UpdateFolder(ctx context.Context, projectID string, folderID string, req UpdateFolderRequest) (*database.ReplayFolder, error) {
	// Verify project exists
	if _, err := s.repo.FindProjectByID(ctx, projectID); err != nil {
		return nil, fmt.Errorf("project not found")
	}

	// Verify the folder exists and belongs to the project
	folders, err := s.repo.FindFoldersByProjectID(ctx, projectID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch folders: %w", err)
	}

	var targetFolder *database.ReplayFolder
	for i := range folders {
		if folders[i].ID == folderID {
			targetFolder = &folders[i]
			break
		}
	}

	if targetFolder == nil {
		return nil, fmt.Errorf("folder not found")
	}

	// Verify the new parent folder exists and is valid (not a child of itself)
	if req.ParentID != nil {
		// Parent can't be itself
		if *req.ParentID == folderID {
			return nil, fmt.Errorf("invalid parent folder")
		}

		// Verify parent exists
		var parentExists bool
		for _, f := range folders {
			if f.ID == *req.ParentID {
				parentExists = true
				break
			}
		}
		if !parentExists {
			return nil, fmt.Errorf("parent folder not found")
		}

		// Check for circular reference by ensuring the new parent isn't a child of this folder
		// A simple way to check is to traverse up from the new parent
		currentWalkID := *req.ParentID
		for currentWalkID != "" {
			var currentParentID *string
			for _, f := range folders {
				if f.ID == currentWalkID {
					currentParentID = f.ParentID
					break
				}
			}
			if currentParentID == nil || *currentParentID == "" {
				break
			}
			if *currentParentID == folderID {
				return nil, fmt.Errorf("circular folder reference detected")
			}
			currentWalkID = *currentParentID
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
