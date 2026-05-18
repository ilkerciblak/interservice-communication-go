package messaging

import (
	"context"
	"time"
)

// Event contacts the unit of event driven communication
type Event struct {
	Name      string
	Payload   []byte
	TimeStamp time.Time
}

// EventHandler typed function handling given single event, returns error
//
// EventHandler function enables the event-bus implementation decides on retry/dead-lettering later.
type EventHandler func(context.Context, Event) error

// Publisher interface contracts event publisher service party
type Publisher interface {
	Publish(context.Context, Event) error
}

// Subscriber interface contracts event consumer service party
type Subscriber interface {
	Subscribe(context.Context, string, EventHandler) error
}

// EventBus contract  embeds both Publisher and Subscriber. A caller can still depends on just single interface or both.
//
// EventBus interface also contracts lifecycle methods as `Start` and `Close` for an event-bus vendor implementation.
type EventBus interface {
	Publisher
	Subscriber
	Start(context.Context) error
	Close(context.Context) error
}
