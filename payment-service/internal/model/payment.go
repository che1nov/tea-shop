package model

import "time"

type Payment struct {
	ID        int64
	OrderID   int64
	Amount    float64
	Status    string
	Method    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ProcessPaymentRequest struct {
	OrderID int64
	Amount  float64
	Method  string
}
