package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("envfile", "EnvFile", parseEnvfile)
}

func parseEnvfile(config *api.TaskConfig, task *api.Task) error {
	envfilePath := ".env"

	check := func(ctx *context.Context) error {
		return nil
	}
	task.AddActionBuilder("", check).SetFeature("envfile", envfilePath)

	return nil
}
