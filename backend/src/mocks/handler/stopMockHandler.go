package handler

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/traefik"
	"mockoon-control-panel/backend_new/src/types"
)

// StopMockHandler handles requests to stop a running mock server
func StopMockHandler(c *gin.Context) {
	var req types.StopMockRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format: " + err.Error()})
		return
	}

	log.Printf("Stopping mock server on port: %d", req.Port)

	// Check if the mock instance exists
	instance, exists := repositories.MockInstanceRepo.Get(req.Port)
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "No mock server running on port " + strconv.Itoa(req.Port),
		})
		return
	}

	// Kill the process
	var err error
	if instance.Process != nil {
		err = instance.Process.Kill()
	} else if instance.Cmd != nil && instance.Cmd.Process != nil {
		err = instance.Cmd.Process.Kill()
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to stop mock server: " + err.Error(),
		})
		return
	}

	// Remove the instance from the repository
	repositories.MockInstanceRepo.Remove(req.Port)

	// Update Traefik configuration
	if err := traefik.GenerateDynamicTraefikConfig(); err != nil {
		log.Printf("Warning: Failed to update Traefik configuration: %v", err)
	}

	// Return success response
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"port":   req.Port,
			"status": "stopped",
		},
	})
}
