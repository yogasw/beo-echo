package auth

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"beo-echo/backend/src/database"
	"beo-echo/backend/src/lib"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in the JWT
type JWTClaims struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// GenerateToken creates a new JWT access token for the given user with 1 day expiry
func GenerateToken(user *database.User) (string, error) {
	// Create the claims
	claims := JWTClaims{
		UserID: user.ID,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			// Token expires in 1 day
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(lib.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func GenerateJWTFromString(plainText string) (string, error) {
	// Create the claims
	claims := JWTClaims{
		UserID: plainText,
		Name:   "",
		RegisteredClaims: jwt.RegisteredClaims{
			// Token expires in 24 hours
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(lib.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken creates a new refresh token for the given user (30 days expiry)
func GenerateRefreshToken(user *database.User) (string, error) {
	// Create refresh token claims with longer expiry
	claims := JWTClaims{
		UserID: user.ID,
		Name:   user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			// Refresh token expires in 30 days
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * 24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	// Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with our secret key
	tokenString, err := token.SignedString(lib.GetJWTSecret())
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates the given token and returns the claims
func ValidateToken(tokenString string) (*JWTClaims, error) {
	// Parse and validate the token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// Validate the signing method
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return lib.GetJWTSecret(), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	// Extract the claims
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	return claims, nil
}

// ValidateRefreshToken validates the given refresh token and returns the claims
func ValidateRefreshToken(tokenString string) (*JWTClaims, error) {
	// Parse and validate the token (same validation as access token)
	return ValidateToken(tokenString)
}

// HashToken hashes the given token string using SHA-256 and returns the hex-encoded hash
func HashToken(tokenString string) string {
	// Create a new SHA-256 hash object
	hash := sha256.New()

	// Write the token string to the hash
	hash.Write([]byte(tokenString))

	// Calculate the hash sum and return the hex-encoded string
	return hex.EncodeToString(hash.Sum(nil))
}

// HashRefreshToken creates a SHA256 hash of the refresh token for database storage
func HashRefreshToken(token string) string {
	hash := sha256.Sum256([]byte(token))
	return hex.EncodeToString(hash[:])
}

// ShouldRotateRefreshToken checks if refresh token should be rotated
// Returns true if token was issued more than 7 days ago
func ShouldRotateRefreshToken(tokenString string) bool {
	claims, err := ValidateRefreshToken(tokenString)
	if err != nil {
		// If token is invalid, we should rotate (though this case shouldn't happen in normal flow)
		return true
	}

	// Check if token was issued more than 7 days ago
	sevenDaysAgo := time.Now().Add(-7 * 24 * time.Hour)
	return claims.IssuedAt.Time.Before(sevenDaysAgo)
}
