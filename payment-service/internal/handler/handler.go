package handler

import (
	"context"

	pb "github.com/che1nov/tea-shop/shared/pb"

	"github.com/che1nov/tea-shop/payment-service/internal/model"
	"github.com/che1nov/tea-shop/payment-service/internal/service"
)

type PaymentsHandler struct {
	service service.PaymentServiceInterface
	pb.UnimplementedPaymentsServiceServer
}

func New(svc service.PaymentServiceInterface) *PaymentsHandler {
	return &PaymentsHandler{
		service: svc,
	}
}

func (h *PaymentsHandler) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest) (*pb.Payment, error) {
	payment, err := h.service.ProcessPayment(ctx, &model.ProcessPaymentRequest{
		OrderID: req.OrderId,
		Amount:  req.Amount,
		Method:  req.Method,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Payment{
		Id:        payment.ID,
		OrderId:   payment.OrderID,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Unix(),
		UpdatedAt: payment.UpdatedAt.Unix(),
	}, nil
}

func (h *PaymentsHandler) GetPayment(ctx context.Context, req *pb.GetPaymentRequest) (*pb.Payment, error) {
	payment, err := h.service.GetPayment(ctx, req.PaymentId)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, nil
	}

	return &pb.Payment{
		Id:        payment.ID,
		OrderId:   payment.OrderID,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Unix(),
		UpdatedAt: payment.UpdatedAt.Unix(),
	}, nil
}

func (h *PaymentsHandler) GetPaymentByOrderID(ctx context.Context, req *pb.GetPaymentByOrderIDRequest) (*pb.Payment, error) {
	payment, err := h.service.GetPaymentByOrderID(ctx, req.OrderId)
	if err != nil {
		return nil, err
	}

	if payment == nil {
		return nil, nil
	}

	return &pb.Payment{
		Id:        payment.ID,
		OrderId:   payment.OrderID,
		Amount:    payment.Amount,
		Status:    payment.Status,
		CreatedAt: payment.CreatedAt.Unix(),
		UpdatedAt: payment.UpdatedAt.Unix(),
	}, nil
}
