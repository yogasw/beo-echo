package services

import (
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/utils"
)

// SyncConfigsToGit synchronizes the configurations with a Git repository
func SyncConfigsToGit() error {
	// Get Git configuration
	config, err := GetGitConfig()
	if err != nil {
		return fmt.Errorf("failed to get Git config: %w", err)
	}

	// Skip if Git sync is disabled
	if !config.Enabled {
		log.Println("Git sync is disabled, skipping")
		return nil
	}

	// Validate configuration
	if config.RepoURL == "" {
		return errors.New("repository URL is required")
	}

	if config.Branch == "" {
		config.Branch = "main"
	}

	// Create or ensure clone folder exists
	cloneFolder := config.CloneFolder
	if cloneFolder == "" {
		cloneFolder = filepath.Join(lib.CONFIGS_DIR, "git-repo")
	}
	if err := utils.EnsureDirectoryExists(cloneFolder); err != nil {
		return fmt.Errorf("failed to create clone directory: %w", err)
	}

	// Check if repository is already cloned
	gitDir := filepath.Join(cloneFolder, ".git")
	repoExists := utils.FileExists(gitDir)

	if repoExists {
		// Pull changes if repository exists
		if config.PullBefore {
			cmd := exec.Command("git", "-C", cloneFolder, "pull", "origin", config.Branch)
			if output, err := cmd.CombinedOutput(); err != nil {
				return fmt.Errorf("git pull failed: %w\n%s", err, output)
			}
		}
	} else {
		// Clone the repository
		var cmd *exec.Cmd
		switch config.AuthType {
		case "basic":
			// Basic authentication via URL
			repoURL := fmt.Sprintf("https://%s:%s@%s",
				config.Username,
				config.Password,
				config.RepoURL[8:]) // Strip "https://"
			cmd = exec.Command("git", "clone", "--single-branch", "--branch", config.Branch, repoURL, cloneFolder)
		case "ssh":
			// SSH authentication (would require SSH key setup)
			// For simplicity, we'll just use the SSH URL
			cmd = exec.Command("git", "clone", "--single-branch", "--branch", config.Branch, config.RepoURL, cloneFolder)
		default:
			// No authentication
			cmd = exec.Command("git", "clone", "--single-branch", "--branch", config.Branch, config.RepoURL, cloneFolder)
		}

		if output, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("git clone failed: %w\n%s", err, output)
		}
	}

	// Copy configuration files from Git to configs directory
	configsInRepo := filepath.Join(cloneFolder, "configs")
	if utils.FileExists(configsInRepo) {
		// Read files in the repo/configs directory
		files, err := os.ReadDir(configsInRepo)
		if err != nil {
			return fmt.Errorf("failed to read configs directory in repo: %w", err)
		}

		// Copy each JSON file to the application configs directory
		for _, file := range files {
			if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
				continue
			}

			src := filepath.Join(configsInRepo, file.Name())
			dst := filepath.Join(lib.CONFIGS_DIR, file.Name())

			if err := utils.CopyFile(src, dst); err != nil {
				log.Printf("Warning: Failed to copy config file %s: %v", file.Name(), err)
			}
		}
	}

	return nil
}
