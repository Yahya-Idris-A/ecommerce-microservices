package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// --- ENTITIES & DTOs ---

type Address struct {
	ID            uuid.UUID `json:"id"`
	UserID        uuid.UUID `json:"user_id"`
	Label         string    `json:"label"` // Misal: "Rumah", "Kantor"
	RecipientName string    `json:"recipient_name"`
	PhoneNumber   string    `json:"phone_number"`
	FullAddress   string    `json:"full_address"`
	City          string    `json:"city"`
	Province      string    `json:"province"`
	PostalCode    string    `json:"postal_code"`
	IsPrimary     bool      `json:"is_primary"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type CreateAddressReq struct {
	Label         string `json:"label" binding:"required"`
	RecipientName string `json:"recipient_name" binding:"required"`
	PhoneNumber   string `json:"phone_number" binding:"required"`
	FullAddress   string `json:"full_address" binding:"required"`
	City          string `json:"city" binding:"required"`
	Province      string `json:"province" binding:"required"`
	PostalCode    string `json:"postal_code" binding:"required"`
	IsPrimary     bool   `json:"is_primary"`
}

type UpdateAddressReq struct {
	Label         string `json:"label"`
	RecipientName string `json:"recipient_name"`
	PhoneNumber   string `json:"phone_number"`
	FullAddress   string `json:"full_address"`
	City          string `json:"city"`
	Province      string `json:"province"`
	PostalCode    string `json:"postal_code"`
}

// --- INTERFACES ---

type AddressRepository interface {
	Create(ctx context.Context, address *Address) error
	GetByUserID(ctx context.Context, userID uuid.UUID) ([]Address, error)
	GetByID(ctx context.Context, id uuid.UUID) (*Address, error)
	Update(ctx context.Context, address *Address) error
	Delete(ctx context.Context, id uuid.UUID) error

	// Fungsi sakti untuk mereset status primary alamat lain milik user yang sama
	RemovePrimaryStatus(ctx context.Context, userID uuid.UUID) error
}

type AddressUsecase interface {
	AddAddress(ctx context.Context, userID uuid.UUID, req *CreateAddressReq) (*Address, error)
	GetMyAddresses(ctx context.Context, userID uuid.UUID) ([]Address, error)
	UpdateAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID, req *UpdateAddressReq) (*Address, error)
	DeleteAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error
	SetPrimaryAddress(ctx context.Context, userID uuid.UUID, addressID uuid.UUID) error
}
