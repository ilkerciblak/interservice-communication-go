package main

import "context"

type Order struct {
	ID       string
	Customer string
	Item     string
	Quantity int
}

type OrderServiceInterface interface {
	PlaceOrder(ctx context.Context, customer, item string, quantity int) (*Order, error)
}

type OrderRepositoryInterface interface {
	CreateOrder(ctx context.Context, order Order) error
	// ReadOrder(ctx context.Context, orderID string) (*Order, error)
}
