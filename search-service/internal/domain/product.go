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

// 2. Define the contract for the Repository (Database layer)
type ProductSearchRepository interface {
	SearchProducts(ctx context.Context, query string) ([]Product, error)
}

// 3. Define the contract for the Usecase (Business logic layer)
type ProductSearchUsecase interface {
	Search(ctx context.Context, query string) ([]Product, error)
}
