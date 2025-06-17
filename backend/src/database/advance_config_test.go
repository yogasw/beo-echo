package database

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAdvanceConfigProject_Validate(t *testing.T) {
	t.Run("Valid delayMs within range", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: 5000, // 5 seconds
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid delayMs at minimum boundary", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: 0, // 0 milliseconds (minimum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid delayMs at maximum boundary", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: 120000, // 2 minutes (maximum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid delayMs with 1 second", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: 1000, // 1 second
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid negative delayMs", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: -1000, // Negative delayMs
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delayMs cannot be negative")
	})

	t.Run("Invalid delayMs above maximum", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: 130000, // Above 2 minutes
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delayMs cannot exceed 120000ms")
	})
}

func TestAdvanceConfigEndpoint_Validate(t *testing.T) {
	t.Run("Valid delayMs within range", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: 10000, // 10 seconds
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid delayMs at minimum boundary", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: 0, // 0 milliseconds (minimum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid delayMs at maximum boundary", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: 120000, // 2 minutes (maximum)
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Valid delayMs with 1 second", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: 1000, // 1 second
		}

		err := config.Validate()
		assert.NoError(t, err)
	})

	t.Run("Invalid negative delayMs", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: -500, // Negative delayMs
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delayMs cannot be negative")
	})

	t.Run("Invalid delayMs above maximum", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: 130000, // Above 2 minutes
		}

		err := config.Validate()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "delayMs cannot exceed 120000ms")
	})
}

func TestParseProjectAdvanceConfig(t *testing.T) {
	t.Run("Valid JSON with delayMs", func(t *testing.T) {
		configJSON := `{"delayMs": 5000}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 5000, config.DelayMs)
	})

	t.Run("Valid empty JSON", func(t *testing.T) {
		configJSON := `{}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.DelayMs)
	})

	t.Run("Valid empty string", func(t *testing.T) {
		configJSON := ""

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.DelayMs)
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		configJSON := `{"delayMs": 5000,}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "invalid JSON format")
	})

	t.Run("Invalid delayMs value with negative value", func(t *testing.T) {
		configJSON := `{"delayMs": -1}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "delayMs cannot be negative")
	})

	t.Run("Invalid delayMs too high", func(t *testing.T) {
		configJSON := `{"delayMs": 400000}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "delayMs cannot exceed 120000ms (2 minutes)")
	})

	t.Run("Valid JSON with extra fields ignored", func(t *testing.T) {
		configJSON := `{"delayMs": 3000, "unknown_field": "value"}`

		config, err := ParseProjectAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 3000, config.DelayMs)
	})
}

func TestParseEndpointAdvanceConfig(t *testing.T) {
	t.Run("Valid JSON with delayMs", func(t *testing.T) {
		configJSON := `{"delayMs": 15000}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 15000, config.DelayMs)
	})

	t.Run("Valid empty JSON", func(t *testing.T) {
		configJSON := `{}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.DelayMs)
	})

	t.Run("Valid empty string", func(t *testing.T) {
		configJSON := ""

		config, err := ParseEndpointAdvanceConfig(configJSON)
		require.NoError(t, err)
		assert.Equal(t, 0, config.DelayMs)
	})

	t.Run("Invalid JSON format", func(t *testing.T) {
		configJSON := `{"delayMs":}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "invalid JSON format")
	})

	t.Run("Invalid delayMs value", func(t *testing.T) {
		configJSON := `{"delayMs": -1}`

		config, err := ParseEndpointAdvanceConfig(configJSON)
		assert.Error(t, err)
		assert.Nil(t, config)
		assert.Contains(t, err.Error(), "delayMs cannot be negative")
	})
}

func TestAdvanceConfigProject_ToJSON(t *testing.T) {
	t.Run("Config with delayMs to JSON", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: 5000,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"delayMs":5000}`, jsonStr)
	})

	t.Run("Empty config to JSON", func(t *testing.T) {
		config := &AdvanceConfigProject{
			DelayMs: 0,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, "", jsonStr)
	})
}

func TestAdvanceConfigEndpoint_ToJSON(t *testing.T) {
	t.Run("Config with delayMs to JSON", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: 8000,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, `{"delayMs":8000}`, jsonStr)
	})

	t.Run("Empty config to JSON", func(t *testing.T) {
		config := &AdvanceConfigEndpoint{
			DelayMs: 0,
		}

		jsonStr, err := config.ToJSON()
		require.NoError(t, err)
		assert.Equal(t, "", jsonStr)
	})
}
