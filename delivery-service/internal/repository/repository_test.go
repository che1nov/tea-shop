package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/che1nov/tea-shop/delivery-service/internal/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDBWithCleanup создает БД и очищает её перед тестом
func setupTestDBWithCleanup(t *testing.T) *sql.DB {
	connStr := "user=user password=password dbname=deliveries_db host=localhost port=5436 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	createTable := `
		CREATE TABLE IF NOT EXISTS deliveries (
			id SERIAL PRIMARY KEY,
			order_id INT NOT NULL UNIQUE,
			address VARCHAR(500) NOT NULL,
			status VARCHAR(50) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE INDEX IF NOT EXISTS idx_deliveries_order ON deliveries(order_id);
	`
	_, err = db.Exec(createTable)
	require.NoError(t, err)

	// Очищаем таблицу перед тестом
	_, err = db.Exec("TRUNCATE TABLE deliveries RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE deliveries RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func TestCreateDelivery_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	delivery := &model.Delivery{
		OrderID: 100,
		Address: "123 Main St",
		Status:  "pending",
	}

	err := repo.CreateDelivery(ctx, delivery)

	assert.NoError(t, err)
	assert.Greater(t, delivery.ID, int64(0))
	assert.NotZero(t, delivery.CreatedAt)
	assert.NotZero(t, delivery.UpdatedAt)
}

func TestCreateDelivery_DuplicateOrderID(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	delivery1 := &model.Delivery{
		OrderID: 100,
		Address: "123 Main St",
		Status:  "pending",
	}

	err := repo.CreateDelivery(ctx, delivery1)
	assert.NoError(t, err)

	delivery2 := &model.Delivery{
		OrderID: 100,
		Address: "456 Oak Ave",
		Status:  "pending",
	}

	err = repo.CreateDelivery(ctx, delivery2)
	assert.Error(t, err)
}

func TestGetDelivery_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	// Создаем доставку напрямую в БД
	var deliveryID int64
	err := db.QueryRow(`
		INSERT INTO deliveries (order_id, address, status, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, 200, "789 Pine Rd", "pending", time.Now(), time.Now()).Scan(&deliveryID)
	require.NoError(t, err)

	delivery, err := repo.GetDelivery(ctx, deliveryID)

	assert.NoError(t, err)
	assert.NotNil(t, delivery)
	assert.Equal(t, deliveryID, delivery.ID)
	assert.Equal(t, int64(200), delivery.OrderID)
	assert.Equal(t, "789 Pine Rd", delivery.Address)
	assert.Equal(t, "pending", delivery.Status)
}

func TestGetDelivery_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	delivery, err := repo.GetDelivery(ctx, 9999)

	assert.NoError(t, err)
	assert.Nil(t, delivery)
}

func TestGetDeliveryByOrderID_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	// Создаем доставку
	delivery := &model.Delivery{
		OrderID: 300,
		Address: "321 Elm St",
		Status:  "in_transit",
	}
	err := repo.CreateDelivery(ctx, delivery)
	require.NoError(t, err)

	found, err := repo.GetDeliveryByOrderID(ctx, 300)

	assert.NoError(t, err)
	assert.NotNil(t, found)
	assert.Equal(t, int64(300), found.OrderID)
	assert.Equal(t, "321 Elm St", found.Address)
	assert.Equal(t, "in_transit", found.Status)
}

func TestGetDeliveryByOrderID_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	delivery, err := repo.GetDeliveryByOrderID(ctx, 9999)

	assert.NoError(t, err)
	assert.Nil(t, delivery)
}

func TestUpdateDeliveryStatus_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	// Создаем доставку
	delivery := &model.Delivery{
		OrderID: 400,
		Address: "555 Oak Ave",
		Status:  "pending",
	}
	err := repo.CreateDelivery(ctx, delivery)
	require.NoError(t, err)

	err = repo.UpdateDeliveryStatus(ctx, delivery.ID, "delivered")
	assert.NoError(t, err)

	// Проверяем обновление
	updated, err := repo.GetDelivery(ctx, delivery.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updated)
	assert.Equal(t, "delivered", updated.Status)
}

func TestUpdateDeliveryStatus_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &DeliveryRepository{db: db}
	ctx := context.Background()

	err := repo.UpdateDeliveryStatus(ctx, 9999, "delivered")
	assert.NoError(t, err) // UPDATE не возвращает ошибку для несуществующих записей
}

