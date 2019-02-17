package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("envfile", "EnvFile", parseEnvfile)
}

func parseEnvfile(config *taskapi.TaskConfig, task *taskapi.Task) error {
	warning := func(ctx *taskapi.Context) error {
		ctx.UI.TaskWarning("the .env file does NOT override existing variables")
		ctx.UI.TaskWarning("the .env file is NOT unloaded when leaving the project")
		return nil
	}
	task.AddActionWithBuilder("", warning).SetFeature("envfile", "")

	return nil
}
