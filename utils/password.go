package utils

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
)

// EncryptPassword hashes the input password using bcrypt.
// An error is returned if hashing fails.
func EncryptPassword(password string) (string, error) {
	bs, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bs), nil
}

// VerifyPassword checks the input password against the stored hash using bcrypt.
// SECURITY: MD5 fallback has been removed - only bcrypt hashes are accepted.
// If you have legacy MD5 passwords, users must reset their passwords.
func VerifyPassword(hash, input string) (bool, string, error) {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(input))
	if err == nil {
		return true, "", nil
	}

	// Reject legacy MD5 hashes - they are cryptographically insecure
	var invalidPrefixErr bcrypt.InvalidHashPrefixError
	if errors.As(err, &invalidPrefixErr) || errors.Is(err, bcrypt.ErrHashTooShort) {
		// Password hash format is invalid (possibly legacy MD5)
		// User must reset their password
		return false, "", errors.New("invalid password hash format - password reset required")
	}

	if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
		return false, "", nil
	}
	return false, "", err
}
