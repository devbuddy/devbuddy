package tasks

import (
	"fmt"
)

type Invalid struct {
	header  string
	name    string
	payload interface{}
	err     error
}

func NewInvalid() Task {
	return &Invalid{header: "Invalid task"}
}

func NewUnknown() Task {
	return &Invalid{header: "Unknown task"}
}

func (u *Invalid) Load(config *taskConfig) (bool, error) {
	u.name = config.name
	u.payload = config.payload
	u.err = nil
	return true, nil
}

func (u *Invalid) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader(u.header, u.name)

	if u.err != nil {
		ctx.ui.TaskError(fmt.Errorf("%s: %+v", u.err, u.payload))
	}
	return nil
}
