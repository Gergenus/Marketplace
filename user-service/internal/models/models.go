package models

import "github.com/google/uuid"

type User struct {
	ID       uuid.UUID `json:"id,omitempty"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Role     string    `json:"role"`
	Password string    `json:"password"`
}
