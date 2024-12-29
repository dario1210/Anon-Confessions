package models

import (
	"time"
)

type Users struct {
	ID            int       `json:"id" gorm:"primaryKey;autoIncrement"`
	AccountNumber string    `json:"account_number" gorm:"type:varchar(255);not null;unique"`
	CreatedAt     time.Time `json:"created_at" gorm:"autoCreateTime"`
}

type UserResponse struct {
	AccountNumber string `json:"accountNumber"`
}
