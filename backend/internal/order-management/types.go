package main

import (
	"context"
	orderpb "ilkerciblak/order-management/shared/proto/order"
)

type Order struct {
	ID       string
	Customer string
	Item     string
	Quantity int
	Status   string
}

type OrderServiceInterface interface {
	PlaceOrder(ctx context.Context, customer, item string, quantity int) (*orderpb.PlaceOrderResponse, error)
	RejectOrder(ctx context.Context, orderID string) error
	ConfirmOrder(ctx context.Context, orderID string) error
	CancelOrder(ctx context.Context, orderID string) error
}

type OrderRepositoryInterface interface {
	CreateOrder(ctx context.Context, order Order) error
	UpdateOrder(ctx context.Context, orderID string, updated Order) error
	GetOrder(ctx context.Context, orderID string) (*Order, error)
}
