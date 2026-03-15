package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// 1. Define the core entity for the User Service
type User struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"` // Disembunyikan dari JSON
	FullName     string    `json:"full_name"`
	PhoneNumber  string    `json:"phone_number"`
	AvatarURL    string    `json:"avatar_url"`
	Role         string    `json:"role"` // "buyer", "merchant", "admin"
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	ID          uuid.UUID `json:"id"`
	Email       string    `json:"email"`
	FullName    string    `json:"full_name"`
	PhoneNumber string    `json:"phone_number"`
	AvatarURL   string    `json:"avatar_url"`
	Role        string    `json:"role"`
	CreatedAt   time.Time `json:"created_at"`
}

type RegisterReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
	FullName string `json:"full_name" binding:"required"`
}

type LoginReq struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string       `json:"token"`
	User  UserResponse `json:"user"`
}

type UpdateProfileReq struct {
	FullName    string `json:"full_name" binding:"required"`
	PhoneNumber string `json:"phone_number"`
	AvatarURL   string `json:"avatar_url"`
}

// 2. Define the contract for the Repository (Database layer)
type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByID(ctx context.Context, id uuid.UUID) (*User, error)
	Update(ctx context.Context, user *User) error
	Delete(ctx context.Context, id uuid.UUID) error
}

// 3. Define the contract for the Usecase (Business logic layer)
type UserUsecase interface {
	Register(ctx context.Context, req *RegisterReq) (*UserResponse, error)
	Login(ctx context.Context, req *LoginReq) (*LoginResponse, error)
	GetProfile(ctx context.Context, userID uuid.UUID) (*UserResponse, error)
	UpdateProfile(ctx context.Context, userID uuid.UUID, req *UpdateProfileReq) (*UserResponse, error)
	DeleteAccount(ctx context.Context, userID uuid.UUID) error
}
