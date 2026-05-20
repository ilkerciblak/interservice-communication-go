// Package eventbus exposes typed Events implements messaging.Event type
package eventbus

import (
	"encoding/json"
	"ilkerciblak/order-management/shared/messaging"
	"time"
)

const (
	OrderPlaced      = "order.placed"
	OrderCancelled   = "order.cancelled"
	OrderConfirmed   = "order.confirmed"
	StockReserved    = "stock.reserved"
	StockNotReserved = "stock.not.reserved"
)

type OrderPlacedPayload struct {
	OrderID  string `json:"order_id"`
	Customer string `json:"customer"`
	Item     string `json:"item"`
	Quantity int32  `json:"quantity"`
}

type StockReservedPayload struct {
	OrderID  string `json:"order_id"`
	Item     string `json:"item"`
	Quantity int32  `json:"quantity"`
	Reserved bool   `json:"reserved"`
	Message  string `json:"message"`
}

type OrderConfirmedPayload struct {
	OrderID string `json:"order_id"`
	Message string `json:"message"`
}

type OrderCancelledPayload struct {
	OrderID string `json:"order_id"`
	Message string `json:"message"`
}

type StockNotReservedPayload struct {
	OrderID string `json:"order_id"`
	Message string `json:"message"`
}

func OrderPlacedEvent(payload OrderPlacedPayload) messaging.Event {
	name := OrderPlaced
	timeStamp := time.Now()
	payloadByte := make([]byte, 0)
	data, err := json.Marshal(payload)
	if err == nil {
		payloadByte = append(payloadByte, data...)
	}

	return messaging.Event{
		Name:      name,
		TimeStamp: timeStamp,
		Payload:   payloadByte,
	}

}

func StockReservedEvent(payload StockReservedPayload) messaging.Event {
	name := StockReserved
	timeStamp := time.Now()
	payloadByte := make([]byte, 0)

	data, err := json.Marshal(payload)
	if err == nil {
		payloadByte = append(payloadByte, data...)
	}

	return messaging.Event{
		Name:      name,
		TimeStamp: timeStamp,
		Payload:   payloadByte,
	}
}

func NewEvent(eventName string, payload any) messaging.Event {
	timeStamp, payloadBytes := time.Now(), make([]byte, 0)
	data, err := json.Marshal(payload)
	if err == nil {
		payloadBytes = append(payloadBytes, data...)
	}

	return messaging.Event{
		Name:      eventName,
		Payload:   payloadBytes,
		TimeStamp: timeStamp,
	}
}
