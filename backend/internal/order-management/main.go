package main

import (
	"context"
	"encoding/json"
	"fmt"
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

	if err := rabbit.Subscribe(ctx, eventbus.OrderConfirmed, func(ctx context.Context, e messaging.Event) error {
		var payload eventbus.OrderConfirmedPayload

		if err := json.Unmarshal(e.Payload, &payload); err != nil {
			return fmt.Errorf("failed to parse event %s: %w", e.Name, err)
		}

		if err := orderService.ConfirmOrder(ctx, payload.OrderID); err != nil {
			return err
		}

		return nil

	}); err != nil {
		log.Fatalf("failed to subscribe event %s: %v", eventbus.OrderConfirmed, err)
	}

	if err := rabbit.Subscribe(ctx, eventbus.OrderCancelled, func(ctx context.Context, e messaging.Event) error {
		var payload eventbus.OrderCancelledPayload

		if err := json.Unmarshal(e.Payload, &payload); err != nil {
			return fmt.Errorf("failed to parse event %s: %w", e.Name, err)
		}

		if err := orderService.CancelOrder(ctx, payload.OrderID); err != nil {
			return err
		}

		return nil

	}); err != nil {
		log.Fatalf("failed to subscribe event %s: %v", eventbus.OrderCancelled, err)
	}

	if err := rabbit.Subscribe(ctx, eventbus.StockNotReserved, func(ctx context.Context, e messaging.Event) error {
		var payload eventbus.StockNotReservedPayload

		if err := json.Unmarshal(e.Payload, &payload); err != nil {
			return fmt.Errorf("failed to parse event %s: %w", e.Name, err)
		}

		if err := orderService.RejectOrder(ctx, payload.OrderID); err != nil {
			return err
		}

		return nil

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
