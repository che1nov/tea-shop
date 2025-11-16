package model

import "time"

type Delivery struct {
	ID        int64
	OrderID   int64
	Address   string
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateDeliveryRequest struct {
	OrderID int64
	Address string
}
