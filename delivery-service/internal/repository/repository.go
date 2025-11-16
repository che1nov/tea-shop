package repository

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/lib/pq"
	
	"github.com/che1nov/tea-shop/delivery-service/internal/model"
)

// DeliveryRepositoryInterface определяет методы репозитория
type DeliveryRepositoryInterface interface {
	CreateDelivery(ctx context.Context, delivery *model.Delivery) error
	GetDelivery(ctx context.Context, id int64) (*model.Delivery, error)
	GetDeliveryByOrderID(ctx context.Context, orderID int64) (*model.Delivery, error)
	UpdateDeliveryStatus(ctx context.Context, id int64, status string) error
	ListDeliveries(ctx context.Context, limit, offset int32, status string) ([]*model.Delivery, error)
	GetTotalDeliveries(ctx context.Context, status string) (int32, error)
}

type DeliveryRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *DeliveryRepository {
	return &DeliveryRepository{db: db}
}

func (r *DeliveryRepository) CreateDelivery(ctx context.Context, delivery *model.Delivery) error {
	query := `
		INSERT INTO deliveries (order_id, address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	now := time.Now()
	return r.db.QueryRowContext(
		ctx,
		query,
		delivery.OrderID,
		delivery.Address,
		delivery.Status,
		now,
		now,
	).Scan(&delivery.ID)
}

func (r *DeliveryRepository) GetDelivery(ctx context.Context, id int64) (*model.Delivery, error) {
	query := `SELECT id, order_id, address, status, created_at, updated_at FROM deliveries WHERE id = $1`

	delivery := &model.Delivery{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&delivery.ID,
		&delivery.OrderID,
		&delivery.Address,
		&delivery.Status,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (r *DeliveryRepository) GetDeliveryByOrderID(ctx context.Context, orderID int64) (*model.Delivery, error) {
	query := `SELECT id, order_id, address, status, created_at, updated_at FROM deliveries WHERE order_id = $1`

	delivery := &model.Delivery{}
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&delivery.ID,
		&delivery.OrderID,
		&delivery.Address,
		&delivery.Status,
		&delivery.CreatedAt,
		&delivery.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return delivery, nil
}

func (r *DeliveryRepository) UpdateDeliveryStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE deliveries SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *DeliveryRepository) ListDeliveries(ctx context.Context, limit, offset int32, statusFilter string) ([]*model.Delivery, error) {
	var query string
	var rows *sql.Rows
	var err error

	if statusFilter != "" {
		query = `SELECT id, order_id, address, status, created_at, updated_at FROM deliveries WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
		rows, err = r.db.QueryContext(ctx, query, statusFilter, limit, offset)
	} else {
		query = `SELECT id, order_id, address, status, created_at, updated_at FROM deliveries ORDER BY created_at DESC LIMIT $1 OFFSET $2`
		rows, err = r.db.QueryContext(ctx, query, limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var deliveries []*model.Delivery
	for rows.Next() {
		delivery := &model.Delivery{}
		if err := rows.Scan(
			&delivery.ID,
			&delivery.OrderID,
			&delivery.Address,
			&delivery.Status,
			&delivery.CreatedAt,
			&delivery.UpdatedAt,
		); err != nil {
			return nil, err
		}
		deliveries = append(deliveries, delivery)
	}

	return deliveries, rows.Err()
}

func (r *DeliveryRepository) GetTotalDeliveries(ctx context.Context, statusFilter string) (int32, error) {
	var query string
	var total int32

	if statusFilter != "" {
		query = `SELECT COUNT(*) FROM deliveries WHERE status = $1`
		err := r.db.QueryRowContext(ctx, query, statusFilter).Scan(&total)
		return total, err
	}

	query = `SELECT COUNT(*) FROM deliveries`
	err := r.db.QueryRowContext(ctx, query).Scan(&total)
	return total, err
}
