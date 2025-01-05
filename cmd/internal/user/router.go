package user

import (
	"github.com/gin-gonic/gin"
)

// RegisterUserRoutes registers all routes related to posts.
func RegisterUsersRoutes(router *gin.RouterGroup, userHandler *UserHandler) {
	userGroup := router.Group("/users")
	{
		userGroup.POST("/register", userHandler.handleUserAccountCreation)
	}
}

// @Summary Create a new user account
// @Description Generate a new 16-digit anonymous account number and return it.
// @Tags users
// @Accept json
// @Produce json
// @Success 200 {object} models.UserResponse
// @Router /users/register [post]
func createUser(c *gin.Context) {}
