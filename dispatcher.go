package event

import (
	"reflect"
	"sync"
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
}

type eventDispatcher struct {
	m                    sync.RWMutex
	eventHolders         map[string][]IListener
	wildcardEventHolders []IListener
}

func NewDispatcher() *eventDispatcher {
	d := new(eventDispatcher)
	d.eventHolders = make(map[string][]IListener)
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
	dispatcher.m.RLock()
	defer dispatcher.m.RUnlock()

	for _, listener := range dispatcher.eventHolders[e.eventType] {
		if listener.AsyncProcess() {
			go func(l IListener) {
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
			go func(l IListener) {
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
