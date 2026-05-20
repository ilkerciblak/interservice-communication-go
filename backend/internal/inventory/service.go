package main

import (
	"context"
	eventbus "ilkerciblak/order-management/internal/event-bus"
	"ilkerciblak/order-management/shared/messaging"
)

type inventoryService struct {
	Repo InventoryRepositoryInterface
	messaging.Publisher
}

func (s *inventoryService) GetInventory(ctx context.Context) (*Inventory, error) {

	return s.Repo.GetInventory(ctx)

}
func (s *inventoryService) ReserveProductRPC(ctx context.Context, productID, orderID string) error {

	return s.Repo.ReserveProduct(ctx, productID)
}
func (s *inventoryService) ReserveProduct(ctx context.Context, productID, orderID string) error {
	if err := s.Repo.ReserveProduct(ctx, productID); err != nil {
		publishErr := s.Publish(ctx, eventbus.NewEvent(
			eventbus.StockNotReserved,
			eventbus.StockNotReservedPayload{
				OrderID: orderID,
				Message: err.Error(),
			},
		))

		if publishErr != nil {
			return publishErr
		}

		return nil
	}

	if err := s.Publish(ctx, eventbus.NewEvent(
		eventbus.StockReserved,
		eventbus.StockReservedPayload{
			OrderID:  orderID,
			Item:     productID,
			Quantity: 0,
			Reserved: true,
			Message:  "stock reserved",
		},
	)); err != nil {
		return err
	}

	return nil
}
