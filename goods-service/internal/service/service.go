package service

import (
	"context"
	"database/sql"

	"github.com/che1nov/tea-shop/goods-service/internal/model"
	"github.com/che1nov/tea-shop/goods-service/internal/repository"
)

// GoodsServiceInterface определяет методы сервиса
type GoodsServiceInterface interface {
	CreateGood(ctx context.Context, req *model.CreateGoodRequest) (*model.Good, error)
	GetGood(ctx context.Context, id int64) (*model.Good, error)
	ListGoods(ctx context.Context, limit, offset int32) ([]*model.Good, error)
	UpdateGood(ctx context.Context, id int64, req *model.UpdateGoodRequest) (*model.Good, error)
	DeleteGood(ctx context.Context, id int64) error
	GetTotalGoods(ctx context.Context) (int32, error)
	CheckStock(ctx context.Context, goodID int64, quantity int32) (bool, error)
	ReserveStock(ctx context.Context, goodID int64, quantity int32, orderID int64) (bool, error)
}

type GoodsService struct {
	repo repository.GoodsRepositoryInterface
}

func New(repo repository.GoodsRepositoryInterface) *GoodsService {
	return &GoodsService{
		repo: repo,
	}
}

func (s *GoodsService) CreateGood(ctx context.Context, req *model.CreateGoodRequest) (*model.Good, error) {
	good := &model.Good{
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.repo.CreateGood(ctx, good); err != nil {
		return nil, err
	}

	return good, nil
}

func (s *GoodsService) GetGood(ctx context.Context, id int64) (*model.Good, error) {
	return s.repo.GetGood(ctx, id)
}

func (s *GoodsService) ListGoods(ctx context.Context, limit, offset int32) ([]*model.Good, error) {
	return s.repo.ListGoods(ctx, limit, offset)
}

func (s *GoodsService) UpdateGood(ctx context.Context, id int64, req *model.UpdateGoodRequest) (*model.Good, error) {
	good, err := s.repo.GetGood(ctx, id)
	if err != nil {
		return nil, err
	}
	if good == nil {
		return nil, nil
	}

	// Обновляем поля
	if req.Name != "" {
		good.Name = req.Name
	}
	if req.Description != "" {
		good.Description = req.Description
	}
	if req.Price > 0 {
		good.Price = req.Price
	}
	if req.Stock >= 0 {
		good.Stock = req.Stock
	}
	if req.SKU != "" {
		good.SKU = req.SKU
	}

	if err := s.repo.UpdateGood(ctx, good); err != nil {
		return nil, err
	}

	return good, nil
}

func (s *GoodsService) DeleteGood(ctx context.Context, id int64) error {
	return s.repo.DeleteGood(ctx, id)
}

func (s *GoodsService) GetTotalGoods(ctx context.Context) (int32, error) {
	return s.repo.GetTotalGoods(ctx)
}

func (s *GoodsService) CheckStock(ctx context.Context, goodID int64, quantity int32) (bool, error) {
	good, err := s.repo.GetGood(ctx, goodID)
	if err != nil {
		return false, err
	}

	if good == nil {
		return false, nil
	}

	return good.Stock >= quantity, nil
}

func (s *GoodsService) ReserveStock(ctx context.Context, goodID int64, quantity int32, orderID int64) (bool, error) {
	err := s.repo.ReserveStock(ctx, goodID, quantity, orderID)
	if err == sql.ErrNoRows {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
