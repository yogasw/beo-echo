package cmd

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/server"
	"mockoon-control-panel/backend_new/src/utils"
)

var apiPort string

// apiCmd represents the api command
var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the API server only",
	Long: `Starts the HTTP server for Mockoon Control Panel API without additional services.
This provides only the REST API endpoints for managing mock instances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runApiServer()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)

	// Add flags specific to the API command
	apiCmd.Flags().StringVarP(&apiPort, "port", "p", "", "Port to run the API server on (overrides env setting)")
}

func runApiServer() error {
	log.Println("🔧 Initializing API-only server mode...")

	// Load environment variables from .env file
	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		log.Println("⚠️  Warning: .env file not found or could not be loaded")
	} else {
		log.Println("✅ Environment variables loaded")
	}

	// Override port if specified as flag
	if apiPort != "" {
		os.Setenv("SERVER_PORT", apiPort)
		lib.SERVER_PORT = apiPort
		log.Printf("ℹ️  Port overridden to: %s", apiPort)
	}

	// Setup required directories
	log.Println("🔧 Setting up required directories...")
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		log.Printf("❌ Failed to create required directories: %v", err)
		return err
	}
	log.Println("✅ Required directories created")

	// Check if mockoon CLI is available
	log.Println("🔍 Checking for Mockoon CLI...")
	mockoonAvailable, err := utils.CheckMockoonCli()
	if err != nil || !mockoonAvailable {
		log.Printf("❌ Mockoon CLI not available: %v", err)
		return err
	}
	log.Println("✅ Mockoon CLI found")

	// Setup database connection without syncing or generating configs
	log.Println("🔧 Setting up database connection...")
	if err := database.CheckAndHandlePrisma(); err != nil {
		log.Printf("❌ Database setup failed: %v", err)
		return err
	}
	log.Println("✅ Database connected")

	log.Println("ℹ️  Running in API-only mode (no git sync, no traefik config)")
	log.Println("🚀 All systems initialized, starting API server...")

	// Start the server
	return server.StartServer()
}
