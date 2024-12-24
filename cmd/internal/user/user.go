package user

import (
	"time"
)

type Users struct {
	AccountNumber string
	CreatedAt     time.Time
}

type UserResponse struct {
	AccountNumber string `json:"accountNumber"`
}
