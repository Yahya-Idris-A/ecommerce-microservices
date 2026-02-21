package usecase

import (
	"errors"
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
	productRepo domain.ProductRepository
}

// NewProductUsecase adalah constructor (Penerapan Dependency Injection).
// Kita "menyuntikkan" repository database ke dalam usecase.
func NewProductUsecase(repo domain.ProductRepository) domain.ProductUsecase {
	return &productUsecase{
		productRepo: repo,
	}
}

// CreateProduct menangani semua logika bisnis sebelum data disimpan.
func (u *productUsecase) CreateProduct(req *domain.CreateProductRequest) (*domain.Product, error) {
	// 1. Validasi Bisnis Dasar
	if req.Price <= 0 {
		return nil, errors.New("price must be greater than zero")
	}

	// 2. Logika Bisnis: Membuat Slug otomatis dari Nama Produk
	// Contoh: "Kopi Susu Gayo" -> "kopi-susu-gayo"
	slug := strings.ToLower(strings.ReplaceAll(req.Name, " ", "-"))

	// 3. Memetakan (Mapping) DTO dari user menjadi Entitas Domain
	product := &domain.Product{
		ID:          uuid.New(), // Generate UUID baru secara otomatis
		MerchantID:  req.MerchantID,
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Slug:        slug,
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
