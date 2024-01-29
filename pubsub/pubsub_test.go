package pubsub

import (
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"
)

func Example() {
	// Create a new topic.
	t := New[int]()

	// Subscribe to changes.
	changes := t.Subscribe(nil)
	go func() {
		for change := range changes {
			fmt.Println("change:", change)
		}
	}()

	// Publish a value.
	t.Publish(1)

	// Publish a value and wait for it to be received.
	t.Publish(2)

	time.Sleep(time.Millisecond * 100)
	// Output:
	// change: 1
	// change: 2
}

func TestPubsub(t *testing.T) {
	pubsub := New[string]()
	ch := make(chan string, 64)
	pubsub.Subscribe(ch)
	pubsub.Publish("hello")
	select {
	case msg := <-ch:
		assert.Equal(t, "hello", msg)

	case <-time.After(time.Millisecond * 100):
		t.Fail()
	}
	_ = pubsub.Close()
	assert.Panics(t, func() { pubsub.Subscribe(ch) })
	assert.Panics(t, func() { pubsub.Unsubscribe(ch) })
	assert.Panics(t, func() { pubsub.Publish("hello") })
}
