package main

import (
	"context"
	"encoding/json"
	"ilkerciblak/order-management/shared/config"
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
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
	defer stop()
	eventBus, err := messaging.RegisterRabbitMQ()
	if err != nil {
		log.Fatalf("failed to initialize rabbitMQ: %v", err)
		return
	}
	defer eventBus.Close(ctx)

	cfg := config.Config()

	repo := inventoryRepository{DataDir: cfg.DataDir}
	service := inventoryService{Repo: &repo, Publisher: eventBus}

	if err = eventBus.Subscribe(ctx, eventbus.OrderPlaced, func(ctx context.Context, e messaging.Event) error {
		var payload eventbus.OrderPlacedPayload

		if err = json.Unmarshal(e.Payload, &payload); err != nil {
			log.Printf("[inventory] failed to decode event %s payload: %v", e.Name, e.Payload)
			return err
		}

		if err = service.ReserveProduct(ctx, payload.Item, payload.OrderID); err != nil {
			log.Printf("[inventory] -- Reservation Rejected for %s", payload.OrderID)
			eventBus.Publish(ctx, eventbus.NewEvent(eventbus.StockNotReserved, eventbus.StockNotReservedPayload{
				OrderID: payload.OrderID,
				Message: err.Error(),
			},
			))
		}

		return nil
	}); err != nil {
		log.Fatal(err)
	}

	log.Printf("[inventory] subscribed to `order.placed` publishes `stock.reserved`")

	grpcServer := grpc.NewServer()
	initInventoryServer(grpcServer, &service)

	ls, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to start inventory server: %v", err)
	}

	go func() {
		if err := eventBus.Start(ctx); err != nil {
			log.Fatal(err)
		}
	}()

	go func() {
		if err := grpcServer.Serve(ls); err != nil {
			log.Fatal(err)
		}
	}()

	select {
	case <-ctx.Done():
		grpcServer.GracefulStop()
		os.Exit(1)
	}

}
