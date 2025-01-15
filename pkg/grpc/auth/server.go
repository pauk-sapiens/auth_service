package auth

import (
	"context"
	"fmt"

	auth1 "github.com/pauk-sapiens/protos/gen/go/auth"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Auth interface {
	Login(ctx context.Context, email, password string, appId int32) (token string, err error)
	RegisterNewUser(ctx context.Context, email, passowrd string) (userID int64, err error)
	IsAdmin(ctx context.Context, userId int64) (bool, error)
}

type serverAPI struct {
	auth1.UnimplementedAuthServer
	auth Auth
}

func Register(gRPC *grpc.Server, auth Auth) {
	auth1.RegisterAuthServer(gRPC, &serverAPI{auth: auth})
}

func (s *serverAPI) Login(ctx context.Context, req *auth1.LoginRequest) (*auth1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	if req.GetAppId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "appId is empty")
	}

	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), req.GetAppId())
	if err != nil {
		return nil, fmt.Errorf("cannot login user, err: %w", err)
	}

	return &auth1.LoginResponse{Token: token}, nil
}

func (s *serverAPI) Register(ctx context.Context, req *auth1.RegisterRequest) (*auth1.RegisterResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "email is empty")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "password is empty")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		return nil, fmt.Errorf("cannot register user, err: %w", err)
	}

	return &auth1.RegisterResponse{UserId: userID}, nil
}

func (s *serverAPI) IsAdmin(ctx context.Context, req *auth1.IsAdminRequest) (*auth1.IsAdminResponse, error) {
	if req.GetUserId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "userID is empty")
	}

	isAdmin, err := s.auth.IsAdmin(ctx, req.GetUserId())
	if err != nil {
		return nil, fmt.Errorf("cannot determine user, err: %w", err)
	}

	return &auth1.IsAdminResponse{IsAdmin: isAdmin}, nil
}
