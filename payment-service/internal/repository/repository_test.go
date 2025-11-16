package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/che1nov/tea-shop/payment-service/internal/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDBWithCleanup создает БД и очищает её перед тестом
func setupTestDBWithCleanup(t *testing.T) *sql.DB {
	connStr := "user=user password=password dbname=payments_db host=localhost port=5435 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	createTable := `
		CREATE TABLE IF NOT EXISTS payments (
			id SERIAL PRIMARY KEY,
			order_id INT NOT NULL,
			amount DECIMAL(10, 2) NOT NULL,
			status VARCHAR(50) NOT NULL,
			method VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`
	_, err = db.Exec(createTable)
	require.NoError(t, err)

	// Очищаем таблицу перед тестом
	_, err = db.Exec("TRUNCATE TABLE payments RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE payments RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func TestCreatePayment_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &PaymentRepository{db: db}
	ctx := context.Background()

	payment := &model.Payment{
		OrderID: 100,
		Amount:  99.99,
		Status:  "pending",
		Method:  "card",
	}

	err := repo.CreatePayment(ctx, payment)

	assert.NoError(t, err)
	assert.Greater(t, payment.ID, int64(0))
}

func TestGetPayment_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &PaymentRepository{db: db}
	ctx := context.Background()

	// Создаем платеж напрямую в БД
	var paymentID int64
	err := db.QueryRow(`
		INSERT INTO payments (order_id, amount, status, method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, 100, 99.99, "completed", "card", time.Now(), time.Now()).Scan(&paymentID)
	require.NoError(t, err)

	// Получаем платеж через репозиторий
	payment, err := repo.GetPayment(ctx, paymentID)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, paymentID, payment.ID)
	assert.Equal(t, int64(100), payment.OrderID)
}

func TestGetPayment_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &PaymentRepository{db: db}
	ctx := context.Background()

	payment, err := repo.GetPayment(ctx, 99999)

	assert.NoError(t, err)
	assert.Nil(t, payment)
}

func TestUpdatePaymentStatus_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &PaymentRepository{db: db}
	ctx := context.Background()

	// Создаем платеж
	var paymentID int64
	err := db.QueryRow(`
		INSERT INTO payments (order_id, amount, status, method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, 100, 99.99, "pending", "card", time.Now(), time.Now()).Scan(&paymentID)
	require.NoError(t, err)

	// Обновляем статус
	err = repo.UpdatePaymentStatus(ctx, paymentID, "completed")
	assert.NoError(t, err)

	// Проверяем изменение
	var status string
	err = db.QueryRow("SELECT status FROM payments WHERE id = $1", paymentID).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "completed", status)
}

func TestGetPaymentByOrderID_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &PaymentRepository{db: db}
	ctx := context.Background()

	// Создаем платеж
	_, err := db.Exec(`
		INSERT INTO payments (order_id, amount, status, method, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
	`, 100, 99.99, "completed", "card", time.Now(), time.Now())
	require.NoError(t, err)

	// Получаем платеж по order_id
	payment, err := repo.GetPaymentByOrderID(ctx, 100)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, int64(100), payment.OrderID)
}

