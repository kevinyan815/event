package domainevent

import (
	"fmt"
)

type UserUpdated struct {
	UserId   int
	UserName string
}

func (*UserUpdated) EventName() string {
	return "UserUpdated"
}

func (e *UserUpdated) EntityID() string {
	return fmt.Sprintf("%d", e.UserId)
}

func (e *UserUpdated) EntityType() string {
	return "User"
}
