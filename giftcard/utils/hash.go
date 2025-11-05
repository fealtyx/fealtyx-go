package utils

import (
	"crypto/sha256"
	"fmt"
	"regexp"
)

// GetSHA256Hash returns the sha256 hash of the input string
// if the input string is already hashed, returns the string as it is
func GetSHA256Hash(s string) string {
	if s == "" || IsSHA256Hash(s) {
		return s
	}

	// Create a new SHA-512 hash object
	hash := sha256.New()

	// Write the input string to the hash object
	hash.Write([]byte(s))

	// Get the SHA-512 hash as a byte slice
	hashedBytes := hash.Sum(nil)

	// Convert the byte slice to a hexadecimal string
	hashedString := fmt.Sprintf("%x", hashedBytes)

	return hashedString
}

// IsSHA256Hash checks if the input string is a SHA-256 hash.
func IsSHA256Hash(s string) bool {
	// SHA-256 hash: 64-character hexadecimal string
	sha256Regex := regexp.MustCompile(`^[a-fA-F0-9]{64}$`)
	return sha256Regex.MatchString(s)
}
