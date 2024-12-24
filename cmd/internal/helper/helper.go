// Package helper provides utility functions, reusable structs, and common helpers
// that can be used across the entire application to reduce redundancy
// and promote code reusability.package helper
package helper

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"golang.org/x/crypto/bcrypt"
)

// GenerateAccountNumber generates a random 16-digit account number.
// It returns the generated account number as a string and an error if any occurs during generation.
func GenerateAccountNumber() (string, error) {
	max := big.NewInt(10_000_000_000_000_000)

	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%016d", n), nil
}

// HashAccountNumber hashes the given account number using bcrypt.
// It returns the hashed account number as a string.
// If an error occurs during hashing, it logs the error
func HashAccountNumber(accountNumber string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(accountNumber), bcrypt.DefaultCost)

	if err != nil {
		log.Printf("Error hashing account number: %s", accountNumber)
	}

	return string(hash)
}

// CompareHashAndPassword compares a hashed password with its possible plaintext equivalent.
// It returns nil on success, or an error on failure.
func CompareHashAndPassword(hashedPassword, password []byte) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
