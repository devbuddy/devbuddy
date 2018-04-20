package tasks

import (
	"fmt"
)

type Unknown struct {
	providedName string
}

func newUnknown(config *taskConfig) (Task, error) {
	return &Unknown{providedName: config.name}, nil
}

func (u *Unknown) name() string {
	return u.providedName
}

func (u *Unknown) header() string {
	return ""
}

func (u *Unknown) perform(ctx *Context) (err error) {
	ctx.ui.TaskError(fmt.Errorf("Unknown task"))
	return nil
}
