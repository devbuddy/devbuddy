package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("env", "Env", parseEnv)
}

func parseEnv(config *api.TaskConfig, task *api.Task) error {
	check := func(ctx *context.Context) error {
		return nil
	}
	task.AddActionBuilder("", check).SetFeature("env", "")

	return nil
}
