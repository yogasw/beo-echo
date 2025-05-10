package cmd

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/utils"
)

// generateCmd represents the generate command for generating necessary configuration files
var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate configuration files",
	Long:  `Generate configuration files for Traefik and Mockoon Control Panel`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Load environment variables from .env file
		if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
			log.Println("Warning: .env file not found or could not be loaded")
		}

		// Ensure required folders and environment variables
		if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
			log.Fatalf("Failed to ensure required folders and environment: %v", err)
			return err
		}

		// Generate static Traefik configuration
		if err := traefik.GenerateStaticTraefikConfig(); err != nil {
			log.Fatalf("Failed to generate static Traefik configuration: %v", err)
			return err
		}

		// Generate dynamic Traefik configuration
		if err := traefik.GenerateDynamicTraefikConfig(); err != nil {
			log.Fatalf("Failed to generate dynamic Traefik configuration: %v", err)
			return err
		}

		log.Println("Configuration generation completed successfully!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(generateCmd)
}
