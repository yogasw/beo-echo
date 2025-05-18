package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "BeoEcho",
	Short: "BeoEcho - Backend server",
	Long: `BeoEcho is a backend server that provides API endpoints
for managing BeoEcho configurations and instances.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Starting BeoEcho server...")
		fmt.Println("Use Ctrl+C to stop the server")
		// By default, run the server
		if err := serverCmd.RunE(cmd, args); err != nil {
			fmt.Println("Error running server:", err)
			os.Exit(1)
		}
	},
}

// Initialize cobra configuration
func init() {
	// Add persistent flags that are global to the application
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is ../configs/app.yaml)")

	// Add subcommands
	rootCmd.AddCommand(serverCmd)

	// Set working directory for consistent path resolution
	dir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
	os.Chdir(dir)
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}
