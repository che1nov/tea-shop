//go:build integration

package repository

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"testing"

	"github.com/che1nov/tea-shop/users-service/internal/model"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
)

func openIntegrationDB(t *testing.T) *sql.DB {
	t.Helper()

	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("set RUN_INTEGRATION_TESTS=1 to run integration tests")
	}

	dsn := os.Getenv("USERS_TEST_DATABASE_URL")
	if dsn == "" {
		dsn = "user=user password=password dbname=users_db host=localhost port=5432 sslmode=disable"
	}

	db, err := sql.Open("pgx", dsn)
	require.NoError(t, err)
	t.Cleanup(func() { _ = db.Close() })

	require.NoError(t, db.Ping())
	return db
}

func migrateUsersIntegrationDB(t *testing.T, db *sql.DB) {
	t.Helper()

	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			email VARCHAR(255) UNIQUE NOT NULL,
			name VARCHAR(255) NOT NULL,
			password_hash VARCHAR(255) NOT NULL,
			created_at TIMESTAMP NOT NULL,
			updated_at TIMESTAMP NOT NULL
		);
		TRUNCATE TABLE users RESTART IDENTITY CASCADE;
	`)
	require.NoError(t, err)
}

func TestUserRepositoryIntegration_CreateAndReadUser(t *testing.T) {
	db := openIntegrationDB(t)
	migrateUsersIntegrationDB(t, db)

	repo := New(db)
	ctx := context.Background()

	user := &model.User{
		Email:        "integration@example.com",
		Name:         "Integration User",
		PasswordHash: "hash",
	}

	require.NoError(t, repo.CreateUser(ctx, user))
	require.NotZero(t, user.ID)

	byID, err := repo.GetUserByID(ctx, user.ID)
	require.NoError(t, err)
	require.NotNil(t, byID)
	require.Equal(t, user.Email, byID.Email)
	require.Equal(t, model.RoleUser, byID.Role)

	byEmail, err := repo.GetUserByEmail(ctx, user.Email)
	require.NoError(t, err)
	require.NotNil(t, byEmail)
	require.Equal(t, user.ID, byEmail.ID)
}

func TestUserRepositoryIntegration_DuplicateEmail(t *testing.T) {
	db := openIntegrationDB(t)
	migrateUsersIntegrationDB(t, db)

	repo := New(db)
	ctx := context.Background()

	user := &model.User{Email: "duplicate@example.com", Name: "User", PasswordHash: "hash"}
	require.NoError(t, repo.CreateUser(ctx, user))

	duplicate := &model.User{Email: user.Email, Name: "Other User", PasswordHash: "hash"}
	err := repo.CreateUser(ctx, duplicate)
	require.True(t, errors.Is(err, ErrEmailAlreadyExists), "expected ErrEmailAlreadyExists, got %v", err)
}
