package handler

import (
	"context"
	"testing"

	"github.com/che1nov/tea-shop/payment-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/che1nov/tea-shop/shared/pb"
)

// MockPaymentService - мок для сервиса
type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) ProcessPayment(ctx context.Context, req *model.ProcessPaymentRequest) (*model.Payment, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentService) GetPayment(ctx context.Context, id int64) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockPaymentService) GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func TestNew(t *testing.T) {
	mockService := new(MockPaymentService)
	handler := New(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestProcessPayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.ProcessPaymentRequest{
		OrderId: 100,
		Amount:  99.99,
		Method:  "card",
	}

	expectedPayment := &model.Payment{
		ID:      1,
		OrderID: 100,
		Amount:  99.99,
		Status:  "completed",
		Method:  "card",
	}

	mockService.On("ProcessPayment", ctx, &model.ProcessPaymentRequest{
		OrderID: 100,
		Amount:  99.99,
		Method:  "card",
	}).Return(expectedPayment, nil)

	resp, err := handler.ProcessPayment(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, int64(100), resp.OrderId)
	assert.Equal(t, 99.99, resp.Amount)
	assert.Equal(t, "completed", resp.Status)
	mockService.AssertExpectations(t)
}

func TestGetPayment_Success(t *testing.T) {
	mockService := new(MockPaymentService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetPaymentRequest{
		PaymentId: 1,
	}

	expectedPayment := &model.Payment{
		ID:      1,
		OrderID: 100,
		Amount:  99.99,
		Status:  "completed",
		Method:  "card",
	}

	mockService.On("GetPayment", ctx, int64(1)).Return(expectedPayment, nil)

	resp, err := handler.GetPayment(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, int64(100), resp.OrderId)
	mockService.AssertExpectations(t)
}

func TestGetPayment_NotFound(t *testing.T) {
	mockService := new(MockPaymentService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetPaymentRequest{
		PaymentId: 999,
	}

	mockService.On("GetPayment", ctx, int64(999)).Return(nil, nil)

	resp, err := handler.GetPayment(ctx, req)

	assert.NoError(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

