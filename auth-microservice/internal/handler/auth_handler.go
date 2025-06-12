package handler

import (
	"context"
	"log"

	authpb "auth-microservice/proto"
)

type AuthServiceHandler struct {
	authpb.UnimplementedAuthServiceServer
}

func NewAuthServiceHandler() *AuthServiceHandler {
	return &AuthServiceHandler{}
}

func (s *AuthServiceHandler) Register(ctx context.Context, req *authpb.RegisterRequest) (*authpb.RegisterResponse, error) {
	log.Printf("Register: email=%s, name=%s", req.Email, req.Name)

	return &authpb.RegisterResponse{
		Id:    "mock-id-123",
		Email: req.Email,
	}, nil
}

func (s *AuthServiceHandler) Login(ctx context.Context, req *authpb.LoginRequest) (*authpb.LoginResponse, error) {
	log.Printf("Login: email=%s", req.Email)

	return &authpb.LoginResponse{
		Token: "mock-jwt-token",
	}, nil
}
