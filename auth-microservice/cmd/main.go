package main

import (
	"log"
	"net"

	"auth-microservice/internal/handler"
	authpb "auth-microservice/proto"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()

	authService := handler.NewAuthServiceHandler()
	authpb.RegisterAuthServiceServer(grpcServer, authService)

	log.Println("gRPC server is running on port 50051...")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
