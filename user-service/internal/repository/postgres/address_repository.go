package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Yahya-idris-A/ecommerce-microservices/user-service/internal/domain"
	"github.com/google/uuid"
)

type addressRepository struct {
	db *sql.DB
}

func NewAddressRepository(db *sql.DB) domain.AddressRepository {
	return &addressRepository{db: db}
}

func (r *addressRepository) Create(ctx context.Context, addr *domain.Address) error {
	query := `
		INSERT INTO user_addresses (id, user_id, label, recipient_name, phone_number, full_address, city, province, postal_code, is_primary, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
	`
	_, err := r.db.ExecContext(ctx, query,
		addr.ID, addr.UserID, addr.Label, addr.RecipientName, addr.PhoneNumber,
		addr.FullAddress, addr.City, addr.Province, addr.PostalCode, addr.IsPrimary,
		addr.CreatedAt, addr.UpdatedAt,
	)
	return err
}

func (r *addressRepository) GetByUserID(ctx context.Context, userID uuid.UUID) ([]domain.Address, error) {
	query := `
		SELECT id, user_id, label, recipient_name, phone_number, full_address, city, province, postal_code, is_primary, created_at, updated_at
		FROM user_addresses WHERE user_id = $1
		ORDER BY is_primary DESC, created_at DESC
	`
	// ORDER BY is_primary DESC memastikan alamat utama selalu tampil paling atas di frontend

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var addresses []domain.Address
	for rows.Next() {
		var a domain.Address
		err := rows.Scan(
			&a.ID, &a.UserID, &a.Label, &a.RecipientName, &a.PhoneNumber,
			&a.FullAddress, &a.City, &a.Province, &a.PostalCode, &a.IsPrimary,
			&a.CreatedAt, &a.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, a)
	}
	return addresses, nil
}

func (r *addressRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Address, error) {
	query := `
		SELECT id, user_id, label, recipient_name, phone_number, full_address, city, province, postal_code, is_primary, created_at, updated_at
		FROM user_addresses WHERE id = $1
	`
	a := &domain.Address{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&a.ID, &a.UserID, &a.Label, &a.RecipientName, &a.PhoneNumber,
		&a.FullAddress, &a.City, &a.Province, &a.PostalCode, &a.IsPrimary,
		&a.CreatedAt, &a.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("address not found")
		}
		return nil, err
	}
	return a, nil
}

func (r *addressRepository) Update(ctx context.Context, addr *domain.Address) error {
	query := `
		UPDATE user_addresses 
		SET label = $1, recipient_name = $2, phone_number = $3, full_address = $4, city = $5, province = $6, postal_code = $7, is_primary = $8, updated_at = $9
		WHERE id = $10
	`
	_, err := r.db.ExecContext(ctx, query,
		addr.Label, addr.RecipientName, addr.PhoneNumber, addr.FullAddress,
		addr.City, addr.Province, addr.PostalCode, addr.IsPrimary, addr.UpdatedAt, addr.ID,
	)
	return err
}

func (r *addressRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM user_addresses WHERE id = $1`
	_, err := r.db.ExecContext(ctx, query, id)
	return err
}

func (r *addressRepository) RemovePrimaryStatus(ctx context.Context, userID uuid.UUID) error {
	query := `UPDATE user_addresses SET is_primary = false WHERE user_id = $1`
	_, err := r.db.ExecContext(ctx, query, userID)
	return err
}
