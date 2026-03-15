package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

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
func (r *postgresProductRepo) Create(ctx context.Context, product *domain.Product) error {
	// Query SQL murni untuk memasukkan data.
	// Kita menggunakan $1, $2, dst untuk mencegah SQL Injection (best practice!).
	query := `
		INSERT INTO products (id, merchant_id, category_id, name, slug, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	// Mengeksekusi query dengan mengirimkan data dari entitas Product
	_, err := r.db.ExecContext(ctx, query,
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
func (r *postgresProductRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Product, error) {
	query := `
		SELECT id, merchant_id, category_id, name, slug, description, price, stock, created_at, updated_at
		FROM products
		WHERE id = $1
	`

	product := &domain.Product{}

	// QueryRow digunakan karena kita hanya ekspektasi 1 baris data
	// Scan memindahkan hasil query (kolom) ke dalam struct product
	err := r.db.QueryRowContext(ctx, query, id).Scan(
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

// GetAll mengambil semua produk dari database
func (r *postgresProductRepo) GetAll(ctx context.Context, keyword string, limit int, cursorCreatedAt string, cursorID string) ([]domain.Product, string, error) {
	// 1. Siapkan kerangka dasar query (WHERE 1=1 adalah trik agar kita bisa menyambung AND dengan mudah)
	baseQuery := `
        SELECT id, merchant_id, category_id, name, slug, description, price, stock, created_at, updated_at
        FROM products
        WHERE 1=1
    `

	// Siapkan wadah untuk parameter ($1, $2, dst)
	var args []interface{}
	argCount := 1

	// 2. Jika ada Keyword, tambahkan filter pencarian
	if keyword != "" {
		// Karena argCount 1, ini akan menjadi: AND (name ILIKE $1 OR description ILIKE $1)
		baseQuery += fmt.Sprintf(` AND (name ILIKE $%d OR description ILIKE $%d)`, argCount, argCount)
		args = append(args, "%"+keyword+"%")
		argCount++ // Naikkan hitungan
	}

	// 3. Jika ada Cursor (halaman selanjutnya), tambahkan filter tuple (created_at, id)
	if cursorCreatedAt != "" && cursorID != "" {
		// Jika argCount sekarang 2, ini akan menjadi: AND (created_at, id) < ($2, $3)
		baseQuery += fmt.Sprintf(` AND (created_at, id) < ($%d, $%d)`, argCount, argCount+1)
		args = append(args, cursorCreatedAt, cursorID)
		argCount += 2
	}

	// 4. Selalu akhiri dengan ORDER BY dan LIMIT
	baseQuery += fmt.Sprintf(` ORDER BY created_at DESC, id DESC LIMIT $%d`, argCount)
	args = append(args, limit)

	// 5. Eksekusi query dengan parameter yang sudah dirakit
	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var products []domain.Product
	var lastCreatedAt time.Time
	var lastID uuid.UUID

	// 6. Looping data seperti biasa
	for rows.Next() {
		var p domain.Product
		err := rows.Scan(&p.ID, &p.MerchantID, &p.CategoryID, &p.Name, &p.Slug, &p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, "", err
		}
		products = append(products, p)

		// Simpan data terakhir untuk dijadikan cursor berikutnya
		lastCreatedAt = p.CreatedAt
		lastID = p.ID
	}

	// 7. Buat string Cursor baru jika data tidak kosong
	var nextCursor string
	if len(products) > 0 {
		// Format cursor: "waktu|id"
		nextCursor = fmt.Sprintf("%s|%s", lastCreatedAt.Format(time.RFC3339Nano), lastID.String())
	}

	return products, nextCursor, nil
}

// GetByMerchantID mengambil semua produk milik satu toko tertentu
func (r *postgresProductRepo) GetByMerchantID(ctx context.Context, merchantID uuid.UUID, limit int, cursorCreatedAt string, cursorID string) ([]domain.Product, string, error) {
	var rows *sql.Rows
	var err error

	// 1. Cek apakah ini halaman selanjutnya (cursor terisi)
	if cursorCreatedAt != "" && cursorID != "" {
		query := `
            SELECT id, merchant_id, category_id, name, slug, description, price, stock, created_at, updated_at
            FROM products
            WHERE merchant_id = $1 AND (created_at, id) < ($2, $3)
            ORDER BY created_at DESC, id DESC
            LIMIT $4
        `
		rows, err = r.db.QueryContext(ctx, query, merchantID, cursorCreatedAt, cursorID, limit)
	} else {
		// 2. Jika tidak ada cursor, ini adalah halaman pertama
		query := `
            SELECT id, merchant_id, category_id, name, slug, description, price, stock, created_at, updated_at
            FROM products
            WHERE merchant_id = $1
            ORDER BY created_at DESC, id DESC
            LIMIT $2
        `
		rows, err = r.db.QueryContext(ctx, query, merchantID, limit)
	}

	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var products []domain.Product
	var lastCreatedAt time.Time
	var lastID uuid.UUID

	for rows.Next() {
		var p domain.Product
		err := rows.Scan(
			&p.ID, &p.MerchantID, &p.CategoryID, &p.Name, &p.Slug,
			&p.Description, &p.Price, &p.Stock, &p.CreatedAt, &p.UpdatedAt,
		)
		if err != nil {
			return nil, "", err
		}
		products = append(products, p)

		// 3. Selalu simpan data terakhir di setiap iterasi
		lastCreatedAt = p.CreatedAt
		lastID = p.ID
	}

	// 4. Rakit cursor baru untuk dikirim ke frontend
	var nextCursor string
	if len(products) > 0 {
		nextCursor = fmt.Sprintf("%s|%s", lastCreatedAt.Format(time.RFC3339Nano), lastID.String())
	}

	return products, nextCursor, nil
}

func (r *postgresProductRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM products WHERE id = $1`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return err // Error dari database (misal koneksi mati)
	}

	// Cek statistiknya: ada berapa baris yang beneran kehapus?
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	// Kalau ternyata 0 baris, berarti ID-nya emang nggak ada!
	if rowsAffected == 0 {
		return errors.New("Product not found")
	}

	return nil
}
