package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	_ "github.com/lib/pq"

	"github.com/che1nov/tea-shop/order-service/internal/model"
)

// OrderRepositoryInterface определяет методы репозитория
type OrderRepositoryInterface interface {
	CreateOrder(ctx context.Context, order *model.Order) error
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status string) error
	ListUserOrders(ctx context.Context, userID int64) ([]*model.Order, error)
}

type OrderRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	itemsJSON, err := json.Marshal(order.Items)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO orders (user_id, items, status, total_price, address, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now()
	return r.db.QueryRowContext(
		ctx,
		query,
		order.UserID,
		itemsJSON,
		order.Status,
		order.TotalPrice,
		order.Address,
		now,
		now,
	).Scan(&order.ID)
}

func (r *OrderRepository) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	query := `SELECT id, user_id, items, status, total_price, address, created_at, updated_at FROM orders WHERE id = $1`

	order := &model.Order{}
	var itemsJSON []byte

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&order.ID,
		&order.UserID,
		&itemsJSON,
		&order.Status,
		&order.TotalPrice,
		&order.Address,
		&order.CreatedAt,
		&order.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
		return nil, err
	}

	return order, nil
}

func (r *OrderRepository) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE orders SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *OrderRepository) ListUserOrders(ctx context.Context, userID int64) ([]*model.Order, error) {
	query := `SELECT id, user_id, items, status, total_price, address, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []*model.Order
	for rows.Next() {
		order := &model.Order{}
		var itemsJSON []byte

		if err := rows.Scan(
			&order.ID,
			&order.UserID,
			&itemsJSON,
			&order.Status,
			&order.TotalPrice,
			&order.Address,
			&order.CreatedAt,
			&order.UpdatedAt,
		); err != nil {
			return nil, err
		}

		if err := json.Unmarshal(itemsJSON, &order.Items); err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, rows.Err()
}
