package usecase

import (
	"context"
	"errors"
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

func (u *categoryUsecase) CreateCategory(ctx context.Context, req *domain.CreateCategoryRequest, merchantID *uuid.UUID) (*domain.Category, error) {
	category := &domain.Category{
		ID:         uuid.New(),
		MerchantID: merchantID,
		Name:       req.Name,
		CreatedAt:  time.Now(),
	}

	if err := u.categoryRepo.Create(ctx, category); err != nil {
		return nil, err
	}

	return category, nil
}

func (u *categoryUsecase) GetCategoryByID(ctx context.Context, id uuid.UUID) (*domain.Category, error) {
	return u.categoryRepo.GetByID(ctx, id)
}

// Tambahkan fungsi ini di bagian bawah file usecase kamu
func (u *categoryUsecase) GetAllCategories(ctx context.Context, merchantID string) ([]domain.Category, error) {
	return u.categoryRepo.GetAll(ctx, merchantID)
}

func (u *categoryUsecase) DeleteCategory(ctx context.Context, role string, merchantID uuid.UUID, categoryID uuid.UUID) error {
	// 1. Cek dulu apakah kategori dengan ID tersebut ada
	category, err := u.categoryRepo.GetByID(ctx, categoryID)
	if err != nil {
		return err // Akan melempar error jika kategori tidak ditemukan
	}

	isGlobalCategory := category.MerchantID == nil

	if isGlobalCategory {
		// ATURAN 1: Kategori Global HANYA boleh dihapus oleh Admin
		if role != "admin" {
			return errors.New("unauthorized: only admin can delete global categories")
		}
	} else {
		// ATURAN 2: Kategori Merchant
		if role == "merchant" {
			// Jika dia merchant, pastikan kategori ini benar-benar miliknya
			if *category.MerchantID != merchantID {
				return errors.New("unauthorized: category does not belong to this merchant")
			}
		} else if role != "admin" {
			// Opsional: Admin biasanya punya hak "Dewa" untuk menghapus kategori merchant juga.
			// Tapi kalau user bukan admin dan bukan pemiliknya, tolak!
			return errors.New("unauthorized: you don't have permission to delete this category")
		}
	}

	// 2. Jika kategori ditemukan, lanjutkan dengan penghapusan
	return u.categoryRepo.Delete(ctx, category.ID)
}
