// Package eventsource provides a pubsub.Topic that also stores the last published value in an atomic.Value.
//
// Updating the value will result in a publish event.
package eventsource

import (
	"sync/atomic"

	"github.com/alecthomas/types/pubsub"
)

// EventSource is a pubsub.Topic that also stores the last published value in an atomic.Value.
//
// Updating the value will result in a publish event.
type EventSource[T any] struct {
	*pubsub.Topic[T]
	value atomic.Value
}

func New[T any]() *EventSource[T] {
	var t T
	e := &EventSource[T]{Topic: pubsub.New[T]()}
	e.value.Store(t)
	changes := make(chan pubsub.Message[T], 64)
	e.SubscribeSync(changes)
	go func() {
		for msg := range changes {
			e.value.Store(msg.Msg)
			msg.Ack()
		}
	}()
	return e
}

func (e *EventSource[T]) Store(value T) {
	e.value.Store(value)
	e.PublishSync(value)
}

func (e *EventSource[T]) Load() T {
	return e.value.Load().(T)
}

func (e *EventSource[T]) Swap(value T) T {
	rv := e.value.Swap(value)
	e.PublishSync(value)
	return rv.(T)
}

func (e *EventSource[T]) CompareAndSwap(old, new T) bool { //nolint:predeclared
	if e.value.CompareAndSwap(old, new) {
		e.PublishSync(new)
		return true
	}
	return false
}
