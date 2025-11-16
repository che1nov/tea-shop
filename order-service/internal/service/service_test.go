package service

import (
	"context"
	"testing"

	"github.com/che1nov/tea-shop/order-service/internal/kafka"
	"github.com/che1nov/tea-shop/order-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"

	pb "github.com/che1nov/tea-shop/shared/pb"
)

// MockRepository - мок для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateOrder(ctx context.Context, order *model.Order) error {
	args := m.Called(ctx, order)
	if args.Error(0) == nil {
		order.ID = 1
	}
	return args.Error(0)
}

func (m *MockRepository) GetOrder(ctx context.Context, id int64) (*model.Order, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Order), args.Error(1)
}

func (m *MockRepository) UpdateOrderStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRepository) ListUserOrders(ctx context.Context, userID int64) ([]*model.Order, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Order), args.Error(1)
}

// MockProducer - мок для Kafka producer
type MockProducer struct {
	mock.Mock
}

func (m *MockProducer) PublishOrderCreated(ctx context.Context, event *kafka.OrderEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockProducer) PublishOrderCompleted(ctx context.Context, event *kafka.OrderEvent) error {
	args := m.Called(ctx, event)
	return args.Error(0)
}

func (m *MockProducer) Close() error {
	args := m.Called()
	return args.Error(0)
}

// MockGoodsServiceClient - мок для gRPC клиента goods service
type MockGoodsServiceClient struct {
	mock.Mock
}

func (m *MockGoodsServiceClient) GetGood(ctx context.Context, req *pb.GetGoodRequest, opts ...grpc.CallOption) (*pb.Good, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.Good), args.Error(1)
}

func (m *MockGoodsServiceClient) ListGoods(ctx context.Context, req *pb.ListGoodsRequest, opts ...grpc.CallOption) (*pb.ListGoodsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.ListGoodsResponse), args.Error(1)
}

func (m *MockGoodsServiceClient) CheckStock(ctx context.Context, req *pb.CheckStockRequest, opts ...grpc.CallOption) (*pb.CheckStockResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.CheckStockResponse), args.Error(1)
}

func (m *MockGoodsServiceClient) ReserveStock(ctx context.Context, req *pb.ReserveStockRequest, opts ...grpc.CallOption) (*pb.ReserveStockResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.ReserveStockResponse), args.Error(1)
}

// MockPaymentsServiceClient - мок для gRPC клиента payment service
type MockPaymentsServiceClient struct {
	mock.Mock
}

func (m *MockPaymentsServiceClient) ProcessPayment(ctx context.Context, req *pb.ProcessPaymentRequest, opts ...grpc.CallOption) (*pb.Payment, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.Payment), args.Error(1)
}

func (m *MockPaymentsServiceClient) GetPayment(ctx context.Context, req *pb.GetPaymentRequest, opts ...grpc.CallOption) (*pb.Payment, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*pb.Payment), args.Error(1)
}

func TestNew(t *testing.T) {
	mockRepo := new(MockRepository)
	mockProducer := new(MockProducer)
	mockGoodsClient := new(MockGoodsServiceClient)
	mockPaymentClient := new(MockPaymentsServiceClient)

	service := New(mockRepo, mockProducer, mockGoodsClient, mockPaymentClient)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, mockProducer, service.producer)
}

func TestGetOrder_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, new(MockProducer), new(MockGoodsServiceClient), new(MockPaymentsServiceClient))
	ctx := context.Background()

	expectedOrder := &model.Order{
		ID:         1,
		UserID:     100,
		Status:     "pending",
		TotalPrice: 99.99,
	}

	mockRepo.On("GetOrder", ctx, int64(1)).Return(expectedOrder, nil)

	order, err := service.GetOrder(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, expectedOrder, order)
	mockRepo.AssertExpectations(t)
}

func TestUpdateOrderStatus_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, new(MockProducer), new(MockGoodsServiceClient), new(MockPaymentsServiceClient))
	ctx := context.Background()

	expectedOrder := &model.Order{
		ID:         1,
		UserID:     100,
		Status:     "completed",
		TotalPrice: 99.99,
	}

	mockRepo.On("UpdateOrderStatus", ctx, int64(1), "completed").Return(nil)
	mockRepo.On("GetOrder", ctx, int64(1)).Return(expectedOrder, nil)

	order, err := service.UpdateOrderStatus(ctx, 1, "completed")

	assert.NoError(t, err)
	assert.NotNil(t, order)
	assert.Equal(t, "completed", order.Status)
	mockRepo.AssertExpectations(t)
}

func TestListUserOrders_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, new(MockProducer), new(MockGoodsServiceClient), new(MockPaymentsServiceClient))
	ctx := context.Background()

	expectedOrders := []*model.Order{
		{ID: 1, UserID: 100, Status: "pending"},
		{ID: 2, UserID: 100, Status: "completed"},
	}

	mockRepo.On("ListUserOrders", ctx, int64(100)).Return(expectedOrders, nil)

	orders, err := service.ListUserOrders(ctx, 100)

	assert.NoError(t, err)
	assert.NotNil(t, orders)
	assert.Equal(t, 2, len(orders))
	mockRepo.AssertExpectations(t)
}

// Note: CreateOrder тест более сложный, так как требует множества моков
// Для полного тестирования нужны моки для goods и payment клиентов

