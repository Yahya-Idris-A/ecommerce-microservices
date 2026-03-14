package domain

import (
	"time"

	"github.com/google/uuid"
)

type Merchant struct {
	ID          uuid.UUID `json:"id"`
	UserID      uuid.UUID `json:"user_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// UserID tidak dikirim via JSON, melainkan diambil dari Token JWT demi keamanan
type CreateMerchantRequest struct {
	Name        string `json:"name" binding:"required"`
	Description string `json:"description"`
}

type MerchantRepository interface {
	Create(merchant *Merchant) error
	GetByUserID(userID uuid.UUID) (*Merchant, error)
	GetByID(id uuid.UUID) (*Merchant, error)
	GetAll(keyword string, limit int, cursorCreatedAt string, cursorID string) ([]Merchant, string, error)
}

type MerchantUsecase interface {
	CreateMerchant(userID uuid.UUID, req *CreateMerchantRequest) (*Merchant, error)
	GetByUserID(userID uuid.UUID) (*Merchant, error)
	GetMerchantByID(id uuid.UUID) (*Merchant, error)
	GetAllMerchants(keyword string, limit int, cursorCreatedAt string, cursorID string) ([]Merchant, string, error)
}
