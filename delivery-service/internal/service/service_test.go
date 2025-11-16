package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/che1nov/tea-shop/delivery-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository - мок для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateDelivery(ctx context.Context, delivery *model.Delivery) error {
	args := m.Called(ctx, delivery)
	if args.Error(0) == nil {
		delivery.ID = 1
	}
	return args.Error(0)
}

func (m *MockRepository) GetDelivery(ctx context.Context, id int64) (*model.Delivery, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Delivery), args.Error(1)
}

func (m *MockRepository) GetDeliveryByOrderID(ctx context.Context, orderID int64) (*model.Delivery, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Delivery), args.Error(1)
}

func (m *MockRepository) UpdateDeliveryStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func TestNew(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestCreateDelivery_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	req := &model.CreateDeliveryRequest{
		OrderID: 1,
		Address: "123 Main St",
	}

	mockRepo.On("CreateDelivery", ctx, mock.AnythingOfType("*model.Delivery")).Return(nil).Run(func(args mock.Arguments) {
		delivery := args.Get(1).(*model.Delivery)
		delivery.ID = 1
		delivery.Status = "pending"
	})

	delivery, err := service.CreateDelivery(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, delivery)
	assert.Equal(t, int64(1), delivery.ID)
	assert.Equal(t, int64(1), delivery.OrderID)
	assert.Equal(t, "123 Main St", delivery.Address)
	assert.Equal(t, "pending", delivery.Status)
	mockRepo.AssertExpectations(t)
}

func TestGetDelivery_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	expectedDelivery := &model.Delivery{
		ID:        1,
		OrderID:   1,
		Address:   "123 Main St",
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetDelivery", ctx, int64(1)).Return(expectedDelivery, nil)

	delivery, err := service.GetDelivery(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, delivery)
	assert.Equal(t, expectedDelivery.ID, delivery.ID)
	assert.Equal(t, expectedDelivery.OrderID, delivery.OrderID)
	mockRepo.AssertExpectations(t)
}

func TestGetDelivery_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetDelivery", ctx, int64(999)).Return(nil, nil)

	delivery, err := service.GetDelivery(ctx, 999)

	assert.NoError(t, err)
	assert.Nil(t, delivery)
	mockRepo.AssertExpectations(t)
}

func TestGetDeliveryByOrderID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	expectedDelivery := &model.Delivery{
		ID:        1,
		OrderID:   100,
		Address:   "456 Oak Ave",
		Status:    "in_transit",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("GetDeliveryByOrderID", ctx, int64(100)).Return(expectedDelivery, nil)

	delivery, err := service.GetDeliveryByOrderID(ctx, 100)

	assert.NoError(t, err)
	assert.NotNil(t, delivery)
	assert.Equal(t, expectedDelivery.OrderID, delivery.OrderID)
	mockRepo.AssertExpectations(t)
}

func TestUpdateDeliveryStatus_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	updatedDelivery := &model.Delivery{
		ID:        1,
		OrderID:   1,
		Address:   "123 Main St",
		Status:    "delivered",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	mockRepo.On("UpdateDeliveryStatus", ctx, int64(1), "delivered").Return(nil)
	mockRepo.On("GetDelivery", ctx, int64(1)).Return(updatedDelivery, nil)

	delivery, err := service.UpdateDeliveryStatus(ctx, 1, "delivered")

	assert.NoError(t, err)
	assert.NotNil(t, delivery)
	assert.Equal(t, "delivered", delivery.Status)
	mockRepo.AssertExpectations(t)
}

func TestUpdateDeliveryStatus_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("UpdateDeliveryStatus", ctx, int64(999), "delivered").Return(nil)
	mockRepo.On("GetDelivery", ctx, int64(999)).Return(nil, nil)

	delivery, err := service.UpdateDeliveryStatus(ctx, 999, "delivered")

	assert.NoError(t, err)
	assert.Nil(t, delivery)
	mockRepo.AssertExpectations(t)
}

func TestUpdateDeliveryStatus_UpdateError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("UpdateDeliveryStatus", ctx, int64(1), "delivered").Return(errors.New("database error"))

	delivery, err := service.UpdateDeliveryStatus(ctx, 1, "delivered")

	assert.Error(t, err)
	assert.Nil(t, delivery)
	mockRepo.AssertExpectations(t)
}

