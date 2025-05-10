package main

import (
	"log"
	"os"
	"path/filepath"

	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/utils"
)

func main() {
	// Determine the path for the .env file
	envPath := filepath.Join("..", ".env")

	// Check if .env file exists and load it
	if _, err := os.Stat(envPath); err == nil {
		log.Println("Loading environment from .env file")
	} else {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	// Ensure required folders and environment variables
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		log.Fatalf("Failed to ensure required folders and environment: %v", err)
	}

	// Generate static Traefik configuration
	if err := traefik.GenerateStaticTraefikConfig(); err != nil {
		log.Fatalf("Failed to generate static Traefik configuration: %v", err)
	}

	// Generate dynamic Traefik configuration with firstInit=true to skip trying to load aliases
	// since the database might not be fully set up during initial configuration
	// if err := traefik.GenerateDynamicTraefikConfigWithInit(true); err != nil {
	// 	log.Fatalf("Failed to generate dynamic Traefik configuration: %v", err)
	// }

	log.Println("Configuration generation completed successfully!")
}
