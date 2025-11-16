package service

import (
	"context"
	"math/rand"

	"github.com/che1nov/tea-shop/payment-service/internal/model"
	"github.com/che1nov/tea-shop/payment-service/internal/repository"
)

// PaymentServiceInterface определяет методы сервиса
type PaymentServiceInterface interface {
	ProcessPayment(ctx context.Context, req *model.ProcessPaymentRequest) (*model.Payment, error)
	GetPayment(ctx context.Context, id int64) (*model.Payment, error)
	GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error)
}

type PaymentService struct {
	repo repository.PaymentRepositoryInterface
}

func New(repo repository.PaymentRepositoryInterface) *PaymentService {
	return &PaymentService{
		repo: repo,
	}
}

func (s *PaymentService) ProcessPayment(ctx context.Context, req *model.ProcessPaymentRequest) (*model.Payment, error) {
	// Имитируем обработку платежа (в реальной системе интегрируем с платёжным шлюзом)
	payment := &model.Payment{
		OrderID: req.OrderID,
		Amount:  req.Amount,
		Method:  req.Method,
		Status:  "pending",
	}

	if err := s.repo.CreatePayment(ctx, payment); err != nil {
		return nil, err
	}

	// Имитируем обработку платежа
	// 90% успешных, 10% неудачных
	// Для тестируемости можно передавать seed или использовать интерфейс
	if rand.Intn(100) < 90 {
		payment.Status = "completed"
	} else {
		payment.Status = "failed"
	}

	if err := s.repo.UpdatePaymentStatus(ctx, payment.ID, payment.Status); err != nil {
		return nil, err
	}

	return payment, nil
}

func (s *PaymentService) GetPayment(ctx context.Context, id int64) (*model.Payment, error) {
	return s.repo.GetPayment(ctx, id)
}

func (s *PaymentService) GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error) {
	return s.repo.GetPaymentByOrderID(ctx, orderID)
}
