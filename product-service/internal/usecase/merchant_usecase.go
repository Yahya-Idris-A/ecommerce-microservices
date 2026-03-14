package usecase

import (
	"strings"
	"time"

	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
	"github.com/google/uuid"
)

type merchantUsecase struct {
	merchantRepo domain.MerchantRepository
}

func NewMerchantUsecase(repo domain.MerchantRepository) domain.MerchantUsecase {
	return &merchantUsecase{
		merchantRepo: repo,
	}
}

func (u *merchantUsecase) CreateMerchant(userID uuid.UUID, req *domain.CreateMerchantRequest) (*domain.Merchant, error) {
	// Membuat slug sederhana (ubah ke lowercase dan ganti spasi dengan strip)
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	merchant := &domain.Merchant{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        req.Name,
		Slug:        slug,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := u.merchantRepo.Create(merchant); err != nil {
		return nil, err
	}

	return merchant, nil
}

func (u *merchantUsecase) GetMerchantByID(id uuid.UUID) (*domain.Merchant, error) {
	return u.merchantRepo.GetByID(id)
}

func (u *merchantUsecase) GetByUserID(userID uuid.UUID) (*domain.Merchant, error) {
	// Di sini kita bisa menambahkan validasi bisnis lain jika perlu,
	// tapi untuk sekarang kita langsung lemparkan ke Repository.
	return u.merchantRepo.GetByUserID(userID)
}

func (u *merchantUsecase) GetAllMerchants(keyword string, limit int, cursorCreatedAt string, cursorID string) ([]domain.Merchant, string, error) {
	return u.merchantRepo.GetAll(keyword, limit, cursorCreatedAt, cursorID)
}
