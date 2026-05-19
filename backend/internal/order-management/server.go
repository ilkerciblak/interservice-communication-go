package main

import (
	"context"
	orderpb "ilkerciblak/order-management/shared/proto/order"

	"google.golang.org/grpc"
)

type orderServer struct {
	orderpb.UnimplementedOrderServiceServer
	service OrderServiceInterface
}

func OrderServer(s *grpc.Server, orderService OrderServiceInterface) *orderServer {
	grpcServerHandler := orderServer{
		service: orderService,
	}

	orderpb.RegisterOrderServiceServer(s, &grpcServerHandler)

	return &grpcServerHandler
}

func (s *orderServer) PlaceOrder(ctx context.Context, req *orderpb.PlaceOrderRequest) (*orderpb.PlaceOrderResponse, error) {
	resp, err := s.service.PlaceOrder(ctx, req.Customer, req.Item, int(req.Quantity))
	if err != nil {
		return nil, err
	}

	return resp, nil

	// o, err := s.service.PlaceOrder(ctx, req.Customer, req.Item, int(req.Quantity))
	// if err != nil {
	// 	if o != nil {
	// 		return &orderpb.PlaceOrderResponse{OrderId: o.ID, Status: "rejected"}, nil
	// 	}
	// 	return nil, err
	// }
	// return &orderpb.PlaceOrderResponse{OrderId: o.ID, Status: "placed"}, nil
}
