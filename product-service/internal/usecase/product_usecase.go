package usecase

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	// Ingat untuk menyesuaikan path import ini dengan nama module kamu
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
)

// productUsecase struct ini sengaja dibuat private (huruf kecil di awal).
// Layer lain tidak bisa langsung memanggil struct ini,
// mereka harus lewat interface domain.ProductUsecase.
type productUsecase struct {
	productRepo  domain.ProductRepository
	merchantRepo domain.MerchantRepository
}

// NewProductUsecase adalah constructor (Penerapan Dependency Injection).
// Kita "menyuntikkan" repository database ke dalam usecase.
func NewProductUsecase(repo domain.ProductRepository, merchantRepo domain.MerchantRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo:  repo,
		merchantRepo: merchantRepo,
	}
}

// CreateProduct menangani semua logika bisnis sebelum data disimpan.
func (u *productUsecase) CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error) {
	// 1. Validasi Bisnis Dasar
	if req.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	// 2. Logika Bisnis: Membuat Slug otomatis dari Nama Produk
	productID := uuid.New()
	baseSlug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))
	uniqueSlug := fmt.Sprintf("%s-%s", baseSlug, productID.String()[:8])

	// 3. Memetakan (Mapping) DTO dari user menjadi Entitas Domain
	product := &domain.Product{
		ID:          productID,
		MerchantID:  req.MerchantID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        uniqueSlug,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	// 4. Memanggil Repository untuk menyimpan data ke database (PostgreSQL nantinya)
	err := u.productRepo.Create(product)
	if err != nil {
		return nil, err // Mengembalikan error ke layer presentasi jika database gagal
	}

	return product, nil
}

// GetProductByID sekadar meneruskan permintaan pencarian ke repository
func (u *productUsecase) GetProductByID(id uuid.UUID) (*domain.Product, error) {
	return u.productRepo.GetByID(id)
}

// Tambahkan fungsi ini di bagian bawah file usecase kamu
func (u *productUsecase) GetAllProducts(merchantID string, keyword string, limit int, cursorCreatedAt string, cursorID string) ([]domain.Product, string, error) {
	if merchantID != "" {
		parsedID, err := uuid.Parse(merchantID)
		if err != nil {
			return nil, "", err // URL tidak valid
		}
		return u.productRepo.GetByMerchantID(parsedID, limit, cursorCreatedAt, cursorID)
	}
	return u.productRepo.GetAll(keyword, limit, cursorCreatedAt, cursorID)
}

func (u *productUsecase) GetMyProducts(userID uuid.UUID, limit int, cursorCreatedAt string, cursorID string) ([]domain.Product, string, error) {
	// 1. Cari dulu profil tokonya berdasarkan userID
	merchant, err := u.merchantRepo.GetByUserID(userID)
	if err != nil {
		return nil, "", err // Akan melempar error jika user belum buat toko
	}

	// 2. Jika toko ketemu, ambil semua produk berdasarkan ID toko tersebut
	return u.productRepo.GetByMerchantID(merchant.ID, limit, cursorCreatedAt, cursorID)
}
