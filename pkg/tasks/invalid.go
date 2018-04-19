package tasks

import (
	"fmt"
)

type Invalid struct {
	definition interface{}
	err        error
}

func newInvalid(definition interface{}, err error) Task {
	return &Invalid{definition: definition, err: err}
}

func (u *Invalid) load(config *taskConfig) (bool, error) {
	return true, nil
}

func (u *Invalid) name() string {
	return "Invalid task"
}

func (u *Invalid) header() string {
	return ""
}

func (u *Invalid) perform(ctx *Context) (err error) {
	ctx.ui.TaskError(fmt.Errorf("%s: %+v", u.err, u.definition))
	return nil
}
