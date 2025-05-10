package traefik

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"

	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/utils"
)

// TraefikStaticConfig represents the structure of the Traefik static configuration
type TraefikStaticConfig struct {
	EntryPoints map[string]EntryPoint `yaml:"entryPoints"`
	Providers   Providers             `yaml:"providers"`
	API         API                   `yaml:"api"`
}

// EntryPoint represents a Traefik entry point
type EntryPoint struct {
	Address string `yaml:"address"`
}

// Providers represents Traefik providers
type Providers struct {
	File File `yaml:"file"`
}

// File represents Traefik file provider
type File struct {
	Directory string `yaml:"directory"`
	Watch     bool   `yaml:"watch"`
}

// API represents Traefik API configuration
type API struct {
	Insecure  bool `yaml:"insecure"`
	Dashboard bool `yaml:"dashboard"`
}

// GenerateStaticTraefikConfig generates the Traefik static configuration
func GenerateStaticTraefikConfig() error {
	// Create the directory for Traefik configuration if it doesn't exist
	traefikDir := filepath.Dir(lib.TRAEFIK_STATIC_CONFIG_PATH)
	if err := utils.EnsureDirectoryExists(traefikDir); err != nil {
		return fmt.Errorf("failed to create Traefik configuration directory: %w", err)
	}

	// Create the configuration structure
	config := TraefikStaticConfig{
		EntryPoints: map[string]EntryPoint{
			"web": {
				Address: ":80",
			},
		},
		Providers: Providers{
			File: File{
				Directory: filepath.Dir(lib.TRAEFIK_DYNAMIC_CONFIG_PATH),
				Watch:     true,
			},
		},
		API: API{
			Insecure:  true,
			Dashboard: true,
		},
	}

	// Convert to YAML
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal Traefik configuration: %w", err)
	}

	// Write to file
	err = os.WriteFile(lib.TRAEFIK_STATIC_CONFIG_PATH, data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write Traefik configuration: %w", err)
	}

	log.Printf("Generated Traefik static configuration at: %s", lib.TRAEFIK_STATIC_CONFIG_PATH)
	return nil
}
