package service

import (
	"context"
	"time"

	"github.com/che1nov/tea-shop/shared/pkg/logger"
	"github.com/che1nov/tea-shop/users-service/internal/model"
	"github.com/che1nov/tea-shop/users-service/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var ErrEmailAlreadyExists = repository.ErrEmailAlreadyExists

// UserServiceInterface определяет методы сервиса
type UserServiceInterface interface {
	CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error)
	GetUser(ctx context.Context, id int64) (*model.User, error)
	GenerateToken(user *model.User) (string, error)
	ValidateToken(tokenString string) (int64, string, error)
	ValidateTokenWithRole(tokenString string) (int64, string, string, error)
	Login(ctx context.Context, email, password string) (string, error)
}

type UserService struct {
	repo          repository.UserRepositoryInterface
	jwtSecret     string
	adminEmail    string
	adminPassword string
}

func New(repo repository.UserRepositoryInterface, jwtSecret string, adminEmail, adminPassword string) *UserService {
	return &UserService{
		repo:          repo,
		jwtSecret:     jwtSecret,
		adminEmail:    adminEmail,
		adminPassword: adminPassword,
	}
}

func (s *UserService) CreateUser(ctx context.Context, req *model.CreateUserRequest) (*model.User, error) {
	// Хешируем пароль
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Email:        req.Email,
		Name:         req.Name,
		PasswordHash: string(hash),
		Role:         model.RoleUser, // Все создаваемые пользователи имеют роль "user"
	}

	if err := s.repo.CreateUser(ctx, user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetUser(ctx context.Context, id int64) (*model.User, error) {
	return s.repo.GetUserByID(ctx, id)
}

func (s *UserService) GenerateToken(user *model.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": user.ID,
		"email":   user.Email,
		"role":    user.Role, // Добавляем роль в токен
		"exp":     time.Now().Add(time.Hour * 24).Unix(),
	})

	return token.SignedString([]byte(s.jwtSecret))
}

func (s *UserService) ValidateToken(tokenString string) (int64, string, error) {
	userID, email, _, err := s.ValidateTokenWithRole(tokenString)
	return userID, email, err
}

// ValidateTokenWithRole возвращает user_id, email и role из токена
func (s *UserService) ValidateTokenWithRole(tokenString string) (int64, string, string, error) {
	token, err := jwt.ParseWithClaims(
		tokenString,
		jwt.MapClaims{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(s.jwtSecret), nil
		},
	)

	if err != nil {
		return 0, "", "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", "", jwt.ErrSignatureInvalid
	}

	// Безопасное извлечение user_id
	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return 0, "", "", jwt.ErrSignatureInvalid
	}
	userID := int64(userIDFloat)

	// Безопасное извлечение email
	email, ok := claims["email"].(string)
	if !ok {
		return 0, "", "", jwt.ErrSignatureInvalid
	}

	// Безопасное извлечение role (если есть, иначе "user")
	role := model.RoleUser
	if roleVal, ok := claims["role"].(string); ok {
		role = roleVal
	}

	return userID, email, role, nil
}

func (s *UserService) Login(ctx context.Context, email, password string) (string, error) {
	// Проверяем админские credentials
	// Временное логирование для отладки
	if s.adminEmail == "" || s.adminPassword == "" {
		// Если переменные окружения не установлены, используем значения по умолчанию
		if s.adminEmail == "" {
			s.adminEmail = "admin@example.com"
		}
		if s.adminPassword == "" {
			s.adminPassword = "admin123"
		}
	}
	
	// Логирование для отладки (только при попытке входа с админским email)
	if email == s.adminEmail {
		logger.Info("Admin login attempt", 
			"input_email", email, 
			"config_email", s.adminEmail,
			"email_match", email == s.adminEmail, 
			"input_password_len", len(password),
			"config_password_len", len(s.adminPassword),
			"password_match", password == s.adminPassword)
	}
	
	if email == s.adminEmail && password == s.adminPassword {
		// Создаем виртуального пользователя-админа для токена
		adminUser := &model.User{
			ID:    0, // Специальный ID для админа
			Email: email,
			Name:  "Administrator",
			Role:  model.RoleAdmin,
		}
		return s.GenerateToken(adminUser)
	}

	// Обычный пользователь
	user, err := s.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", err
	}

	if user == nil {
		return "", jwt.ErrSignatureInvalid
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", err
	}

	// Убеждаемся, что у обычного пользователя роль "user"
	user.Role = model.RoleUser
	return s.GenerateToken(user)
}
