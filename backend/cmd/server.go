package cmd

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"beo-echo/backend/src"
	"beo-echo/backend/src/database"
	systemConfig "beo-echo/backend/src/systemConfigs"
	"beo-echo/backend/src/utils"
)

var port string
var hostname string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the BeoEcho server",
	Long: `Starts the HTTP server for BeoEcho.
This provides API endpoints for managing mock instances.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runServer()
	},
}

func init() {
	// Add flags specific to the server command
	serverCmd.Flags().StringVarP(&port, "port", "p", "", "Port to run the server on (overrides env setting)")
	serverCmd.Flags().StringVarP(&hostname, "host", "H", "", "Hostname to bind the server to (overrides env setting)")
}

func runServer() error {
	log.Println("🔧 Initializing BeoEcho server...")

	// Load environment variables from .env file
	if err := godotenv.Load(filepath.Join("..", ".env")); err != nil {
		log.Println("⚠️  Warning: .env file not found or could not be loaded")
	} else {
		log.Println("✅ Environment variables loaded")
	}

	// Setup required directories
	log.Println("🔧 Setting up required directories...")
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		log.Printf("❌ Failed to create required directories: %v", err)
		return err
	}
	log.Println("✅ Required directories created")

	// Setup database connection
	log.Println("🔧 Setting up database connection...")
	if err := database.CheckAndHandle(); err != nil {
		log.Printf("❌ Database setup failed: %v", err)
		return err
	}
	log.Println("✅ Database connected")

	log.Println("🚀 All systems initialized, starting HTTP server...")

	// Initialize default system configuration
	// Initialize default system configuration
	if err := systemConfig.InitializeDefaultConfig(); err != nil {
		log.Printf("❌ Failed to initialize default system configuration: %v", err)
	}

	// Start the server (this will block until the server is stopped)
	return src.StartServer()
}
