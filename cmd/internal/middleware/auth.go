package middleware

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/user"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Authentication is a middleware function that authenticates a user based on the account number provided in the request header.
// It checks if the account number exists in the database and sets the user information in the context if authenticated.
func Authentication(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("Starting authentication process...")

		accNum := c.GetHeader("X-Account-Number")
		if accNum == "" {
			log.Println("Authentication failed: Account number missing.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Missing account number"})
			return
		}

		var users []user.Users
		if err := db.Find(&users).Error; err != nil {
			log.Printf("Authentication failed: Database error: %v", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}

		var authenticatedUser *user.Users
		for _, user := range users {
			err := helper.CompareHashAndPassword([]byte(user.AccountNumber), []byte(accNum))
			if err == nil {
				authenticatedUser = &user
				break
			}
		}

		if authenticatedUser == nil {
			log.Println("Authentication failed: Invalid account number.")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid account number"})
			return
		}

		log.Printf("User authenticated successfully.")
		c.Set("user", authenticatedUser)
		c.Next()
	}
}