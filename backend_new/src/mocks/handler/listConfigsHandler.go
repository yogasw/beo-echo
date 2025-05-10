package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"mockoon-control-panel/backend_new/src/lib"
	"mockoon-control-panel/backend_new/src/mocks/repositories"
)

// ListConfigsHandler handles requests to list available configuration files
func ListConfigsHandler(c *gin.Context) {
	configs, err := repositories.FileRepo.ListConfigs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to list configuration files: " + err.Error(),
		})
		return
	}

	// Update running status and generate URLs for each config
	for i := range configs {
		running, port := repositories.MockInstanceRepo.IsRunningWithConfig(configs[i].ConfigFile)
		configs[i].InUse = running
		if running && configs[i].Port == 0 {
			// Only override port if it wasn't already set and the instance is running
			configs[i].Port = port
		}

		// Only generate URL if there's a port (either from the file or from running status)
		if configs[i].Port > 0 {
			// Generate URL just like in TypeScript implementation
			url := ""
			if lib.PROXY_MODE {
				if lib.PROXY_BASE_URL != "" {
					url = lib.PROXY_BASE_URL + "/" + strconv.Itoa(configs[i].Port)
				} else {
					host := c.Request.Host
					protocol := "http"
					if c.Request.TLS != nil {
						protocol = "https"
					}
					url = protocol + "://" + host + "/" + strconv.Itoa(configs[i].Port)
				}
			} else {
				url = "http://" + lib.SERVER_HOSTNAME + ":" + strconv.Itoa(configs[i].Port)
			}
			configs[i].URL = url
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": configs,
	})
}
