package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	userRepo  domain.UserRepository
	jwtSecret []byte
}

func NewUserUsecase(repo domain.UserRepository, secret string) domain.UserUsecase {
	return &userUsecase{
		userRepo:  repo,
		jwtSecret: []byte(secret),
	}
}

func (u *userUsecase) Register(ctx context.Context, email, password, role string) (*domain.User, error) {
	// 1. Basic validation
	if email == "" || password == "" {
		return nil, errors.New("[ERROR] Email and password cannot be empty")
	}

	// Default role to customer if not provided
	if role == "" {
		role = "customer"
	}

	// 2. Hash the password using Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("[ERROR] Failed to hash password")
	}

	// 3. Prepare the user entity
	user := &domain.User{
		Email:    email,
		Password: string(hashedPassword),
		Role:     role,
	}

	// 4. Save to database via Repository
	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *userUsecase) Login(ctx context.Context, email, password string) (string, error) {
	// 1. Find user by email
	user, err := u.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("[ERROR] Invalid email or password")
	}

	// 2. Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("[ERROR] Invalid email or password")
	}

	// 3. Generate JWT Token
	// The payload contains user ID, Role, and expiration time
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"role":    user.Role,
		"exp":     time.Now().Add(time.Hour * 72).Unix(), // Token expires in 72 hours
	})

	// Sign the token with our secret key
	tokenString, err := token.SignedString(u.jwtSecret)
	if err != nil {
		return "", errors.New("[ERROR] Failed to generate token")
	}

	return tokenString, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, id string) (*domain.User, error) {
	return u.userRepo.GetByID(ctx, id)
}
