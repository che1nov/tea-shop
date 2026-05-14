//go:build integration

package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/che1nov/tea-shop/order-service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_INTEGRATION_TESTS=1 to run integration tests")
	}

	dsn := os.Getenv("ORDERS_TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "user=user password=password dbname=orders_db host=localhost port=5434 sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	require.NoError(t, db.Ping())
	return db
}

func migrateOrdersIntegrationDB(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			items JSONB NOT NULL,
			status VARCHAR(50) NOT NULL,
			total_price DECIMAL(10, 2) NOT NULL,
			address TEXT,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		TRUNCATE TABLE orders RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func TestOrderRepositoryIntegration_CreateReadUpdateAndList(t *testing.T) {
	db := openIntegrationDB(t)
	migrateOrdersIntegrationDB(t, db)

	repo := New(db)
	ctx := context.Background()

	order := &model.Order{
		UserID:     42,
		Status:     "pending",
		TotalPrice: 25,
		Address:    "Samara",
		Items: []model.OrderItem{
			{GoodID: 1, Quantity: 2, Price: 12.5},
		},
	}

	require.NoError(t, repo.CreateOrder(ctx, order))
	require.NotZero(t, order.ID)

	stored, err := repo.GetOrder(ctx, order.ID)
	require.NoError(t, err)
	require.NotNil(t, stored)
	require.Equal(t, int64(42), stored.UserID)
	require.Len(t, stored.Items, 1)
	require.Equal(t, "pending", stored.Status)

	require.NoError(t, repo.UpdateOrderStatus(ctx, order.ID, "paid"))

	updated, err := repo.GetOrder(ctx, order.ID)
	require.NoError(t, err)
	require.NotNil(t, updated)
	require.Equal(t, "paid", updated.Status)

	orders, err := repo.ListUserOrders(ctx, 42)
	require.NoError(t, err)
	require.Len(t, orders, 1)
	require.Equal(t, order.ID, orders[0].ID)
}
