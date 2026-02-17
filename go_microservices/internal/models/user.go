package models

import (
	"time"
)

type User struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name" binding:"required"`
	Email     string    `json:"email" db:"email" binding:"required,email"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

type CreateUserRequest struct {
	Name  string `json:"name" binding:"required"`
	Email string `json:"email" binding:"required,email"`
}

type UpdateUserRequest struct {
	Name  string `json:"name" binding:"omitempty"`
	Email string `json:"email" binding:"omitempty,email"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
	Status  int    `json:"status"`
}
