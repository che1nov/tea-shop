//go:build integration

package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/che1nov/tea-shop/delivery-service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_INTEGRATION_TESTS=1 to run integration tests")
	}

	dsn := os.Getenv("DELIVERY_TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "user=user password=password dbname=deliveries_db host=localhost port=5436 sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	require.NoError(t, db.Ping())
	return db
}

func migrateDeliveryIntegrationDB(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS deliveries (
			id SERIAL PRIMARY KEY,
			order_id INT NOT NULL UNIQUE,
			address VARCHAR(500) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		TRUNCATE TABLE deliveries RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func TestDeliveryRepositoryIntegration_CreateReadUpdateAndList(t *testing.T) {
	db := openIntegrationDB(t)
	migrateDeliveryIntegrationDB(t, db)

	repo := New(db)
	ctx := context.Background()

	delivery := &model.Delivery{
		OrderID: 200,
		Address: "Samara",
		Status:  "pending",
	}

	require.NoError(t, repo.CreateDelivery(ctx, delivery))
	require.NotZero(t, delivery.ID)

	byID, err := repo.GetDelivery(ctx, delivery.ID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	require.Equal(t, int64(200), byID.OrderID)

	byOrderID, err := repo.GetDeliveryByOrderID(ctx, delivery.OrderID)
	require.NoError(t, err)
	require.NotNil(t, byOrderID)
	require.Equal(t, delivery.ID, byOrderID.ID)

	require.NoError(t, repo.UpdateDeliveryStatus(ctx, delivery.ID, "shipped"))

	deliveries, err := repo.ListDeliveries(ctx, 10, 0, "shipped")
	require.NoError(t, err)
	require.Len(t, deliveries, 1)

	total, err := repo.GetTotalDeliveries(ctx, "shipped")
	require.NoError(t, err)
	require.Equal(t, int32(1), total)
}
