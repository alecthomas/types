// Package pubsub provides a simple publish/subscribe mechanism.
//
// It supports both synchronous and asynchronous subscriptions.
package pubsub

import (
	"errors"
	"fmt"
	"sync"
	"time"
)

// AckTimeout is the time to wait for an ack before panicking.
//
// This is a last-ditch effort to avoid deadlocks.
const AckTimeout = time.Second * 30

// Message is a message that must be acknowledge by the receiver.
type Message[T any] struct {
	Msg T
	ack chan error
}

func (a *Message[T]) Ack() { close(a.ack) }
func (a *Message[T]) Nack(err error) {
	if err == nil {
		err = errors.New("nack")
	}
	a.ack <- err
	close(a.ack)
}

// Control messages for the topic.
type control[T any] interface{ control() }

type subscribe[T any] chan Message[T]

func (subscribe[T]) control() {}

type unsubscribe[T any] chan Message[T]

func (unsubscribe[T]) control() {}

type stop struct{}

func (stop) control() {}

type Topic[T any] struct {
	// This map is used by Unsubscribe() because the non-ackable channel is not
	// the same as the ackable channel.
	//
	// If this were typed it would be map[chan T]chan Message[T]
	rawChannelMap sync.Map
	publish       chan Message[T]
	control       chan control[T]
	// Closed when the Topic is closed.
	close chan struct{}
}

// New creates a new topic that can be used to publish and subscribe to messages.
func New[T any]() *Topic[T] {
	s := &Topic[T]{
		publish: make(chan Message[T], 16384),
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

// Publish a message to the topic.
func (s *Topic[T]) Publish(t T) {
	s.publish <- Message[T]{Msg: t, ack: make(chan error, 1)}
}

// PublishSync publishes a message to the topic and blocks until all subscriber
// channels have acked the message.
func (s *Topic[T]) PublishSync(t T) error {
	ack := make(chan error, 1)
	s.publish <- Message[T]{Msg: t, ack: ack}
	timer := time.NewTimer(AckTimeout)
	defer timer.Stop()
	select {
	case err := <-ack:
		return err
	case <-timer.C:
		return nil
	}
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
	forward := make(chan Message[T], cap(c))
	go func() {
		for msg := range forward {
			c <- msg.Msg
			msg.Ack()
		}
		close(c)
	}()
	s.rawChannelMap.Store(c, forward)
	s.control <- subscribe[T](forward)
	return c
}

// SubscribeSync creates a synchronous subscription to the topic.
//
// Each message must be acked by the subscriber.
//
// A synchronous publish will block until the message has been acked by
// all subscribers.
//
// The channel will be closed when the topic is closed.
// If "c" is nil a new channel of size 16 will be created.
func (s *Topic[T]) SubscribeSync(c chan Message[T]) chan Message[T] {
	if c == nil {
		c = make(chan Message[T], 16)
	}
	s.control <- subscribe[T](c)
	return c
}

// Unsubscribe a channel from the topic, closing the channel.
func (s *Topic[T]) Unsubscribe(c chan T) {
	ackch, ok := s.rawChannelMap.Load(c)
	if !ok { // This should never happen in practice.
		panic("channel not subscribed")
	}
	s.rawChannelMap.Delete(c)
	s.control <- unsubscribe[T](ackch.(chan Message[T]))
}

// UnsubscribeSync a synchronised subscription from the topic, closing the channel.
func (s *Topic[T]) UnsubscribeSync(c chan Message[T]) {
	s.control <- unsubscribe[T](c)
}

// Close the topic, blocking until all subscribers have been closed.
func (s *Topic[T]) Close() error {
	s.control <- stop{}
	<-s.close
	return nil
}

func (s *Topic[T]) run() {
	subscriptions := map[chan Message[T]]struct{}{}
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
				s.rawChannelMap.Range(func(k, v interface{}) bool {
					s.rawChannelMap.Delete(k)
					return true
				})
				close(s.close)
				return

			default:
				panic(fmt.Sprintf("unknown control message: %T", msg))
			}

		case msg := <-s.publish:
			errs := []error{}
			for ch := range subscriptions {
				smsg := Message[T]{Msg: msg.Msg, ack: make(chan error, 1)}
				ch <- smsg
				timer := time.NewTimer(AckTimeout)
				select {
				case err := <-smsg.ack:
					errs = append(errs, err)
				case <-timer.C:
					panic("ack timeout")
				}
				timer.Stop()
			}
			msg.ack <- errors.Join(errs...)
			close(msg.ack)
		}
	}
}
