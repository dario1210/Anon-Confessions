package models

import "time"

// Users represents a user in the database.
type Users struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountNumber string    `json:"account_number" gorm:"type:varchar(255);not null;unique"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

// UserResponse is a minimal representation of a user used in API responses.
type UserResponse struct {
	AccountNumber string `json:"accountNumber"`
}
