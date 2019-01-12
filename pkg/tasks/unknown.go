package tasks

import (
	"fmt"
)

func parseUnknown(config *TaskConfig, task *Task) error {
	builder := actionBuilder("", func(ctx *Context) error {
		ctx.UI.TaskWarning(fmt.Sprintf("Unknown task: \"%s\"", config.name))
		return nil
	})
	task.AddAction(builder.Build())
	return nil
}
