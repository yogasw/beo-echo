package services

import (
	"testing"
	"time"

	"beo-echo/backend/src/database"

	"github.com/stretchr/testify/assert"
)

func TestMockService_applyDelay(t *testing.T) {
	service := &MockService{}

	t.Run("No delay when all parameters are nil or empty", func(t *testing.T) {
		project := &database.Project{AdvanceConfig: ""}

		start := time.Now()
		service.applyDelay(project, nil, nil)
		elapsed := time.Since(start)

		// Should complete almost immediately (less than 10ms)
		assert.Less(t, elapsed.Milliseconds(), int64(10))
	})

	t.Run("Project level delay only", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 50}`,
		}

		start := time.Now()
		service.applyDelay(project, nil, nil)
		elapsed := time.Since(start)

		// Should delay approximately 50ms (allow some tolerance)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(45))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(70))
	})

	t.Run("Endpoint level delay overrides project delay", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 100}`,
		}
		endpoint := &database.MockEndpoint{
			AdvanceConfig: `{"delayMs": 30}`,
		}

		start := time.Now()
		service.applyDelay(project, endpoint, nil)
		elapsed := time.Since(start)

		// Should delay approximately 30ms (endpoint delay), not 100ms (project delay)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(25))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(50))
	})

	t.Run("Response level delay has highest priority", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 100}`,
		}
		endpoint := &database.MockEndpoint{
			AdvanceConfig: `{"delayMs": 80}`,
		}
		response := &database.MockResponse{
			DelayMS: 20,
		}

		start := time.Now()
		service.applyDelay(project, endpoint, response)
		elapsed := time.Since(start)

		// Should delay approximately 20ms (response delay), not project or endpoint delay
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(15))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(35))
	})

	t.Run("Invalid project config JSON should be ignored", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"invalid": "json"`, // Invalid JSON
		}

		start := time.Now()
		service.applyDelay(project, nil, nil)
		elapsed := time.Since(start)

		// Should complete almost immediately since invalid config is ignored
		assert.Less(t, elapsed.Milliseconds(), int64(10))
	})

	t.Run("Invalid endpoint config JSON should be ignored", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 50}`,
		}
		endpoint := &database.MockEndpoint{
			AdvanceConfig: `{"invalid": "json"`, // Invalid JSON
		}

		start := time.Now()
		service.applyDelay(project, endpoint, nil)
		elapsed := time.Since(start)

		// Should use project delay since endpoint config is invalid
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(45))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(70))
	})

	t.Run("Zero delay values should be ignored", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 0}`,
		}
		endpoint := &database.MockEndpoint{
			AdvanceConfig: `{"delayMs": 0}`,
		}
		response := &database.MockResponse{
			DelayMS: 0,
		}

		start := time.Now()
		service.applyDelay(project, endpoint, response)
		elapsed := time.Since(start)

		// Should complete almost immediately since all delays are zero
		assert.Less(t, elapsed.Milliseconds(), int64(10))
	})

	t.Run("Endpoint with zero delay should fall back to project delay", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 40}`,
		}
		endpoint := &database.MockEndpoint{
			AdvanceConfig: `{"delayMs": 0}`, // Zero delay
		}

		start := time.Now()
		service.applyDelay(project, endpoint, nil)
		elapsed := time.Since(start)

		// Should use project delay since endpoint delay is zero
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(35))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(60))
	})

	t.Run("Empty endpoint config should fall back to project delay", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 35}`,
		}
		endpoint := &database.MockEndpoint{
			AdvanceConfig: "", // Empty config
		}

		start := time.Now()
		service.applyDelay(project, endpoint, nil)
		elapsed := time.Since(start)

		// Should use project delay since endpoint config is empty
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(30))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(55))
	})

	t.Run("Complex scenario with all levels configured", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 200}`,
		}
		endpoint := &database.MockEndpoint{
			AdvanceConfig: `{"delayMs": 150}`,
		}
		response := &database.MockResponse{
			DelayMS: 25,
		}

		start := time.Now()
		service.applyDelay(project, endpoint, response)
		elapsed := time.Since(start)

		// Should use response delay (25ms) as it has highest priority
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(20))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(40))
	})
}

func TestMockService_applyDelay_EdgeCases(t *testing.T) {
	service := &MockService{}

	t.Run("Nil project should not panic", func(t *testing.T) {
		start := time.Now()
		service.applyDelay(nil, nil, nil)
		elapsed := time.Since(start)

		// Should complete almost immediately
		assert.Less(t, elapsed.Milliseconds(), int64(10))
	})

	t.Run("Project with nil AdvanceConfig field", func(t *testing.T) {
		project := &database.Project{} // AdvanceConfig will be empty string by default

		start := time.Now()
		service.applyDelay(project, nil, nil)
		elapsed := time.Since(start)

		// Should complete almost immediately
		assert.Less(t, elapsed.Milliseconds(), int64(10))
	})

	t.Run("Very small delay should still work", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 1}`,
		}

		start := time.Now()
		service.applyDelay(project, nil, nil)
		elapsed := time.Since(start)

		// Even 1ms delay should be detectable (with some tolerance)
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(0))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(15))
	})

	t.Run("Response with negative DelayMS should be ignored", func(t *testing.T) {
		project := &database.Project{
			AdvanceConfig: `{"delayMs": 50}`,
		}
		response := &database.MockResponse{
			DelayMS: -10, // Negative delay
		}

		start := time.Now()
		service.applyDelay(project, nil, response)
		elapsed := time.Since(start)

		// Should use project delay since response delay is negative
		assert.GreaterOrEqual(t, elapsed.Milliseconds(), int64(45))
		assert.LessOrEqual(t, elapsed.Milliseconds(), int64(70))
	})
}
