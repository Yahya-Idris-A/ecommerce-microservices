package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/Yahya-idris-A/ecommerce-microservices/order-service/internal/domain"
	"github.com/google/uuid"
)

type cartRepository struct {
	db *sql.DB
}

func NewCartRepository(db *sql.DB) domain.CartRepository {
	return &cartRepository{db: db}
}

// ==========================================
// 1. URUSAN PAYUNG CART
// ==========================================

func (r *cartRepository) GetCartByUserID(ctx context.Context, userID uuid.UUID) (*domain.Cart, error) {
	query := `SELECT id, user_id, created_at, updated_at FROM carts WHERE user_id = $1`

	cart := &domain.Cart{}
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&cart.ID, &cart.UserID, &cart.CreatedAt, &cart.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("cart not found")
		}
		return nil, err
	}

	return cart, nil
}

func (r *cartRepository) CreateCart(ctx context.Context, cart *domain.Cart) error {
	query := `INSERT INTO carts (id, user_id, created_at, updated_at) VALUES ($1, $2, $3, $4)`

	_, err := r.db.ExecContext(ctx, query, cart.ID, cart.UserID, cart.CreatedAt, cart.UpdatedAt)
	return err
}

// ==========================================
// 2. URUSAN ISI CART (ITEMS)
// ==========================================

func (r *cartRepository) GetItemByCartAndProduct(ctx context.Context, cartID uuid.UUID, productID uuid.UUID) (*domain.CartItem, error) {
	query := `
		SELECT id, cart_id, product_id, merchant_id, quantity, created_at, updated_at 
		FROM cart_items 
		WHERE cart_id = $1 AND product_id = $2
	`

	item := &domain.CartItem{}
	err := r.db.QueryRowContext(ctx, query, cartID, productID).Scan(
		&item.ID, &item.CartID, &item.ProductID, &item.MerchantID,
		&item.Quantity, &item.CreatedAt, &item.UpdatedAt,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil // Return nil, nil artinya barang belum ada di keranjang
		}
		return nil, err
	}

	return item, nil
}

func (r *cartRepository) GetItemsByCartID(ctx context.Context, cartID uuid.UUID) ([]domain.CartItem, error) {
	// ORDER BY merchant_id adalah trik agar grouping di Usecase nanti lebih gampang dan cepat!
	query := `
		SELECT id, cart_id, product_id, merchant_id, quantity, created_at, updated_at 
		FROM cart_items 
		WHERE cart_id = $1
		ORDER BY merchant_id, created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, cartID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []domain.CartItem
	for rows.Next() {
		var i domain.CartItem
		err := rows.Scan(
			&i.ID, &i.CartID, &i.ProductID, &i.MerchantID,
			&i.Quantity, &i.CreatedAt, &i.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		items = append(items, i)
	}

	return items, nil
}

func (r *cartRepository) AddItem(ctx context.Context, item *domain.CartItem) error {
	query := `
		INSERT INTO cart_items (id, cart_id, product_id, merchant_id, quantity, created_at, updated_at) 
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.ExecContext(ctx, query,
		item.ID, item.CartID, item.ProductID, item.MerchantID,
		item.Quantity, item.CreatedAt, item.UpdatedAt,
	)
	return err
}

func (r *cartRepository) UpdateItemQuantity(ctx context.Context, itemID uuid.UUID, quantity int) error {
	query := `UPDATE cart_items SET quantity = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2`

	_, err := r.db.ExecContext(ctx, query, quantity, itemID)
	return err
}

func (r *cartRepository) DeleteItem(ctx context.Context, itemID uuid.UUID) error {
	query := `DELETE FROM cart_items WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, itemID)
	return err
}

func (r *cartRepository) DeleteItemsByCartID(ctx context.Context, cartID uuid.UUID) error {
	// Fungsi sapu bersih ini akan dipakai kalau user berhasil Checkout atau bayar
	query := `DELETE FROM cart_items WHERE cart_id = $1`

	_, err := r.db.ExecContext(ctx, query, cartID)
	return err
}
