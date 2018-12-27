package tasks

import (
	"fmt"
)

func parseUnknown(config *taskConfig, task *Task) error {
	builder := actionBuilder("", func(ctx *Context) error {
		ctx.ui.TaskWarning(fmt.Sprintf("Unknown task: \"%s\"", config.name))
		return nil
	})
	task.addAction(builder.Build())
	return nil
}
