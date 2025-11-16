package model

import "time"

type OrderItem struct {
	GoodID   int64
	Quantity int32
	Price    float64
}

type Order struct {
	ID         int64
	UserID     int64
	Items      []OrderItem
	Status     string
	TotalPrice float64
	Address    string
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateOrderRequest struct {
	UserID  int64
	Items   []OrderItem
	Address string
}
