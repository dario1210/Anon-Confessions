package user

import (
	"anon-confessions/cmd/internal/models"

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
	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}
