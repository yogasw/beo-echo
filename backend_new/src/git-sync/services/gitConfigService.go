package services

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/utils"
)

// GitConfig represents the Git configuration
type GitConfig struct {
	Enabled     bool   `json:"enabled"`
	RepoURL     string `json:"repoUrl"`
	Branch      string `json:"branch"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	SSHKey      string `json:"sshKey"`
	AuthType    string `json:"authType"` // "none", "basic", "ssh"
	CommitMsg   string `json:"commitMsg"`
	PullBefore  bool   `json:"pullBefore"`
	CloneFolder string `json:"cloneFolder"`
}

const gitConfigFile = "git-config.json"

// SaveGitConfig saves Git configuration to a file
func SaveGitConfig(config GitConfig) error {
	configDir := filepath.Join(lib.CONFIGS_DIR, "git")
	if err := utils.EnsureDirectoryExists(configDir); err != nil {
		return fmt.Errorf("failed to create git config directory: %w", err)
	}

	configPath := filepath.Join(configDir, gitConfigFile)
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal git config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write git config file: %w", err)
	}

	return nil
}

// GetGitConfig loads Git configuration from a file
func GetGitConfig() (GitConfig, error) {
	var config GitConfig
	configPath := filepath.Join(lib.CONFIGS_DIR, "git", gitConfigFile)

	// If the file doesn't exist, return default config
	if !utils.FileExists(configPath) {
		return GitConfig{
			Enabled:    false,
			Branch:     "main",
			AuthType:   "none",
			CommitMsg:  "Update mock configurations",
			PullBefore: true,
		}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read git config file: %w", err)
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return config, fmt.Errorf("failed to parse git config: %w", err)
	}

	return config, nil
}
