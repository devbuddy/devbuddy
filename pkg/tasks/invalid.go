package tasks

import (
	"fmt"
)

type Invalid struct {
	header     string
	name       string
	definition interface{}
	err        error
}

func NewInvalid() Task {
	return &Invalid{header: "Invalid task"}
}

func NewUnknown() Task {
	return &Invalid{header: "Unknown task"}
}

func (u *Invalid) Load(definition interface{}) (bool, error) {
	u.name, u.err = extractTaskName(definition)
	u.definition = definition
	return true, nil
}

func (u *Invalid) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader(u.header, u.name)

	if u.err != nil {
		ctx.ui.TaskError(fmt.Errorf("%s: %+v", u.err, u.definition))
	}
	return nil
}
