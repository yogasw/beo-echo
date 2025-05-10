package types

import (
	"os"
	"os/exec"
	"time"
)

// ApiResponse represents a standard API response structure
type ApiResponse struct {
	Error string      `json:"error,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

// StartMockRequest represents a request to start a mock server
type StartMockRequest struct {
	Port       int    `json:"port" binding:"required"`
	ConfigFile string `json:"configFile" binding:"required"`
	UUID       string `json:"uuid,omitempty"`
}

// StartMockResponse represents a response from starting a mock server
type StartMockResponse struct {
	Port       int       `json:"port"`
	ConfigFile string    `json:"configFile"`
	StartTime  time.Time `json:"startTime"`
	Status     string    `json:"status"`
}

// StopMockRequest represents a request to stop a mock server
type StopMockRequest struct {
	Port int `json:"port" binding:"required"`
}

// MockStatusResponse represents status information about mock instances
type MockStatusResponse struct {
	Running   []MockInstance `json:"running"`
	Available []string       `json:"available"`
}

// MockInstance represents a running mock instance
type MockInstance struct {
	Process    *os.Process `json:"-"`
	Cmd        *exec.Cmd   `json:"-"`
	ConfigFile string      `json:"configFile"`
	Port       int         `json:"port"`
	StartTime  time.Time   `json:"startTime"`
	LogFile    *os.File    `json:"-"`
	UUID       string      `json:"uuid,omitempty"`
}

// ConfigFile represents a configuration file
type ConfigFile struct {
	UUID       string    `json:"uuid,omitempty"`
	Name       string    `json:"name"`
	ConfigFile string    `json:"configFile"`
	Port       int       `json:"port,omitempty"`
	Size       string    `json:"size"`
	Modified   time.Time `json:"modified"`
	InUse      bool      `json:"inUse"`
	URL        string    `json:"url,omitempty"`
}
