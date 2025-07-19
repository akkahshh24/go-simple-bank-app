package util

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// HashPassword hashes the given password using a secure hashing algorithm.
func HashPassword(password string) (string, error) {
	// Internally, bcrypt generates a random salt and hashes the password
	// using the bcrypt algorithm with a cost factor.
	// The cost factor determines the computational complexity of the hashing process.
	// A higher cost factor increases security but also increases the time it takes to hash the password.
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("failed to hash password: %w", err)
	}
	return string(hashedPassword), nil
}

// CheckPassword checks if the provided password matches the hashed password.
func CheckPassword(password string, hashedPassword string) error {
	// Internally, bcrypt decodes the hashed password to extract the salt and cost factor,
	// then hashes the provided password using the same salt and cost factor.
	// It then compares the newly hashed password with the original hashed password.
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
