package auth

import (
	"errors"
	"fmt"
	"os"
	"time"

	"mockoon-control-panel/backend_new/src/database"

	"github.com/golang-jwt/jwt/v5"
)

// JWTClaims represents the claims in the JWT
type JWTClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
	jwt.RegisteredClaims
}

// Default JWT secret key - will be overridden by environment variable if available
var jwtSecretKey = []byte("default_jwt_secret_key_change_in_production")

func init() {
	// Use environment variable for JWT secret if available
	if secretKey := os.Getenv("JWT_SECRET"); secretKey != "" {
		jwtSecretKey = []byte(secretKey)
	}
}

// GenerateToken creates a new JWT token for the given user
func GenerateToken(user *database.User) (string, error) {
	// Create the claims
	claims := JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Name:   user.Name,
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
	tokenString, err := token.SignedString(jwtSecretKey)
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
		return jwtSecretKey, nil
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
