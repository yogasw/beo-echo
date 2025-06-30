package database

import (
	"encoding/json"
	"errors"
)

// AdvanceConfigProject defines advance configuration structure for projects
type AdvanceConfigProject struct {
	DelayMs int `json:"delayMs,omitempty"` // Response delay in milliseconds (0-120000)
}

// AdvanceConfigEndpoint defines advance configuration structure for endpoints
type AdvanceConfigEndpoint struct {
	DelayMs int `json:"delayMs,omitempty"` // Response delay in milliseconds (0-120000)
}

// Validate validates the project advance configuration
func (a *AdvanceConfigProject) Validate() error {
	if a.DelayMs < 0 {
		return errors.New("delayMs cannot be negative")
	}
	if a.DelayMs > 120000 {
		return errors.New("delayMs cannot exceed 120000ms (2 minutes)")
	}
	return nil
}

// Validate validates the endpoint advance configuration
func (a *AdvanceConfigEndpoint) Validate() error {
	if a.DelayMs < 0 {
		return errors.New("delayMs cannot be negative")
	}
	if a.DelayMs > 120000 {
		return errors.New("delayMs cannot exceed 120000ms (2 minutes)")
	}
	return nil
}

// ParseProjectAdvanceConfig parses JSON string to AdvanceConfigProject struct
func ParseProjectAdvanceConfig(configJSON string) (*AdvanceConfigProject, error) {
	if configJSON == "" {
		return &AdvanceConfigProject{}, nil
	}

	var config AdvanceConfigProject
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, errors.New("invalid JSON format in advance_config")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// ParseEndpointAdvanceConfig parses JSON string to AdvanceConfigEndpoint struct
func ParseEndpointAdvanceConfig(configJSON string) (*AdvanceConfigEndpoint, error) {
	if configJSON == "" {
		return &AdvanceConfigEndpoint{}, nil
	}

	var config AdvanceConfigEndpoint
	if err := json.Unmarshal([]byte(configJSON), &config); err != nil {
		return nil, errors.New("invalid JSON format in advance_config")
	}

	if err := config.Validate(); err != nil {
		return nil, err
	}

	return &config, nil
}

// ToJSON converts AdvanceConfigProject to JSON string
func (a *AdvanceConfigProject) ToJSON() (string, error) {
	if a.DelayMs == 0 {
		return "", nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// ToJSON converts AdvanceConfigEndpoint to JSON string
func (a *AdvanceConfigEndpoint) ToJSON() (string, error) {
	if a.DelayMs == 0 {
		return "", nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
