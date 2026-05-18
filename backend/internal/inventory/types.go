package main

import "context"

type Inventory struct {
	Products []Product `json:"products"`
}

type Product struct {
	ID            string
	Title         string
	Quantity      int32
	ReservedCount int32 `json:"reserved_count"`
}

type InventoryServiceInterface interface {
	GetInventory(ctx context.Context) (*Inventory, error)
	ReserveProduct(ctx context.Context, productID string) error
}

type InventoryRepositoryInterface interface {
	GetInventory(ctx context.Context) (*Inventory, error)
	ReserveProduct(ctx context.Context, productID string) error
}
