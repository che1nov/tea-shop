package service

import (
	"context"

	pb "github.com/che1nov/tea-shop/shared/pb"

	"github.com/che1nov/tea-shop/order-service/internal/kafka"
	"github.com/che1nov/tea-shop/order-service/internal/model"
	"github.com/che1nov/tea-shop/order-service/internal/repository"
)

// KafkaProducerInterface определяет методы для Kafka producer
type KafkaProducerInterface interface {
	PublishOrderCreated(ctx context.Context, event *kafka.OrderEvent) error
	PublishOrderCompleted(ctx context.Context, event *kafka.OrderEvent) error
	Close() error
}

// OrderServiceInterface определяет методы сервиса
type OrderServiceInterface interface {
	CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.Order, error)
	GetOrder(ctx context.Context, id int64) (*model.Order, error)
	UpdateOrderStatus(ctx context.Context, id int64, status string) (*model.Order, error)
	ListUserOrders(ctx context.Context, userID int64) ([]*model.Order, error)
}

type OrderService struct {
	repo               repository.OrderRepositoryInterface
	producer           KafkaProducerInterface
	goodsServiceConn   pb.GoodsServiceClient
	paymentServiceConn pb.PaymentsServiceClient
	deliveryServiceConn pb.DeliveryServiceClient
}

func New(
	repo repository.OrderRepositoryInterface,
	producer KafkaProducerInterface,
	goodsServiceConn pb.GoodsServiceClient,
	paymentServiceConn pb.PaymentsServiceClient,
	deliveryServiceConn pb.DeliveryServiceClient,
) *OrderService {
	return &OrderService{
		repo:                repo,
		producer:            producer,
		goodsServiceConn:    goodsServiceConn,
		paymentServiceConn:  paymentServiceConn,
		deliveryServiceConn: deliveryServiceConn,
	}
}

func (s *OrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.Order, error) {
	// Расчет общей суммы
	var totalPrice float64

	// Проверяем наличие всех товаров
	for _, item := range req.Items {
		good, err := s.goodsServiceConn.GetGood(ctx, &pb.GetGoodRequest{GoodId: item.GoodID})
		if err != nil {
			return nil, err
		}

		if good == nil {
			return nil, nil
		}

		item.Price = good.Price
		totalPrice += good.Price * float64(item.Quantity)

		// Проверяем наличие товара
		checkResp, err := s.goodsServiceConn.CheckStock(ctx, &pb.CheckStockRequest{
			GoodId:   item.GoodID,
			Quantity: item.Quantity,
		})
		if err != nil {
			return nil, err
		}

		if !checkResp.Available {
			return nil, nil // Товара недостаточно
		}
	}

	// Создаём заказ
	order := &model.Order{
		UserID:     req.UserID,
		Items:      req.Items,
		Status:     "pending",
		TotalPrice: totalPrice,
		Address:    req.Address,
	}

	if err := s.repo.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	// Зарезервировали товары
	for _, item := range req.Items {
		_, err := s.goodsServiceConn.ReserveStock(ctx, &pb.ReserveStockRequest{
			GoodId:   item.GoodID,
			Quantity: item.Quantity,
			OrderId:  order.ID,
		})
		if err != nil {
			return nil, err
		}
	}

	// Обрабатываем платёж
	paymentResp, err := s.paymentServiceConn.ProcessPayment(ctx, &pb.ProcessPaymentRequest{
		OrderId: order.ID,
		Amount:  totalPrice,
		Method:  "card",
	})
	if err != nil {
		return nil, err
	}

	if paymentResp.Status == "completed" {
		order.Status = "paid"
		s.repo.UpdateOrderStatus(ctx, order.ID, "paid")
		
		// После успешной оплаты автоматически создаем доставку
		if order.Address != "" {
			_, err := s.deliveryServiceConn.CreateDelivery(ctx, &pb.CreateDeliveryRequest{
				OrderId: order.ID,
				Address: order.Address,
			})
			if err != nil {
				// Логируем ошибку, но не прерываем создание заказа
				// Доставка может быть создана позже вручную
			}
		}
	} else {
		order.Status = "payment_failed"
		s.repo.UpdateOrderStatus(ctx, order.ID, "payment_failed")
	}

	// Публикуем событие в Kafka
	s.producer.PublishOrderCreated(ctx, &kafka.OrderEvent{
		OrderID:    order.ID,
		UserID:     order.UserID,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
	})

	return order, nil
}

func (s *OrderService) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	return s.repo.GetOrder(ctx, id)
}

func (s *OrderService) UpdateOrderStatus(ctx context.Context, id int64, status string) (*model.Order, error) {
	if err := s.repo.UpdateOrderStatus(ctx, id, status); err != nil {
		return nil, err
	}

	return s.repo.GetOrder(ctx, id)
}

func (s *OrderService) ListUserOrders(ctx context.Context, userID int64) ([]*model.Order, error) {
	return s.repo.ListUserOrders(ctx, userID)
}
