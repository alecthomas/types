package eventsource

import (
	"fmt"
	"log"
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
	err := e.PublishSync(1)
	if err != nil {
		log.Fatal(err)
	}

	// Set and publish a value.
	err = e.Store(2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(e.Load())

	// Output:
	// change: 1
	// change: 2
	// 2
}
