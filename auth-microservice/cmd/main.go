package main

import (
	"log"
	"net"
	"os"

	"auth-microservice/internal/db"

	"google.golang.org/grpc"
)

func main() {
	//เชื่อมต่อ MongoDB
	if err := db.InitMongoDB(os.Getenv("MONGO_URI")); err != nil {
		log.Fatalf("MongoDB connection error: %v", err)
	}

	//gRPC server
	grpcServer := grpc.NewServer()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Println("🚀 gRPC server started at :50051")
	grpcServer.Serve(lis)
}
