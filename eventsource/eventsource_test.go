package eventsource

import (
	"fmt"
	"time"
)

func Example() {

	// Create a new event source.
	e := New[int]()

	// Subscribe to changes.
	changes := e.Subscribe(nil)
	go func() {
		for change := range changes {
			fmt.Println("change:", change)
		}
	}()

	// Publish a set a value.
	e.Publish(1)

	// Set and publish a value.
	e.Store(2)

	time.Sleep(time.Millisecond * 100)

	fmt.Println(e.Load())

	// Output:
	// change: 1
	// change: 2
	// 2
}
