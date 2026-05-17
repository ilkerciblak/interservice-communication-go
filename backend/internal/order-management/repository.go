package main

import (
	"context"
	"encoding/json"
	"os"
)

type OrderRepository struct {
}

func (r *OrderRepository) CreateOrder(ctx context.Context, order Order) error {

	data, err := json.MarshalIndent(order, "", "")
	if err != nil {
		return err
	}
	

	file, err := os.OpenFile("../../data/order.json", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o666)
	if err != nil {
		return err
	}

	if _, err := file.Write(data); err != nil {
		return err
	}

	// if err := os.WriteFile("/data/order.json", data, 0644); err != nil {
	// 	return fmt.Errorf("failed to write to file: %w", err)
	// }

	return nil
}
