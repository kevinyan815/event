package domainevent

import (
	"fmt"

	"github.com/kevinyan815/event"
)

type LogEventListener struct {
}

func (*LogEventListener) EventHandler() event.Handler {
	return func(e *event.Event) error {
		fmt.Printf("LogEventListener executed! eventId: %s, eventName: %s, entityId: %s, entityType: %s, eventTime: %s, eventData: %+v\n", e.ID, e.Name, e.EntityID, e.EntityType, e.EventTime, e.ConcreteEvent)
		return nil
	}
}

func (*LogEventListener) AsyncProcess() bool {
	return false
}
