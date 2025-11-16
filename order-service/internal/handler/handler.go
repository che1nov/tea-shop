package handler

import (
	"context"

	"github.com/che1nov/tea-shop/order-service/internal/model"
	"github.com/che1nov/tea-shop/order-service/internal/service"
	pb "github.com/che1nov/tea-shop/shared/pb"
)

type OrdersHandler struct {
	service service.OrderServiceInterface
	pb.UnimplementedOrdersServiceServer
}

func New(svc service.OrderServiceInterface) *OrdersHandler {
	return &OrdersHandler{
		service: svc,
	}
}

func (h *OrdersHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderRequest) (*pb.Order, error) {
	items := make([]model.OrderItem, len(req.Items))
	for i, item := range req.Items {
		items[i] = model.OrderItem{
			GoodID:   item.GoodId,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	order, err := h.service.CreateOrder(ctx, &model.CreateOrderRequest{
		UserID:  req.UserId,
		Items:   items,
		Address: req.Address,
	})
	if err != nil {
		return nil, err
	}

	return h.orderToProto(order), nil
}

func (h *OrdersHandler) GetOrder(ctx context.Context, req *pb.GetOrderRequest) (*pb.Order, error) {
	order, err := h.service.GetOrder(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	if order == nil {
		return nil, nil
	}

	return h.orderToProto(order), nil
}

func (h *OrdersHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusRequest) (*pb.Order, error) {
	order, err := h.service.UpdateOrderStatus(ctx, req.OrderId, req.Status)
	if err != nil {
		return nil, err
	}

	return h.orderToProto(order), nil
}

func (h *OrdersHandler) orderToProto(order *model.Order) *pb.Order {
	items := make([]*pb.OrderItem, len(order.Items))
	for i, item := range order.Items {
		items[i] = &pb.OrderItem{
			GoodId:   item.GoodID,
			Quantity: item.Quantity,
			Price:    item.Price,
		}
	}

	return &pb.Order{
		Id:         order.ID,
		UserId:     order.UserID,
		Items:      items,
		Status:     order.Status,
		TotalPrice: order.TotalPrice,
		Address:    order.Address,
		CreatedAt:  order.CreatedAt.Unix(),
		UpdatedAt:  order.UpdatedAt.Unix(),
	}
}
