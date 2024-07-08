package auth

import (
	"context"
	"errors"

	"github.com/TiZir/gRPC_service/internal/services/auth"
	"github.com/TiZir/gRPC_service/internal/storage"
	gRPC_servicev1 "github.com/TiZir/gRPC_service_protos/protos/gen/go/gRPC_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email string, password string, appID int) (token string, err error)
	RegisterNewUser(ctx context.Context, email string, password string) (userID int64, err error)
	IsAdmin(ctx context.Context, userID int64) (isAdmin bool, err error)
}

type serverAPI struct {
	gRPC_servicev1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	gRPC_servicev1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

// Валидация через пакет должна быть/отдельная функция
func (s *serverAPI) Login(ctx context.Context, req *gRPC_servicev1.LoginRequest) (*gRPC_servicev1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "app id is required")
	}

	// Работа со структурой serverAPI
	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppId()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid credentials")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}

	return &gRPC_servicev1.LoginResponse{
		Token: token,
	}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *gRPC_servicev1.RegisterRequest) (*gRPC_servicev1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is required")
	}

	// Работа со структурой serverAPI
	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExists) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gRPC_servicev1.RegisterResponse{
		UserId: userID,
	}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *gRPC_servicev1.IsAdminRequest) (*gRPC_servicev1.IsAdminResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "email is required")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		if errors.Is(err, storage.ErrUserNotFound) {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Error(codes.Internal, "internal error")
	}
	return &gRPC_servicev1.IsAdminResponse{
		IsAdmin: isAdmin,
	}, nil

}
