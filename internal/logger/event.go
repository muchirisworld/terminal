package logger

import (
	"context"
	"maps"
	"sync"
)

type eventKey struct{}

// WideEvent holds the context for a single request's wide event.
type WideEvent struct {
	mu   sync.Mutex
	data map[string]any
}

func NewWideEvent() *WideEvent {
	return &WideEvent{
		data: make(map[string]any),
	}
}

func (e *WideEvent) Add(key string, value any) {
	if e == nil {
		return
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	e.data[key] = value
}

func (e *WideEvent) GetAll() map[string]any {
	if e == nil {
		return nil
	}
	e.mu.Lock()
	defer e.mu.Unlock()
	copied := make(map[string]any, len(e.data))
	maps.Copy(copied, e.data)
	return copied
}

// WithEvent returns a new context with the given WideEvent.
func WithEvent(ctx context.Context, e *WideEvent) context.Context {
	return context.WithValue(ctx, eventKey{}, e)
}

// GetEvent retrieves the WideEvent from the context.
func GetEvent(ctx context.Context) *WideEvent {
	if e, ok := ctx.Value(eventKey{}).(*WideEvent); ok {
		return e
	}
	return nil
}

// Add adds a key-value pair to the WideEvent in the context.
func Add(ctx context.Context, key string, value any) {
	if e := GetEvent(ctx); e != nil {
		e.Add(key, value)
	}
}
