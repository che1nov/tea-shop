package handler

import (
	"context"
	"errors"
	"testing"

	"github.com/che1nov/tea-shop/users-service/internal/model"
	"github.com/che1nov/tea-shop/users-service/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	pb "github.com/che1nov/tea-shop/shared/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// MockUserService - мок для сервиса
type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserService) ValidateToken(tokenString string) (int64, string, error) {
	args := m.Called(tokenString)
	return args.Get(0).(int64), args.String(1), args.Error(2)
}

func (m *MockUserService) GenerateToken(user *model.User) (string, error) {
	args := m.Called(user)
	return args.String(0), args.Error(1)
}

func (m *MockUserService) Login(ctx context.Context, email, password string) (string, error) {
	args := m.Called(ctx, email, password)
	return args.String(0), args.Error(1)
}

func TestNew(t *testing.T) {
	mockService := new(MockUserService)
	
	handler := New(mockService)
	
	assert.NotNil(t, handler)
	assert.Equal(t, mockService, handler.service)
}

func TestCreateUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	
	expectedUser := &model.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashed_password",
	}
	
	mockService.On("CreateUser", ctx, &model.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}).Return(expectedUser, nil)
	
	resp, err := handler.CreateUser(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "test@example.com", resp.Email)
	assert.Equal(t, "Test User", resp.Name)
	mockService.AssertExpectations(t)
}

func TestCreateUser_MissingEmail(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Name:     "Test User",
		Password: "password123",
	}
	
	resp, err := handler.CreateUser(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "email is required")
	mockService.AssertNotCalled(t, "CreateUser")
}

func TestCreateUser_MissingName(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Email:    "test@example.com",
		Password: "password123",
	}
	
	resp, err := handler.CreateUser(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "name is required")
}

func TestCreateUser_MissingPassword(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Email: "test@example.com",
		Name:  "Test User",
	}
	
	resp, err := handler.CreateUser(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	assert.Contains(t, st.Message(), "password is required")
}

func TestCreateUser_EmailAlreadyExists(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}
	
	mockService.On("CreateUser", ctx, mock.AnythingOfType("*model.CreateUserRequest")).
		Return(nil, service.ErrEmailAlreadyExists)
	
	resp, err := handler.CreateUser(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.AlreadyExists, st.Code())
	assert.Contains(t, st.Message(), "already exists")
}

func TestGetUser_Success(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.GetUserRequest{
		UserId: 1,
	}
	
	expectedUser := &model.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashed_password",
	}
	
	mockService.On("GetUser", ctx, int64(1)).Return(expectedUser, nil)
	
	resp, err := handler.GetUser(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, int64(1), resp.Id)
	assert.Equal(t, "test@example.com", resp.Email)
	mockService.AssertExpectations(t)
}

func TestGetUser_InvalidUserId(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.GetUserRequest{
		UserId: 0,
	}
	
	resp, err := handler.GetUser(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.InvalidArgument, st.Code())
	mockService.AssertNotCalled(t, "GetUser")
}

func TestGetUser_NotFound(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.GetUserRequest{
		UserId: 999,
	}
	
	mockService.On("GetUser", ctx, int64(999)).Return(nil, nil)
	
	resp, err := handler.GetUser(ctx, req)
	
	assert.Error(t, err)
	assert.Nil(t, resp)
	
	st, ok := status.FromError(err)
	assert.True(t, ok)
	assert.Equal(t, codes.NotFound, st.Code())
	assert.Contains(t, st.Message(), "not found")
}

func TestValidateToken_ValidToken(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.ValidateTokenRequest{
		Token: "valid-token",
	}
	
	mockService.On("ValidateToken", "valid-token").Return(int64(1), "test@example.com", nil)
	
	resp, err := handler.ValidateToken(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.True(t, resp.Valid)
	assert.Equal(t, int64(1), resp.UserId)
	assert.Equal(t, "test@example.com", resp.Email)
	mockService.AssertExpectations(t)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.ValidateTokenRequest{
		Token: "invalid-token",
	}
	
	mockService.On("ValidateToken", "invalid-token").Return(int64(0), "", errors.New("invalid token"))
	
	resp, err := handler.ValidateToken(ctx, req)
	
	assert.NoError(t, err) // Handler не возвращает ошибку, а возвращает Valid: false
	assert.NotNil(t, resp)
	assert.False(t, resp.Valid)
}

func TestValidateToken_EmptyToken(t *testing.T) {
	mockService := new(MockUserService)
	handler := New(mockService)
	ctx := context.Background()
	
	req := &pb.ValidateTokenRequest{
		Token: "",
	}
	
	resp, err := handler.ValidateToken(ctx, req)
	
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.False(t, resp.Valid)
	mockService.AssertNotCalled(t, "ValidateToken")
}

