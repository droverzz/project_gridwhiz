package handler

import (
	"auth-microservice/internal/service"
	"context"
	"log"
	"strings"

	authpb "auth-microservice/proto"

	"auth-microservice/internal/model"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type AuthServiceHandler struct {
	authpb.UnimplementedAuthServiceServer
	authService service.AuthService
}

func NewAuthServiceHandler(authService service.AuthService) *AuthServiceHandler {
	return &AuthServiceHandler{authService: authService}
}

func grpcErrorFromService(err error) error {
	switch err {
	case nil:
		return nil
	case service.ErrInvalidCredentials, service.ErrUnauthenticated:
		return status.Error(codes.Unauthenticated, err.Error())
	case service.ErrForbidden:
		return status.Error(codes.PermissionDenied, err.Error())
	case service.ErrUserExists:
		return status.Error(codes.AlreadyExists, err.Error())
	case service.ErrNotFound:
		return status.Error(codes.NotFound, err.Error())
	case service.ErrInvalidArgument:
		return status.Error(codes.InvalidArgument, err.Error())
	default:
		return status.Error(codes.Internal, err.Error())
	}
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
		return nil, grpcErrorFromService(err)
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
		return nil, grpcErrorFromService(err)
	}
	log.Printf("Login success")
	return &authpb.LoginResponse{
		Token: token,
	}, nil
}

func (s *AuthServiceHandler) Logout(ctx context.Context, req *authpb.LogoutRequest) (*authpb.LogoutResponse, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "missing metadata")
	}

	authHeaders := md["authorization"]
	if len(authHeaders) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token not provided")
	}

	tokenString := authHeaders[0]
	const bearerPrefix = "Bearer "
	if strings.HasPrefix(tokenString, bearerPrefix) {
		tokenString = tokenString[len(bearerPrefix):]
	} else {
		return nil, status.Errorf(codes.Unauthenticated, "invalid authorization header format")
	}

	err := s.authService.Logout(ctx, tokenString)
	if err != nil {
		log.Printf("Logout failed: %v", err)
		return nil, grpcErrorFromService(err)
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
		return nil, grpcErrorFromService(err)
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

	err = s.authService.AddRole(ctx, adminUserID, targetUserID, req.NewRole)
	if err != nil {
		return nil, grpcErrorFromService(err)
	}

	return &authpb.AddRoleResponse{Success: true}, nil
}

func (s *AuthServiceHandler) ListUsers(ctx context.Context, req *authpb.ListUsersRequest) (*authpb.ListUsersResponse, error) {
	filter := &model.UserFilter{
		Name:  req.Name,
		Email: req.Email,
		Page:  req.Page,
		Limit: req.Limit,
	}

	users, total, err := s.authService.ListUsers(ctx, filter)
	if err != nil {
		return nil, grpcErrorFromService(err)
	}

	var userProtos []*authpb.User
	for _, u := range users {
		userProtos = append(userProtos, &authpb.User{
			Id:    u.ID.Hex(),
			Name:  u.Name,
			Email: u.Email,
			Role:  u.Role,
		})
	}

	return &authpb.ListUsersResponse{
		Users: userProtos,
		Total: total,
	}, nil
}

func (s *AuthServiceHandler) UpdateProfile(ctx context.Context, req *authpb.UpdateProfileRequest) (*authpb.UpdateProfileResponse, error) {
	userIDHex, ok := ctx.Value("user_id").(string)
	if !ok || userIDHex == "" {
		return nil, status.Errorf(codes.Unauthenticated, "missing user id in context")
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	err = s.authService.UpdateProfile(ctx, userID, req.Name, req.Email)
	if err != nil {
		return nil, grpcErrorFromService(err)
	}

	return &authpb.UpdateProfileResponse{Success: true}, nil
}

func (s *AuthServiceHandler) DeleteProfile(ctx context.Context, req *authpb.DeleteProfileRequest) (*authpb.DeleteProfileResponse, error) {
	userIDHex, ok := ctx.Value("user_id").(string)
	if !ok || userIDHex == "" {
		return nil, status.Errorf(codes.Unauthenticated, "missing user id in context")
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	err = s.authService.DeleteProfile(ctx, userID)
	if err != nil {
		return nil, grpcErrorFromService(err)
	}

	return &authpb.DeleteProfileResponse{Success: true}, nil
}

func (s *AuthServiceHandler) GeneratePasswordResetToken(ctx context.Context, req *authpb.GeneratePasswordResetTokenRequest) (*authpb.GeneratePasswordResetTokenResponse, error) {
	userIDHex, ok := ctx.Value("user_id").(string)
	if !ok || userIDHex == "" {
		return nil, status.Errorf(codes.Unauthenticated, "missing user id in context")
	}

	userID, err := primitive.ObjectIDFromHex(userIDHex)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user id")
	}

	token, err := s.authService.GeneratePasswordResetToken(ctx, userID)
	if err != nil {
		return nil, grpcErrorFromService(err)
	}

	return &authpb.GeneratePasswordResetTokenResponse{
		ResetToken: token,
	}, nil
}

func (s *AuthServiceHandler) ResetPassword(ctx context.Context, req *authpb.ResetPasswordRequest) (*authpb.ResetPasswordResponse, error) {
	err := s.authService.ResetPassword(ctx, req.ResetToken, req.NewPassword)
	if err != nil {
		return nil, grpcErrorFromService(err)
	}

	return &authpb.ResetPasswordResponse{
		Success: true,
	}, nil
}
