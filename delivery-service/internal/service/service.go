package service

import (
	"context"

	"github.com/che1nov/tea-shop/delivery-service/internal/model"
	"github.com/che1nov/tea-shop/delivery-service/internal/repository"
)

// DeliveryServiceInterface определяет методы сервиса
type DeliveryServiceInterface interface {
	CreateDelivery(ctx context.Context, req *model.CreateDeliveryRequest) (*model.Delivery, error)
	GetDelivery(ctx context.Context, id int64) (*model.Delivery, error)
	GetDeliveryByOrderID(ctx context.Context, orderID int64) (*model.Delivery, error)
	UpdateDeliveryStatus(ctx context.Context, id int64, status string) (*model.Delivery, error)
	ListDeliveries(ctx context.Context, limit, offset int32, status string) ([]*model.Delivery, int32, error)
}

type DeliveryService struct {
	repo repository.DeliveryRepositoryInterface
}

func New(repo repository.DeliveryRepositoryInterface) *DeliveryService {
	return &DeliveryService{
		repo: repo,
	}
}

func (s *DeliveryService) CreateDelivery(ctx context.Context, req *model.CreateDeliveryRequest) (*model.Delivery, error) {
	delivery := &model.Delivery{
		OrderID: req.OrderID,
		Address: req.Address,
		Status:  "pending",
	}

	if err := s.repo.CreateDelivery(ctx, delivery); err != nil {
		return nil, err
	}

	return delivery, nil
}

func (s *DeliveryService) GetDelivery(ctx context.Context, id int64) (*model.Delivery, error) {
	return s.repo.GetDelivery(ctx, id)
}

func (s *DeliveryService) GetDeliveryByOrderID(ctx context.Context, orderID int64) (*model.Delivery, error) {
	return s.repo.GetDeliveryByOrderID(ctx, orderID)
}

func (s *DeliveryService) UpdateDeliveryStatus(ctx context.Context, id int64, status string) (*model.Delivery, error) {
	if err := s.repo.UpdateDeliveryStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return s.repo.GetDelivery(ctx, id)
}

func (s *DeliveryService) ListDeliveries(ctx context.Context, limit, offset int32, status string) ([]*model.Delivery, int32, error) {
	deliveries, err := s.repo.ListDeliveries(ctx, limit, offset, status)
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.GetTotalDeliveries(ctx, status)
	if err != nil {
		return nil, 0, err
	}

	return deliveries, total, nil
}
