package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *UserService
}

func NewUserHandler(userService *UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

func (h *UserHandler) handleUserAccountCreation(c *gin.Context) {
	accNumber, err := h.userService.createUser()

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	user := UserResponse{
		AccountNumber: accNumber,
	}

	c.JSON(http.StatusOK, user)
}
