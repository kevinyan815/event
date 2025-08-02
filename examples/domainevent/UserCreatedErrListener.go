package domainevent

import (
	"fmt"

	"github.com/go-study-lab/event"
)

type UserCreatedErrListener struct {
}

func (*UserCreatedErrListener) EventHandler() event.Handler {
	// this listener will return error to test the error handling in event dispatcher
	return func(e *event.Event) error {
		return fmt.Errorf("can not convert event data to type *UserCreated")
	}
}

func (*UserCreatedErrListener) AsyncProcess() bool {
	return false
}
