package service

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/che1nov/tea-shop/goods-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository - мок для репозитория
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateGood(ctx context.Context, good *model.Good) error {
	args := m.Called(ctx, good)
	if args.Error(0) == nil {
		good.ID = 1
	}
	return args.Error(0)
}

func (m *MockRepository) GetGood(ctx context.Context, id int64) (*model.Good, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Good), args.Error(1)
}

func (m *MockRepository) ListGoods(ctx context.Context, limit, offset int32) ([]*model.Good, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Good), args.Error(1)
}

func (m *MockRepository) ReserveStock(ctx context.Context, goodID int64, quantity int32, orderID int64) error {
	args := m.Called(ctx, goodID, quantity, orderID)
	return args.Error(0)
}

func (m *MockRepository) GetTotalGoods(ctx context.Context) (int32, error) {
	args := m.Called(ctx)
	return args.Get(0).(int32), args.Error(1)
}

func TestNew(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
}

func TestCreateGood_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	req := &model.CreateGoodRequest{
		Name:        "Test Good",
		Description: "Test Description",
		Price:       99.99,
		Stock:       100,
	}

	mockRepo.On("CreateGood", ctx, mock.AnythingOfType("*model.Good")).Return(nil).Run(func(args mock.Arguments) {
		good := args.Get(1).(*model.Good)
		good.ID = 1
	})

	good, err := service.CreateGood(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, good)
	assert.Equal(t, "Test Good", good.Name)
	assert.Equal(t, "Test Description", good.Description)
	assert.Equal(t, 99.99, good.Price)
	assert.Equal(t, int32(100), good.Stock)
	mockRepo.AssertExpectations(t)
}

func TestGetGood_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	expectedGood := &model.Good{
		ID:          1,
		Name:        "Test Good",
		Description: "Test Description",
		Price:       99.99,
		Stock:       100,
	}

	mockRepo.On("GetGood", ctx, int64(1)).Return(expectedGood, nil)

	good, err := service.GetGood(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, good)
	assert.Equal(t, expectedGood, good)
	mockRepo.AssertExpectations(t)
}

func TestListGoods_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	expectedGoods := []*model.Good{
		{ID: 1, Name: "Good 1", Price: 10.0, Stock: 50},
		{ID: 2, Name: "Good 2", Price: 20.0, Stock: 30},
	}

	mockRepo.On("ListGoods", ctx, int32(10), int32(0)).Return(expectedGoods, nil)

	goods, err := service.ListGoods(ctx, 10, 0)

	assert.NoError(t, err)
	assert.NotNil(t, goods)
	assert.Equal(t, 2, len(goods))
	mockRepo.AssertExpectations(t)
}

func TestGetTotalGoods_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetTotalGoods", ctx).Return(int32(100), nil)

	total, err := service.GetTotalGoods(ctx)

	assert.NoError(t, err)
	assert.Equal(t, int32(100), total)
	mockRepo.AssertExpectations(t)
}

func TestCheckStock_Available(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	good := &model.Good{
		ID:    1,
		Stock: 100,
	}

	mockRepo.On("GetGood", ctx, int64(1)).Return(good, nil)

	available, err := service.CheckStock(ctx, 1, 50)

	assert.NoError(t, err)
	assert.True(t, available)
	mockRepo.AssertExpectations(t)
}

func TestCheckStock_NotAvailable(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	good := &model.Good{
		ID:    1,
		Stock: 10,
	}

	mockRepo.On("GetGood", ctx, int64(1)).Return(good, nil)

	available, err := service.CheckStock(ctx, 1, 50)

	assert.NoError(t, err)
	assert.False(t, available)
	mockRepo.AssertExpectations(t)
}

func TestCheckStock_GoodNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("GetGood", ctx, int64(999)).Return(nil, nil)

	available, err := service.CheckStock(ctx, 999, 10)

	assert.NoError(t, err)
	assert.False(t, available)
	mockRepo.AssertExpectations(t)
}

func TestReserveStock_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("ReserveStock", ctx, int64(1), int32(10), int64(100)).Return(nil)

	success, err := service.ReserveStock(ctx, 1, 10, 100)

	assert.NoError(t, err)
	assert.True(t, success)
	mockRepo.AssertExpectations(t)
}

func TestReserveStock_InsufficientStock(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("ReserveStock", ctx, int64(1), int32(100), int64(100)).Return(sql.ErrNoRows)

	success, err := service.ReserveStock(ctx, 1, 100, 100)

	assert.NoError(t, err)
	assert.False(t, success)
	mockRepo.AssertExpectations(t)
}

func TestReserveStock_Error(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo)
	ctx := context.Background()

	mockRepo.On("ReserveStock", ctx, int64(1), int32(10), int64(100)).Return(errors.New("database error"))

	success, err := service.ReserveStock(ctx, 1, 10, 100)

	assert.Error(t, err)
	assert.False(t, success)
	mockRepo.AssertExpectations(t)
}

