package tasks

import (
	"fmt"
)

type Unknown struct {
	providedName string
}

func NewUnknown() Task {
	return &Unknown{}
}

func (u *Unknown) Load(config *taskConfig) (bool, error) {
	u.providedName = config.name
	return true, nil
}

func (u *Unknown) name() string {
	return u.providedName
}

func (u *Unknown) header() string {
	return ""
}

func (u *Unknown) Perform(ctx *Context) (err error) {
	ctx.ui.TaskError(fmt.Errorf("Unknown task"))
	return nil
}
