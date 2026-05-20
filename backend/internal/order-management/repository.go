package main

import (
	"context"
	"encoding/json"
	"os"
	"slices"
)

type OrderRepository struct {
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order Order) error {

	var orders struct {
		Orders []Order `json:"orders"`
	}
	file, err := os.OpenFile("../../data/order.json", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return err
	}

	orders.Orders = append(orders.Orders, order)
	data, err := json.MarshalIndent(orders, "", "")
	if err != nil {
		return err
	}

	if _, err := file.WriteAt(data, 0); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) UpdateOrder(ctx context.Context, orderID string, updated Order) error {

	var orders struct {
		Orders []Order `json:"orders"`
	}

	file, err := os.OpenFile("../../data/order.json", os.O_RDWR|os.O_CREATE, 0o666)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return err
	}

	orderIDX := slices.IndexFunc(orders.Orders, func(o Order) bool {
		return o.ID == orderID
	})

	orders.Orders[orderIDX] = updated

	data, err := json.MarshalIndent(orders, "", "")
	if err != nil {
		return err
	}

	if _, err := file.WriteAt(data, 0); err != nil {
		return err
	}

	return nil
}

func (r *OrderRepository) GetOrder(ctx context.Context, orderID string) (*Order, error) {
	var orders struct {
		Orders []Order `json:"orders"`
	}

	file, err := os.OpenFile("../../data/order.json", os.O_RDONLY, 0o444)
	if err != nil {
		return nil, err
	}

	if err := json.NewDecoder(file).Decode(&orders); err != nil {
		return nil, err
	}

	orderIDX := slices.IndexFunc(orders.Orders, func(o Order) bool { return o.ID == orderID })

	return &orders.Orders[orderIDX], nil
}
