package handler

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "github.com/che1nov/tea-shop/shared/pb"
	"github.com/che1nov/tea-shop/users-service/internal/model"
	"github.com/che1nov/tea-shop/users-service/internal/service"
)

type UsersHandler struct {
	service service.UserServiceInterface
	pb.UnimplementedUsersServiceServer
}

func New(svc service.UserServiceInterface) *UsersHandler {
	return &UsersHandler{
		service: svc,
	}
}

func (h *UsersHandler) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.User, error) {
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if req.Name == "" {
		return nil, status.Errorf(codes.InvalidArgument, "name is required")
	}
	if req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	user, err := h.service.CreateUser(ctx, &model.CreateUserRequest{
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	})
	if err != nil {
		if err == service.ErrEmailAlreadyExists {
			return nil, status.Errorf(codes.AlreadyExists, "user with email %s already exists", req.Email)
		}
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	return &pb.User{
		Id:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.Unix(),
	}, nil
}

func (h *UsersHandler) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.User, error) {
	if req.UserId <= 0 {
		return nil, status.Errorf(codes.InvalidArgument, "user_id must be greater than 0")
	}

	user, err := h.service.GetUser(ctx, req.UserId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	if user == nil {
		return nil, status.Errorf(codes.NotFound, "user with id %d not found", req.UserId)
	}

	return &pb.User{
		Id:           user.ID,
		Email:        user.Email,
		Name:         user.Name,
		PasswordHash: user.PasswordHash,
		Role:         user.Role,
		CreatedAt:    user.CreatedAt.Unix(),
	}, nil
}

func (h *UsersHandler) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	if req.Email == "" {
		return nil, status.Errorf(codes.InvalidArgument, "email is required")
	}
	if req.Password == "" {
		return nil, status.Errorf(codes.InvalidArgument, "password is required")
	}

	token, err := h.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid email or password")
	}

	// Получаем данные пользователя из токена (включая role)
	userID, email, role, err := h.service.ValidateTokenWithRole(token)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "failed to validate token: %v", err)
	}

	// Если userID = 0, это админ (виртуальный пользователь)
	if userID == 0 {
		return &pb.LoginResponse{
			Token: token,
			User: &pb.User{
				Id:           0,
				Email:        email,
				Name:         "Administrator",
				PasswordHash: "", // Не возвращаем пароль для админа
				Role:         role,
				CreatedAt:    0,
			},
		}, nil
	}

	// Для обычных пользователей получаем данные из БД
	user, err := h.service.GetUser(ctx, userID)
	if err != nil || user == nil {
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	return &pb.LoginResponse{
		Token: token,
		User: &pb.User{
			Id:           user.ID,
			Email:        user.Email,
			Name:         user.Name,
			PasswordHash: user.PasswordHash,
			Role:         user.Role,
			CreatedAt:    user.CreatedAt.Unix(),
		},
	}, nil
}

func (h *UsersHandler) ValidateToken(ctx context.Context, req *pb.ValidateTokenRequest) (*pb.ValidateTokenResponse, error) {
	if req.Token == "" {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	userID, email, role, err := h.service.ValidateTokenWithRole(req.Token)
	if err != nil {
		return &pb.ValidateTokenResponse{Valid: false}, nil
	}

	return &pb.ValidateTokenResponse{
		Valid:  true,
		UserId: userID,
		Email:  email,
		Role:   role,
	}, nil
}
