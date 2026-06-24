// Package pat implements long-lived Personal Access Tokens (PAT) used by
// external clients (the MCP server, CLIs) to authenticate against the Beo Echo
// REST API. Tokens are opaque strings of the form "beo_pat_<random>"; only
// their SHA-256 hash is persisted.
package pat

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"strings"
)

// TokenPrefix is the fixed marker that identifies a Beo Echo PAT. The auth
// middleware uses it to distinguish PATs from JWTs.
const TokenPrefix = "beo_pat_"

// randomBytesLen is the number of random bytes packed into the secret part of a
// token (256 bits of entropy).
const randomBytesLen = 32

// GenerateToken returns a new plaintext PAT and its SHA-256 hash. The plaintext
// is only ever returned here — callers must show it to the user once and store
// the hash.
func GenerateToken() (plaintext string, hash string, err error) {
	buf := make([]byte, randomBytesLen)
	if _, err = rand.Read(buf); err != nil {
		return "", "", err
	}
	// URL-safe base64 without padding keeps the token copy-paste friendly.
	secret := base64.RawURLEncoding.EncodeToString(buf)
	plaintext = TokenPrefix + secret
	return plaintext, HashToken(plaintext), nil
}

// HashToken returns the hex-encoded SHA-256 hash of a plaintext token. The same
// function is used at creation time and on every auth lookup.
func HashToken(plaintext string) string {
	sum := sha256.Sum256([]byte(plaintext))
	return hex.EncodeToString(sum[:])
}

// IsPAT reports whether the given credential looks like a Beo Echo PAT.
func IsPAT(token string) bool {
	return strings.HasPrefix(token, TokenPrefix)
}

// DisplayPrefix returns the first chars of a token suitable for showing in a
// list UI without revealing the secret (e.g. "beo_pat_AbCd").
func DisplayPrefix(plaintext string) string {
	const shown = len(TokenPrefix) + 4
	if len(plaintext) <= shown {
		return plaintext
	}
	return plaintext[:shown]
}

// ConstantTimeEqual compares two hashes without leaking timing information.
func ConstantTimeEqual(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}
