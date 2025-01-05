// Package helper provides utility functions, reusable structs, and common helpers
// that can be used across the entire application to reduce redundancy
// and promote code reusability.package helper
package helper

import (
	"anon-confessions/cmd/internal/models"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

type SuccessMessage struct {
	Message string `json:"msg"`
}

type ErrorMessage struct {
	Message string `json:"error"`
}

// RetrieveLoggedInUserId retrieves the logged-in user's ID from the Gin context.
// If the user ID is not found in the context or is of an invalid type,
// the function immediately aborts the HTTP request with an appropriate unauthorized error response.
func RetrieveLoggedInUserId(c *gin.Context) int {
	userID, exists := c.Get("userID")

	if !exists {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorMessage{Message: "Authentication failed."})
		return 0
	}

	intUserId, ok := userID.(int)
	if !ok {
		c.AbortWithStatusJSON(http.StatusUnauthorized, ErrorMessage{Message: "Invalid user ID format."})
		return 0
	}

	return intUserId
}

// ParseIDParam retrieves the parameter specified from the route parameter as an integer.
// If the format is invalid, it aborts the HTTP request with a 400 Bad Request status.
func ParseIDParam(c *gin.Context, param string) int {
	id := c.Param(param)
	intID, err := strconv.Atoi(id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ErrorMessage{Message: "Invalid ID format."})
		return 0
	}
	return intID
}

// GenerateSortOrder generates an SQL ORDER BY clause based on the provided sorting preferences.
// It will give priority to the SortByLikes field if it is set.
func GenerateOrderClause(postQueryParam models.PostQueryParams) string {
	var orderClause string

	if postQueryParam.SortByLikes != "" {
		orderClause = "total_likes " + postQueryParam.SortByLikes
	}
	if postQueryParam.SortByCreationDate != "" {
		orderClause = "created_at " + postQueryParam.SortByCreationDate
	}

	log.Println("Generated ORDER BY clause:", orderClause)
	return orderClause
}
