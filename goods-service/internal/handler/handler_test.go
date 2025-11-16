package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/che1nov/tea-shop/goods-service/internal/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/che1nov/tea-shop/shared/pb"
)

// MockGoodsService - мок для сервиса
type MockGoodsService struct {
	mock.Mock
}

func (m *MockGoodsService) CreateGood(ctx context.Context, req *model.CreateGoodRequest) (*model.Good, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Good), args.Error(1)
}

func (m *MockGoodsService) GetGood(ctx context.Context, id int64) (*model.Good, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.Good), args.Error(1)
}

func (m *MockGoodsService) ListGoods(ctx context.Context, limit, offset int32) ([]*model.Good, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*model.Good), args.Error(1)
}

func (m *MockGoodsService) GetTotalGoods(ctx context.Context) (int32, error) {
	args := m.Called(ctx)
	return args.Get(0).(int32), args.Error(1)
}

func (m *MockGoodsService) CheckStock(ctx context.Context, goodID int64, quantity int32) (bool, error) {
	args := m.Called(ctx, goodID, quantity)
	return args.Bool(0), args.Error(1)
}

func (m *MockGoodsService) ReserveStock(ctx context.Context, goodID int64, quantity int32, orderID int64) (bool, error) {
	args := m.Called(ctx, goodID, quantity, orderID)
	return args.Bool(0), args.Error(1)
}

func TestNew(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)

	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestGetGood_Success(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetGoodRequest{
		GoodId: 1,
	}

	expectedGood := &model.Good{
		ID:          1,
		Name:        "Test Good",
		Description: "Test Description",
		Price:       99.99,
		Stock:       100,
	}

	mockService.On("GetGood", ctx, int64(1)).Return(expectedGood, nil)

	resp, err := handler.GetGood(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "Test Good", resp.Name)
	mockService.AssertExpectations(t)
}

func TestGetGood_NotFound(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.GetGoodRequest{
		GoodId: 999,
	}

	mockService.On("GetGood", ctx, int64(999)).Return(nil, nil)

	resp, err := handler.GetGood(ctx, req)

	assert.NoError(t, err)
	assert.Nil(t, resp)
	mockService.AssertExpectations(t)
}

func TestListGoods_Success(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.ListGoodsRequest{
		Limit:  10,
		Offset: 0,
	}

	goods := []*model.Good{
		{ID: 1, Name: "Good 1", Price: 10.0, Stock: 50},
		{ID: 2, Name: "Good 2", Price: 20.0, Stock: 30},
	}

	mockService.On("ListGoods", ctx, int32(10), int32(0)).Return(goods, nil)
	mockService.On("GetTotalGoods", ctx).Return(int32(2), nil)

	resp, err := handler.ListGoods(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, 2, len(resp.Goods))
	assert.Equal(t, int32(2), resp.Total)
	mockService.AssertExpectations(t)
}

func TestCheckStock_Available(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.CheckStockRequest{
		GoodId:   1,
		Quantity: 10,
	}

	mockService.On("CheckStock", ctx, int64(1), int32(10)).Return(true, nil)

	resp, err := handler.CheckStock(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Available)
	mockService.AssertExpectations(t)
}

func TestCheckStock_NotAvailable(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.CheckStockRequest{
		GoodId:   1,
		Quantity: 100,
	}

	mockService.On("CheckStock", ctx, int64(1), int32(100)).Return(false, nil)

	resp, err := handler.CheckStock(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Available)
	mockService.AssertExpectations(t)
}

func TestReserveStock_Success(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.ReserveStockRequest{
		GoodId:   1,
		Quantity: 10,
		OrderId:  100,
	}

	mockService.On("ReserveStock", ctx, int64(1), int32(10), int64(100)).Return(true, nil)

	resp, err := handler.ReserveStock(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Success)
	assert.Empty(t, resp.Error)
	mockService.AssertExpectations(t)
}

func TestReserveStock_InsufficientStock(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.ReserveStockRequest{
		GoodId:   1,
		Quantity: 100,
		OrderId:  100,
	}

	mockService.On("ReserveStock", ctx, int64(1), int32(100), int64(100)).Return(false, nil)

	resp, err := handler.ReserveStock(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "insufficient stock", resp.Error)
	mockService.AssertExpectations(t)
}

func TestReserveStock_Error(t *testing.T) {
	mockService := new(MockGoodsService)
	handler := New(mockService)
	ctx := context.Background()

	req := &pb.ReserveStockRequest{
		GoodId:   1,
		Quantity: 10,
		OrderId:  100,
	}

	mockService.On("ReserveStock", ctx, int64(1), int32(10), int64(100)).Return(false, errors.New("database error"))

	resp, err := handler.ReserveStock(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Success)
	assert.Equal(t, "database error", resp.Error)
	mockService.AssertExpectations(t)
}

