package service

import (
	"context"
	"errors"
	"testing"

	"github.com/che1nov/tea-shop/payment-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository - мок для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreatePayment(ctx context.Context, payment *model.Payment) error {
	args := m.Called(ctx, payment)
	if args.Error(0) == nil {
		payment.ID = 1
	}
	return args.Error(0)
}

func (m *MockRepository) GetPayment(ctx context.Context, id int64) (*model.Payment, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func (m *MockRepository) UpdatePaymentStatus(ctx context.Context, id int64, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRepository) GetPaymentByOrderID(ctx context.Context, orderID int64) (*model.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Payment), args.Error(1)
}

func TestNew(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestProcessPayment_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	req := &model.ProcessPaymentRequest{
		OrderID: 1,
		Amount:  100.50,
		Method:  "card",
	}

	mockRepo.On("CreatePayment", ctx, mock.AnythingOfType("*model.Payment")).Return(nil).Run(func(args mock.Arguments) {
		payment := args.Get(1).(*model.Payment)
		payment.ID = 1
		payment.Status = "pending"
	})
	mockRepo.On("UpdatePaymentStatus", ctx, int64(1), mock.AnythingOfType("string")).Return(nil)

	payment, err := service.ProcessPayment(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, int64(1), payment.ID)
	assert.Equal(t, int64(1), payment.OrderID)
	assert.Equal(t, 100.50, payment.Amount)
	assert.Equal(t, "card", payment.Method)
	assert.Contains(t, []string{"completed", "failed"}, payment.Status)
	mockRepo.AssertExpectations(t)
}

func TestProcessPayment_CreateError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	req := &model.ProcessPaymentRequest{
		OrderID: 1,
		Amount:  100.50,
		Method:  "card",
	}

	mockRepo.On("CreatePayment", ctx, mock.AnythingOfType("*model.Payment")).Return(errors.New("database error"))

	payment, err := service.ProcessPayment(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, payment)
	mockRepo.AssertExpectations(t)
}

func TestGetPayment_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	expectedPayment := &model.Payment{
		ID:      1,
		OrderID: 100,
		Amount:  99.99,
		Status:  "completed",
		Method:  "card",
	}

	mockRepo.On("GetPayment", ctx, int64(1)).Return(expectedPayment, nil)

	payment, err := service.GetPayment(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, expectedPayment, payment)
	mockRepo.AssertExpectations(t)
}

func TestGetPayment_NotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetPayment", ctx, int64(999)).Return(nil, nil)

	payment, err := service.GetPayment(ctx, 999)

	assert.NoError(t, err)
	assert.Nil(t, payment)
	mockRepo.AssertExpectations(t)
}

func TestGetPaymentByOrderID_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	expectedPayment := &model.Payment{
		ID:      1,
		OrderID: 100,
		Amount:  99.99,
		Status:  "completed",
		Method:  "card",
	}

	mockRepo.On("GetPaymentByOrderID", ctx, int64(100)).Return(expectedPayment, nil)

	payment, err := service.GetPaymentByOrderID(ctx, 100)

	assert.NoError(t, err)
	assert.NotNil(t, payment)
	assert.Equal(t, expectedPayment, payment)
	mockRepo.AssertExpectations(t)
}

