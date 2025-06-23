package systemConfig

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/utils"
)

// Test constants - only used in tests
const (
	TEST_ARRAY_CONFIG = "TEST_ARRAY_CONFIG" // Test array configuration for testing array type configs
)

// TestSystemConfigOperations tests the system config operations
func TestSystemConfigOperations(t *testing.T) {
	// Setup test environment
	utils.SetupFolderConfigForTest()

	// Add test array config to DefaultConfigSettings for testing
	DefaultConfigSettings[TEST_ARRAY_CONFIG] = ConfigSetting{
		Type:        TypeArray,
		Value:       `["test1","test2","test3"]`,
		Description: "Test array configuration for testing array type functionality",
		Category:    "Testing",
	}

	// Initialize database
	err := database.CheckAndHandle()
	if err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	// Cleanup after test
	t.Cleanup(func() {
		// Remove test config from DefaultConfigSettings
		delete(DefaultConfigSettings, TEST_ARRAY_CONFIG)
		utils.CleanupTestFolders()
	})

	// Test SetSystemConfig
	t.Run("SetSystemConfig", func(t *testing.T) {
		// Test setting string config using existing key
		err := SetSystemConfig(CUSTOM_SUBDOMAIN_DOMAIN, "*.example.com")
		assert.NoError(t, err, "Setting string config should not error")

		// Test setting boolean config
		err = SetSystemConfig(FEATURE_SHOW_PASSWORD_REQUIREMENTS, "true")
		assert.NoError(t, err, "Setting boolean config should not error")

		// Test setting array config using the predefined test array config
		err = SetSystemConfig(TEST_ARRAY_CONFIG, `["item1","item2","item3"]`)
		assert.NoError(t, err, "Setting array config should not error")

		// Test setting config with hidden value
		err = SetSystemConfig(JWT_SECRET, "secret-jwt-key")
		assert.NoError(t, err, "Setting config with hidden value should not error")

		// Test setting config with non-existent key (should error)
		err = SetSystemConfig("NON_EXISTENT_KEY", "some_value")
		assert.Error(t, err, "Setting non-existent config key should error")
		assert.Contains(t, err.Error(), "configuration key NON_EXISTENT_KEY not found", "Error should mention key not found")
	})

	// Test GetSystemConfig
	t.Run("GetSystemConfig", func(t *testing.T) {
		// Test getting string config
		val, err := GetSystemConfig(CUSTOM_SUBDOMAIN_DOMAIN)
		assert.NoError(t, err, "Getting string config should not error")
		assert.Equal(t, "*.example.com", val, "Retrieved value should match set value")

		// Test getting boolean config
		val, err = GetSystemConfig(FEATURE_SHOW_PASSWORD_REQUIREMENTS)
		assert.NoError(t, err, "Getting boolean config should not error")
		assert.Equal(t, true, val, "Retrieved value should match set value")

		// Test getting array config
		val, err = GetSystemConfig(TEST_ARRAY_CONFIG)
		assert.NoError(t, err, "Getting array config should not error")
		arrVal, ok := val.([]string)
		assert.True(t, ok, "Value should be a string array")
		assert.Equal(t, []string{"item1", "item2", "item3"}, arrVal, "Retrieved array should match set value")

		// Test getting non-existent config key (should error)
		_, err = GetSystemConfig("NON_EXISTENT_KEY")
		assert.Error(t, err, "Getting non-existent config key should error")
		assert.Contains(t, err.Error(), "configuration key NON_EXISTENT_KEY not found", "Error should mention key not found")
	})

	// Test GetConfig with type inference
	t.Run("GetConfig with type inference", func(t *testing.T) {
		// Test string config
		strVal, err := GetSystemConfigWithType[string](CUSTOM_SUBDOMAIN_DOMAIN)
		assert.NoError(t, err, "Getting typed string config should not error")
		assert.Equal(t, "*.example.com", strVal, "Retrieved string should match set value")

		// Test boolean config
		boolVal, err := GetSystemConfigWithType[bool](FEATURE_SHOW_PASSWORD_REQUIREMENTS)
		assert.NoError(t, err, "Getting typed boolean config should not error")
		assert.Equal(t, true, boolVal, "Retrieved boolean should match set value")

		// Test array config
		arrVal, err := GetSystemConfigWithType[[]string](TEST_ARRAY_CONFIG)
		assert.NoError(t, err, "Getting typed array config should not error")
		assert.Equal(t, []string{"item1", "item2", "item3"}, arrVal, "Retrieved array should match set value")

		// Test getting non-existent config key with type (should error)
		_, err = GetSystemConfigWithType[string]("NON_EXISTENT_KEY")
		assert.Error(t, err, "Getting non-existent config key with type should error")
		assert.Contains(t, err.Error(), "configuration key NON_EXISTENT_KEY not found", "Error should mention key not found")
	})

	// Test GetAllSystemConfigs
	t.Run("GetAllSystemConfigs", func(t *testing.T) {
		configs, err := GetAllSystemConfigs()
		assert.NoError(t, err, "Getting all configs should not error")
		assert.GreaterOrEqual(t, len(configs), 3, "Should retrieve at least the 3 configs we set")

		// Count our test configs
		foundCustomDomain := false
		foundPasswordReq := false
		foundTestArray := false

		for _, cfg := range configs {
			if cfg.Key == CUSTOM_SUBDOMAIN_DOMAIN {
				foundCustomDomain = true
				assert.Equal(t, "*.example.com", cfg.Value, "CUSTOM_SUBDOMAIN_DOMAIN value should match")
			}
			if cfg.Key == FEATURE_SHOW_PASSWORD_REQUIREMENTS {
				foundPasswordReq = true
				assert.Equal(t, "true", cfg.Value, "FEATURE_SHOW_PASSWORD_REQUIREMENTS value should match")
			}
			if cfg.Key == TEST_ARRAY_CONFIG {
				foundTestArray = true
				assert.Equal(t, "array", cfg.Type, "TEST_ARRAY_CONFIG should have array type")
			}
		}

		assert.True(t, foundCustomDomain, "CUSTOM_SUBDOMAIN_DOMAIN config should be found")
		assert.True(t, foundPasswordReq, "FEATURE_SHOW_PASSWORD_REQUIREMENTS config should be found")
		assert.True(t, foundTestArray, "TEST_ARRAY_CONFIG config should be found")
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
		// Test getting a config that was already set in previous test
		val, err := GetSystemConfig(CUSTOM_SUBDOMAIN_DOMAIN)
		assert.NoError(t, err, "Getting config should not error")
		assert.Equal(t, "*.example.com", val, "Should return the value that was set in earlier test")
		
		// Test getting a config that hasn't been explicitly set (should return default)
		val, err = GetSystemConfig(CUSTOM_SUBDOMAIN_ENABLED)
		assert.NoError(t, err, "Getting default config should not error")
		assert.Equal(t, false, val, "Default boolean value should be returned")
	})

	// Test AddConfig and SetConfigByID
	t.Run("AddConfig and SetConfigByID", func(t *testing.T) {
		// Test AddConfig with existing key - this should work as AddConfig allows any key
		key := "TEST_ADD_CONFIG_KEY"
		value := "test_value"
		description := "Test config for adding"
		
		// Add a new config (AddConfig allows creating any key)
		config, err := AddConfig(key, value, description, "string", false)
		assert.NoError(t, err, "Adding config should not error")
		assert.NotNil(t, config, "Config should be created")
		assert.Equal(t, key, config.Key, "Config key should match")

		// Get the ID for updating
		c := &database.SystemConfig{}
		err = database.DB.Model(&database.SystemConfig{}).Where("key = ?", key).Scan(&c).Error
		assert.NoError(t, err, "Getting config ID should not error")
		id := c.ID

		// Update the config using SetConfigByID
		updatedConfig, err := SetConfigByID(id, key, "updated", "Updated description")
		assert.NoError(t, err, "Updating config should not error")
		assert.Equal(t, "updated", updatedConfig.Value, "Value should be updated")
		assert.Equal(t, "Updated description", updatedConfig.Description, "Description should be updated")

		// Test SetConfigByID with non-existent ID (should error)
		_, err = SetConfigByID("99999", "some_key", "some_value", "some description")
		assert.Error(t, err, "Setting config with non-existent ID should error")
	})

	// Test SetSystemConfig error scenarios
	t.Run("SetSystemConfig Error Scenarios", func(t *testing.T) {
		// Test setting config with empty key (should error) 
		err := SetSystemConfig("", "some_value")
		assert.Error(t, err, "Setting config with empty key should error")
		assert.Contains(t, err.Error(), "configuration key  not found", "Error should mention key not found")

		// Test setting config with key that doesn't exist in DefaultConfigSettings (should error)
		err = SetSystemConfig("UNKNOWN_KEY", "some_value")
		assert.Error(t, err, "Setting config with unknown key should error")
		assert.Contains(t, err.Error(), "configuration key UNKNOWN_KEY not found", "Error should mention key not found")
	})
}
