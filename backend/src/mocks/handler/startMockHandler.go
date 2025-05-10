package handler

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/types"
	"mockoon-control-panel/backend_new/src/utils"
)

// UseUnsafePort determines if we should allow ports outside the safe range
const UseUnsafePort = true

// StartMockHandler handles requests to start a mock server
func StartMockHandler(c *gin.Context) {
	var req types.StartMockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	log.Printf("Starting mock server with config: %s on port: %d", req.ConfigFile, req.Port)

	// Validate port if not using unsafe port
	if !UseUnsafePort && !utils.IsPortSafe(req.Port) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid port. Port must be between 9001 and 9999.",
		})
		return
	}

	// Check if port is already in use
	if utils.IsPortInUse(req.Port) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Port " + strconv.Itoa(req.Port) + " is already in use.",
		})
		return
	}

	// Check if config file exists
	configPath := repositories.FileRepo.GetConfigPath(req.ConfigFile)
	if !repositories.FileRepo.ConfigExists(req.ConfigFile) {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Configuration file not found",
		})
		return
	}

	// Create log directory if it doesn't exist
	if err := utils.EnsureDirectoryExists(lib.LOGS_DIR); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create log directory: " + err.Error(),
		})
		return
	}

	// Create log file
	logFilePath := filepath.Join(lib.LOGS_DIR, "mock-"+strconv.Itoa(req.Port)+".log")
	logFile, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create log file: " + err.Error(),
		})
		return
	}

	// Start mockoon-cli process
	cmd := exec.Command("mockoon-cli",
		"start",
		"--data", configPath,
		"--port", strconv.Itoa(req.Port),
	)

	// Redirect stdout and stderr to log file
	cmd.Stdout = logFile
	cmd.Stderr = logFile

	// Start the process
	if err := cmd.Start(); err != nil {
		logFile.Close()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to start mock server: " + err.Error(),
		})
		return
	}

	// Record the instance
	now := time.Now()
	repositories.MockInstanceRepo.Add(req.Port, types.MockInstance{
		Process:    cmd.Process,
		Cmd:        cmd,
		ConfigFile: req.ConfigFile,
		Port:       req.Port,
		StartTime:  now,
		LogFile:    logFile,
		UUID:       req.UUID,
	})

	// Update Traefik configuration
	if err := traefik.GenerateDynamicTraefikConfig(); err != nil {
		log.Printf("Warning: Failed to update Traefik configuration: %v", err)
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"data": types.StartMockResponse{
			Port:       req.Port,
			ConfigFile: req.ConfigFile,
			StartTime:  now,
			Status:     "running",
		},
	})
}
