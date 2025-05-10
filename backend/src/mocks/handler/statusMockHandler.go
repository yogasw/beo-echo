package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/mocks/repositories"
	"mockoon-control-panel/backend_new/src/types"
)

// StatusMockHandler handles requests to get the status of running mock servers
func StatusMockHandler(c *gin.Context) {
	// Get all running mock instances
	runningInstances := repositories.MockInstanceRepo.GetAll()

	// Get a list of available config files
	configs, err := repositories.FileRepo.ListConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list configuration files: " + err.Error(),
		})
		return
	}

	// Create a list of available config file names
	availableConfigs := make([]string, 0, len(configs))
	for _, config := range configs {
		// Check if this config is running
		running, port := repositories.MockInstanceRepo.IsRunningWithConfig(config.Name)

		// Add the config to the available list
		availableConfigs = append(availableConfigs, config.Name)

		// Update the config's running status
		if running {
			for i := range configs {
				if configs[i].Name == config.Name {
					configs[i].InUse = true
					configs[i].Port = port
					break
				}
			}
		}
	}

	// Create running instances for response
	runningForResponse := make([]types.MockInstance, 0, len(runningInstances))
	for _, instance := range runningInstances {
		// Create a copy without sensitive fields
		runningForResponse = append(runningForResponse, types.MockInstance{
			ConfigFile: instance.ConfigFile,
			Port:       instance.Port,
			StartTime:  instance.StartTime,
			UUID:       instance.UUID,
		})
	}

	// Return the status response
	c.JSON(http.StatusOK, gin.H{
		"data": types.MockStatusResponse{
			Running:   runningForResponse,
			Available: availableConfigs,
		},
	})
}
