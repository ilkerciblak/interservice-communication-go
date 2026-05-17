package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to start order service: %v", err)
	}

	grpcServer := grpc.NewServer()

	orderRepository := OrderRepository{}
	orderService := OrderService{Repository: &orderRepository}
	OrderServer(grpcServer, &orderService)
	log.Fatal(grpcServer.Serve(lis))
}
