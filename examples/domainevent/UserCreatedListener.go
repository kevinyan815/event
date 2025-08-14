package domainevent

import (
	"fmt"

	"github.com/kevinyan815/event"
)

type UserCreatedListener struct {
}

func (*UserCreatedListener) EventHandler() event.Handler {
	return func(e *event.Event) error {
		var eData *UserCreated
		var ok bool

		if eData, ok = e.EventData().(*UserCreated); !ok {
			return fmt.Errorf("can not convert event data to type *UserCreated, eventId: %s, eventName: %s, entityId: %s, entityType: %s",
				e.ID, e.Name, e.EntityID, e.EntityType)
		}

		fmt.Printf("event data, userId:%d, userName:%s \n", eData.UserId, eData.UserName)
		return nil
	}
}

// When we want to use a same
// DB transaction in event trigger and listener to  get ACID assurance, then
// AsyncProcess must return false.
func (*UserCreatedListener) AsyncProcess() bool {
	return false
}
