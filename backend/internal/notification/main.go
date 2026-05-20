package main

import (
	"context"
	"encoding/json"
	"fmt"
	eventbus "ilkerciblak/order-management/internal/event-bus"
	"ilkerciblak/order-management/shared/messaging"
	"log"
	"math/rand/v2"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	eventBus, err := messaging.RegisterRabbitMQ()
	if err != nil {
		log.Fatalf("failed to start eventBus in notificaiton.main: %v", err)
	}

	defer eventBus.Close(ctx)

	if err := eventBus.Subscribe(ctx, eventbus.StockReserved, func(ctx context.Context, e messaging.Event) error {

		var payload eventbus.StockReservedPayload

		if err := json.Unmarshal(e.Payload, &payload); err != nil {
			return fmt.Errorf("failed to parse event %s: %w", e.Name, err)
		}

		if payload.Reserved {
			userConfirmation := rand.IntN(2) == 1
			if userConfirmation {
				_ = eventBus.Publish(ctx, eventbus.NewEvent(eventbus.OrderConfirmed, eventbus.OrderConfirmedPayload{OrderID: payload.OrderID}))
				return nil
			}
			_ = eventBus.Publish(ctx, eventbus.NewEvent(eventbus.OrderCancelled, eventbus.OrderCancelledPayload{OrderID: payload.OrderID, Message: "user cancelled"}))
		}
		_ = eventBus.Publish(ctx, eventbus.NewEvent(eventbus.OrderCancelled, eventbus.OrderCancelledPayload{OrderID: payload.OrderID, Message: "stock not reserved"}))

		return nil

	}); err != nil {
		log.Fatal(err)
	}

	InitNotificationGrpcServer(grpcServer)
	go func() {
		if err := eventBus.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("notification service consuming order.placed")
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		os.Exit(1)
	}
}
