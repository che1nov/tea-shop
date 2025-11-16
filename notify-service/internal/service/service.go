package service

import (
	"fmt"

	"github.com/che1nov/tea-shop/notify-service/internal/kafka"
	"github.com/che1nov/tea-shop/shared/pkg/logger"
)

type NotifyService struct {
	emailFrom string
}

func New(emailFrom string) *NotifyService {
	return &NotifyService{
		emailFrom: emailFrom,
	}
}

func (s *NotifyService) HandleOrderCreated(event *kafka.OrderEvent) error {
	// Имитируем отправку email
	logger.Info("Sending email notification for order created", "order_id", event.OrderID, "total_price", event.TotalPrice)

	// В реальной системе здесь будет:
	// - подключение к SMTP серверу
	// - формирование HTML письма
	// - отправка email

	return nil
}

func (s *NotifyService) HandleOrderCompleted(event *kafka.OrderEvent) error {
	// Имитируем отправку email
	logger.Info("Sending email notification for order completed", "order_id", event.OrderID, "status", event.Status)

	// В реальной системе здесь будет отправка email с информацией о завершении заказа

	return nil
}

func (s *NotifyService) HandleOrderPaymentFailed(event *kafka.OrderEvent) error {
	logger.Info("Sending email notification for order payment failed", "order_id", event.OrderID)

	return nil
}

func (s *NotifyService) HandleEvent(event *kafka.OrderEvent) error {
	switch event.EventType {
	case "order.created":
		return s.HandleOrderCreated(event)
	case "order.completed":
		return s.HandleOrderCompleted(event)
	case "order.payment_failed":
		return s.HandleOrderPaymentFailed(event)
	default:
		return fmt.Errorf("unknown event type: %s", event.EventType)
	}
}
