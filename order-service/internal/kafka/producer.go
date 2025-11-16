package kafka

import (
	"context"
	"encoding/json"

	"github.com/segmentio/kafka-go"
)

type OrderEvent struct {
	OrderID    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	EventType  string  `json:"event_type"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
}

type Producer struct {
	writer *kafka.Writer
}

func NewProducer(brokers []string) *Producer {
	return &Producer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...),
			Topic:    "order-events",
			Balancer: &kafka.LeastBytes{},
		},
	}
}

func (p *Producer) PublishOrderCreated(ctx context.Context, event *OrderEvent) error {
	event.EventType = "order.created"

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("order_" + string(rune(event.OrderID))),
		Value: data,
	})
}

func (p *Producer) PublishOrderCompleted(ctx context.Context, event *OrderEvent) error {
	event.EventType = "order.completed"

	data, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte("order_" + string(rune(event.OrderID))),
		Value: data,
	})
}

func (p *Producer) Close() error {
	return p.writer.Close()
}
