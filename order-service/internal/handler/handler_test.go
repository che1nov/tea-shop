package handler

import (
	"context"
	"testing"

	"github.com/che1nov/tea-shop/order-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/che1nov/tea-shop/shared/pb"
)

// MockOrderService - мок для сервиса
type MockOrderService struct {
	mock.Mock
}

func (m *MockOrderService) CreateOrder(ctx context.Context, req *model.CreateOrderRequest) (*model.Order, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderService) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderService) UpdateOrderStatus(ctx context.Context, id int64, status string) (*model.Order, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockOrderService) ListUserOrders(ctx context.Context, userID int64) ([]*model.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

func TestNew(t *testing.T) {
	mockService := new(MockOrderService)
	handler := New(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestGetOrder_Success(t *testing.T) {
	mockService := new(MockOrderService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetOrderRequest{
		OrderId: 1,
	}

	expectedOrder := &model.Order{
		ID:         1,
		UserID:     100,
		Status:     "pending",
		TotalPrice: 99.99,
		Items: []model.OrderItem{
			{GoodID: 1, Quantity: 2, Price: 49.99},
		},
	}

	mockService.On("GetOrder", ctx, int64(1)).Return(expectedOrder, nil)

	resp, err := handler.GetOrder(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, int64(100), resp.UserId)
	assert.Equal(t, "pending", resp.Status)
	mockService.AssertExpectations(t)
}

func TestGetOrder_NotFound(t *testing.T) {
	mockService := new(MockOrderService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetOrderRequest{
		OrderId: 999,
	}

	mockService.On("GetOrder", ctx, int64(999)).Return(nil, nil)

	resp, err := handler.GetOrder(ctx, req)

	assert.NoError(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	mockService := new(MockOrderService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.UpdateOrderStatusRequest{
		OrderId: 1,
		Status:  "completed",
	}

	expectedOrder := &model.Order{
		ID:         1,
		UserID:     100,
		Status:     "completed",
		TotalPrice: 99.99,
	}

	mockService.On("UpdateOrderStatus", ctx, int64(1), "completed").Return(expectedOrder, nil)

	resp, err := handler.UpdateOrderStatus(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, "completed", resp.Status)
	mockService.AssertExpectations(t)
}

func TestOrderToProto(t *testing.T) {
	mockService := new(MockOrderService)
	handler := New(mockService)

	order := &model.Order{
		ID:         1,
		UserID:     100,
		Status:     "pending",
		TotalPrice: 99.99,
		Items: []model.OrderItem{
			{GoodID: 1, Quantity: 2, Price: 49.99},
			{GoodID: 2, Quantity: 1, Price: 50.0},
		},
	}

	pbOrder := handler.orderToProto(order)

	assert.NotNil(t, pbOrder)
	assert.Equal(t, int64(1), pbOrder.Id)
	assert.Equal(t, int64(100), pbOrder.UserId)
	assert.Equal(t, 2, len(pbOrder.Items))
	assert.Equal(t, int64(1), pbOrder.Items[0].GoodId)
	assert.Equal(t, int32(2), pbOrder.Items[0].Quantity)
}

