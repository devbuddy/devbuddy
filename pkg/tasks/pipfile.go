package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	taskapi.RegisterTaskDefinition("pipfile", "Pipfile", parserPipfile).AddRequiredTask(pythonTaskName)
}

func parserPipfile(config *taskapi.TaskConfig, task *taskapi.Task) error {
	builder := actionBuilder("install pipfile command", func(ctx *taskapi.Context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pipenv: %s", result.Error)
		}
		return nil
	})
	builder.OnFunc(func(ctx *taskapi.Context) *actionResult {
		featureInfo := ctx.features["python"]
		name := helpers.VirtualenvName(ctx.proj, featureInfo.Param)
		venv := helpers.NewVirtualenv(ctx.cfg, name)
		pipenvCmd := venv.Which("pipenv")
		if !utils.PathExists(pipenvCmd) {
			return actionNeeded("Pipenv is not installed in the virtualenv")
		}
		return actionNotNeeded()
	})
	task.AddAction(builder.Build())

	builder = actionBuilder("install dependencies from the Pipfile", func(ctx *taskapi.Context) error {
		result := command(ctx, "pipenv", "install", "--system", "--dev").SetEnvVar("PIPENV_QUIET", "1").Run()
		if result.Error != nil {
			return fmt.Errorf("pipenv failed: %s", result.Error)
		}
		return nil
	})
	builder.OnFileChange("Pipfile")
	builder.OnFileChange("Pipfile.lock")
	task.AddAction(builder.Build())

	return nil
}
