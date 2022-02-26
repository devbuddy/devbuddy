package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("envfile", "EnvFile", parseEnvfile)
}

func parseEnvfile(config *taskapi.TaskConfig, task *taskapi.Task) error {
	envfilePath := ".env"

	check := func(ctx *context.Context) error {
		return nil
	}
	task.AddActionBuilder("", check).SetFeature("envfile", envfilePath)

	return nil
}
