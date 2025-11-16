package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/che1nov/tea-shop/payment-service/internal/model"
	_ "github.com/lib/pq"
)

// PaymentRepositoryInterface определяет методы репозитория
type PaymentRepositoryInterface interface {
	CreatePayment(ctx context.Context, payment *model.Payment) error
	GetPayment(ctx context.Context, id int64) (*model.Payment, error)
	UpdatePaymentStatus(ctx context.Context, id int64, status string) error
	GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error)
}

type PaymentRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) CreatePayment(ctx context.Context, payment *model.Payment) error {
	query := `
		INSERT INTO payments (order_id, amount, status, method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	return r.db.QueryRowContext(
		ctx,
		query,
		payment.OrderID,
		payment.Amount,
		payment.Status,
		payment.Method,
		now,
		now,
	).Scan(&payment.ID)
}

func (r *PaymentRepository) GetPayment(ctx context.Context, id int64) (*model.Payment, error) {
	query := `SELECT id, order_id, amount, status, method, created_at, updated_at FROM payments WHERE id = $1`

	payment := &model.Payment{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Status,
		&payment.Method,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}

func (r *PaymentRepository) UpdatePaymentStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE payments SET status = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.ExecContext(ctx, query, status, time.Now(), id)
	return err
}

func (r *PaymentRepository) GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error) {
	query := `SELECT id, order_id, amount, status, method, created_at, updated_at FROM payments WHERE order_id = $1`

	payment := &model.Payment{}
	err := r.db.QueryRowContext(ctx, query, orderID).Scan(
		&payment.ID,
		&payment.OrderID,
		&payment.Amount,
		&payment.Status,
		&payment.Method,
		&payment.CreatedAt,
		&payment.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return payment, nil
}
