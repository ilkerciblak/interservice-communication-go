package main

import (
	"context"
	"encoding/json"
	eventbus "ilkerciblak/order-management/internal/event-bus"
	"ilkerciblak/order-management/shared/messaging"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	defer stop()

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("failed to start order service: %v", err)
	}

	// replaced with eventBus subscriptions

	// conn, err := grpc.NewClient("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatalf("%v", err)
	// }
	// inventoryClient := inventorypb.NewInventoryServiceClient(conn)
	//
	// notificationConn, err := grpc.NewClient("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//
	// notificationClient := notificationpb.NewNotificationServiceClient(notificationConn)
	//

	grpcServer := grpc.NewServer()

	rabbit, err := messaging.RegisterRabbitMQ()
	if err != nil {
		log.Fatal(err)
	}
	defer rabbit.Close(context.Background())

	orderRepository := OrderRepository{}
	orderService := OrderService{
		Repository: &orderRepository,
		Publisher:  rabbit,
	}

	if err := rabbit.Subscribe(ctx, eventbus.StockReserved, func(ctx context.Context, e messaging.Event) error {
		var payload eventbus.StockReservedPayload

		if err := json.Unmarshal(e.Payload, &payload); err != nil {
			return err
		}

		if payload.Reserved {
			return orderService.ConfirmOrder(ctx, payload.OrderID)
		}

		return orderService.RejectOrder(ctx, payload.OrderID)
	}); err != nil {
		log.Fatal(err)
	}

	OrderServer(grpcServer, &orderService)
	go func() {
		log.Fatal(grpcServer.Serve(lis))
	}()

	go func() {
		if err := rabbit.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		os.Exit(1)
	}
}
