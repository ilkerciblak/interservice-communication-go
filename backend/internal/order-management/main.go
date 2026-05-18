package main

import (
	inventorypb "ilkerciblak/order-management/shared/proto/inventory"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to start order service: %v", err)
	}

	conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("%v", err)
	}
	inventoryClient := inventorypb.NewInventoryServiceClient(conn)

	grpcServer := grpc.NewServer()

	orderRepository := OrderRepository{}
	orderService := OrderService{Repository: &orderRepository, inventoryClient: inventoryClient}
	OrderServer(grpcServer, &orderService)

	log.Fatal(grpcServer.Serve(lis))
}
