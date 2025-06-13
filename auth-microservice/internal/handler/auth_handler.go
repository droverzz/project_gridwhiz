package handler

import (
	"auth-microservice/internal/service"
	"context"
	"log"

	authpb "auth-microservice/proto"

	"auth-microservice/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
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
		Role:     "user",
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
		Role:  user.Role,
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

func (s *AuthServiceHandler) GetUserByID(ctx context.Context, req *authpb.GetUserByIDRequest) (*authpb.GetUserByIDResponse, error) {
	objectID, err := primitive.ObjectIDFromHex(req.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	user, err := s.authService.GetUserByID(ctx, objectID)
	if err != nil {
		return nil, status.Errorf(codes.NotFound, "user not found")
	}

	return &authpb.GetUserByIDResponse{
		Id:    user.ID.Hex(),
		Name:  user.Name,
		Email: user.Email,
		Role:  user.Role,
	}, nil
}

func (s *AuthServiceHandler) AddRole(ctx context.Context, req *authpb.AddRoleRequest) (*authpb.AddRoleResponse, error) {

	adminUserIDHex, ok := ctx.Value("user_id").(string)
	if !ok || adminUserIDHex == "" {
		return nil, status.Errorf(codes.Unauthenticated, "missing user id in context")
	}

	adminUserID, err := primitive.ObjectIDFromHex(adminUserIDHex)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid admin user id")
	}

	targetUserID, err := primitive.ObjectIDFromHex(req.TargetUserId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid target user id")
	}

	// เรียก service logic ตรวจสอบ admin และ update role
	err = s.authService.AddRole(ctx, adminUserID, targetUserID, req.NewRole)
	if err != nil {
		if err.Error() == "forbidden: only admin can update role" {
			return nil, status.Errorf(codes.PermissionDenied, err.Error())
		}
		return nil, status.Errorf(codes.Internal, err.Error())
	}

	return &authpb.AddRoleResponse{Success: true}, nil
}
