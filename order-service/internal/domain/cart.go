package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// ==========================================
// 1. ENTITIES (Representasi Tabel Database)
// ==========================================

type Cart struct {
	ID        uuid.UUID `json:"id"`
	UserID    uuid.UUID `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CartItem struct {
	ID         uuid.UUID `json:"id"`
	CartID     uuid.UUID `json:"cart_id"`
	ProductID  uuid.UUID `json:"product_id"`
	MerchantID uuid.UUID `json:"merchant_id"`
	Quantity   int       `json:"quantity"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// ==========================================
// 2. DTOs (Data Transfer Objects untuk Frontend)
// ==========================================

type AddCartItemReq struct {
	ProductID  string `json:"product_id" binding:"required"`
	MerchantID string `json:"merchant_id" binding:"required"`
	Quantity   int    `json:"quantity" binding:"required,min=1"`
}

type UpdateCartItemReq struct {
	Quantity int `json:"quantity" binding:"required,min=1"`
}

// CartItemDetail adalah gabungan data CartItem di database dengan data asli dari Product Service
type CartItemDetail struct {
	ItemID      uuid.UUID `json:"item_id"`
	ProductID   uuid.UUID `json:"product_id"`
	ProductName string    `json:"product_name"` // Didapat dari Product Service
	ProductImg  string    `json:"product_img"`  // Didapat dari Product Service
	Price       float64   `json:"price"`        // Harga real-time dari Product Service
	Quantity    int       `json:"quantity"`
	SubTotal    float64   `json:"sub_total"` // Price * Quantity
}

// MerchantCartGroup mengelompokkan item berdasarkan tokonya
type MerchantCartGroup struct {
	MerchantID   uuid.UUID        `json:"merchant_id"`
	MerchantName string           `json:"merchant_name"` // Didapat dari Product Service
	Items        []CartItemDetail `json:"items"`
}

// CartResponse adalah output final yang dikirim ke Next.js
type CartResponse struct {
	CartID      uuid.UUID           `json:"cart_id"`
	UserID      uuid.UUID           `json:"user_id"`
	Groups      []MerchantCartGroup `json:"groups"`       // Data terkelompok rapi
	TotalAmount float64             `json:"total_amount"` // Grand total seluruh keranjang
}

// ==========================================
// 3. INTERFACES (Kontrak Kerja)
// ==========================================

type CartRepository interface {
	// Urusan Payung Cart
	GetCartByUserID(ctx context.Context, userID uuid.UUID) (*Cart, error)
	CreateCart(ctx context.Context, cart *Cart) error

	// Urusan Isi Cart (Items)
	GetItemByCartAndProduct(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) (*CartItem, error)
	GetItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]CartItem, error)
	AddItem(ctx context.Context, item *CartItem) error
	UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error
	DeleteItem(ctx context.Context, itemID uuid.UUID) error
	DeleteItemsByCartID(ctx context.Context, cartID uuid.UUID) error // Dipakai saat checkout berhasil
}

type CartUsecase interface {
	AddItemToCart(ctx context.Context, userID uuid.UUID, req *AddCartItemReq) error
	GetMyCart(ctx context.Context, userID uuid.UUID) (*CartResponse, error)
	UpdateItemQuantity(ctx context.Context, userID uuid.UUID, itemID uuid.UUID, req *UpdateCartItemReq) error
	RemoveItemFromCart(ctx context.Context, userID uuid.UUID, itemID uuid.UUID) error
}
