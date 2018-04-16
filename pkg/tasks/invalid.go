package tasks

import (
	"fmt"
)

type Invalid struct {
	definition interface{}
	err        error
}

func NewInvalid(definition interface{}, err error) Task {
	return &Invalid{definition: definition, err: err}
}

func (u *Invalid) Load(config *taskConfig) (bool, error) {
	return true, nil
}

func (u *Invalid) name() string {
	return "Invalid task"
}

func (u *Invalid) header() string {
	return ""
}

func (u *Invalid) Perform(ctx *Context) (err error) {
	ctx.ui.TaskError(fmt.Errorf("%s: %+v", u.err, u.definition))
	return nil
}
