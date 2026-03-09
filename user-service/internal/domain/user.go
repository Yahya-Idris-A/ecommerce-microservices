package domain

import (
	"context"
	"time"
)

// 1. Define the core entity for the User Service
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"`    // Hidden from JSON to prevent accidental credential leaks
	Role      string    `json:"role"` // e.g., "admin", "merchant", "customer"
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// 2. Define the contract for the Repository (Database layer)
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id string) (*User, error)
}

// 3. Define the contract for the Usecase (Business logic layer)
type UserUsecase interface {
	Register(ctx context.Context, email, password, role string) (*User, error)
	Login(ctx context.Context, email, password string) (string, error) // Returns a JWT token
	GetProfile(ctx context.Context, id string) (*User, error)
}
