package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/che1nov/tea-shop/payment-service/internal/model"
	"github.com/stretchr/testify/require"
)

func setupPaymentRepo(t *testing.T) (*PaymentRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return New(db), mock, func() { _ = db.Close() }
}

func TestCreatePaymentSuccess(t *testing.T) {
	repo, mock, cleanup := setupPaymentRepo(t)
	defer cleanup()

	p := &model.Payment{OrderID: 1, Amount: 10, Status: "pending", Method: "card"}
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO payments (order_id, amount, status, method, created_at, updated_at)")).
		WithArgs(p.OrderID, p.Amount, p.Status, p.Method, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(3)))

	err := repo.CreatePayment(context.Background(), p)
	require.NoError(t, err)
	require.Equal(t, int64(3), p.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPaymentNotFound(t *testing.T) {
	repo, mock, cleanup := setupPaymentRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, amount, status, method, created_at, updated_at FROM payments WHERE id = $1")).
		WithArgs(int64(123)).
		WillReturnError(sql.ErrNoRows)

	p, err := repo.GetPayment(context.Background(), 123)
	require.NoError(t, err)
	require.Nil(t, p)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetPaymentByOrderIDSuccess(t *testing.T) {
	repo, mock, cleanup := setupPaymentRepo(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, amount, status, method, created_at, updated_at FROM payments WHERE order_id = $1")).
		WithArgs(int64(99)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "amount", "status", "method", "created_at", "updated_at"}).
			AddRow(int64(9), int64(99), 19.5, "completed", "card", now, now))

	p, err := repo.GetPaymentByOrderID(context.Background(), 99)
	require.NoError(t, err)
	require.NotNil(t, p)
	require.Equal(t, int64(99), p.OrderID)
	require.NoError(t, mock.ExpectationsWereMet())
}
