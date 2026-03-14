package usecase

import (
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
	"github.com/google/uuid"
)

type categoryUsecase struct {
	categoryRepo domain.CategoryRepository
}

func NewCategoryUsecase(repo domain.CategoryRepository) domain.CategoryUsecase {
	return &categoryUsecase{
		categoryRepo: repo,
	}
}

func (u *categoryUsecase) CreateCategory(req *domain.CreateCategoryRequest, merchantID *uuid.UUID) (*domain.Category, error) {
	category := &domain.Category{
		ID:         uuid.New(),
		MerchantID: merchantID,
		Name:       req.Name,
		CreatedAt:  time.Now(),
	}

	if err := u.categoryRepo.Create(category); err != nil {
		return nil, err
	}

	return category, nil
}

func (u *categoryUsecase) GetCategoryByID(id uuid.UUID) (*domain.Category, error) {
	return u.categoryRepo.GetByID(id)
}

// Tambahkan fungsi ini di bagian bawah file usecase kamu
func (u *categoryUsecase) GetAllCategories(merchantID string) ([]domain.Category, error) {
	return u.categoryRepo.GetAll(merchantID)
}
