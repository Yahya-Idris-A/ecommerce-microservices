package domain

import "context"

// 1. Define the core entity
type Product struct {
	ID          string `json:"id"`
	CategoryID  string `json:"category_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	Price       string `json:"price"`
	Stock       int32  `json:"stock"`
	MerchantID  string `json:"merchant_id"`
}

type Merchant struct {
	ID          string `json:"id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
}

type GlobalSearchResult struct {
	Products  []Product  `json:"products"`
	Merchants []Merchant `json:"merchants"`
}

// 2. Define the contract for the Repository
type ProductSearchRepository interface {
	// Ubah return type menjadi GlobalSearchResult
	SearchGlobal(ctx context.Context, query string, limit int, cursor string) (*GlobalSearchResult, string, error)
}

// 3. Define the contract for the Usecase
type ProductSearchUsecase interface {
	// Ubah return type menjadi GlobalSearchResult
	Search(ctx context.Context, query string, limit int, cursor string) (*GlobalSearchResult, string, error)
}
