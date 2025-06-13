package main

import (
	"log"
	"net"
	"os"

	"auth-microservice/internal/db"
	"auth-microservice/internal/handler"
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

	// เชื่อมต่อ MongoDB
	if err := db.InitMongoDB(os.Getenv("MONGO_URI")); err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	// สร้าง gRPC server
	grpcServer := grpc.NewServer()

	// สร้าง AuthService และ Handler
	authService := service.NewAuthService()
	authHandler := handler.NewAuthServiceHandler(authService)

	// ผูก AuthService กับ gRPC server
	authpb.RegisterAuthServiceServer(grpcServer, authHandler)

	// เริ่มฟังพอร์ต
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("🚀 gRPC server started at :50051")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
