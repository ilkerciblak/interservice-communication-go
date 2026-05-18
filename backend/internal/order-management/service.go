package main

import (
	"context"
	"crypto/rand"
	"fmt"
	inventorypb "ilkerciblak/order-management/shared/proto/inventory"
	notificationpb "ilkerciblak/order-management/shared/proto/notification"
	"log"
	"time"
)

type OrderService struct {
	Repository         OrderRepositoryInterface
	inventoryClient    inventorypb.InventoryServiceClient
	notificationClient notificationpb.NotificationServiceClient
}

func (s *OrderService) PlaceOrder(ctx context.Context, customer, item string, quantity int) (*Order, error) {

	order := Order{
		ID:       fmt.Sprintf("order-%s-%d", rand.Text(), time.Now().Unix()),
		Customer: customer,
		Item:     item,
		Quantity: quantity,
	}

	invresp, err := s.inventoryClient.ReserveStock(ctx, &inventorypb.ReserveStockRequest{
		OrderId: order.ID,
		Item:    order.Item,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to reserve stock: %v", err)
	}

	if !invresp.Reserved {
		return &order, fmt.Errorf("no sufficient stock")
	}

	if err := s.Repository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	if _, err := s.notificationClient.SendConfirmation(ctx, &notificationpb.SendConfirmationRequest{
		Customer: order.Customer,
		OrderId:  order.ID,
	}); err != nil {
		return nil, fmt.Errorf("notification failed: %w", err)
	}

	log.Printf("[order] placed %s for %s (%dx %s)\n", order.ID, order.Customer, order.Quantity, order.Item)

	return &order, nil
}
