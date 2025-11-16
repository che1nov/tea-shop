package model

import "time"

type Good struct {
	ID          int64
	SKU         string
	Name        string
	Description string
	Price       float64
	Stock       int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type CreateGoodRequest struct {
	Name        string
	Description string
	Price       float64
	Stock       int32
}

type UpdateGoodRequest struct {
	Name        string
	Description string
	Price       float64
	Stock       int32
	SKU         string
}

type StockReservation struct {
	ID        int64
	GoodID    int64
	OrderID   int64
	Quantity  int32
	CreatedAt time.Time
}
