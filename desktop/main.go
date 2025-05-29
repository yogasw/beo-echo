package main

import (
	"context"
	"embed"
	"fmt"
	"io"
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
	"beo-echo/backend/src/lib"
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

	log.Println("🔄 OnStartup called...")
	log.Printf("App startup context: %v", ctx)

	// Set desktop mode for backend to use proper paths
	log.Println("🔄 Setting desktop mode...")
	lib.SetDesktopMode(true)

	// Initialize application directories and configurations
	log.Println("🔄 Setting up desktop environment...")
	if err := setupDesktopEnvironment(); err != nil {
		log.Printf("❌ Failed to setup desktop environment: %v", err)
		// Don't return early - try to continue with backend startup
	}

	// Initialize and start backend server in a goroutine
	log.Println("🔄 Starting backend server...")
	backendCtx, cancel := context.WithCancel(ctx)
	a.backendCancel = cancel

	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ Backend server panic: %v", r)
			}
		}()

		if err := startBackendServer(backendCtx); err != nil {
			log.Printf("❌ Backend server error: %v", err)
		}
	}()

	// Give backend a moment to start
	time.Sleep(2 * time.Second)
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
	log.Printf("Executable path: %s", execPath)

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

	// No need to change working directory anymore since backend uses absolute paths
	log.Printf("✅ Desktop environment initialized in: %s", appDataDir)
	log.Printf("Current working directory remains: %s", getCurrentWorkingDir())

	return nil
}

// setupLogging configures logging to both console and file for better debugging
func setupLogging() {
	// Get user home directory for log file
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Warning: Could not get home directory for logging: %v", err)
		return
	}

	// Create logs directory
	logDir := filepath.Join(homeDir, ".beoecho", "logs")
	if err := os.MkdirAll(logDir, 0755); err != nil {
		log.Printf("Warning: Could not create log directory: %v", err)
		return
	}

	// Create log file with timestamp
	logFile := filepath.Join(logDir, fmt.Sprintf("desktop-%s.log", time.Now().Format("2006-01-02")))
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Warning: Could not create log file: %v", err)
		return
	}

	// Setup multi-writer to write to both file and console
	multiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	log.Printf("📝 Logging initialized. Log file: %s", logFile)
}

// getCurrentWorkingDir returns current working directory safely
func getCurrentWorkingDir() string {
	if wd, err := os.Getwd(); err == nil {
		return wd
	}
	return "unknown"
}

// getExecutablePath returns executable path safely
func getExecutablePath() string {
	if exec, err := os.Executable(); err == nil {
		return exec
	}
	return "unknown"
}

// startBackendServer initializes and starts the backend server
func startBackendServer(ctx context.Context) error {
	log.Println("🚀 Starting backend server...")

	// Log current environment for debugging
	log.Printf("Current working directory: %s", getCurrentWorkingDir())
	log.Printf("Backend CURRENT_DIR(): %s", lib.CURRENT_DIR())
	log.Printf("HOME: %s", os.Getenv("HOME"))
	log.Printf("PATH: %s", os.Getenv("PATH"))

	// Since we set desktop mode, the backend will automatically use user home/.beoecho
	// No need to change working directory manually anymore

	// Setup required directories using backend utilities
	log.Println("🔄 Ensuring required folders and environment...")
	if err := utils.EnsureRequiredFoldersAndEnv(); err != nil {
		log.Printf("❌ Failed to setup required folders: %v", err)
		return fmt.Errorf("failed to setup required folders: %w", err)
	}

	// Setup database connection
	log.Println("🔄 Setting up database connection...")
	if err := database.CheckAndHandle(); err != nil {
		log.Printf("❌ Failed to setup database: %v", err)
		return fmt.Errorf("failed to setup database: %w", err)
	}

	// Initialize default system configuration
	log.Println("🔄 Initializing system configuration...")
	if err := systemConfig.InitializeDefaultConfig(); err != nil {
		log.Printf("⚠️  Warning: Failed to initialize default system configuration: %v", err)
	}

	// Initialize log service
	log.Println("🔄 Initializing log service...")
	handlerLogs.InitLogService()

	// Create a channel to monitor server shutdown
	serverDone := make(chan error, 1)

	// Start the server in a separate goroutine
	go func() {
		defer close(serverDone)
		defer func() {
			if r := recover(); r != nil {
				log.Printf("❌ Server startup panic: %v", r)
				serverDone <- fmt.Errorf("server panic: %v", r)
			}
		}()

		log.Println("🔄 Starting HTTP server...")
		if err := src.StartServer(); err != nil {
			log.Printf("❌ Server startup error: %v", err)
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
			log.Printf("❌ Backend server error: %v", err)
			return fmt.Errorf("backend server error: %w", err)
		}
		log.Println("✅ Backend server started successfully")
		return nil
	}
}

func main() {
	// Setup logging to file for desktop app debugging
	setupLogging()

	log.Println("🚀 Starting BeoEcho Desktop Application...")
	log.Printf("Current working directory: %s", getCurrentWorkingDir())
	log.Printf("Executable path: %s", getExecutablePath())
	log.Printf("Environment PATH: %s", os.Getenv("PATH"))

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
		log.Printf("❌ Application error: %v", err)
		println("Error:", err.Error())
	}
}
