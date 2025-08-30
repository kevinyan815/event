package event

import (
	"context"
	"reflect"
	"sync"
	"time"

	"github.com/google/uuid"
)

type IEvent interface {
	EventName() string
	EntityID() string
	EntityType() string
}

// Event : Type of domainevent
type Event struct {
	ID            string          `json:"event_id"` // unique event id
	Name          string          `json:"event_name"`
	EntityID      string          `json:"entity_id"`   // entity id like: user_id, order_id, order_no, etc. using string to cover all cases not only numeric identiry
	EntityType    string          `json:"entity_type"` // entity type like: user, order, etc.
	EventTime     time.Time       `json:"occurred_on"` // event occurred time
	ConcreteEvent IEvent          `json:"event_data"`  // event data
	context       context.Context `json:"-"`
	eventType     string
}

func NewEvent(concreteEvent IEvent) *Event {
	return NewEventWithContext(context.Background(), concreteEvent)
}

// NewEventWithContext : Factory of Event obj with context
func NewEventWithContext(ctx context.Context, concreteEvent IEvent) *Event {
	eventId := uuid.New().String()
	event := &Event{
		ID:            eventId,
		Name:          concreteEvent.EventName(),
		EntityID:      concreteEvent.EntityID(),
		EntityType:    concreteEvent.EntityType(),
		EventTime:     time.Now(),
		ConcreteEvent: concreteEvent,
		context:       ctx,
		eventType:     reflect.TypeOf(concreteEvent).String(),
	}

	return event
}

func (e *Event) EventName() string {
	return e.Name
}

func (e *Event) OccurredOn() time.Time {
	return e.EventTime
}

func (e *Event) EventData() IEvent {
	return e.ConcreteEvent
}

func (e *Event) Context() context.Context {
	return e.context
}

var _Dispatcher *eventDispatcher
var once sync.Once

func Dispatcher() *eventDispatcher {
	once.Do(func() {
		_Dispatcher = NewDispatcher()
	})
	return _Dispatcher
}
