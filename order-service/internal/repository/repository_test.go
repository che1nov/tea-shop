package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/che1nov/tea-shop/order-service/internal/model"
	"github.com/stretchr/testify/require"
)

func setupOrderRepo(t *testing.T) (*OrderRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return New(db), mock, func() { _ = db.Close() }
}

func TestCreateOrderSuccess(t *testing.T) {
	repo, mock, cleanup := setupOrderRepo(t)
	defer cleanup()

	order := &model.Order{
		UserID:     1,
		Items:      []model.OrderItem{{GoodID: 2, Quantity: 3, Price: 9.9}},
		Status:     "pending",
		TotalPrice: 29.7,
		Address:    "addr",
	}

	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO orders (user_id, items, status, total_price, address, created_at, updated_at)")).
		WithArgs(order.UserID, sqlmock.AnyArg(), order.Status, order.TotalPrice, order.Address, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(55)))

	err := repo.CreateOrder(context.Background(), order)
	require.NoError(t, err)
	require.Equal(t, int64(55), order.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetOrderNotFound(t *testing.T) {
	repo, mock, cleanup := setupOrderRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, items, status, total_price, address, created_at, updated_at FROM orders WHERE id = $1")).
		WithArgs(int64(77)).
		WillReturnError(sql.ErrNoRows)

	order, err := repo.GetOrder(context.Background(), 77)
	require.NoError(t, err)
	require.Nil(t, order)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListUserOrdersSuccess(t *testing.T) {
	repo, mock, cleanup := setupOrderRepo(t)
	defer cleanup()

	now := time.Now()
	itemsJSON := `[{"good_id":1,"quantity":2,"price":5.5}]`
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, user_id, items, status, total_price, address, created_at, updated_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC")).
		WithArgs(int64(9)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "user_id", "items", "status", "total_price", "address", "created_at", "updated_at"}).
			AddRow(int64(1), int64(9), []byte(itemsJSON), "pending", 11.0, "a", now, now))

	orders, err := repo.ListUserOrders(context.Background(), 9)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	require.Len(t, orders[0].Items, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

