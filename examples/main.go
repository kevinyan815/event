package main

import (
	"context"
	"fmt"
	"time"

	"github.com/go-study-lab/event"
	"github.com/go-study-lab/event/examples/domainevent"
)

type User struct {
	UserId   int
	UserName string
}

func PublishCreateUserEvent() {
	user := &User{
		UserId:   1,
		UserName: "KK",
	}
	// assume use tx to create user, then publish a userCreated event
	ctx := context.WithValue(context.Background(), "trace", "11111111")
	event.Dispatcher().Dispatch(event.NewEventWithContext(ctx, &domainevent.UserCreated{
		UserId:   user.UserId,
		UserName: user.UserName,
	}))
}

func HasEventListener() {
	e := new(domainevent.UserUpdated)
	ok := event.Dispatcher().HasEventListener(event.NewEvent(e))
	fmt.Printf("Event: %s have listener? %t \n", e.EventName(), ok)

	if ok {
		event.Dispatcher().Dispatch(event.NewEvent(new(domainevent.UserUpdated)))
	}
}

func RemoveEventListener() {
	listener := new(domainevent.UserUpdatedListener)

	event.Dispatcher().Subscribe(
		event.NewEvent(&domainevent.UserUpdated{}),
		listener,
	)

	event.Dispatcher().RemoveEventListener(
		event.NewEvent(&domainevent.UserUpdated{}),
		listener,
	)

}

func main() {

	PublishCreateUserEvent()
	HasEventListener()
	RemoveEventListener()

	// sleep 10 seconds to wait for async process
	time.Sleep(10 * time.Second)
}
