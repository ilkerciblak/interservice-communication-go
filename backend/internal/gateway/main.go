package main

import (
	"encoding/json"
	"fmt"
	orderpb "ilkerciblak/order-management/shared/proto/order"
	"log"
	"net/http"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.NewClient(
		"localhost:50051",
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalf("failed to start gateway: %v", err)
	}
	defer conn.Close()

	g := gateway{
		orderClient: orderpb.NewOrderServiceClient(conn),
	}
	
	

	mux := http.NewServeMux()

	server := &http.Server{
		Addr:    fmt.Sprintf(":8080"),
		Handler: mux,
	}

	mux.HandleFunc("POST /orders", g.handlePlaceOrder)

	log.Fatal(server.ListenAndServe())
}

type gateway struct {
	orderClient orderpb.OrderServiceClient
}

type placeOrderRequest struct {
	Customer string `json:"customer"`
	Item     string `json:"item"`
	Quantity int    `json:"quantity"`
}

func (g *gateway) handlePlaceOrder(w http.ResponseWriter, r *http.Request) {
	var req placeOrderRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "check request body", http.StatusBadRequest)
		return
	}

	resp, err := g.orderClient.PlaceOrder(
		r.Context(),
		&orderpb.PlaceOrderRequest{
			Customer: req.Customer,
			Item:     req.Item,
			Quantity: int32(req.Quantity),
		},
	)
	if err != nil {
		http.Error(w, "failed to place order: "+err.Error(), http.StatusBadGateway)
		return
	}
	w.Header().Set("content-type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"order_id": resp.OrderId,
		"status":   resp.Status,
	})

}
