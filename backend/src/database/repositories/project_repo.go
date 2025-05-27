package repositories

import (
	"beo-echo/backend/src/database"
	"context"

	"gorm.io/gorm"
)

// projectRepository implements the ProjectRepository interface
type projectRepository struct {
	db *gorm.DB
}

// NewProjectRepository creates a new project repository
func NewProjectRepository(db *gorm.DB) *projectRepository {
	return &projectRepository{db: db}
}

// GetWorkspaceProjects returns all projects in a workspace
func (r *projectRepository) GetWorkspaceProjects(ctx context.Context, workspaceID string) ([]database.Project, error) {
	var projects []database.Project
	err := r.db.Where("workspace_id = ?", workspaceID).Find(&projects).Error
	return projects, err
}
