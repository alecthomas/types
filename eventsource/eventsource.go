// Package eventsource provides a pubsub.Topic that also stores the last published value in an atomic.Value.
//
// Updating the value will result in a publish event.
//
// Example:
//
//	package main
//
//	import (
//		"fmt"
//		"time"
//
//		"github.com/alecthomas/types/eventsource"
//	)
//
//	func main() {
//		// Create a new event source.
//		e := eventsource.New[int]()
//
//		// Subscribe to changes.
//		changes := e.Subscribe(nil)
//		go func() {
//			for change := range changes {
//				fmt.Println("change:", change)
//			}
//		}()
//
//		// Publish a value.
//		e.Store(1)
//
//		// Publish a value and wait for it to be received.
//		e.Store(2)
//		time.Sleep(time.Millisecond)
//	}
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
	changes := make(chan T, 64)
	e.Subscribe(changes)
	go func() {
		for value := range changes {
			e.value.Store(value)
		}
	}()
	return e
}

func (e *EventSource[T]) Store(value T) {
	e.value.Store(value)
	e.Publish(value)
}

func (e *EventSource[T]) Load() T {
	return e.value.Load().(T)
}

func (e *EventSource[T]) Swap(value T) T {
	rv := e.value.Swap(value)
	e.Publish(value)
	return rv.(T)
}

func (e *EventSource[T]) CompareAndSwap(old, new T) bool { //nolint:predeclared
	if e.value.CompareAndSwap(old, new) {
		e.Publish(new)
		return true
	}
	return false
}
