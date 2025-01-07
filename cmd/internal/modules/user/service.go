package user

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/models"
	"log/slog"
	"time"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) createUser() (string, error) {
	slog.Info("Generating account number for new user")

	accNumber, err := helper.GenerateAccountNumber()
	if err != nil {
		slog.Error("Failed to generate account number", slog.String("error", err.Error()))
		return "", err
	}

	hashedAccNumber := helper.HashAccountNumber(accNumber)
	user := models.Users{
		AccountNumber: hashedAccNumber,
		CreatedAt:     time.Now(),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		slog.Error("Failed to insert user", slog.String("error", err.Error()))
		return "", err
	}

	slog.Info("User created successfully", slog.String("accountNumber", accNumber))
	return accNumber, nil
}
