package domainevent

import "github.com/go-study-lab/event"

func init() {
	eventDispatcher := event.Dispatcher()
	eventDispatcher.Subscribe(event.NewEvent(&UserCreated{}), &UserCreatedListener{}, &UserCreatedErrListener{})
	eventDispatcher.Subscribe(event.NewEvent(&UserUpdated{}), &UserUpdatedListener{})
	eventDispatcher.SubscribeWildcard(&LogEventListener{})
}
