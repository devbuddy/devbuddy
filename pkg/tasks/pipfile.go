package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	Register("pipfile", "Pipfile", parserPipfile).SetRequiredTask(pythonTaskName)
}

func parserPipfile(config *TaskConfig, task *Task) error {
	builder := actionBuilder("install pipfile command", func(ctx *Context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pipenv: %s", result.Error)
		}
		return nil
	})
	builder.OnFunc(func(ctx *Context) *ActionResult {
		featureInfo := ctx.Features["python"]
		name := helpers.VirtualenvName(ctx.Project, featureInfo.Param)
		venv := helpers.NewVirtualenv(ctx.Cfg, name)
		pipenvCmd := venv.Which("pipenv")
		if !utils.PathExists(pipenvCmd) {
			return ActionNeeded("Pipenv is not installed in the virtualenv")
		}
		return ActionNotNeeded()
	})
	task.AddAction(builder.Build())

	builder = actionBuilder("install dependencies from the Pipfile", func(ctx *Context) error {
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
