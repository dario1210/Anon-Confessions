package user

import (
	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(Users) error
}

type SQLiteUserRepository struct {
	db *gorm.DB
}

func NewSQLiteUserRepository(db *gorm.DB) *SQLiteUserRepository {
	return &SQLiteUserRepository{db: db}
}

func (repo *SQLiteUserRepository) CreateUser(user Users) error {
	if err := repo.db.Create(&user).Error; err != nil {
		return err
	}
	return nil
}