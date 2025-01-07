package helper

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestAccountNumber(t *testing.T) {
	accountNumber, err := GenerateAccountNumber()
	if err != nil {
		t.Fatalf("Error generating account number: %v", err)
	}

	if len(accountNumber) != 16 {
		t.Fatalf("Expected account number to be 16 digits, got %d", len(accountNumber))
	}
}

func TestHashedAccountNumber(t *testing.T) {
	hashedAccountNumber := HashAccountNumber("1234567890123456")
	if len(hashedAccountNumber) == 0 {
		t.Fatalf("Expected hashed account number to be non-empty, got empty string")
	}

}
func TestRetrieveLoggedInUserId(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Define test cases
	tests := []struct {
		name       string
		setupCtx   func(c *gin.Context)
		expectedID int
		statusCode int
	}{
		{
			name: "No user ID in context",
			setupCtx: func(c *gin.Context) {
				// Do not set userID
			},
			expectedID: 0,
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "Invalid user ID type",
			setupCtx: func(c *gin.Context) {
				c.Set("userID", "invalidType")
			},
			expectedID: 0,
			statusCode: http.StatusUnauthorized,
		},
		{
			name: "Valid user ID",
			setupCtx: func(c *gin.Context) {
				c.Set("userID", 42)
			},
			expectedID: 42,
			statusCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test context and recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up the context
			tt.setupCtx(c)

			userID := RetrieveLoggedInUserId(c)

			// Assert the results
			if userID != tt.expectedID {
				t.Errorf("Expected userID %d, got %d", tt.expectedID, userID)
			}
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}

func TestParseIDParam(t *testing.T) {
	gin.SetMode(gin.TestMode)

	// Define test cases
	tests := []struct {
		name       string
		param      string
		setupCtx   func(c *gin.Context)
		expectedID int
		statusCode int
	}{
		{
			name:  "Valid ID parameter",
			param: "123",
			setupCtx: func(c *gin.Context) {
				c.Params = gin.Params{gin.Param{Key: "id", Value: "123"}}
			},
			expectedID: 123,
			statusCode: http.StatusOK,
		},
		{
			name:  "Invalid ID parameter (non-integer)",
			param: "abc",
			setupCtx: func(c *gin.Context) {
				c.Params = gin.Params{gin.Param{Key: "id", Value: "abc"}}
			},
			expectedID: 0,
			statusCode: http.StatusBadRequest,
		},
		{
			name:  "Missing ID parameter",
			param: "",
			setupCtx: func(c *gin.Context) {
				c.Params = gin.Params{}
			},
			expectedID: 0,
			statusCode: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test context and recorder
			w := httptest.NewRecorder()
			c, _ := gin.CreateTestContext(w)

			// Set up the context
			tt.setupCtx(c)

			// Call the function being tested
			parsedID := ParseIDParam(c, "id")

			// Assert the results
			if parsedID != tt.expectedID {
				t.Errorf("Expected parsedID %d, got %d", tt.expectedID, parsedID)
			}
			if w.Code != tt.statusCode {
				t.Errorf("Expected status code %d, got %d", tt.statusCode, w.Code)
			}
		})
	}
}
