package usecase

import (
	"context"
	"errors"
	"strings"

	"github.com/Yahya-idris-A/ecommerce-microservices/search-service/internal/domain"
)

type productSearchUsecase struct {
	searchRepo domain.ProductSearchRepository
}

func NewProductSearchUsecase(repo domain.ProductSearchRepository) domain.ProductSearchUsecase {
	return &productSearchUsecase{
		searchRepo: repo,
	}
}

func (u *productSearchUsecase) Search(ctx context.Context, query string) ([]domain.Product, error) {
	// Business Logic / Validation: Clean up spaces and check if empty
	cleanQuery := strings.TrimSpace(query)
	if cleanQuery == "" {
		return nil, errors.New("[ERROR] Search query cannot be empty")
	}

	// If valid, pass the task to the Repository layer
	return u.searchRepo.SearchProducts(ctx, cleanQuery)
}
