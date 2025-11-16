package repository

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/che1nov/tea-shop/users-service/internal/model"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// setupTestDB создает тестовую БД и таблицы
func setupTestDB(t *testing.T) *sql.DB {
	connStr := "user=user password=password dbname=users_db host=localhost port=5432 sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	require.NoError(t, err)
	
	// Используем существующую таблицу users
	// Создаем её если не существует
	createTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
	`
	_, err = db.Exec(createTable)
	require.NoError(t, err)
	
	return db
}

// cleanupTestDB очищает тестовые данные
func cleanupTestDB(t *testing.T, db *sql.DB) {
	_, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
}

// setupTestDBWithCleanup создает БД и очищает её перед тестом
func setupTestDBWithCleanup(t *testing.T) *sql.DB {
	db := setupTestDB(t)
	// Очищаем таблицу перед тестом
	_, err := db.Exec("TRUNCATE TABLE users RESTART IDENTITY CASCADE")
	require.NoError(t, err)
	return db
}

func TestCreateUser_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)
	
	repo := &UserRepository{db: db}
	ctx := context.Background()
	
	user := &model.User{
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashed_password",
	}
	
	err := repo.CreateUser(ctx, user)
	
	assert.NoError(t, err)
	assert.Greater(t, user.ID, int64(0))
	
	// Проверяем, что пользователь создан
	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM users WHERE id = $1", user.ID).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestCreateUser_DuplicateEmail(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)
	
	repo := &UserRepository{db: db}
	ctx := context.Background()
	
	user1 := &model.User{
		Email:        "duplicate@example.com",
		Name:         "First User",
		PasswordHash: "hash1",
	}
	
	err := repo.CreateUser(ctx, user1)
	assert.NoError(t, err)
	
	// Пытаемся создать пользователя с тем же email
	user2 := &model.User{
		Email:        "duplicate@example.com",
		Name:         "Second User",
		PasswordHash: "hash2",
	}
	
	err = repo.CreateUser(ctx, user2)
	assert.Error(t, err)
	assert.Equal(t, ErrEmailAlreadyExists, err)
}

func TestGetUserByID_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)
	
	repo := &UserRepository{db: db}
	ctx := context.Background()
	
	// Создаем пользователя напрямую в БД для теста
	var userID int64
	err := db.QueryRow(`
		INSERT INTO users (email, name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`, "test@example.com", "Test User", "hash", time.Now(), time.Now()).Scan(&userID)
	require.NoError(t, err)
	
	// Получаем пользователя через репозиторий
	user, err := repo.GetUserByID(ctx, userID)
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
}

func TestGetUserByID_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)
	
	repo := &UserRepository{db: db}
	ctx := context.Background()
	
	user, err := repo.GetUserByID(ctx, 99999)
	
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestGetUserByEmail_Success(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)
	
	repo := &UserRepository{db: db}
	ctx := context.Background()
	
	// Создаем пользователя
	_, err := db.Exec(`
		INSERT INTO users (email, name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`, "email@example.com", "Email User", "hash", time.Now(), time.Now())
	require.NoError(t, err)
	
	// Получаем пользователя по email
	user, err := repo.GetUserByEmail(ctx, "email@example.com")
	
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "email@example.com", user.Email)
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	db := setupTestDBWithCleanup(t)
	defer db.Close()
	defer cleanupTestDB(t, db)
	
	repo := &UserRepository{db: db}
	ctx := context.Background()
	
	user, err := repo.GetUserByEmail(ctx, "notfound@example.com")
	
	assert.NoError(t, err)
	assert.Nil(t, user)
}

