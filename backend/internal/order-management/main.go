package main

import (
	"context"
	"ilkerciblak/order-management/shared/messaging"
	inventorypb "ilkerciblak/order-management/shared/proto/inventory"
	notificationpb "ilkerciblak/order-management/shared/proto/notification"
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

	notificationConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	notificationClient := notificationpb.NewNotificationServiceClient(notificationConn)

	grpcServer := grpc.NewServer()

	rabbit, err := messaging.RegisterRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}

	defer rabbit.Close(context.Background())

	orderRepository := OrderRepository{}
	orderService := OrderService{
		Repository:         &orderRepository,
		inventoryClient:    inventoryClient,
		notificationClient: notificationClient,
		Publisher:          rabbit,
	}
	OrderServer(grpcServer, &orderService)

	log.Fatal(grpcServer.Serve(lis))
}
