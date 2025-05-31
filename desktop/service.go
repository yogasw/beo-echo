package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	// Backend imports
	"beo-echo/backend/src"
	"beo-echo/backend/src/database"
	"beo-echo/backend/src/lib"
	handlerLogs "beo-echo/backend/src/logs/handlers"
	systemConfig "beo-echo/backend/src/systemConfigs"
	"beo-echo/backend/src/utils"

	"github.com/wailsapp/wails/v3/pkg/application"
	// Wails imports
)

// NewBackendService creates a new BackendService struct
func NewBackendService() *BackendService {
	return &BackendService{}
}

// BackendService struct
type BackendService struct {
	ctx           context.Context
	backendCancel context.CancelFunc
}

func (b *BackendService) ServiceName() string {
	return "BeoEcho Backend Service"
}

// ServiceStartup is called when the app starts up. It's used to setup the application context
func (b *BackendService) OnStartup(ctx context.Context, options application.ServiceOptions) error {
	b.ctx = ctx

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
	b.backendCancel = cancel

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

	return nil
}

// OnShutdown is called when the application is shutting down
func (b *BackendService) OnShutdown() error {
	log.Println("🔄 OnShutdown called...")

	// Shutdown backend server gracefully
	if b.backendCancel != nil {
		log.Println("🔄 Shutting down backend server...")
		b.backendCancel()
	}
	log.Println("✅ Desktop application shutdown complete")

	return nil
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
