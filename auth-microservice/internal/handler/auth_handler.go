package handler

import (
	"auth-microservice/internal/service"
	"context"
	"log"

	authpb "auth-microservice/proto"

	"auth-microservice/internal/model"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthServiceHandler struct {
	authpb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthServiceHandler(authService service.AuthService) *AuthServiceHandler {
	return &AuthServiceHandler{authService: authService}
}

func (s *AuthServiceHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	user := &model.User{
		Email:    req.Email,
		Password: req.Password,
		Name:     req.Name,
	}

	err := s.authService.Register(ctx, user)
	if err != nil {
		log.Printf("Register failed: %v", err)
		return nil, status.Errorf(codes.Internal, "registration failed")
	}
	log.Printf("Register success")
	return &authpb.RegisterResponse{
		Id:    user.ID.Hex(),
		Email: user.Email,
	}, nil
}

func (s *AuthServiceHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	token, err := s.authService.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("Login failed: %v", err)
		return nil, status.Errorf(codes.Unauthenticated, "invalid credentials")
	}
	log.Printf("Login success")
	return &authpb.LoginResponse{
		Token: token,
	}, nil
}

func (s *AuthServiceHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	err := s.authService.Logout(ctx, req.Token)
	if err != nil {
		log.Printf("Logout failed: %v", err)
		return nil, status.Errorf(codes.Internal, "logout failed")
	}
	log.Printf("Logout success")
	return &authpb.LogoutResponse{
		Success: true,
	}, nil
}
