package handler

import (
	"context"
	"testing"
	"time"

	"github.com/che1nov/tea-shop/delivery-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	pb "github.com/che1nov/tea-shop/shared/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockDeliveryService - мок для сервиса
type MockDeliveryService struct {
	mock.Mock
}

func (m *MockDeliveryService) CreateDelivery(ctx context.Context, req *model.CreateDeliveryRequest) (*model.Delivery, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Delivery), args.Error(1)
}

func (m *MockDeliveryService) GetDelivery(ctx context.Context, id int64) (*model.Delivery, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Delivery), args.Error(1)
}

func (m *MockDeliveryService) GetDeliveryByOrderID(ctx context.Context, orderID int64) (*model.Delivery, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Delivery), args.Error(1)
}

func (m *MockDeliveryService) UpdateDeliveryStatus(ctx context.Context, id int64, status string) (*model.Delivery, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Delivery), args.Error(1)
}

func TestCreateDelivery_Success(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.CreateDeliveryRequest{
		OrderId: 1,
		Address: "123 Main St",
	}

	expectedDelivery := &model.Delivery{
		ID:        1,
		OrderID:   1,
		Address:   "123 Main St",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("CreateDelivery", ctx, &model.CreateDeliveryRequest{
		OrderID: 1,
		Address: "123 Main St",
	}).Return(expectedDelivery, nil)

	resp, err := handler.CreateDelivery(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, int64(1), resp.OrderId)
	assert.Equal(t, "123 Main St", resp.Address)
	assert.Equal(t, "pending", resp.Status)
	mockService.AssertExpectations(t)
}

func TestCreateDelivery_InvalidOrderID(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.CreateDeliveryRequest{
		OrderId: 0,
		Address: "123 Main St",
	}

	resp, err := handler.CreateDelivery(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertNotCalled(t, "CreateDelivery")
}

func TestCreateDelivery_InvalidAddress(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.CreateDeliveryRequest{
		OrderId: 1,
		Address: "",
	}

	resp, err := handler.CreateDelivery(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertNotCalled(t, "CreateDelivery")
}

func TestGetDelivery_Success(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetDeliveryRequest{
		DeliveryId: 1,
	}

	expectedDelivery := &model.Delivery{
		ID:        1,
		OrderID:   1,
		Address:   "123 Main St",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("GetDelivery", ctx, int64(1)).Return(expectedDelivery, nil)

	resp, err := handler.GetDelivery(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, int64(1), resp.OrderId)
	mockService.AssertExpectations(t)
}

func TestGetDelivery_NotFound(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetDeliveryRequest{
		DeliveryId: 999,
	}

	mockService.On("GetDelivery", ctx, int64(999)).Return(nil, nil)

	resp, err := handler.GetDelivery(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	mockService.AssertExpectations(t)
}

func TestGetDelivery_InvalidID(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetDeliveryRequest{
		DeliveryId: 0,
	}

	resp, err := handler.GetDelivery(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertNotCalled(t, "GetDelivery")
}

func TestUpdateDeliveryStatus_Success(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.UpdateDeliveryStatusRequest{
		DeliveryId: 1,
		Status:     "delivered",
	}

	expectedDelivery := &model.Delivery{
		ID:        1,
		OrderID:   1,
		Address:   "123 Main St",
		Status:    "delivered",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockService.On("UpdateDeliveryStatus", ctx, int64(1), "delivered").Return(expectedDelivery, nil)

	resp, err := handler.UpdateDeliveryStatus(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "delivered", resp.Status)
	mockService.AssertExpectations(t)
}

func TestUpdateDeliveryStatus_InvalidID(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.UpdateDeliveryStatusRequest{
		DeliveryId: 0,
		Status:     "delivered",
	}

	resp, err := handler.UpdateDeliveryStatus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertNotCalled(t, "UpdateDeliveryStatus")
}

func TestUpdateDeliveryStatus_InvalidStatus(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.UpdateDeliveryStatusRequest{
		DeliveryId: 1,
		Status:     "",
	}

	resp, err := handler.UpdateDeliveryStatus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertNotCalled(t, "UpdateDeliveryStatus")
}

func TestUpdateDeliveryStatus_NotFound(t *testing.T) {
	mockService := new(MockDeliveryService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.UpdateDeliveryStatusRequest{
		DeliveryId: 999,
		Status:     "delivered",
	}

	mockService.On("UpdateDeliveryStatus", ctx, int64(999), "delivered").Return(nil, nil)

	resp, err := handler.UpdateDeliveryStatus(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, resp)
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	mockService.AssertExpectations(t)
}

