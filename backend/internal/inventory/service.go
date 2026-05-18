package main

import "context"

type inventoryService struct {
	Repo InventoryRepositoryInterface
}

func (s *inventoryService) GetInventory(ctx context.Context) (*Inventory, error) {

	return s.Repo.GetInventory(ctx)

}
func (s *inventoryService) ReserveProduct(ctx context.Context, productID string) error {

	return s.Repo.ReserveProduct(ctx, productID)
}
