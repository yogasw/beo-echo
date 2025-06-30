package repositories

import (
	"context"
	"strings"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/echo/services"
	"gorm.io/gorm"
)

// projectSearchRepository implements services.ProjectSearchRepository
type projectSearchRepository struct {
	db *gorm.DB
}

// NewProjectSearchRepository creates a new project search repository
func NewProjectSearchRepository(db *gorm.DB) services.ProjectSearchRepository {
	return &projectSearchRepository{db: db}
}

// CheckAliasAvailability checks if a project alias is available globally
func (r *projectSearchRepository) CheckAliasAvailability(ctx context.Context, alias string) (bool, error) {
	var count int64
	err := r.db.Model(&database.Project{}).Where("alias = ?", alias).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count == 0, nil
}

// SearchProjectsByNameInUserWorkspaces searches for projects by name in workspaces the user has joined
func (r *projectSearchRepository) SearchProjectsByNameInUserWorkspaces(ctx context.Context, userID, searchQuery string) ([]services.ProjectSearchResult, error) {
	var results []services.ProjectSearchResult
	
	query := `
		SELECT 
			p.id,
			p.name,
			p.alias,
			p.workspace_id,
			w.name as workspace_name
		FROM projects p
		INNER JOIN workspaces w ON p.workspace_id = w.id
		INNER JOIN user_workspaces uw ON w.id = uw.workspace_id AND uw.user_id = ?
		WHERE LOWER(p.name) LIKE LOWER(?)
		ORDER BY p.name ASC
	`
	
	searchPattern := "%" + strings.ToLower(searchQuery) + "%"
	err := r.db.Raw(query, userID, searchPattern).Scan(&results).Error
	
	return results, err
}
