package user

import (
	"anon-confessions/cmd/internal/models"
	"log/slog"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(models.Users) error
}

type SQLiteUserRepository struct {
	db *gorm.DB
}

func NewSQLiteUserRepository(db *gorm.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (repo *SQLiteUserRepository) CreateUser(user models.Users) error {
	slog.Info("Creating a new user", slog.String("accountNumber", user.AccountNumber))

	if err := repo.db.Create(&user).Error; err != nil {
		slog.Error("Failed to create user", slog.String("error", err.Error()), slog.String("accountNumber", user.AccountNumber))
		return err
	}

	slog.Info("User created successfully", slog.String("accountNumber", user.AccountNumber))
	return nil
}
