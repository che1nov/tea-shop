package service

import (
	"testing"

	"github.com/che1nov/tea-shop/notify-service/internal/kafka"
	"github.com/stretchr/testify/require"
)

func TestHandleEventKnownTypes(t *testing.T) {
	svc := New("noreply@test.local")

	for _, eventType := range []string{"order.created", "order.completed", "order.payment_failed"} {
		err := svc.HandleEvent(&kafka.OrderEvent{EventType: eventType, OrderID: 1})
		require.NoError(t, err)
	}
}

func TestHandleEventUnknownType(t *testing.T) {
	svc := New("noreply@test.local")
	err := svc.HandleEvent(&kafka.OrderEvent{EventType: "unknown"})
	require.Error(t, err)
}
