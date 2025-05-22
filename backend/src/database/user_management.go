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
// This method is maintained here for model method functionality
func (user *User) VerifyPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// Note: All other user management functions have been moved to repositories/user_repo.go
// and are accessed via the users module services.
