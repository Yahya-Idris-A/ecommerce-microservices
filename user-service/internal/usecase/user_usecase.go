package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/domain"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
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

func (u *userUsecase) Register(ctx context.Context, req *domain.RegisterReq) (*domain.UserResponse, error) {
	// 1. Cek apakah email sudah terdaftar
	existingUser, _ := u.userRepo.GetByEmail(ctx, req.Email)
	if existingUser != nil {
		return nil, errors.New("email already registered")
	}

	// 2. Hash the password using Bcrypt
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, errors.New("[ERROR] Failed to hash password")
	}

	// 3. Mapping ke Entitas
	user := &domain.User{
		ID:           uuid.New(),
		Email:        req.Email,
		PasswordHash: string(hashedPassword),
		FullName:     req.FullName,
		Role:         "buyer", // Default role untuk pendaftar baru
		CreatedAt:    time.Now().UTC(),
		UpdatedAt:    time.Now().UTC(),
	}

	// 4. Simpan ke DB
	err = u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:        user.ID,
		Email:     user.Email,
		FullName:  user.FullName,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
	}, nil
}

func (u *userUsecase) Login(ctx context.Context, req *domain.LoginReq) (*domain.LoginResponse, error) {
	// 1. Find user by email
	user, err := u.userRepo.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, errors.New("[ERROR] Invalid email or password")
	}

	// 2. Compare the provided password with the hashed password in the database
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, errors.New("[ERROR] Invalid email or password")
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
		return nil, errors.New("[ERROR] Failed to generate token")
	}

	return &domain.LoginResponse{
		Token: tokenString,
		User: domain.UserResponse{
			ID:          user.ID,
			Email:       user.Email,
			FullName:    user.FullName,
			PhoneNumber: user.PhoneNumber,
			AvatarURL:   user.AvatarURL,
			Role:        user.Role,
			CreatedAt:   user.CreatedAt,
		},
	}, nil
}

func (u *userUsecase) GetProfile(ctx context.Context, userID uuid.UUID) (*domain.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (u *userUsecase) UpdateProfile(ctx context.Context, userID uuid.UUID, req *domain.UpdateProfileReq) (*domain.UserResponse, error) {
	user, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	user.FullName = req.FullName
	user.PhoneNumber = req.PhoneNumber
	user.AvatarURL = req.AvatarURL
	user.UpdatedAt = time.Now().UTC()

	err = u.userRepo.Update(ctx, user)
	if err != nil {
		return nil, err
	}

	return &domain.UserResponse{
		ID:          user.ID,
		Email:       user.Email,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
		AvatarURL:   user.AvatarURL,
		Role:        user.Role,
		CreatedAt:   user.CreatedAt,
	}, nil
}

func (u *userUsecase) DeleteAccount(ctx context.Context, userID uuid.UUID) error {
	// Pastikan user-nya ada sebelum dihapus
	_, err := u.userRepo.GetByID(ctx, userID)
	if err != nil {
		return errors.New("user not found")
	}

	// Hapus user.
	// Catatan Arsitektur: Di sistem aslinya, karena kita pakai PostgreSQL dengan
	// foreign key ON DELETE CASCADE (yang akan kita set saat migrasi nanti),
	// menghapus user di sini otomatis akan menghapus semua alamatnya di tabel user_addresses juga!
	return u.userRepo.Delete(ctx, userID)
}
