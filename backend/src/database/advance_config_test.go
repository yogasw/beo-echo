package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvanceConfigProject_Validate(t *testing.T) {
	t.Run("Valid timeout within range", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: 5000, // 5 seconds
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid timeout at minimum boundary", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: 0, // 0 milliseconds (minimum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid timeout at maximum boundary", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: 120000, // 2 minutes (maximum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid timeout with 1 second", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: 1000, // 1 second
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid negative timeout", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: -1000, // Negative timeout
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout cannot be negative")
	})

	t.Run("Invalid timeout above maximum", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: 130000, // Above 2 minutes
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout cannot exceed 120000ms")
	})
}

func TestAdvanceConfigEndpoint_Validate(t *testing.T) {
	t.Run("Valid timeout within range", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: 10000, // 10 seconds
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid timeout at minimum boundary", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: 0, // 0 milliseconds (minimum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid timeout at maximum boundary", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: 120000, // 2 minutes (maximum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid timeout with 1 second", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: 1000, // 1 second
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid negative timeout", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: -500, // Negative timeout
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout cannot be negative")
	})

	t.Run("Invalid timeout above maximum", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: 130000, // Above 2 minutes
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "timeout cannot exceed 120000ms")
	})
}

func TestParseProjectAdvanceConfig(t *testing.T) {
	t.Run("Valid JSON with timeout", func(t *testing.T) {
		configJSON := `{"timeout": 5000}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 5000, config.Timeout)
	})

	t.Run("Valid empty JSON", func(t *testing.T) {
		configJSON := `{}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.Timeout)
	})

	t.Run("Valid empty string", func(t *testing.T) {
		configJSON := ""

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.Timeout)
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		configJSON := `{"timeout": 5000,}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "invalid JSON format")
	})

	t.Run("Invalid timeout value with negative value", func(t *testing.T) {
		configJSON := `{"timeout": -1}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "timeout cannot be negative")
	})

	t.Run("Invalid timeout too high", func(t *testing.T) {
		configJSON := `{"timeout": 400000}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "timeout cannot exceed 120000ms (2 minutes)")
	})

	t.Run("Valid JSON with extra fields ignored", func(t *testing.T) {
		configJSON := `{"timeout": 3000, "unknown_field": "value"}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 3000, config.Timeout)
	})
}

func TestParseEndpointAdvanceConfig(t *testing.T) {
	t.Run("Valid JSON with timeout", func(t *testing.T) {
		configJSON := `{"timeout": 15000}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 15000, config.Timeout)
	})

	t.Run("Valid empty JSON", func(t *testing.T) {
		configJSON := `{}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.Timeout)
	})

	t.Run("Valid empty string", func(t *testing.T) {
		configJSON := ""

		config, err := ParseEndpointAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.Timeout)
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		configJSON := `{"timeout":}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "invalid JSON format")
	})

	t.Run("Invalid timeout value", func(t *testing.T) {
		configJSON := `{"timeout": -1}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "timeout cannot be negative")
	})
}

func TestAdvanceConfigProject_ToJSON(t *testing.T) {
	t.Run("Config with timeout to JSON", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: 5000,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"timeout":5000}`, jsonStr)
	})

	t.Run("Empty config to JSON", func(t *testing.T) {
		config := &AdvanceConfigProject{
			Timeout: 0,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, "", jsonStr)
	})
}

func TestAdvanceConfigEndpoint_ToJSON(t *testing.T) {
	t.Run("Config with timeout to JSON", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: 8000,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"timeout":8000}`, jsonStr)
	})

	t.Run("Empty config to JSON", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			Timeout: 0,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, "", jsonStr)
	})
}
