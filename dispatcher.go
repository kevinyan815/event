package event

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"time"
)

// IEventDispatcher  interface of event dispatcher
type IEventDispatcher interface {
	// Subscribe for listener subscribe specified event
	Subscribe(event *Event, listener ...IListener)
	// SubscribeWildcard for listener subscribe all events using wildcard
	SubscribeWildcard(listener ...IListener)
	// RemoveEventListener remove event listener held by dispatcher
	RemoveEventListener(event *Event, listener IListener)
	// RemoveWildcardListener remove wildcard listener
	RemoveWildcardListener(listener IListener)
	//HasEventListener verify whether dispatcher held listener for specified event or not
	HasEventListener(event *Event) bool
	// Dispatch publish event and execute event listeners
	Dispatch(event *Event)
	// Shutdown dispatcher and wait for all event listeners to finish processing,
	// timeout is the max time to wait for all event listeners to finish processing
	Shutdown(timeout time.Duration) error
}

type eventDispatcher struct {
	m                    sync.RWMutex
	eventHolders         map[string][]IListener
	wildcardEventHolders []IListener
	wg                   sync.WaitGroup // track async events
	isShutdown           atomic.Value
}

func NewDispatcher() *eventDispatcher {
	d := new(eventDispatcher)
	d.eventHolders = make(map[string][]IListener)
	d.isShutdown.Store(false)
	return d
}

func (dispatcher *eventDispatcher) Subscribe(e *Event, l ...IListener) {
	dispatcher.m.Lock()
	defer dispatcher.m.Unlock()

	dispatcher.eventHolders[e.eventType] = append(dispatcher.eventHolders[e.eventType], l...)
}

func (dispatcher *eventDispatcher) SubscribeWildcard(l ...IListener) {
	dispatcher.m.Lock()
	defer dispatcher.m.Unlock()

	dispatcher.wildcardEventHolders = append(dispatcher.wildcardEventHolders, l...)
}

func (dispatcher *eventDispatcher) Dispatch(e *Event) {
	if dispatcher.isShutdown.Load().(bool) {
		GetLogger().Error(e.Context(), "Dispatcher is shutting down, rejecting new events", P("event", e))
		return
	}

	dispatcher.m.RLock()
	defer dispatcher.m.RUnlock()

	for _, listener := range dispatcher.eventHolders[e.eventType] {
		if listener.AsyncProcess() {
			dispatcher.wg.Add(1)
			go func(l IListener) {
				defer dispatcher.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						GetLogger().Error(e.Context(), "Panic in event listener", P("error", r), P("event", e))
					}
				}()
				if err := l.EventHandler()(e); err != nil {
					GetLogger().Error(e.Context(), "Error in event listener", P("error", err), P("event", e))
				}
			}(listener)
		} else {
			if err := listener.EventHandler()(e); err != nil {
				GetLogger().Error(e.Context(), "Error in event listener", P("error", err), P("event", e))
			}
		}
	}
	for _, listener := range dispatcher.wildcardEventHolders {
		if listener.AsyncProcess() {
			dispatcher.wg.Add(1)
			go func(l IListener) {
				defer dispatcher.wg.Done()
				defer func() {
					if r := recover(); r != nil {
						GetLogger().Error(e.Context(), "Panic in event listener", P("error", r), P("event", e))
					}
				}()
				if err := l.EventHandler()(e); err != nil {
					GetLogger().Error(e.Context(), "Error in event listener", P("error", err), P("event", e))
				}
			}(listener)
		} else {
			if err := listener.EventHandler()(e); err != nil {
				GetLogger().Error(e.Context(), "Error in event listener", P("error", err), P("event", e))
			}
		}
	}
}

func (dispatcher *eventDispatcher) HasEventListener(e *Event) bool {
	dispatcher.m.RLock()
	defer dispatcher.m.RUnlock()

	_, ok := dispatcher.eventHolders[e.eventType]
	return ok
}

func (dispatcher *eventDispatcher) RemoveEventListener(e *Event, l IListener) {
	dispatcher.m.Lock()
	defer dispatcher.m.Unlock()

	ptr := reflect.ValueOf(l).Pointer()
	listeners := dispatcher.eventHolders[e.eventType]
	for idx, listener := range listeners {
		if reflect.ValueOf(listener).Pointer() == ptr {
			dispatcher.eventHolders[e.eventType] = append(listeners[:idx], listeners[idx+1:]...)
			return
		}
	}
}

func (dispatcher *eventDispatcher) RemoveWildcardListener(l IListener) {
	dispatcher.m.Lock()
	defer dispatcher.m.Unlock()

	ptr := reflect.ValueOf(l).Pointer()
	for idx, listener := range dispatcher.wildcardEventHolders {
		if reflect.ValueOf(listener).Pointer() == ptr {
			dispatcher.wildcardEventHolders = append(dispatcher.wildcardEventHolders[:idx], dispatcher.wildcardEventHolders[idx+1:]...)
			return
		}
	}
}

func (dispatcher *eventDispatcher) Shutdown(timeout time.Duration) error {
	if !dispatcher.compareAndSwapShutdown(false, true) {
		return fmt.Errorf("dispatcher is already shutting down")
	}
	dispatcher.isShutdown.Store(true)
	if timeout == 0 || timeout > 30*time.Second {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		dispatcher.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return fmt.Errorf("shutdown timeout after %v", timeout)
	}
}

func (dispatcher *eventDispatcher) compareAndSwapShutdown(old, new bool) bool {
	for {
		v := dispatcher.isShutdown.Load()
		if v == nil || v.(bool) != old {
			return false
		}
		if dispatcher.isShutdown.CompareAndSwap(v, new) {
			return true
		}
	}
}
