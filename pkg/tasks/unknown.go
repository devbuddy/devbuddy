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

func (u *Unknown) perform(ctx *context) (err error) {
	ctx.ui.TaskWarning(fmt.Sprintf("Unknown task"))
	return nil
}

func (u *Unknown) actions(ctx *context) (actions []taskAction) {
	return
}
