package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/Yahya-idris-A/ecommerce-microservices/product-service/internal/domain"
)

type postgresMerchantRepo struct {
	db *sql.DB
}

func NewPostgresMerchantRepository(db *sql.DB) domain.MerchantRepository {
	return &postgresMerchantRepo{
		db: db,
	}
}

func (r *postgresMerchantRepo) Create(merchant *domain.Merchant) error {
	query := `
		INSERT INTO merchants (id, user_id, name, slug, description, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.Exec(query,
		merchant.ID,
		merchant.UserID,
		merchant.Name,
		merchant.Slug,
		merchant.Description,
		merchant.CreatedAt,
		merchant.UpdatedAt,
	)
	return err
}

func (r *postgresMerchantRepo) GetByUserID(userID uuid.UUID) (*domain.Merchant, error) {
	query := `
		SELECT id, user_id, name, slug, description, created_at, updated_at
		FROM merchants
		WHERE user_id = $1
	`

	merchant := &domain.Merchant{}
	err := r.db.QueryRow(query, userID).Scan(
		&merchant.ID, &merchant.UserID, &merchant.Name, &merchant.Slug,
		&merchant.Description, &merchant.CreatedAt, &merchant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("merchant profile not found")
		}
		return nil, err
	}

	return merchant, nil
}

// GetByID mengambil profil toko spesifik berdasarkan UUID toko
func (r *postgresMerchantRepo) GetByID(id uuid.UUID) (*domain.Merchant, error) {
	query := `SELECT id, user_id, name, slug, description, created_at, updated_at FROM merchants WHERE id = $1`

	merchant := &domain.Merchant{}
	err := r.db.QueryRow(query, id).Scan(
		&merchant.ID, &merchant.UserID, &merchant.Name, &merchant.Slug,
		&merchant.Description, &merchant.CreatedAt, &merchant.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("merchant not found")
		}
		return nil, err
	}

	return merchant, nil
}

// GetAll mengambil semua toko, dan akan melakukan pencarian jika ada keyword
func (r *postgresMerchantRepo) GetAll(keyword string, limit int, cursorCreatedAt string, cursorID string) ([]domain.Merchant, string, error) {
	baseQuery := `
        SELECT id, user_id, name, slug, description, created_at, updated_at 
		FROM merchants
        WHERE 1=1
    `

	var args []interface{}
	argCount := 1

	if keyword != "" {
		// Karena argCount 1, ini akan menjadi: AND (name ILIKE $1 OR description ILIKE $1)
		baseQuery += fmt.Sprintf(` AND (name ILIKE $%d OR description ILIKE $%d)`, argCount, argCount)
		args = append(args, "%"+keyword+"%")
		argCount++ // Naikkan hitungan
	}

	if cursorCreatedAt != "" && cursorID != "" {
		// Jika argCount sekarang 2, ini akan menjadi: AND (created_at, id) < ($2, $3)
		baseQuery += fmt.Sprintf(` AND (created_at, id) < ($%d, $%d)`, argCount, argCount+1)
		args = append(args, cursorCreatedAt, cursorID)
		argCount += 2
	}

	baseQuery += fmt.Sprintf(` ORDER BY created_at DESC, id DESC LIMIT $%d`, argCount)
	args = append(args, limit)

	rows, err := r.db.Query(baseQuery, args...)
	if err != nil {
		return nil, "", err
	}
	defer rows.Close()

	var merchants []domain.Merchant
	var lastCreatedAt time.Time
	var lastID uuid.UUID

	// 6. Looping data seperti biasa
	for rows.Next() {
		var p domain.Merchant
		err := rows.Scan(&p.ID, &p.UserID, &p.Name, &p.Slug, &p.Description, &p.CreatedAt, &p.UpdatedAt)
		if err != nil {
			return nil, "", err
		}
		merchants = append(merchants, p)

		// Simpan data terakhir untuk dijadikan cursor berikutnya
		lastCreatedAt = p.CreatedAt
		lastID = p.ID
	}

	// 7. Buat string Cursor baru jika data tidak kosong
	var nextCursor string
	if len(merchants) > 0 {
		// Format cursor: "waktu|id"
		nextCursor = fmt.Sprintf("%s|%s", lastCreatedAt.Format(time.RFC3339Nano), lastID.String())
	}

	return merchants, nextCursor, nil
}
