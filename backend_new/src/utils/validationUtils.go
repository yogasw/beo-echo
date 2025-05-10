package utils

import (
	"regexp"
	"strings"
)

// IsValidEmail checks if a string is a valid email address
func IsValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`)
	return emailRegex.MatchString(email)
}

// IsValidSshUrl checks if a string is a valid SSH URL
func IsValidSshUrl(url string) bool {
	return strings.HasPrefix(url, "git@")
}

// IsValidSshKey checks if a string is a valid SSH private key
func IsValidSshKey(key string) bool {
	return strings.HasPrefix(key, "-----BEGIN OPENSSH PRIVATE KEY-----") &&
		strings.HasSuffix(key, "-----END OPENSSH PRIVATE KEY-----")
}
