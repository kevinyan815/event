package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/kevinyan815/event"
)

// SlowAsyncListener Simulate a slow async listener
type SlowAsyncListener struct {
	processTime time.Duration
}

func (l *SlowAsyncListener) EventHandler() event.Handler {
	return func(e *event.Event) error {
		fmt.Printf("Start processing slow async event, will take %v\n", l.processTime)
		time.Sleep(l.processTime) // simulate time-consuming operation
		fmt.Printf("Slow async event processing completed\n")
		return nil
	}
}

func (l *SlowAsyncListener) AsyncProcess() bool {
	return true
}

// FastAsyncListener Simulate a fast async listener
type FastAsyncListener struct{}

func (l *FastAsyncListener) EventHandler() event.Handler {
	return func(e *event.Event) error {
		fmt.Println("Start processing fast async event")
		time.Sleep(100 * time.Millisecond) // 模拟短暂操作
		fmt.Println("Fast async event processing completed")
		return nil
	}
}

func (l *FastAsyncListener) AsyncProcess() bool {
	return true
}

// SyncListener Simulate a sync listener
type SyncListener struct{}

func (l *SyncListener) EventHandler() event.Handler {
	return func(e *event.Event) error {
		fmt.Println("Processing sync event")
		return nil
	}
}

func (l *SyncListener) AsyncProcess() bool {
	return false
}

// DummyEvent Test event
type DummyEvent struct{}

func (e *DummyEvent) EventName() string {
	return "dummy_event"
}

func (e *DummyEvent) EntityID() string {
	return "dummy_1"
}

func (e *DummyEvent) EntityType() string {
	return "dummy"
}

func main() {
	// Register listeners
	dispatcher := event.Dispatcher()
	dummyEvent := event.NewEvent(&DummyEvent{})

	dispatcher.Subscribe(dummyEvent,
		&SlowAsyncListener{processTime: 3 * time.Second},
		&FastAsyncListener{},
		&SyncListener{},
	)

	// Send some events
	fmt.Println("Start sending events...")

	// Send first batch of events
	for i := 0; i < 3; i++ {
		ctx := context.Background()
		dispatcher.Dispatch(event.NewEventWithContext(ctx, &DummyEvent{}))
	}

	// Wait for 1 second, let some events start processing
	time.Sleep(1 * time.Second)

	fmt.Println("\nStart shutting down event dispatcher...")
	// Try graceful shutdown, set 5 seconds timeout
	shutdownTimeout := 5 * time.Second
	start := time.Now()

	if err := dispatcher.Shutdown(shutdownTimeout); err != nil {
		log.Printf("Error shutting down event dispatcher: %v\n", err)
	} else {
		fmt.Printf("Event dispatcher successfully closed, took: %v\n", time.Since(start))
	}

	// Try to send events after shutdown
	fmt.Println("\nTry to send events after shutdown...")
	dispatcher.Dispatch(event.NewEventWithContext(context.Background(), &DummyEvent{}))
}
