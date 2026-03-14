package repository

import (
	"context"
	"database/sql"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/che1nov/tea-shop/users-service/internal/model"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupRepo(t *testing.T) (*UserRepository, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return New(db), mock, func() { _ = db.Close() }
}

func TestCreateUserSuccess(t *testing.T) {
	repo, mock, cleanup := setupRepo(t)
	defer cleanup()

	user := &model.User{Email: "test@example.com", Name: "Test", PasswordHash: "hash"}
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (email, name, password_hash, created_at, updated_at)")).
		WithArgs(user.Email, user.Name, user.PasswordHash, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(int64(1)))

	err := repo.CreateUser(context.Background(), user)
	require.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUserDuplicateEmail(t *testing.T) {
	repo, mock, cleanup := setupRepo(t)
	defer cleanup()

	user := &model.User{Email: "dup@example.com", Name: "Dup", PasswordHash: "hash"}
	mock.ExpectQuery(regexp.QuoteMeta("INSERT INTO users (email, name, password_hash, created_at, updated_at)")).
		WithArgs(user.Email, user.Name, user.PasswordHash, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnError(&pq.Error{Code: "23505"})

	err := repo.CreateUser(context.Background(), user)
	require.Error(t, err)
	assert.ErrorIs(t, err, ErrEmailAlreadyExists)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIDSuccess(t *testing.T) {
	repo, mock, cleanup := setupRepo(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE id = $1")).
		WithArgs(int64(7)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password_hash", "created_at", "updated_at"}).
			AddRow(int64(7), "u@example.com", "User", "hash", now, now))

	user, err := repo.GetUserByID(context.Background(), 7)
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, model.RoleUser, user.Role)
	assert.Equal(t, int64(7), user.ID)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByIDNotFound(t *testing.T) {
	repo, mock, cleanup := setupRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE id = $1")).
		WithArgs(int64(100)).
		WillReturnError(sql.ErrNoRows)

	user, err := repo.GetUserByID(context.Background(), 100)
	require.NoError(t, err)
	assert.Nil(t, user)
	require.NoError(t, mock.ExpectationsWereMet())
}

func TestGetUserByEmailSuccess(t *testing.T) {
	repo, mock, cleanup := setupRepo(t)
	defer cleanup()

	now := time.Now()
	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE email = $1")).
		WithArgs("mail@example.com").
		WillReturnRows(sqlmock.NewRows([]string{"id", "email", "name", "password_hash", "created_at", "updated_at"}).
			AddRow(int64(3), "mail@example.com", "Name", "hash", now, now))

	user, err := repo.GetUserByEmail(context.Background(), "mail@example.com")
	require.NoError(t, err)
	require.NotNil(t, user)
	assert.Equal(t, model.RoleUser, user.Role)
	assert.Equal(t, "mail@example.com", user.Email)
	require.NoError(t, mock.ExpectationsWereMet())
}

