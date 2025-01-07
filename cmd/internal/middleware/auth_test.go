package middleware

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/helper/testutils"
	"anon-confessions/cmd/internal/models"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAuthMiddleware(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Set up the mock database
	db := testutils.SetupMockDB()

	// Insert mock users into the database
	mockUsers := []models.Users{
		{ID: 1, AccountNumber: helper.HashAccountNumber("3998442793406687")},
		{ID: 2, AccountNumber: helper.HashAccountNumber("1234567891234567")},
	}
	for _, user := range mockUsers {
		db.Create(&user)
	}

	// Define test cases
	tests := []struct {
		name           string
		accountNumber  string
		expectedStatus int
		expectedUserID int
	}{
		{
			name:           "Valid account number",
			accountNumber:  "3998442793406687",
			expectedStatus: http.StatusOK,
			expectedUserID: 1,
		},
		{
			name:           "Invalid account number",
			accountNumber:  "1234567891234566",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
		{
			name:           "Missing account number",
			accountNumber:  "",
			expectedStatus: http.StatusUnauthorized,
			expectedUserID: 0,
		},
	}

	// Iterate through test cases
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			router := gin.New()
			router.Use(Authentication(db))
			router.GET("/test", func(c *gin.Context) {
				userID, exists := c.Get("userID")
				if exists {
					c.JSON(http.StatusOK, gin.H{"userID": userID})
				} else {
					c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
				}
			})

			req := httptest.NewRequest("GET", "/test", nil)
			req.Header.Set("X-Account-Number", tt.accountNumber)
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			if w.Code != tt.expectedStatus {
				t.Errorf("Expected status %d, got %d", tt.expectedStatus, w.Code)
			}

			// Assert user ID if status is OK
			if tt.expectedStatus == http.StatusOK {
				var resp map[string]interface{}
				if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
					t.Fatalf("Failed to unmarshal response: %v", err)
				}
				userID, ok := resp["userID"].(float64)
				if !ok {
					t.Fatalf("Expected userID in response, but none found")
				}
				if int(userID) != tt.expectedUserID {
					t.Errorf("Expected userID %d, got %d", tt.expectedUserID, int(userID))
				}
			}
		})
	}
}
