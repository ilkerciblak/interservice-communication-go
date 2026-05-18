package main

import (
	"log"
	"net"

	"google.golang.org/grpc"
)

func main() {

	repo := inventoryRepository{}
	service := inventoryService{Repo: &repo}

	grpcServer := grpc.NewServer()
	initInventoryServer(grpcServer, &service)

	ls, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("failed to start inventory server: %v", err)
	}

	log.Fatal(grpcServer.Serve(ls))

}
