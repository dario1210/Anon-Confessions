package user

import (
	"anon-confessions/cmd/internal/helper"
	"anon-confessions/cmd/internal/models"
	"log"
	"time"
)

type UserService struct {
	userRepo UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) createUser() (string, error) {
	accNumber, err := helper.GenerateAccountNumber()

	if err != nil {
		log.Print(err)
	}

	hashedAccNumber := helper.HashAccountNumber(accNumber)
	user := models.Users{
		AccountNumber: hashedAccNumber,
		CreatedAt:     time.Now(),
	}

	err = s.userRepo.CreateUser(user)
	if err != nil {
		log.Print(err)
		return "", err
	}

	return accNumber, nil
}
