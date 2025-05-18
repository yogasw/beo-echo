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
		ID:    "123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := GenerateToken(user)
	assert.NoError(t, err, "Failed to generate token")
	assert.NotEmpty(t, token, "Generated token is empty")
	assert.Contains(t, string(lib.GetJWTSecret()), "test_secret2", "JWT secret does not contain expected value")
}

func TestValidateToken(t *testing.T) {
	// Mock user data
	user := &database.User{
		ID:    "123",
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := GenerateToken(user)
	assert.NoError(t, err, "Failed to generate token")

	claims, err := ValidateToken(token)
	assert.NoError(t, err, "Failed to validate token")
	assert.Equal(t, user.ID, claims.UserID, "UserID does not match")
	assert.Equal(t, user.Email, claims.Email, "Email does not match")
	assert.Equal(t, user.Name, claims.Name, "Name does not match")
	assert.True(t, claims.ExpiresAt.Time.After(time.Now()), "Token is already expired")
}
