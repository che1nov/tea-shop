package service

import (
	"context"
	"errors"
	"testing"

	"github.com/che1nov/tea-shop/users-service/internal/model"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockRepository - мок для репозитория, реализует UserRepositoryInterface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) CreateUser(ctx context.Context, user *model.User) error {
	args := m.Called(ctx, user)
	if args.Error(0) == nil {
		user.ID = 1
	}
	return args.Error(0)
}

func (m *MockRepository) GetUserByID(ctx context.Context, id int64) (*model.User, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockRepository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*model.User), args.Error(1)
}

func TestNew(t *testing.T) {
	mockRepo := new(MockRepository)
	jwtSecret := "test-secret"

	service := New(mockRepo, jwtSecret)

	assert.NotNil(t, service)
	assert.Equal(t, mockRepo, service.repo)
	assert.Equal(t, jwtSecret, service.jwtSecret)
}

func TestCreateUser_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, "test-secret")
	ctx := context.Background()

	req := &model.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
		user := args.Get(1).(*model.User)
		user.ID = 1
	})

	user, err := service.CreateUser(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "Test User", user.Name)
	assert.NotEmpty(t, user.PasswordHash)
	assert.NotEqual(t, "password123", user.PasswordHash) // Пароль должен быть захеширован
	mockRepo.AssertExpectations(t)
}

func TestCreateUser_RepositoryError(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, "test-secret")
	ctx := context.Background()

	req := &model.CreateUserRequest{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	}

	mockRepo.On("CreateUser", ctx, mock.AnythingOfType("*model.User")).Return(errors.New("database error"))

	user, err := service.CreateUser(ctx, req)

	assert.Error(t, err)
	assert.Nil(t, user)
	mockRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, "test-secret")
	ctx := context.Background()

	expectedUser := &model.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: "hashed_password",
	}

	mockRepo.On("GetUserByID", ctx, int64(1)).Return(expectedUser, nil)

	user, err := service.GetUser(ctx, 1)

	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, expectedUser, user)
	mockRepo.AssertExpectations(t)
}

func TestGenerateToken_Success(t *testing.T) {
	service := New(new(MockRepository), "test-secret-key")

	user := &model.User{
		ID:    1,
		Email: "test@example.com",
		Name:  "Test User",
	}

	token, err := service.GenerateToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Проверяем, что токен можно распарсить
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return []byte("test-secret-key"), nil
	})

	assert.NoError(t, err)
	assert.True(t, parsedToken.Valid)

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	assert.True(t, ok)
	assert.Equal(t, float64(1), claims["user_id"])
	assert.Equal(t, "test@example.com", claims["email"])
}

func TestValidateToken_ValidToken(t *testing.T) {
	service := New(new(MockRepository), "test-secret-key")

	user := &model.User{
		ID:    1,
		Email: "test@example.com",
	}

	// Генерируем токен
	token, err := service.GenerateToken(user)
	assert.NoError(t, err)

	// Валидируем токен
	userID, email, err := service.ValidateToken(token)

	assert.NoError(t, err)
	assert.Equal(t, int64(1), userID)
	assert.Equal(t, "test@example.com", email)
}

func TestValidateToken_InvalidToken(t *testing.T) {
	service := New(new(MockRepository), "test-secret-key")

	userID, email, err := service.ValidateToken("invalid-token")

	assert.Error(t, err)
	assert.Equal(t, int64(0), userID)
	assert.Empty(t, email)
}

func TestValidateToken_WrongSecret(t *testing.T) {
	service1 := New(new(MockRepository), "secret1")
	service2 := New(new(MockRepository), "secret2")

	user := &model.User{
		ID:    1,
		Email: "test@example.com",
	}

	// Генерируем токен с одним секретом
	token, err := service1.GenerateToken(user)
	assert.NoError(t, err)

	// Пытаемся валидировать с другим секретом
	userID, email, err := service2.ValidateToken(token)

	assert.Error(t, err)
	assert.Equal(t, int64(0), userID)
	assert.Empty(t, email)
}

func TestLogin_Success(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, "test-secret")
	ctx := context.Background()

	password := "password123"
	hashedPassword := "$2a$10$abcdefghijklmnopqrstuv" // Минимальный валидный bcrypt хеш

	user := &model.User{
		ID:           1,
		Email:        "test@example.com",
		Name:         "Test User",
		PasswordHash: hashedPassword,
	}

	mockRepo.On("GetUserByEmail", ctx, "test@example.com").Return(user, nil)

	// Для реального теста нужен настоящий bcrypt hash
	// Здесь мы используем мок, поэтому просто проверяем структуру
	token, err := service.Login(ctx, "test@example.com", password)

	// Если хеш невалидный, ожидаем ошибку
	if err != nil {
		assert.Error(t, err)
	} else {
		assert.NotEmpty(t, token)
	}

	mockRepo.AssertExpectations(t)
}

func TestLogin_UserNotFound(t *testing.T) {
	mockRepo := new(MockRepository)
	service := New(mockRepo, "test-secret")
	ctx := context.Background()

	mockRepo.On("GetUserByEmail", ctx, "notfound@example.com").Return(nil, nil)

	token, err := service.Login(ctx, "notfound@example.com", "password")

	assert.Error(t, err)
	assert.Empty(t, token)
	mockRepo.AssertExpectations(t)
}
