package pubsub_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/alecthomas/assert/v2"

	. "github.com/alecthomas/types/pubsub" //nolint
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
	select {
	case _, ok := <-ch:
		assert.True(t, !ok, "channel should be closed")

	case <-time.After(time.Millisecond * 100):
		t.Fatal("channel should have been closed")
	}
	assert.Panics(t, func() { pubsub.Subscribe(ch) })
	assert.Panics(t, func() { pubsub.Unsubscribe(ch) })
	assert.Panics(t, func() { pubsub.Publish("hello") })
}

func TestSyncPubSub(t *testing.T) {
	pubsub := New[string]()
	defer pubsub.Close() //nolint
	order := make(chan string, 64)
	finished := make(chan struct{})
	ch := pubsub.SubscribeSync(nil)
	go func() {
		err := pubsub.PublishSync("hello")
		order <- "published"
		assert.NoError(t, err)

		err = pubsub.PublishSync("world")
		order <- "published"
		assert.EqualError(t, err, "nack")

		close(finished)
	}()
	select {
	case msg := <-ch:
		assert.Equal(t, "hello", msg.Msg)
		order <- "received"
		msg.Ack()

		// Test nack
		select {
		case msg := <-ch:
			assert.Equal(t, "world", msg.Msg)
			order <- "received"
			msg.Nack(errors.New("nack"))

		case <-time.After(time.Millisecond * 500):
			t.Fatal("timeout")
		}

	case <-time.After(time.Millisecond * 500):
		t.Fatal("timeout")
	}
	<-finished
	close(order)

	// Ensure that the message was received before it was published and thus
	// acked.
	actual := []string{}
	for o := range order {
		actual = append(actual, o)
	}
	assert.Equal(t, []string{"received", "published", "received", "published"}, actual)
}

func TestPubSubPanicAfterUnsubscribe(t *testing.T) {
	t.Skip("This test is slow")
	topic := New[string]()
	for range 100 {
		go func() {
			foo := topic.Subscribe(make(chan string, 1))
			<-time.After(time.Second)
			topic.Unsubscribe(foo)
		}()
	}
	go func() {
		for {
			select {
			case <-time.After(time.Millisecond * 100):
				topic.Publish("foo")
			}
		}
	}()
	<-time.After(time.Minute)
}
