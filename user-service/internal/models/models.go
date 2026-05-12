package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Verified bool      `json:"verified"`
	Password string    `json:"password"`
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type RefreshSession struct {
	ID            int
	UserID        uuid.UUID
	ResfreshToken uuid.UUID
	Fingerprint   string
	IP            string
	ExpiresIn     int64
	CreatedAt     time.Time
}

type ConfirmationRequest struct {
	Email string `json:"email"`
}
