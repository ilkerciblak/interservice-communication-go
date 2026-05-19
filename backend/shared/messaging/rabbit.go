package messaging

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/rabbitmq/amqp091-go"
)

const exchangeName string = "events"

type rabbitMQEventBus struct {
	connection     *amqp091.Connection
	messageChannel *amqp091.Channel
	handlers       map[string][]EventHandler
}

func RegisterRabbitMQ() (EventBus, error) {
	url := os.Getenv("AMQP_URL")
	if strings.TrimSpace(url) == "" {
		return nil, fmt.Errorf("amqp_url is empty")
	}
	conn, err := amqp091.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to dial amqp01: %w", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("failed to open concurrent message channel: %w", err)
	}

	if err := ch.ExchangeDeclare(
		exchangeName,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	eventBus := &rabbitMQEventBus{
		connection:     conn,
		messageChannel: ch,
		handlers:       make(map[string][]EventHandler),
	}

	return eventBus,
		nil
}

func (eb *rabbitMQEventBus) Publish(ctx context.Context, event Event) error {
	if err := eb.messageChannel.PublishWithContext(
		ctx,
		exchangeName,
		event.Name,
		false,
		false,
		amqp091.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp091.Persistent,
			Body:         event.Payload,
			Timestamp:    event.TimeStamp,
		},
	); err != nil {
		return fmt.Errorf("failed to publish event %s: %w", err)
	}

	return nil

}

func (eb *rabbitMQEventBus) Subscribe(ctx context.Context, eventName string, event EventHandler) error {
	queue, err := eb.messageChannel.QueueDeclare(eventName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare queue %s: %w", eventName, err)
	}

	if err := eb.messageChannel.QueueBind(queue.Name, eventName, exchangeName, false, nil); err != nil {
		return fmt.Errorf("failed to bind queue: %w", err)
	}

	eb.handlers[eventName] = append(eb.handlers[eventName], event)

	return nil
}

func (b *rabbitMQEventBus) Start(ctx context.Context) error {
	for eventName, handlers := range b.handlers {
		msgs, err := b.messageChannel.Consume(eventName, "", true /* auto-ack for now */, false, false, false, nil)
		if err != nil {
			return fmt.Errorf("consume %s: %w", eventName, err)
		}
		handlers := handlers
		go func(name string) {
			for d := range msgs {
				for _, handler := range handlers {
					_ = handler(ctx, Event{Name: name, Payload: d.Body, TimeStamp: d.Timestamp})
				}
			}
		}(eventName)
	}
	return nil
}

func (b *rabbitMQEventBus) Close(ctx context.Context) error {
	if b.messageChannel != nil {
		_ = b.messageChannel.Close()
	}
	if b.connection != nil {
		return b.connection.Close()
	}
	return nil
}

// The var _ EventBus = (*adapterType)(nil) line is a Go idiom: it makes the compiler verify the type implements the interface, with zero runtime cost. Use it for every implementation.
var _ EventBus = (*rabbitMQEventBus)(nil)
