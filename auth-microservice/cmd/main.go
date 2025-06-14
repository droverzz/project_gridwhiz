package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"auth-microservice/internal/db"
	"auth-microservice/internal/handler"
	"auth-microservice/internal/middleware"

	"auth-microservice/internal/redis"
	"auth-microservice/internal/service"
	authpb "auth-microservice/proto"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
)

func main() {
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	err = redis.Ping(context.Background())
	if err != nil {
		log.Fatalf(" Redis not connected: %v", err)
	}
	fmt.Println("Redis connected")

	if err := db.InitMongoDB(os.Getenv("MONGO_URI")); err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(middleware.AuthInterceptor),
	)

	authService := service.NewAuthService()
	authHandler := handler.NewAuthServiceHandler(authService)

	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("🚀 gRPC server started at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
