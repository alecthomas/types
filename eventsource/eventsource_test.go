package eventsource

import (
	"fmt"
)

func Example() {
	// Create a new event source.
	e := New[int]()

	// Subscribe to changes.
	changes := e.SubscribeSync(nil)
	go func() {
		for change := range changes {
			fmt.Println("change:", change.Msg)
			change.Ack()
		}
	}()

	// Publish a set a value.
	e.PublishSync(1)

	// Set and publish a value.
	e.Store(2)

	fmt.Println(e.Load())

	// Output:
	// change: 1
	// change: 2
	// 2
}
