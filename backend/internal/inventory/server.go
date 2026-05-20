package main

import (
	"context"
	inventorypb "ilkerciblak/order-management/shared/proto/inventory"

	"google.golang.org/grpc"
)

type inventoryServer struct {
	inventorypb.UnimplementedInventoryServiceServer
	service InventoryServiceInterface
}

func initInventoryServer(s *grpc.Server, service InventoryServiceInterface) error {
	handler := inventoryServer{
		service: service,
	}

	inventorypb.RegisterInventoryServiceServer(s, &handler)

	return nil
}

func (i *inventoryServer) GetInventory(ctx context.Context, in *inventorypb.GetInventoryRequest) (*inventorypb.GetInventoryResponse, error) {

	inv, err := i.service.GetInventory(ctx)
	if err != nil {
		return nil, err
	}

	arr := make([]*inventorypb.Product, 0, len(inv.Products))
	for _, product := range inv.Products {
		arr = append(arr, &inventorypb.Product{
			Id:            product.ID,
			Title:         product.Title,
			Quantity:      product.Quantity,
			ReservedCount: product.ReservedCount,
		})
	}
	return &inventorypb.GetInventoryResponse{
		Products: arr,
	}, nil
}

func (i *inventoryServer) ReserveStock(ctx context.Context, in *inventorypb.ReserveStockRequest) (*inventorypb.ReserveStockResponse, error) {
	if err := i.service.ReserveProduct(ctx, in.Item, in.OrderId); err != nil {
		return &inventorypb.ReserveStockResponse{Reserved: false}, err
	}

	return &inventorypb.ReserveStockResponse{Reserved: true}, nil
}

func InventoryServer(s *grpc.Server, service InventoryServiceInterface) *inventoryServer {
	handler := &inventoryServer{
		service: service,
	}

	inventorypb.RegisterInventoryServiceServer(s, handler)

	return handler
}
