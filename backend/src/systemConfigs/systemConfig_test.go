package systemConfig

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/utils"
)

var TEST_ARRAY = "TEST_ARRAY"
var BASE_URL = "BASE_URL"

// TestSystemConfigOperations tests the system config operations
func TestSystemConfigOperations(t *testing.T) {
	// Setup test environment
	utils.SetupFolderConfigForTest()

	// Initialize database
	err := database.CheckAndHandle()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		utils.CleanupTestFolders()
	})

	// Test SetSystemConfig
	t.Run("SetSystemConfig", func(t *testing.T) {
		// Test setting string config
		err := SetSystemConfig(BASE_URL, "https://example.com")
		assert.NoError(t, err, "Setting string config should not error")

		// Test setting boolean config
		err = SetSystemConfig(FEATURE_SHOW_PASSWORD_REQUIREMENTS, "true")
		assert.NoError(t, err, "Setting boolean config should not error")

		// Add config with array type manually for testing
		_, err = AddConfig(TEST_ARRAY, `["item1","item2"]`, "Test array config", "array", false)
		assert.NoError(t, err, "Adding array config should not error")

		// Test setting config with hidden value using AddConfig
		_, err = AddConfig("ssh_key", "secret-ssh-key", "SSH key for Git authentication", "string", true)
		assert.NoError(t, err, "Adding config with hidden value should not error")
	})

	// Test GetSystemConfig
	t.Run("GetSystemConfig", func(t *testing.T) {
		// Test getting string config
		val, err := GetSystemConfig(BASE_URL)
		assert.NoError(t, err, "Getting string config should not error")
		assert.Equal(t, "https://example.com", val, "Retrieved value should match set value")

		// Test getting boolean config
		val, err = GetSystemConfig(FEATURE_SHOW_PASSWORD_REQUIREMENTS)
		assert.NoError(t, err, "Getting boolean config should not error")
		assert.Equal(t, true, val, "Retrieved value should match set value")

		// Test getting array config
		val, err = GetSystemConfig(TEST_ARRAY)
		assert.NoError(t, err, "Getting array config should not error")
		arrVal, ok := val.([]string)
		assert.True(t, ok, "Value should be a string array")
		assert.Equal(t, []string{"item1", "item2"}, arrVal, "Retrieved array should match set value")
	})

	// Test GetConfig with type inference
	t.Run("GetConfig with type inference", func(t *testing.T) {
		// Test string config
		strVal, err := GetSystemConfigWithType[string]("BASE_URL")
		assert.NoError(t, err, "Getting typed string config should not error")
		assert.Equal(t, "https://example.com", strVal, "Retrieved string should match set value")

		// Test boolean config
		boolVal, err := GetSystemConfigWithType[bool]("FEATURE_SHOW_PASSWORD_REQUIREMENTS")
		assert.NoError(t, err, "Getting typed boolean config should not error")
		assert.Equal(t, true, boolVal, "Retrieved boolean should match set value")

		// Test array config
		arrVal, err := GetSystemConfigWithType[[]string]("TEST_ARRAY")
		assert.NoError(t, err, "Getting typed array config should not error")
		assert.Equal(t, []string{"item1", "item2"}, arrVal, "Retrieved array should match set value")
	})

	// Test GetAllSystemConfigs
	t.Run("GetAllSystemConfigs", func(t *testing.T) {
		configs, err := GetAllSystemConfigs()
		assert.NoError(t, err, "Getting all configs should not error")
		assert.GreaterOrEqual(t, len(configs), 3, "Should retrieve at least the 3 configs we set")

		// Count our test configs
		foundBaseURL := false
		foundPasswordReq := false
		foundTestArray := false

		for _, cfg := range configs {
			if cfg.Key == "BASE_URL" {
				foundBaseURL = true
				assert.Equal(t, "https://example.com", cfg.Value, "BASE_URL value should match")
			}
			if cfg.Key == "FEATURE_SHOW_PASSWORD_REQUIREMENTS" {
				foundPasswordReq = true
				assert.Equal(t, "true", cfg.Value, "FEATURE_SHOW_PASSWORD_REQUIREMENTS value should match")
			}
			if cfg.Key == "TEST_ARRAY" {
				foundTestArray = true
				assert.Equal(t, "array", cfg.Type, "TEST_ARRAY should have array type")
			}
		}

		assert.True(t, foundBaseURL, "BASE_URL config should be found")
		assert.True(t, foundPasswordReq, "FEATURE_SHOW_PASSWORD_REQUIREMENTS config should be found")
		assert.True(t, foundTestArray, "TEST_ARRAY config should be found")
	})

	// Test GetFeatureFlags
	t.Run("GetFeatureFlags", func(t *testing.T) {
		flags, err := GetFeatureFlags()
		assert.NoError(t, err, "Getting feature flags should not error")
		assert.NotNil(t, flags, "Feature flags should not be nil")

		// Check our test feature flag
		assert.Equal(t, true, flags["FEATURE_SHOW_PASSWORD_REQUIREMENTS"], "Feature flag value should match")
	})

	// Test Default Values
	t.Run("Default Values", func(t *testing.T) {
		// Test getting a default config that hasn't been explicitly set
		val, err := GetSystemConfig(CUSTOM_SUBDOMAIN_DOMAIN)
		assert.NoError(t, err, "Getting default config should not error")
		assert.Equal(t, "*.beo-echo.com", val, "Default value should be returned")
	})

	// Test AddConfig and SetConfigByID
	t.Run("AddConfig and SetConfigByID", func(t *testing.T) {
		key := "TEST_ADD_CONFIG"
		value := "test_value"
		description := "Test config for adding"
		// Add a new config
		config, err := AddConfig(key, value, description, "string", false)
		assert.NoError(t, err, "Adding config should not error")
		assert.NotNil(t, config, "Config should be created")
		assert.Equal(t, key, config.Key, "Config key should match")

		// Get the ID for updating
		c := &database.SystemConfig{}
		err = database.DB.Model(&database.SystemConfig{}).Where("key = ?", key).Scan(&c).Error
		assert.NoError(t, err, "Getting config ID should not error")
		id := c.ID

		// Update the config
		updatedConfig, err := SetConfigByID(id, key, "updated", "Updated description")
		assert.NoError(t, err, "Updating config should not error")
		assert.Equal(t, "updated", updatedConfig.Value, "Value should be updated")
		assert.Equal(t, "Updated description", updatedConfig.Description, "Description should be updated")
	})
}
