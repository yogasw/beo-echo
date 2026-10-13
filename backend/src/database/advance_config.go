package database

import (
	"encoding/json"
	"errors"
)

// CorsMode selects how a project answers CORS preflight requests on its mock routes.
type CorsMode string

const (
	CorsModeStatic  CorsMode = "static"  // Fixed Access-Control-Allow-Headers list (legacy behaviour)
	CorsModeReflect CorsMode = "reflect" // Echo Access-Control-Request-Headers back
	CorsModeCustom  CorsMode = "custom"  // Use CorsConfig.AllowHeaders
)

// maxCorsMaxAge caps Access-Control-Max-Age (24h).
const maxCorsMaxAge = 86400

// CorsConfig is the per-project CORS configuration for mock routes.
// An empty Mode behaves as CorsModeStatic so existing projects keep their behaviour.
type CorsConfig struct {
	Mode         CorsMode `json:"mode,omitempty"`
	AllowHeaders []string `json:"allowHeaders,omitempty"` // Used when Mode is custom
	AllowOrigins []string `json:"allowOrigins,omitempty"` // Defaults to the server-wide CORS_ORIGIN
	MaxAge       int      `json:"maxAge,omitempty"`       // Access-Control-Max-Age in seconds (0 = omit)
}

// Validate validates the CORS configuration
func (c *CorsConfig) Validate() error {
	switch c.Mode {
	case "", CorsModeStatic, CorsModeReflect, CorsModeCustom:
	default:
		return errors.New("cors.mode must be one of: static, reflect, custom")
	}
	if c.MaxAge < 0 || c.MaxAge > maxCorsMaxAge {
		return errors.New("cors.maxAge must be between 0 and 86400 seconds")
	}
	return nil
}

// SequenceOnEnd decides what a "sequence" endpoint serves once every response was used.
type SequenceOnEnd string

const (
	SequenceStayLast SequenceOnEnd = "stay_last" // Keep serving the last response (default)
	SequenceLoop     SequenceOnEnd = "loop"      // Start over from the first response
)

// SequenceConfig configures the "sequence" response mode of an endpoint.
type SequenceConfig struct {
	OnEnd SequenceOnEnd `json:"onEnd,omitempty"`
}

// AdvanceConfigProject defines advance configuration structure for projects
type AdvanceConfigProject struct {
	DelayMs int         `json:"delayMs,omitempty"` // Response delay in milliseconds (0-120000)
	Cors    *CorsConfig `json:"cors,omitempty"`    // CORS preflight behaviour for mock routes
}

// AdvanceConfigEndpoint defines advance configuration structure for endpoints
type AdvanceConfigEndpoint struct {
	DelayMs  int             `json:"delayMs,omitempty"`  // Response delay in milliseconds (0-120000)
	Sequence *SequenceConfig `json:"sequence,omitempty"` // Options for response_mode "sequence"
}

// DefaultProjectAdvanceConfig is the advance config given to newly created projects.
// Existing projects keep an empty config, which means static CORS.
func DefaultProjectAdvanceConfig() string {
	return `{"cors":{"mode":"reflect"}}`
}

// Validate validates the project advance configuration
func (a *AdvanceConfigProject) Validate() error {
	if a.DelayMs < 0 {
		return errors.New("delayMs cannot be negative")
	}
	if a.DelayMs > 120000 {
		return errors.New("delayMs cannot exceed 120000ms (2 minutes)")
	}
	if a.Cors != nil {
		return a.Cors.Validate()
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
	if a.Sequence != nil {
		switch a.Sequence.OnEnd {
		case "", SequenceStayLast, SequenceLoop:
		default:
			return errors.New("sequence.onEnd must be one of: stay_last, loop")
		}
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
	if a.DelayMs == 0 && a.Cors == nil {
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
	if a.DelayMs == 0 && a.Sequence == nil {
		return "", nil
	}

	data, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
