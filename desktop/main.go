package main

import (
	"context"
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"

	// Backend imports
	"beo-echo/backend/src"
	"beo-echo/backend/src/database"
	handlerLogs "beo-echo/backend/src/logs/handlers"
	systemConfig "beo-echo/backend/src/systemConfigs"
	"beo-echo/backend/src/utils"
)

//go:embed all:frontend
var assets embed.FS

// App struct
type App struct {
	ctx           context.Context
	backendCancel context.CancelFunc
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// OnStartup is called when the app starts up. It's used to setup the application context
func (a *App) OnStartup(ctx context.Context) {
	a.ctx = ctx

	// Initialize application directories and configurations
	if err := setupDesktopEnvironment(); err != nil {
		log.Printf("Failed to setup desktop environment: %v", err)
		return
	}

	// Initialize and start backend server in a goroutine
	backendCtx, cancel := context.WithCancel(ctx)
	a.backendCancel = cancel

	go func() {
		if err := startBackendServer(backendCtx); err != nil {
			log.Printf("Backend server error: %v", err)
		}
	}()

	// Give backend a moment to start
	time.Sleep(1 * time.Second)
	log.Println("✅ Desktop application initialized successfully")
}

// OnDomReady is called after front-end resources have been loaded
func (a *App) OnDomReady(ctx context.Context) {
	// Add your action here
}

// OnBeforeClose is called when the application is about to quit,
// either by clicking the window close button or calling runtime.Quit.
// Returning true will cause the application to continue,
// false will continue shutdown as normal.
func (a *App) OnBeforeClose(ctx context.Context) (prevent bool) {
	return false
}

// OnShutdown is called when the application is shutting down
func (a *App) OnShutdown(ctx context.Context) {
	// Shutdown backend server gracefully
	if a.backendCancel != nil {
		log.Println("🔄 Shutting down backend server...")
		a.backendCancel()
	}
	log.Println("✅ Desktop application shutdown complete")
}

// setupDesktopEnvironment initializes the desktop application environment
func setupDesktopEnvironment() error {
	log.Println("🔧 Setting up desktop environment...")

	// Get the executable directory for desktop app
	execPath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	_ = filepath.Dir(execPath) // execDir not used for now

	// Setup application folders in user's home directory for desktop app
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	appDataDir := filepath.Join(homeDir, ".beoecho")
	if err := os.MkdirAll(appDataDir, 0755); err != nil {
		return fmt.Errorf("failed to create app data directory: %w", err)
	}

	// Setup required folders
	folders := []string{
		filepath.Join(appDataDir, "configs"),
		filepath.Join(appDataDir, "logs"),
		filepath.Join(appDataDir, "uploads"),
		filepath.Join(appDataDir, "db"),
	}

	for _, folder := range folders {
		if err := os.MkdirAll(folder, 0755); err != nil {
			return fmt.Errorf("failed to setup folder %s: %w", folder, err)
		}
	}

	// Change working directory to app data directory
	if err := os.Chdir(appDataDir); err != nil {
		return fmt.Errorf("failed to change working directory: %w", err)
	}

	log.Printf("✅ Desktop environment initialized in: %s", appDataDir)
	return nil
}

// startBackendServer initializes and starts the backend server
func startBackendServer(ctx context.Context) error {
	log.Println("🚀 Starting backend server...")

	// Setup required directories using backend utilities
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		return fmt.Errorf("failed to setup required folders: %w", err)
	}

	// Setup database connection
	if err := database.CheckAndHandle(); err != nil {
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Initialize default system configuration
	if err := systemConfig.InitializeDefaultConfig(); err != nil {
		log.Printf("⚠️  Warning: Failed to initialize default system configuration: %v", err)
	}

	// Initialize log service
	handlerLogs.InitLogService()

	// Create a channel to monitor server shutdown
	serverDone := make(chan error, 1)

	// Start the server in a separate goroutine
	go func() {
		defer close(serverDone)
		if err := src.StartServer(); err != nil {
			serverDone <- fmt.Errorf("server startup failed: %w", err)
		}
	}()

	// Wait for either context cancellation or server error
	select {
	case <-ctx.Done():
		log.Println("🔄 Backend server shutdown requested")
		return ctx.Err()
	case err := <-serverDone:
		if err != nil {
			return fmt.Errorf("backend server error: %w", err)
		}
		return nil
	}
}

func main() {
	// Create an instance of the app structure
	app := NewApp()

	// Create application with options
	err := wails.Run(&options.App{
		Title:             "BeoEcho Desktop",
		Width:             1200,
		Height:            800,
		MinWidth:          800,
		MinHeight:         600,
		MaxWidth:          1920,
		MaxHeight:         1080,
		DisableResize:     false,
		Fullscreen:        false,
		Frameless:         false,
		StartHidden:       false,
		HideWindowOnClose: false,
		BackgroundColour:  &options.RGBA{R: 31, G: 41, B: 55, A: 1}, // Tailwind gray-800
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		OnStartup:        app.OnStartup,
		OnDomReady:       app.OnDomReady,
		OnBeforeClose:    app.OnBeforeClose,
		OnShutdown:       app.OnShutdown,
		WindowStartState: options.Normal,
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
