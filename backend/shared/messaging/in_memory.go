package messaging

import (
	"context"
	"log"
	"sync"
)

type inMemoryBus struct {
	mu       *sync.Mutex
	handlers map[string][]EventHandler
}

func NewInMemoryBus() *inMemoryBus {
	return &inMemoryBus{handlers: make(map[string][]EventHandler)}
}

func (b *inMemoryBus) Publish(ctx context.Context, e Event) error {

	b.mu.Lock()
	defer b.mu.Unlock()

	hs, k := b.handlers[e.Name]
	if !k {
		log.Printf("no handler found for: %s", e.Name)
	}

	for _, h := range hs {
		h := h
		go func() {
			if err := h(ctx, e); err != nil {
				log.Printf("[eventbus] handler error for %s: %v", e.Name, err)
			}
		}()
	}

	return nil
}

func (b *inMemoryBus) Subscribe(ctx context.Context, eventName string, h EventHandler) error {

	b.mu.Lock()
	defer b.mu.Unlock()

	b.handlers[eventName] = append(b.handlers[eventName], h)

	return nil
}

func (b *inMemoryBus) Start(ctx context.Context) error { return nil }
func (b *inMemoryBus) Close(ctx context.Context) error { return nil }
