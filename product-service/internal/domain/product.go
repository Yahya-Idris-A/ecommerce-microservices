package domain

import (
	"context"
	"time"

	"github.com/google/uuid"
)

// Category merepresentasikan kategori dari sebuah produk
type Category struct {
	ID         uuid.UUID  `json:"id"`
	MerchantID *uuid.UUID `json:"merchant_id,omitempty"` // Tambahkan pointer ini (omitempty agar tidak tampil di JSON jika nil)
	Name       string     `json:"name"`
	CreatedAt  time.Time  `json:"created_at"`
}

// Product adalah entitas utama sebagai "Source of Truth"
type Product struct {
	ID          uuid.UUID `json:"id"`
	MerchantID  uuid.UUID `json:"merchant_id"`
	CategoryID  uuid.UUID `json:"category_id"`
	Name        string    `json:"name"`
	Slug        string    `json:"slug"`
	Description string    `json:"description"`
	Price       float64   `json:"price"` // Untuk portofolio kita pakai float64, di sistem perbankan asli biasanya pakai custom Decimal package
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// CreateProductRequest adalah DTO (Data Transfer Object)
// yang berisi payload dari user saat membuat produk baru.
type CreateProductRequest struct {
	MerchantID  uuid.UUID `json:"merchant_id" binding:"required"`
	CategoryID  uuid.UUID `json:"category_id" binding:"required"`
	Name        string    `json:"name" binding:"required"`
	Description string    `json:"description"`
	Price       float64   `json:"price" binding:"required,gt=0"`
	Stock       int       `json:"stock" binding:"required,min=0"`
}

// ProductRepository adalah kontrak (interface) untuk komunikasi dengan database.
// Layer logika bisnis (Usecase) hanya akan memanggil fungsi di interface ini,
// sehingga Usecase tidak perlu tahu apakah kita pakai PostgreSQL atau MySQL.
type ProductRepository interface {
	Create(ctx context.Context, product *Product) error
	GetByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetAll(ctx context.Context, keyword string, limit int, cursorCreatedAt string, cursorID string) ([]Product, string, error)
	GetByMerchantID(ctx context.Context, merchantID uuid.UUID, limit int, cursorCreatedAt string, cursorID string) ([]Product, string, error)
	Delete(ctx context.Context, id uuid.UUID) error
	// Kita akan tambahkan fungsi lain (Update, Delete) perlahan nanti
}

// ProductUsecase adalah kontrak untuk logika bisnis kita.
// Handler/API (Delivery) akan memanggil interface ini.
type ProductUsecase interface {
	CreateProduct(ctx context.Context, req *CreateProductRequest) (*Product, error)
	GetProductByID(ctx context.Context, id uuid.UUID) (*Product, error)
	GetAllProducts(ctx context.Context, merchantID string, keyword string, limit int, cursorCreatedAt string, cursorID string) ([]Product, string, error)
	GetMyProducts(ctx context.Context, userID uuid.UUID, limit int, cursorCreatedAt string, cursorID string) ([]Product, string, error)
	DeleteProduct(ctx context.Context, userID uuid.UUID, productID uuid.UUID) error
}

// --- Tambahan untuk Category ---

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required"`
}

type CategoryRepository interface {
	Create(ctx context.Context, category *Category) error
	GetByID(ctx context.Context, id uuid.UUID) (*Category, error)
	GetAll(ctx context.Context, merchantID string) ([]Category, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CategoryUsecase interface {
	CreateCategory(ctx context.Context, req *CreateCategoryRequest, merchantID *uuid.UUID) (*Category, error)
	GetCategoryByID(ctx context.Context, id uuid.UUID) (*Category, error)
	GetAllCategories(ctx context.Context, merchantID string) ([]Category, error)
	DeleteCategory(ctx context.Context, role string, merchantID uuid.UUID, categoryID uuid.UUID) error
}
