package domainevent

import (
	"fmt"
)

type UserCreated struct {
	UserId   int
	UserName string
}

func (*UserCreated) EventName() string {
	return "UserCreated"
}

func (e *UserCreated) EntityID() string {
	return fmt.Sprintf("%d", e.UserId)
}

func (e *UserCreated) EntityType() string {
	return "User"
}
