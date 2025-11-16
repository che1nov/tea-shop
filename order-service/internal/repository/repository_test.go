package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"testing"
	"time"

	"github.com/che1nov/tea-shop/order-service/internal/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDBWithCleanup создает БД и очищает её перед тестом
func setupTestDBWithCleanup(t *testing.T) *sql.DB {
	connStr := "user=user password=password dbname=orders_db host=localhost port=5434 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	createTable := `
		CREATE TABLE IF NOT EXISTS orders (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL,
			items JSONB NOT NULL,
			status VARCHAR(50) NOT NULL,
			total_price DECIMAL(10, 2) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`
	_, err = db.Exec(createTable)
	require.NoError(t, err)

	// Очищаем таблицу перед тестом
	_, err = db.Exec("TRUNCATE TABLE orders RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE orders RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func TestCreateOrder_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &OrderRepository{db: db}
	ctx := context.Background()

	order := &model.Order{
		UserID: 100,
		Items: []model.OrderItem{
			{GoodID: 1, Quantity: 2, Price: 49.99},
		},
		Status:     "pending",
		TotalPrice: 99.98,
	}

	err := repo.CreateOrder(ctx, order)

	assert.NoError(t, err)
	assert.Greater(t, order.ID, int64(0))
}

func TestGetOrder_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &OrderRepository{db: db}
	ctx := context.Background()

	// Создаем заказ напрямую в БД
	items := []model.OrderItem{
		{GoodID: 1, Quantity: 2, Price: 49.99},
	}
	itemsJSON, _ := json.Marshal(items)

	var orderID int64
	err := db.QueryRow(`
		INSERT INTO orders (user_id, items, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, 100, itemsJSON, "pending", 99.98, time.Now(), time.Now()).Scan(&orderID)
	require.NoError(t, err)

	// Получаем заказ через репозиторий
	order, err := repo.GetOrder(ctx, orderID)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, orderID, order.ID)
	assert.Equal(t, int64(100), order.UserID)
	assert.Equal(t, 1, len(order.Items))
}

func TestGetOrder_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &OrderRepository{db: db}
	ctx := context.Background()

	order, err := repo.GetOrder(ctx, 99999)

	assert.NoError(t, err)
	assert.Nil(t, order)
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &OrderRepository{db: db}
	ctx := context.Background()

	// Создаем заказ
	itemsJSON, _ := json.Marshal([]model.OrderItem{})
	var orderID int64
	err := db.QueryRow(`
		INSERT INTO orders (user_id, items, status, total_price, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, 100, itemsJSON, "pending", 99.99, time.Now(), time.Now()).Scan(&orderID)
	require.NoError(t, err)

	// Обновляем статус
	err = repo.UpdateOrderStatus(ctx, orderID, "completed")
	assert.NoError(t, err)

	// Проверяем изменение
	var status string
	err = db.QueryRow("SELECT status FROM orders WHERE id = $1", orderID).Scan(&status)
	assert.NoError(t, err)
	assert.Equal(t, "completed", status)
}

func TestListUserOrders_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &OrderRepository{db: db}
	ctx := context.Background()

	// Создаем заказы для пользователя
	itemsJSON, _ := json.Marshal([]model.OrderItem{})
	_, err := db.Exec(`
		INSERT INTO orders (user_id, items, status, total_price, created_at, updated_at)
		VALUES 
			($1, $2, 'pending', 50.0, NOW(), NOW()),
			($1, $2, 'completed', 100.0, NOW(), NOW())
	`, 100, itemsJSON)
	require.NoError(t, err)

	orders, err := repo.ListUserOrders(ctx, 100)

	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.GreaterOrEqual(t, len(orders), 2)
}

