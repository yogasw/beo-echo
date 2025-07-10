package auth

import (
	"testing"
	"time"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/lib"

	"github.com/stretchr/testify/assert"
)

func TestGenerateToken(t *testing.T) {
	lib.SetJWTSecret("test_secret2")
	// Mock user data
	user := &database.User{
		ID:   "123",
		Name: "Test User",
	}

	token, err := GenerateToken(user)
	assert.NoError(t, err, "Failed to generate token")
	assert.NotEmpty(t, token, "Generated token is empty")
	assert.Contains(t, string(lib.GetJWTSecret()), "test_secret2", "JWT secret does not contain expected value")
}

func TestValidateToken(t *testing.T) {
	// Mock user data
	user := &database.User{
		ID:   "123",
		Name: "Test User",
	}

	token, err := GenerateToken(user)
	assert.NoError(t, err, "Failed to generate token")

	claims, err := ValidateToken(token)
	assert.NoError(t, err, "Failed to validate token")
	assert.Equal(t, user.ID, claims.UserID, "UserID does not match")
	assert.Equal(t, user.Name, claims.Name, "Name does not match")
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()), "Token is already expired")
}

func TestGenerateRefreshToken(t *testing.T) {
	lib.SetJWTSecret("test_secret2")
	// Create a test user
	user := &database.User{
		ID:   "test-user-id",
		Name: "Test User",
	}

	// Generate refresh token
	refreshToken, err := GenerateRefreshToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, refreshToken)

	// Validate the refresh token
	claims, err := ValidateRefreshToken(refreshToken)
	assert.NoError(t, err)
	assert.Equal(t, user.ID, claims.UserID)
	assert.Equal(t, user.Name, claims.Name)

	// Check that token expires in approximately 30 days
	expectedExpiry := time.Now().Add(30 * 24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time
	timeDiff := actualExpiry.Sub(expectedExpiry)

	// Allow some tolerance (1 minute) for test execution time
	assert.True(t, timeDiff < time.Minute && timeDiff > -time.Minute,
		"Token expiry should be approximately 30 days from now")
}

func TestHashRefreshToken(t *testing.T) {
	token := "test-refresh-token-123"

	// Hash the token
	hash1 := HashRefreshToken(token)
	hash2 := HashRefreshToken(token)

	// Same token should produce same hash
	assert.Equal(t, hash1, hash2)
	assert.NotEqual(t, token, hash1) // Hash should be different from original
	assert.Len(t, hash1, 64)         // SHA256 hex string should be 64 characters
}

func TestAccessTokenExpiry(t *testing.T) {
	lib.SetJWTSecret("test_secret2")
	// Create a test user
	user := &database.User{
		ID:   "test-user-id",
		Name: "Test User",
	}

	// Generate access token
	accessToken, err := GenerateToken(user)
	assert.NoError(t, err)
	assert.NotEmpty(t, accessToken)

	// Validate the access token
	claims, err := ValidateToken(accessToken)
	assert.NoError(t, err)

	// Check that token expires in approximately 24 hours
	expectedExpiry := time.Now().Add(24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time
	timeDiff := actualExpiry.Sub(expectedExpiry)

	// Allow some tolerance (1 minute) for test execution time
	assert.True(t, timeDiff < time.Minute && timeDiff > -time.Minute,
		"Access token expiry should be approximately 15 minutes from now")
}

func TestShouldRotateRefreshToken(t *testing.T) {
	lib.SetJWTSecret("test_secret2")

	// Create a test user
	user := &database.User{
		ID:   "test-user-id",
		Name: "Test User",
	}

	// Generate fresh refresh token (should not rotate)
	freshToken, err := GenerateRefreshToken(user)
	assert.NoError(t, err)

	shouldRotate := ShouldRotateRefreshToken(freshToken)
	assert.False(t, shouldRotate, "Fresh token should not need rotation")

	// Test with invalid token (should rotate)
	invalidToken := "invalid.jwt.token"
	shouldRotate = ShouldRotateRefreshToken(invalidToken)
	assert.True(t, shouldRotate, "Invalid token should trigger rotation")
}
