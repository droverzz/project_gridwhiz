package middleware

import (
	"context"

	"auth-microservice/internal/utils"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func AuthInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

	skipAuthMethods := map[string]bool{
		"/auth.AuthService/Register": true,
		"/auth.AuthService/Login":    true,
	}

	if skipAuthMethods[info.FullMethod] {
		return handler(ctx, req)
	}

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
	if len(tokenString) > len(bearerPrefix) && tokenString[:len(bearerPrefix)] == bearerPrefix {
		tokenString = tokenString[len(bearerPrefix):]
	}

	userID, err := utils.ExtractUserIDFromJWT(tokenString)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid token: %v", err)
	}

	newCtx := context.WithValue(ctx, "user_id", userID)
	return handler(newCtx, req)
}
