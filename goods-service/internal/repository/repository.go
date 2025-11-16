package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/lib/pq"

	"github.com/che1nov/tea-shop/goods-service/internal/model"
)

// GoodsRepositoryInterface определяет методы репозитория
type GoodsRepositoryInterface interface {
	CreateGood(ctx context.Context, good *model.Good) error
	GetGood(ctx context.Context, id int64) (*model.Good, error)
	ListGoods(ctx context.Context, limit, offset int32) ([]*model.Good, error)
	UpdateGood(ctx context.Context, good *model.Good) error
	DeleteGood(ctx context.Context, id int64) error
	ReserveStock(ctx context.Context, goodID int64, quantity int32, orderID int64) error
	GetTotalGoods(ctx context.Context) (int32, error)
}

type GoodsRepository struct {
	db *sql.DB
}

func New(db *sql.DB) *GoodsRepository {
	return &GoodsRepository{db: db}
}

// generateSKU генерирует уникальный артикул для товара
func (r *GoodsRepository) generateSKU(ctx context.Context) (string, error) {
	var maxID int64
	err := r.db.QueryRowContext(ctx, "SELECT COALESCE(MAX(id), 0) FROM goods").Scan(&maxID)
	if err != nil {
		return "", err
	}
	// Генерируем SKU в формате GOOD-000001
	return fmt.Sprintf("GOOD-%06d", maxID+1), nil
}

func (r *GoodsRepository) CreateGood(ctx context.Context, good *model.Good) error {
	// Если SKU не указан, генерируем автоматически
	if good.SKU == "" {
		sku, err := r.generateSKU(ctx)
		if err != nil {
			return err
		}
		good.SKU = sku
	}

	query := `
		INSERT INTO goods (sku, name, description, price, stock, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id
	`
	now := time.Now()
	return r.db.QueryRowContext(
		ctx,
		query,
		good.SKU,
		good.Name,
		good.Description,
		good.Price,
		good.Stock,
		now,
		now,
	).Scan(&good.ID)
}

func (r *GoodsRepository) GetGood(ctx context.Context, id int64) (*model.Good, error) {
	query := `SELECT id, sku, name, description, price, stock, created_at, updated_at FROM goods WHERE id = $1`

	good := &model.Good{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&good.ID,
		&good.SKU,
		&good.Name,
		&good.Description,
		&good.Price,
		&good.Stock,
		&good.CreatedAt,
		&good.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return good, nil
}

func (r *GoodsRepository) ListGoods(ctx context.Context, limit, offset int32) ([]*model.Good, error) {
	query := `
		SELECT id, sku, name, description, price, stock, created_at, updated_at 
		FROM goods 
		ORDER BY id 
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var goods []*model.Good
	for rows.Next() {
		good := &model.Good{}
		if err := rows.Scan(
			&good.ID,
			&good.SKU,
			&good.Name,
			&good.Description,
			&good.Price,
			&good.Stock,
			&good.CreatedAt,
			&good.UpdatedAt,
		); err != nil {
			return nil, err
		}
		goods = append(goods, good)
	}

	return goods, rows.Err()
}

func (r *GoodsRepository) ReserveStock(ctx context.Context, goodID int64, quantity int32, orderID int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Проверяем, есть ли товар и достаточно ли его
	var currentStock int32
	err = tx.QueryRowContext(ctx, "SELECT stock FROM goods WHERE id = $1 FOR UPDATE", goodID).Scan(&currentStock)
	if err != nil {
		return err
	}

	if currentStock < quantity {
		return sql.ErrNoRows // Товара недостаточно
	}

	// Уменьшаем остаток
	_, err = tx.ExecContext(
		ctx,
		"UPDATE goods SET stock = stock - $1 WHERE id = $2",
		quantity,
		goodID,
	)
	if err != nil {
		return err
	}

	// Сохраняем информацию о резервировании
	_, err = tx.ExecContext(
		ctx,
		"INSERT INTO stock_reservations (good_id, order_id, quantity, created_at) VALUES ($1, $2, $3, $4)",
		goodID,
		orderID,
		quantity,
		time.Now(),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *GoodsRepository) UpdateGood(ctx context.Context, good *model.Good) error {
	query := `
		UPDATE goods 
		SET sku = $1, name = $2, description = $3, price = $4, stock = $5, updated_at = $6
		WHERE id = $7
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		good.SKU,
		good.Name,
		good.Description,
		good.Price,
		good.Stock,
		time.Now(),
		good.ID,
	)
	return err
}

func (r *GoodsRepository) DeleteGood(ctx context.Context, id int64) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Сначала удаляем все резервации для этого товара
	_, err = tx.ExecContext(ctx, "DELETE FROM stock_reservations WHERE good_id = $1", id)
	if err != nil {
		return err
	}

	// Затем удаляем сам товар
	_, err = tx.ExecContext(ctx, "DELETE FROM goods WHERE id = $1", id)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *GoodsRepository) GetTotalGoods(ctx context.Context) (int32, error) {
	var total int32
	err := r.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM goods").Scan(&total)
	return total, err
}
