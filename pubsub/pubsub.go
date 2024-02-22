// Package pubsub provides a simple publish/subscribe mechanism.
package pubsub

import "fmt"

// Control messages for the topic.
type control[T any] interface{ control() }

type subscribe[T any] chan T

func (subscribe[T]) control() {}

type unsubscribe[T any] chan T

func (unsubscribe[T]) control() {}

type stop struct{}

func (stop) control() {}

type Topic[T any] struct {
	publish chan T
	control chan control[T]
	close   chan struct{}
}

// New creates a new topic that can be used to publish and subscribe to messages.
func New[T any]() *Topic[T] {
	s := &Topic[T]{
		publish: make(chan T, 64),
		control: make(chan control[T]),
		close:   make(chan struct{}),
	}
	go s.run()
	return s
}

// Wait that returns a channel that will be closed when the Topic is closed.
func (s *Topic[T]) Wait() chan struct{} {
	return s.close
}

func (s *Topic[T]) Publish(t T) {
	s.publish <- t
}

// Subscribe a channel to the topic.
//
// The channel will be closed when the topic is closed.
//
// If "c" is nil a new channel of size 16 will be created.
func (s *Topic[T]) Subscribe(c chan T) chan T {
	if c == nil {
		c = make(chan T, 16)
	}
	s.control <- subscribe[T](c)
	return c
}

// Unsubscribe a channel from the topic, closing the channel.
func (s *Topic[T]) Unsubscribe(c chan T) {
	s.control <- unsubscribe[T](c)
}

// Close the topic, blocking until all subscribers have been closed.
func (s *Topic[T]) Close() error {
	s.control <- stop{}
	<-s.close
	return nil
}

func (s *Topic[T]) run() {
	subscriptions := map[chan T]struct{}{}
	for {
		select {
		case msg := <-s.control:
			switch msg := msg.(type) {
			case subscribe[T]:
				subscriptions[msg] = struct{}{}

			case unsubscribe[T]:
				delete(subscriptions, msg)
				close(msg)

			case stop:
				for ch := range subscriptions {
					close(ch)
				}
				close(s.control)
				close(s.publish)
				close(s.close)
				return

			default:
				panic(fmt.Sprintf("unknown control message: %T", msg))
			}

		case msg := <-s.publish:
			for ch := range subscriptions {
				ch <- msg
			}
		}
	}
}
