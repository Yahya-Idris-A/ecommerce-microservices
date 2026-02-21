package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"
	// Ingat sesuaikan path import ini
	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
)

// postgresProductRepo adalah implementasi konkrit dari interface domain.ProductRepository
type postgresProductRepo struct {
	db *sql.DB // Kita menggunakan standar library database/sql bawaan Golang
}

// NewPostgresProductRepository adalah constructor untuk menginisialisasi repository
func NewPostgresProductRepository(db *sql.DB) domain.ProductRepository {
	return &postgresProductRepo{
		db: db,
	}
}

// Create menyimpan produk baru ke PostgreSQL
func (r *postgresProductRepo) Create(product *domain.Product) error {
	// Query SQL murni untuk memasukkan data.
	// Kita menggunakan $1, $2, dst untuk mencegah SQL Injection (best practice!).
	query := `
		INSERT INTO products (id, merchant_id, category_id, name, slug, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Mengeksekusi query dengan mengirimkan data dari entitas Product
	_, err := r.db.Exec(query,
		product.ID,
		product.MerchantID,
		product.CategoryID,
		product.Name,
		product.Slug,
		product.Description,
		product.Price,
		product.Stock,
		product.CreatedAt,
		product.UpdatedAt,
	)

	if err != nil {
		// Nanti kita bisa tambahkan custom error logger di sini
		return err
	}

	return nil
}

// GetByID mengambil satu produk berdasarkan UUID
func (r *postgresProductRepo) GetByID(id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, merchant_id, category_id, name, slug, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &domain.Product{}

	// QueryRow digunakan karena kita hanya ekspektasi 1 baris data
	// Scan memindahkan hasil query (kolom) ke dalam struct product
	err := r.db.QueryRow(query, id).Scan(
		&product.ID,
		&product.MerchantID,
		&product.CategoryID,
		&product.Name,
		&product.Slug,
		&product.Description,
		&product.Price,
		&product.Stock,
		&product.CreatedAt,
		&product.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Jika data tidak ditemukan, kembalikan custom error
			return nil, errors.New("product not found")
		}
		return nil, err
	}

	return product, nil
}
