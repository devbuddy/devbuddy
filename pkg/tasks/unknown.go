package tasks

import (
	"fmt"
)

type Unknown struct {
	definition interface{}
}

func NewUnknown() Task {
	return &Unknown{}
}

func (u *Unknown) Load(definition interface{}) (bool, error) {
	u.definition = definition
	return true, nil
}

func (u *Unknown) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader("Unknown task", "")

	ctx.ui.TaskError(fmt.Errorf("invalid task definition: %+v", u.definition))
	return nil
}
