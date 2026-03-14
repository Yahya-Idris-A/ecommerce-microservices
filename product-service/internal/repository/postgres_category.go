package repository

import (
	"database/sql"
	"errors"

	"github.com/google/uuid"

	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
)

type postgresCategoryRepo struct {
	db *sql.DB
}

func NewPostgresCategoryRepository(db *sql.DB) domain.CategoryRepository {
	return &postgresCategoryRepo{
		db: db,
	}
}

func (r *postgresCategoryRepo) Create(category *domain.Category) error {
	query := `
		INSERT INTO categories (id, merchant_id, name, created_at)
		VALUES ($1, $2, $3, $4)
	`
	_, err := r.db.Exec(query, category.ID, category.MerchantID, category.Name, category.CreatedAt)
	return err
}

// GetByID mengambil satu kategori berdasarkan UUID
func (r *postgresCategoryRepo) GetByID(id uuid.UUID) (*domain.Category, error) {
	query := `
		SELECT id, merchant_id, name, created_at
		FROM categories
		WHERE id = $1
	`

	category := &domain.Category{}

	// QueryRow digunakan karena kita hanya ekspektasi 1 baris data
	// Scan memindahkan hasil query (kolom) ke dalam struct category
	err := r.db.QueryRow(query, id).Scan(
		&category.ID,
		&category.MerchantID,
		&category.Name,
		&category.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			// Jika data tidak ditemukan, kembalikan custom error
			return nil, errors.New("category not found")
		}
		return nil, err
	}

	return category, nil
}

// GetAll mengambil semua kategori dari database
func (r *postgresCategoryRepo) GetAll(merchantID string) ([]domain.Category, error) {
	var query string
	var rows *sql.Rows
	var err error

	if merchantID != "" {
		// Jika melihat toko spesifik: Ambil Kategori Global (NULL) + Kategori Toko Tersebut
		query = `SELECT id, merchant_id, name, created_at FROM categories WHERE merchant_id = $1 ORDER BY name ASC`
		rows, err = r.db.Query(query, merchantID)
	} else {
		// Jika di halaman utama: Hanya ambil Kategori Global (NULL)
		query = `SELECT id, merchant_id, name, created_at FROM categories WHERE merchant_id IS NULL ORDER BY name ASC`
		rows, err = r.db.Query(query)
	}

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var categories []domain.Category
	for rows.Next() {
		var c domain.Category
		err := rows.Scan(
			&c.ID,
			&c.MerchantID,
			&c.Name,
			&c.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		categories = append(categories, c)
	}

	return categories, nil
}
