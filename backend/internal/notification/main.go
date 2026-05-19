package main

import (
	"context"
	"ilkerciblak/order-management/shared/messaging"
	"log"
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

	// if err := eventBus.Subscribe(ctx, eventbus.OrderPlaced, func(ctx context.Context, e messaging.Event) error {
	// 	var payload eventbus.OrderPlacedPayload
	// 	if err := json.Unmarshal(e.Payload, &payload); err != nil {
	// 		log.Printf("failed to decode event (%s) payload: %v", e.Name, err)
	// 	}
	//
	// 	log.Printf("[notification]: event arrived: %s | order %s confirmed", e.Name, payload.OrderID)
	//
	// 	return nil
	// }); err != nil {
	// 	log.Fatal(err)
	// }

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
