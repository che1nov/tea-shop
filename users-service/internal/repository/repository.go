package repository

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/lib/pq"
	"github.com/che1nov/tea-shop/users-service/internal/model"
)

var ErrEmailAlreadyExists = errors.New("email already exists")

// UserRepositoryInterface определяет методы репозитория
type UserRepositoryInterface interface {
	CreateUser(ctx context.Context, user *model.User) error
	GetUserByID(ctx context.Context, id int64) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
}

type UserRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) CreateUser(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (email, name, password_hash, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
	`
	now := time.Now()
	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Email,
		user.Name,
		user.PasswordHash,
		now,
		now,
	).Scan(&user.ID)
	
	if err != nil {
		// Проверяем ошибку уникального ограничения
		if pqErr, ok := err.(*pq.Error); ok {
			if pqErr.Code == "23505" { // unique_violation
				return ErrEmailAlreadyExists
			}
		}
		return err
	}
	
	return nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	query := `SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE id = $1`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Все пользователи из БД имеют роль "user"
	user.Role = model.RoleUser
	return user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `SELECT id, email, name, password_hash, created_at, updated_at FROM users WHERE email = $1`

	user := &model.User{}
	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.PasswordHash,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	// Все пользователи из БД имеют роль "user"
	user.Role = model.RoleUser
	return user, nil
}
