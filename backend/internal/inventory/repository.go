package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"slices"
)

// type InventoryRepositoryInterface interface {
// 	GetInventory(ctx context.Context) (*Inventory, error)
// 	ReserveProduct(ctx context.Context, productID string) error
// }

type inventoryRepository struct {
}

func (r *inventoryRepository) GetInventory(ctx context.Context) (*Inventory, error) {

	file, err := os.OpenFile("../../data/inventory_2.json", os.O_RDONLY, 0o666)
	if err != nil {
		return nil, fmt.Errorf("failed to read inventory.json: %w", err)
	}

	defer file.Close()

	var inventory Inventory

	if err := json.NewDecoder(file).Decode(&inventory); err != nil {
		return nil, fmt.Errorf("failed to decode inventory data: %w", err)
	}

	return &inventory, nil
}

func (r *inventoryRepository) ReserveProduct(ctx context.Context, productID string) error {
	file, err := os.OpenFile("../../data/inventory_2.json", os.O_CREATE|os.O_RDWR, 0o666)
	if err != nil {
		return fmt.Errorf("failed to read inventory.json: %w", err)
	}
	defer file.Close()

	var inventory Inventory

	if err := json.NewDecoder(file).Decode(&inventory); err != nil {
		return fmt.Errorf("failed to decode inventory data: %w", err)
	}

	// find the item

	productIdx := slices.IndexFunc(inventory.Products, func(p Product) bool {
		return p.ID == productID
	})

	product := inventory.Products[productIdx]
	fmt.Fprintf(os.Stdout, "product: %s - reserved: %d", product.Title, product.ReservedCount)
	if product.ReservedCount >= product.Quantity {
		return fmt.Errorf("in-sufficient stock")
	}
	product.ReservedCount++

	inventory.Products[productIdx] = product

	data, err := json.MarshalIndent(inventory, "", "")
	if err != nil {
		return err
	}

	if _, err := file.WriteAt(data, 0); err != nil {
		return fmt.Errorf("error while updating inventory: %w", err)
	}

	return nil
}
