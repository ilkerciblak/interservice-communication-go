package main

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"time"
)

type OrderService struct {
	Repository OrderRepositoryInterface
}

func (s *OrderService) PlaceOrder(ctx context.Context, customer, item string, quantity int) (*Order, error) {

	order := Order{
		ID:       fmt.Sprintf("order-%s-%d", rand.Text(), time.Now().Unix()),
		Customer: customer,
		Item:     item,
		Quantity: quantity,
	}

	if err := s.Repository.CreateOrder(ctx, order); err != nil {
		return nil, err
	}

	log.Printf("[order] placed %s for %s (%dx %s)\n", order.ID, order.Customer, order.Quantity, order.Item)

	return &order, nil
}
