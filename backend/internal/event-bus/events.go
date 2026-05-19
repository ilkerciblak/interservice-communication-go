package eventbus

import (
	"encoding/json"
	"ilkerciblak/order-management/shared/messaging"
	"time"
)

const (
	OrderPlaced   = "order.placed"
	StockReserved = "stock.reserved"
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
