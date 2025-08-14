package domainevent

import (
	"fmt"
	"github.com/kevinyan815/event"
)

type UserUpdatedListener struct {
}

func (*UserUpdatedListener) EventHandler() event.Handler {
	return func(e *event.Event) error {
		fmt.Println("UserUpdatedListener executed!")
		return nil
	}
}

func (*UserUpdatedListener) AsyncProcess() bool {
	return true
}
