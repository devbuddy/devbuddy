package tasks

import (
	"fmt"
)

func parseUnknown(config *taskConfig, task *Task) error {
	warning := fmt.Sprintf("Unknown task: \"%s\"", config.name)

	task.addAction(newAction("", func(ctx *context) error {
		ctx.ui.TaskWarning(warning)
		return nil
	}))

	return nil
}
