package main

import (
	"context"
	"crypto/rand"
	"fmt"
	eventbus "ilkerciblak/order-management/internal/event-bus"
	"ilkerciblak/order-management/shared/messaging"
	orderpb "ilkerciblak/order-management/shared/proto/order"
	"log"
	"time"
)

type OrderService struct {
	Repository OrderRepositoryInterface
	messaging.Publisher
}

func (s *OrderService) PlaceOrder(ctx context.Context, customer, item string, quantity int) (*orderpb.PlaceOrderResponse, error) {
	order := Order{
		ID:       fmt.Sprintf("order-%s-%d", rand.Text(), time.Now().Unix()),
		Customer: customer,
		Item:     item,
		Quantity: quantity,
		Status:   "pending",
	}

	if err := s.Repository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	if err := s.Publisher.Publish(ctx, eventbus.OrderPlacedEvent(eventbus.OrderPlacedPayload{
		OrderID:  order.ID,
		Customer: order.Customer,
		Item:     order.Item,
		Quantity: int32(order.Quantity),
	})); err != nil {
		log.Printf("[order-service] failed to publish event %s: %v", eventbus.OrderPlaced, err)
	}

	return &orderpb.PlaceOrderResponse{
		OrderId: order.ID,
		Status:  order.Status,
	}, nil
}

func (s *OrderService) ConfirmOrder(ctx context.Context, orderID string) error {

	order, err := s.Repository.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = "confirmed"

	if err := s.Repository.UpdateOrder(ctx, orderID, *order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) CancelOrder(ctx context.Context, orderID string) error {
	order, err := s.Repository.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = "cancelled"

	if err := s.Repository.UpdateOrder(ctx, orderID, *order); err != nil {
		return err
	}

	return nil
}

func (s *OrderService) RejectOrder(ctx context.Context, orderID string) error {
	order, err := s.Repository.GetOrder(ctx, orderID)
	if err != nil {
		return err
	}

	order.Status = "rejected"

	if err := s.Repository.UpdateOrder(ctx, orderID, *order); err != nil {
		return err
	}

	return nil
}

// func (s *OrderService) PlaceOrder(ctx context.Context, customer, item string, quantity int) (*Order, error) {
//
// 	order := Order{
// 		ID:       fmt.Sprintf("order-%s-%d", rand.Text(), time.Now().Unix()),
// 		Customer: customer,
// 		Item:     item,
// 		Quantity: quantity,
// 	}
//
// 	invresp, err := s.inventoryClient.ReserveStock(ctx, &inventorypb.ReserveStockRequest{
// 		OrderId: order.ID,
// 		Item:    order.Item,
// 	})
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to reserve stock: %v", err)
// 	}
//
// 	if !invresp.Reserved {
// 		return &order, fmt.Errorf("no sufficient stock")
// 	}
//
// 	if err := s.Repository.CreateOrder(ctx, order); err != nil {
// 		return nil, err
// 	}
// 	// replaced with event publisher
// 	//
// 	// if _, err := s.notificationClient.SendConfirmation(ctx, &notificationpb.SendConfirmationRequest{
// 	// 	Customer: order.Customer,
// 	// 	OrderId:  order.ID,
// 	// }); err != nil {
// 	// 	return nil, fmt.Errorf("notification failed: %w", err)
// 	// }
//
// 	if err := s.Publish(ctx, eventbus.OrderPlacedEvent(eventbus.OrderPlacedPayload{
// 		OrderID:  order.ID,
// 		Customer: order.Customer,
// 		Item:     order.Item,
// 		Quantity: int32(order.Quantity),
// 	})); err != nil {
// 		log.Println("[order] failed to publish order.placed: %v", err)
// 	}
//
// 	log.Printf("[order] placed %s for %s (%dx %s)\n", order.ID, order.Customer, order.Quantity, order.Item)
//
// 	return &order, nil
// }
