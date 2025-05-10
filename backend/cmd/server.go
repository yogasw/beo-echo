package cmd

import (
	"log"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"

	"mockoon-control-panel/backend_new/src"
	"mockoon-control-panel/backend_new/src/database"
	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/utils"
)

var port string
var hostname string

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start the Mockoon Control Panel server",
	Long: `Starts the HTTP server for Mockoon Control Panel.
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
	log.Println("🔧 Initializing Mockoon Control Panel server...")

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
	if err := database.CheckAndHandlePrisma(); err != nil {
		log.Printf("❌ Database setup failed: %v", err)
		return err
	}
	log.Println("✅ Database connected")

	// Generate Traefik config
	// log.Println("🔧 Generating Traefik configuration...")
	// if err := traefik.GenerateDynamicTraefikConfig(); err != nil {
	// 	log.Printf("❌ Failed to generate dynamic Traefik config: %v", err)
	// 	return err
	// }
	// log.Println("✅ Dynamic Traefik configuration generated")

	if err := traefik.GenerateStaticTraefikConfig(); err != nil {
		log.Printf("❌ Failed to generate static Traefik config: %v", err)
		return err
	}
	log.Println("✅ Static Traefik configuration generated")

	log.Println("🚀 All systems initialized, starting HTTP server...")

	// Start the server (this will block until the server is stopped)
	return src.StartServer()
}
