//go:build integration

package repository

import (
	"context"
	"database/sql"
	"os"
	"testing"

	"github.com/che1nov/tea-shop/payment-service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_INTEGRATION_TESTS=1 to run integration tests")
	}

	dsn := os.Getenv("PAYMENTS_TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "user=user password=password dbname=payments_db host=localhost port=5435 sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	require.NoError(t, db.Ping())
	return db
}

func migratePaymentsIntegrationDB(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			order_id INT NOT NULL UNIQUE,
			amount DECIMAL(10, 2) NOT NULL,
			status VARCHAR(50) NOT NULL,
			method VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		TRUNCATE TABLE payments RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func TestPaymentRepositoryIntegration_CreateReadAndUpdate(t *testing.T) {
	db := openIntegrationDB(t)
	migratePaymentsIntegrationDB(t, db)

	repo := New(db)
	ctx := context.Background()

	payment := &model.Payment{
		OrderID: 100,
		Amount:  50,
		Status:  "pending",
		Method:  "card",
	}

	require.NoError(t, repo.CreatePayment(ctx, payment))
	require.NotZero(t, payment.ID)

	byID, err := repo.GetPayment(ctx, payment.ID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	require.Equal(t, int64(100), byID.OrderID)

	require.NoError(t, repo.UpdatePaymentStatus(ctx, payment.ID, "completed"))

	byOrderID, err := repo.GetPaymentByOrderID(ctx, payment.OrderID)
	require.NoError(t, err)
	require.NotNil(t, byOrderID)
	require.Equal(t, "completed", byOrderID.Status)
}
