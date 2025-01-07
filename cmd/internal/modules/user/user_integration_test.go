// Package user contains integration tests for the user creation functionality in the application.
// These tests validate the behavior of the user registration endpoint and ensure all layers
// (handler → service → repository) work together seamlessly under real-world conditions.
package user

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/helper/testutils"
	"anon-confessions/cmd/internal/models"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
)

// TestUserIntegration validates if a user is created successfully.
func TestUserIntegration(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Step 1: Setup Mock Database
	db := testutils.SetupMockDB()

	// Step 2: Initialize Components
	repo := NewSQLiteUserRepository(db)
	service := NewUserService(repo)
	handler := NewUserHandler(service)

	// Step 3: Setup Router
	router := gin.Default()
	RegisterUsersRoutes(router.Group("/api/v1"), handler)

	t.Run("Test Create User Handler", func(t *testing.T) {
		// Create a POST request to /api/v1/users/register
		w, req := testutils.HTTPTestRequest(http.MethodPost, "/api/v1/users/register", nil)
		router.ServeHTTP(w, req)

		// Assert status code
		if w.Code != http.StatusOK {
			t.Fatalf("Expected status code %d, got %d", http.StatusOK, w.Code)
		}

		// Parse the response
		var resp models.UserResponse
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("Failed to unmarshal response: %v", err)
		}

		// Assert the response contains an account number
		if resp.AccountNumber == "" {
			t.Fatalf("Expected non-empty account number, got '%s'", resp.AccountNumber)
		}

		t.Logf("Created account number: %s", resp.AccountNumber)

		// Query the database to ensure the user exists
		var createdUser models.Users
		if err := db.Where("account_number IS NOT NULL").First(&createdUser).Error; err != nil {
			t.Fatalf("Failed to find created user in database: %v", err)
		}

		if err := helper.CompareHashAndPassword([]byte(createdUser.AccountNumber), []byte(resp.AccountNumber)); err != nil {
			t.Fatalf("Stored account number hash does not match the response account number")
		}

		t.Logf("Hashed account number stored correctly in the database")
	})
}
