package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("env", "Env", parseEnv)
}

func parseEnv(config *taskapi.TaskConfig, task *taskapi.Task) error {
	check := func(ctx *context.Context) error {
		return nil
	}
	task.AddActionWithBuilder("", check).SetFeature("env", "")

	return nil
}
