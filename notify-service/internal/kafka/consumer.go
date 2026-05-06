package kafka

import (
	"context"
	"encoding/json"

	"github.com/che1nov/tea-shop/shared/pkg/logger"
	appmetrics "github.com/che1nov/tea-shop/shared/pkg/metrics"
	"github.com/segmentio/kafka-go"
)

type OrderEvent struct {
	OrderID    int64   `json:"order_id"`
	UserID     int64   `json:"user_id"`
	EventType  string  `json:"event_type"`
	Status     string  `json:"status"`
	TotalPrice float64 `json:"total_price"`
}

type Consumer struct {
	reader *kafka.Reader
}

func NewConsumer(brokers []string, groupID string) *Consumer {
	return &Consumer{
		reader: kafka.NewReader(kafka.ReaderConfig{
			Brokers: brokers,
			Topic:   "order-events",
			GroupID: groupID,
		}),
	}
}

func (c *Consumer) Start(ctx context.Context, handleEvent func(*OrderEvent) error) error {
	for {
		msg, err := c.reader.ReadMessage(ctx)
		if err != nil {
			appmetrics.ObserveKafkaMessage("notify-service", "order-events", "unknown", "consumed", "read_error")
			return err
		}

		event := &OrderEvent{}
		if err := json.Unmarshal(msg.Value, event); err != nil {
			appmetrics.ObserveKafkaMessage("notify-service", "order-events", "unknown", "consumed", "decode_error")
			logger.Error("Failed to unmarshal event", "error", err)
			continue
		}

		appmetrics.ObserveKafkaMessage("notify-service", "order-events", event.EventType, "consumed", "received")
		logger.Info("Received event", "event_type", event.EventType, "order_id", event.OrderID)

		if err := handleEvent(event); err != nil {
			appmetrics.ObserveKafkaMessage("notify-service", "order-events", event.EventType, "consumed", "handler_error")
			logger.Error("Failed to handle event", "error", err, "event_type", event.EventType, "order_id", event.OrderID)
			continue
		}

		appmetrics.ObserveKafkaMessage("notify-service", "order-events", event.EventType, "consumed", "success")
	}
}

func (c *Consumer) Close() error {
	return c.reader.Close()
}
