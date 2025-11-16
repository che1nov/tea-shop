package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/che1nov/tea-shop/goods-service/internal/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDBWithCleanup создает БД и очищает её перед тестом
func setupTestDBWithCleanup(t *testing.T) *sql.DB {
	connStr := "user=user password=password dbname=goods_db host=localhost port=5433 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)

	createTable := `
		CREATE TABLE IF NOT EXISTS goods (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			description TEXT,
			price DECIMAL(10, 2) NOT NULL,
			stock INT NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		CREATE TABLE IF NOT EXISTS stock_reservations (
			id SERIAL PRIMARY KEY,
			good_id INT NOT NULL,
			order_id INT NOT NULL,
			quantity INT NOT NULL,
			created_at TIMESTAMP NOT NULL
		);
	`
	_, err = db.Exec(createTable)
	require.NoError(t, err)

	// Очищаем таблицы перед тестом
	_, err = db.Exec("TRUNCATE TABLE goods, stock_reservations RESTART IDENTITY CASCADE")
	require.NoError(t, err)

	return db
}

func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE goods, stock_reservations RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

func TestCreateGood_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &GoodsRepository{db: db}
	ctx := context.Background()

	good := &model.Good{
		Name:        "Test Good",
		Description:  "Test Description",
		Price:       99.99,
		Stock:       100,
	}

	err := repo.CreateGood(ctx, good)

	assert.NoError(t, err)
	assert.Greater(t, good.ID, int64(0))
}

func TestGetGood_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &GoodsRepository{db: db}
	ctx := context.Background()

	// Создаем товар напрямую в БД
	var goodID int64
	err := db.QueryRow(`
		INSERT INTO goods (name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`, "Test Good", "Description", 99.99, 100, time.Now(), time.Now()).Scan(&goodID)
	require.NoError(t, err)

	// Получаем товар через репозиторий
	good, err := repo.GetGood(ctx, goodID)

	assert.NoError(t, err)
	assert.NotNil(t, good)
	assert.Equal(t, goodID, good.ID)
	assert.Equal(t, "Test Good", good.Name)
}

func TestGetGood_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &GoodsRepository{db: db}
	ctx := context.Background()

	good, err := repo.GetGood(ctx, 99999)

	assert.NoError(t, err)
	assert.Nil(t, good)
}

func TestListGoods_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &GoodsRepository{db: db}
	ctx := context.Background()

	// Создаем несколько товаров
	_, err := db.Exec(`
		INSERT INTO goods (name, description, price, stock, created_at, updated_at)
		VALUES 
			('Good 1', 'Desc 1', 10.0, 50, NOW(), NOW()),
			('Good 2', 'Desc 2', 20.0, 30, NOW(), NOW())
	`)
	require.NoError(t, err)

	goods, err := repo.ListGoods(ctx, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, goods)
	assert.GreaterOrEqual(t, len(goods), 2)
}

func TestGetTotalGoods_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)

	repo := &GoodsRepository{db: db}
	ctx := context.Background()

	// Создаем товары
	_, err := db.Exec(`
		INSERT INTO goods (name, description, price, stock, created_at, updated_at)
		VALUES 
			('Good 1', 'Desc 1', 10.0, 50, NOW(), NOW()),
			('Good 2', 'Desc 2', 20.0, 30, NOW(), NOW())
	`)
	require.NoError(t, err)

	total, err := repo.GetTotalGoods(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int32(2), total)
}

