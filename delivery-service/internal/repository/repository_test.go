package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/che1nov/tea-shop/delivery-service/internal/model"
	"github.com/stretchr/testify/require"
)

func setupDeliveryRepo(t *testing.T) (*DeliveryRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return New(db), mock, func() { _ = db.Close() }
}

func TestCreateDeliverySuccess(t *testing.T) {
	repo, mock, cleanup := setupDeliveryRepo(t)
	defer cleanup()

	delivery := &model.Delivery{OrderID: 1, Address: "addr", Status: "pending"}
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO deliveries (order_id, address, status, created_at, updated_at)")).
		WithArgs(delivery.OrderID, delivery.Address, delivery.Status, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(10)))

	err := repo.CreateDelivery(context.Background(), delivery)
	require.NoError(t, err)
	require.Equal(t, int64(10), delivery.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetDeliveryNotFound(t *testing.T) {
	repo, mock, cleanup := setupDeliveryRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, address, status, created_at, updated_at FROM deliveries WHERE id = $1")).
		WithArgs(int64(5)).
		WillReturnError(sql.ErrNoRows)

	got, err := repo.GetDelivery(context.Background(), 5)
	require.NoError(t, err)
	require.Nil(t, got)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListDeliveriesWithFilter(t *testing.T) {
	repo, mock, cleanup := setupDeliveryRepo(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, order_id, address, status, created_at, updated_at FROM deliveries WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3")).
		WithArgs("pending", int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "order_id", "address", "status", "created_at", "updated_at"}).
			AddRow(int64(1), int64(2), "addr", "pending", now, now))

	list, err := repo.ListDeliveries(context.Background(), 10, 0, "pending")
	require.NoError(t, err)
	require.Len(t, list, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

