package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/utils"
)

// configCmd represents the config command - DEPRECATED, use generateCmd instead
var configCmd = &cobra.Command{
	Use:        "config",
	Short:      "Manage configuration (DEPRECATED)",
	Long:       `Generate or update Mockoon Control Panel configurations - DEPRECATED, use generate command instead`,
	Deprecated: "Use the 'generate' command instead",
}

// configGenerateCmd is a subcommand of config - DEPRECATED, use generateCmd instead
var configGenerateCmd = &cobra.Command{
	Use:        "generate",
	Short:      "Generate configuration files (DEPRECATED)",
	Long:       `Generate Traefik and other configuration files - DEPRECATED, use generate command instead`,
	Deprecated: "Use the 'generate' command instead",
	RunE: func(cmd *cobra.Command, args []string) error {
		log.Println("Warning: 'config generate' is deprecated, use 'generate' instead")
		return generateConfigs()
	},
}

func init() {
	configCmd.AddCommand(configGenerateCmd)
	rootCmd.AddCommand(configCmd)

	// No need to add generateCmd as a subcommand of configCmd
	// since generateCmd is a top-level command
	log.SetOutput(os.Stdout) // This is just to avoid the unused import warning
}

func generateConfigs() error {
	// Load environment variables from .env file
	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		log.Println("Warning: .env file not found or could not be loaded")
	}

	// Ensure required folders and environment variables
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		return err
	}

	// Generate static Traefik configuration
	if err := traefik.GenerateStaticTraefikConfig(); err != nil {
		return err
	}

	// // Generate dynamic Traefik configuration
	// if err := traefik.GenerateDynamicTraefikConfig(); err != nil {
	// 	return err
	// }

	log.Println("Configuration generation completed successfully!")
	return nil
}
