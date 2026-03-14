package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/che1nov/tea-shop/goods-service/internal/model"
	"github.com/stretchr/testify/require"
)

func setupGoodsRepo(t *testing.T) (*GoodsRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return New(db), mock, func() { _ = db.Close() }
}

func TestCreateGoodAutoSKU(t *testing.T) {
	repo, mock, cleanup := setupGoodsRepo(t)
	defer cleanup()

	good := &model.Good{Name: "Tea", Description: "Green", Price: 100, Stock: 5}
	mock.ExpectQuery(regexp.QuoteMeta("SELECT COALESCE(MAX(id), 0) FROM goods")).
		WillReturnRows(sqlmock.NewRows([]string{"max"}).AddRow(int64(15)))
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO goods (sku, name, description, price, stock, created_at, updated_at)")).
		WithArgs("GOOD-000016", good.Name, good.Description, good.Price, good.Stock, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(16)))

	err := repo.CreateGood(context.Background(), good)
	require.NoError(t, err)
	require.Equal(t, int64(16), good.ID)
	require.Equal(t, "GOOD-000016", good.SKU)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetGoodNotFound(t *testing.T) {
	repo, mock, cleanup := setupGoodsRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, description, price, stock, created_at, updated_at FROM goods WHERE id = $1")).
		WithArgs(int64(1)).
		WillReturnError(sql.ErrNoRows)

	good, err := repo.GetGood(context.Background(), 1)
	require.NoError(t, err)
	require.Nil(t, good)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestListGoodsSuccess(t *testing.T) {
	repo, mock, cleanup := setupGoodsRepo(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, sku, name, description, price, stock, created_at, updated_at FROM goods ORDER BY id LIMIT $1 OFFSET $2")).
		WithArgs(int32(10), int32(0)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "sku", "name", "description", "price", "stock", "created_at", "updated_at"}).
			AddRow(int64(1), "GOOD-000001", "Tea", "Desc", 10.5, int32(3), now, now))

	goods, err := repo.ListGoods(context.Background(), 10, 0)
	require.NoError(t, err)
	require.Len(t, goods, 1)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestReserveStockInsufficient(t *testing.T) {
	repo, mock, cleanup := setupGoodsRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT stock FROM goods WHERE id = $1 FOR UPDATE")).
		WithArgs(int64(10)).
		WillReturnRows(sqlmock.NewRows([]string{"stock"}).AddRow(int32(1)))
	mock.ExpectRollback()

	err := repo.ReserveStock(context.Background(), 10, 2, 100)
	require.ErrorIs(t, err, sql.ErrNoRows)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteGoodSuccess(t *testing.T) {
	repo, mock, cleanup := setupGoodsRepo(t)
	defer cleanup()

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM stock_reservations WHERE good_id = $1")).
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectExec(regexp.QuoteMeta("DELETE FROM goods WHERE id = $1")).
		WithArgs(int64(5)).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err := repo.DeleteGood(context.Background(), 5)
	require.NoError(t, err)
	require.NoError(t, mock.ExpectationsWereMet())
}
