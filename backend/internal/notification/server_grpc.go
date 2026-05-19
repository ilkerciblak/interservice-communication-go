package main

import (
	"context"
	"fmt"
	notificationpb "ilkerciblak/order-management/shared/proto/notification"
	"log"

	"google.golang.org/grpc"
)

type notificationServer struct {
	 notificationpb.UnimplementedNotificationServiceServer
}

func InitNotificationGrpcServer(s *grpc.Server) {
	notificationpb.RegisterNotificationServiceServer(
		s,
		&notificationServer{},
	)
}

func (n *notificationServer) SendConfirmation(ctx context.Context, in *notificationpb.SendConfirmationRequest) (*notificationpb.SendConfirmationResponse, error) {
	log.Printf(
		"[notification] emailing %s: order %s confirmed",
		in.Customer,
		in.OrderId,
	)
	return &notificationpb.SendConfirmationResponse{
		Sent: true,

		Message: fmt.Sprintf(
			"[notification] emailing %s: order %s confirmed",
			in.Customer,
			in.OrderId,
		),
	}, nil
}
