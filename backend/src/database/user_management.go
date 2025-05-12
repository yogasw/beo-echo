package database

import (
	"log"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// InitializeDefaultUserAndWorkspace checks if any users exist,
// and if not, creates a default admin user with a demo workspace and project
func InitializeDefaultUserAndWorkspace(db *gorm.DB) error {
	// Check if any users exist in the system
	var userCount int64
	if err := db.Model(&User{}).Count(&userCount).Error; err != nil {
		return err
	}

	// If users already exist, no action needed
	if userCount > 0 {
		log.Println("Users already exist in the system, skipping default user creation")
		return nil
	}

	// Create default admin user
	log.Println("No users found, creating default admin user...")

	// Generate hashed password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin"), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	defaultUser := User{
		Email:    "admin@admin.com",
		Name:     "Admin",
		Password: string(hashedPassword),
		IsOwner:  true,
	}

	// Create user in transaction to ensure all related objects are created
	tx := db.Begin()
	if err := tx.Create(&defaultUser).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create demo workspace
	demoWorkspace := Workspace{
		Name: "Demo Workspace",
	}

	if err := tx.Create(&demoWorkspace).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Link user to workspace as admin
	userWorkspace := UserWorkspace{
		UserID:      defaultUser.ID,
		WorkspaceID: demoWorkspace.ID,
		Role:        "admin",
	}

	if err := tx.Create(&userWorkspace).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Create demo project in the workspace
	demoProject := Project{
		Name:        "Demo Project",
		WorkspaceID: demoWorkspace.ID,
		Mode:        ModeMock,
		Status:      "running",
		Alias:       "demo",
	}

	if err := tx.Create(&demoProject).Error; err != nil {
		tx.Rollback()
		return err
	}

	// Commit transaction
	if err := tx.Commit().Error; err != nil {
		return err
	}

	log.Println("Successfully created default admin user, workspace and project")
	return nil
}

// VerifyPassword checks if the provided password matches the hashed password in the database
func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// HashPassword creates a bcrypt hash of the password
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// GetUserByEmail retrieves a user by their email address
func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.Where("email = ?", email).First(&user)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

// GetUserWorkspaces returns all workspaces that the user is a member of
func GetUserWorkspaces(userID string) ([]Workspace, error) {
	var workspaces []Workspace
	err := DB.Joins("JOIN user_workspaces ON user_workspaces.workspace_id = workspaces.id").
		Where("user_workspaces.user_id = ?", userID).
		Find(&workspaces).Error

	return workspaces, err
}

// GetWorkspaceProjects returns all projects in a workspace
func GetWorkspaceProjects(workspaceID string) ([]Project, error) {
	var projects []Project
	err := DB.Where("workspace_id = ?", workspaceID).Find(&projects).Error
	return projects, err
}

// IsUserWorkspaceAdmin checks if a user is an admin in a specific workspace
func IsUserWorkspaceAdmin(userID string, workspaceID string) (bool, error) {
	var userWorkspace UserWorkspace
	result := DB.Where("user_id = ? AND workspace_id = ?", userID, workspaceID).First(&userWorkspace)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return false, nil
		}
		return false, result.Error
	}

	return userWorkspace.Role == "admin", nil
}
