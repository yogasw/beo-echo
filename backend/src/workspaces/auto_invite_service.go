package workspaces

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"

	"beo-echo/backend/src/database"

	"github.com/rs/zerolog/log"
)

// AutoInviteService handles automatic invitations to workspaces based on email domains
type AutoInviteService struct {
	db *gorm.DB
}

// NewAutoInviteService creates a new auto-invite service
func NewAutoInviteService(db *gorm.DB) *AutoInviteService {
	return &AutoInviteService{db: db}
}

// ProcessUserAutoInvite checks if a user should be automatically invited to any workspaces
// based on their email domain, and creates the necessary UserWorkspace records if so
func (s *AutoInviteService) ProcessUserAutoInvite(ctx context.Context, user *database.User) error {

	if user == nil || user.Email == "" {
		return nil // No user or email, nothing to do
	}

	// Extract the domain from the user's email
	parts := strings.Split(user.Email, "@")
	if len(parts) != 2 {
		return nil // Invalid email format, nothing to do
	}
	userDomain := strings.ToLower(parts[1]) // Convert to lowercase for case-insensitive comparison

	// Find all workspaces with auto-invite enabled
	var workspaces []database.Workspace
	if err := s.db.Where("auto_invite_enabled = ?", true).Find(&workspaces).Error; err != nil {
		log.Error().Err(err).Str("user_id", user.ID).Msg("failed to get workspaces for auto-invite")
		return err
	}

	for _, workspace := range workspaces {
		// Check if user is already a member of this workspace
		var existingMembership database.UserWorkspace
		existingMembershipCount := s.db.Where("user_id = ? AND workspace_id = ?", user.ID, workspace.ID).
			First(&existingMembership).RowsAffected

		if existingMembershipCount > 0 {
			// User already has a membership record for this workspace, skip
			continue
		}

		// Check if the user's domain matches any in the workspace's auto-invite domains
		if workspace.AutoInviteDomains == "" {
			continue // No domains configured
		}

		domains := strings.Split(workspace.AutoInviteDomains, ",")
		for _, domain := range domains {
			// Trim whitespace and convert to lowercase for comparison
			domainToCheck := strings.ToLower(strings.TrimSpace(domain))

			if domainToCheck == userDomain {
				// Create a new UserWorkspace record
				newMembership := database.UserWorkspace{
					UserID:      user.ID,
					WorkspaceID: workspace.ID,
					Role:        workspace.AutoInviteRole,
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}

				if err := s.db.Create(&newMembership).Error; err != nil {
					log.Error().
						Err(err).
						Str("user_id", user.ID).
						Str("workspace_id", workspace.ID).
						Str("domain", userDomain).
						Msg("failed to create auto-invite membership")
					return err
				}

				log.Info().
					Str("user_id", user.ID).
					Str("email", user.Email).
					Str("workspace_id", workspace.ID).
					Str("workspace_name", workspace.Name).
					Str("role", workspace.AutoInviteRole).
					Msg("user auto-invited to workspace based on email domain")

				break // No need to check other domains for this workspace
			}
		}
	}

	return nil
}
