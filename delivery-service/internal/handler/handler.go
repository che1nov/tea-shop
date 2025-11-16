package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/che1nov/tea-shop/shared/pb"

	"github.com/che1nov/tea-shop/delivery-service/internal/model"
	"github.com/che1nov/tea-shop/delivery-service/internal/service"
)

type DeliveryHandler struct {
	service service.DeliveryServiceInterface
	pb.UnimplementedDeliveryServiceServer
}

func New(svc service.DeliveryServiceInterface) *DeliveryHandler {
	return &DeliveryHandler{
		service: svc,
	}
}

func (h *DeliveryHandler) CreateDelivery(ctx context.Context, req *pb.CreateDeliveryRequest) (*pb.Delivery, error) {
	if req.OrderId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "order_id is required")
	}
	if req.Address == "" {
		return nil, status.Errorf(codes.InvalidArgument, "address is required")
	}

	delivery, err := h.service.CreateDelivery(ctx, &model.CreateDeliveryRequest{
		OrderID: req.OrderId,
		Address: req.Address,
	})
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to create delivery: %v", err)
	}

	return &pb.Delivery{
		Id:        delivery.ID,
		OrderId:   delivery.OrderID,
		Address:   delivery.Address,
		Status:    delivery.Status,
		CreatedAt: delivery.CreatedAt.Unix(),
		UpdatedAt: delivery.UpdatedAt.Unix(),
	}, nil
}

func (h *DeliveryHandler) GetDelivery(ctx context.Context, req *pb.GetDeliveryRequest) (*pb.Delivery, error) {
	if req.DeliveryId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "delivery_id is required")
	}

	delivery, err := h.service.GetDelivery(ctx, req.DeliveryId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get delivery: %v", err)
	}

	if delivery == nil {
		return nil, status.Errorf(codes.NotFound, "delivery with id %d not found", req.DeliveryId)
	}

	return &pb.Delivery{
		Id:        delivery.ID,
		OrderId:   delivery.OrderID,
		Address:   delivery.Address,
		Status:    delivery.Status,
		CreatedAt: delivery.CreatedAt.Unix(),
		UpdatedAt: delivery.UpdatedAt.Unix(),
	}, nil
}

func (h *DeliveryHandler) UpdateDeliveryStatus(ctx context.Context, req *pb.UpdateDeliveryStatusRequest) (*pb.Delivery, error) {
	if req.DeliveryId == 0 {
		return nil, status.Errorf(codes.InvalidArgument, "delivery_id is required")
	}
	if req.Status == "" {
		return nil, status.Errorf(codes.InvalidArgument, "status is required")
	}

	delivery, err := h.service.UpdateDeliveryStatus(ctx, req.DeliveryId, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to update delivery status: %v", err)
	}

	if delivery == nil {
		return nil, status.Errorf(codes.NotFound, "delivery with id %d not found", req.DeliveryId)
	}

	return &pb.Delivery{
		Id:        delivery.ID,
		OrderId:   delivery.OrderID,
		Address:   delivery.Address,
		Status:    delivery.Status,
		CreatedAt: delivery.CreatedAt.Unix(),
		UpdatedAt: delivery.UpdatedAt.Unix(),
	}, nil
}

func (h *DeliveryHandler) ListDeliveries(ctx context.Context, req *pb.ListDeliveriesRequest) (*pb.ListDeliveriesResponse, error) {
	limit := req.Limit
	if limit <= 0 {
		limit = 100
	}

	deliveries, total, err := h.service.ListDeliveries(ctx, limit, req.Offset, req.Status)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to list deliveries: %v", err)
	}

	pbDeliveries := make([]*pb.Delivery, len(deliveries))
	for i, delivery := range deliveries {
		pbDeliveries[i] = &pb.Delivery{
			Id:        delivery.ID,
			OrderId:   delivery.OrderID,
			Address:   delivery.Address,
			Status:    delivery.Status,
			CreatedAt: delivery.CreatedAt.Unix(),
			UpdatedAt: delivery.UpdatedAt.Unix(),
		}
	}

	return &pb.ListDeliveriesResponse{
		Deliveries: pbDeliveries,
		Total:      total,
	}, nil
}
