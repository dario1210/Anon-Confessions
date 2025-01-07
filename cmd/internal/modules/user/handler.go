package user

import (
	"anon-confessions/cmd/internal/models"
	"log/slog"
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
	slog.Info("Starting user account creation")

	accNumber, err := h.userService.createUser()
	if err != nil {
		slog.Error("Failed to create user account", slog.String("error", err.Error()))
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}

	slog.Info("User account created successfully", slog.String("accountNumber", accNumber))

	user := models.UserResponse{
		AccountNumber: accNumber,
	}

	c.JSON(http.StatusOK, user)
}