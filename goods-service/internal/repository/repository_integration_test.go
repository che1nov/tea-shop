//go:build integration

package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/che1nov/tea-shop/goods-service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_INTEGRATION_TESTS=1 to run integration tests")
	}

	dsn := os.Getenv("GOODS_TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "user=user password=password dbname=goods_db host=localhost port=5433 sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	require.NoError(t, db.Ping())
	return db
}

func migrateGoodsIntegrationDB(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS goods (
			id SERIAL PRIMARY KEY,
			sku VARCHAR(50) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			stock INT NOT NULL DEFAULT 0,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE TABLE IF NOT EXISTS stock_reservations (
			id SERIAL PRIMARY KEY,
			good_id INT NOT NULL REFERENCES goods(id),
			order_id INT NOT NULL,
			quantity INT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
		TRUNCATE TABLE stock_reservations, goods RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func TestGoodsRepositoryIntegration_CreateListUpdateAndDeleteGood(t *testing.T) {
	db := openIntegrationDB(t)
	migrateGoodsIntegrationDB(t, db)

	repo := New(db)
	ctx := context.Background()

	good := &model.Good{Name: "Sencha", Description: "Green tea", Price: 12.5, Stock: 10}
	require.NoError(t, repo.CreateGood(ctx, good))
	require.NotZero(t, good.ID)
	require.Equal(t, "GOOD-000001", good.SKU)

	stored, err := repo.GetGood(ctx, good.ID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.Equal(t, good.Name, stored.Name)

	stored.Name = "Gyokuro"
	stored.Stock = 7
	require.NoError(t, repo.UpdateGood(ctx, stored))

	goods, err := repo.ListGoods(ctx, 10, 0)
	require.NoError(t, err)
	require.Len(t, goods, 1)
	require.Equal(t, "Gyokuro", goods[0].Name)

	total, err := repo.GetTotalGoods(ctx)
	require.NoError(t, err)
	require.Equal(t, int32(1), total)

	require.NoError(t, repo.DeleteGood(ctx, good.ID))
	deleted, err := repo.GetGood(ctx, good.ID)
	require.NoError(t, err)
	require.Nil(t, deleted)
}

func TestGoodsRepositoryIntegration_ReserveStock(t *testing.T) {
	db := openIntegrationDB(t)
	migrateGoodsIntegrationDB(t, db)

	repo := New(db)
	ctx := context.Background()

	good := &model.Good{Name: "Assam", Description: "Black tea", Price: 8, Stock: 5}
	require.NoError(t, repo.CreateGood(ctx, good))

	require.NoError(t, repo.ReserveStock(ctx, good.ID, 3, 1001))

	stored, err := repo.GetGood(ctx, good.ID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.Equal(t, int32(2), stored.Stock)

	err = repo.ReserveStock(ctx, good.ID, 3, 1002)
	require.True(t, errors.Is(err, sql.ErrNoRows), "expected sql.ErrNoRows, got %v", err)
}
