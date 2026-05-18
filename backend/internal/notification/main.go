package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()
	InitNotificationGrpcServer(grpcServer)
	log.Fatal(grpcServer.Serve(lis))
}
