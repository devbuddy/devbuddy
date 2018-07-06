package tasks

import (
	"fmt"
)

func parseUnknown(config *taskConfig, task *Task) error {
	task.perform = func(ctx *context) (err error) {
		ctx.ui.TaskWarning(fmt.Sprintf("Unknown task: \"%s\"", config.name))
		return nil
	}
	return nil
}
