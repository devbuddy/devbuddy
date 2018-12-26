package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	t := registerTaskDefinition("pipfile")
	t.name = "Pipfile"
	t.requiredTask = pythonTaskName
	t.parser = parserPipfile
}

func parserPipfile(config *taskConfig, task *Task) error {
	builder := actionBuilder("install pipfile command", func(ctx *Context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pipenv: %s", result.Error)
		}
		return nil
	})
	builder.OnFunc(func(ctx *Context) *actionResult {
		pythonParam := ctx.features["python"]
		name := helpers.VirtualenvName(ctx.proj, pythonParam)
		venv := helpers.NewVirtualenv(ctx.cfg, name)
		pipenvCmd := venv.Which("pipenv")
		if !utils.PathExists(pipenvCmd) {
			return actionNeeded("Pipenv is not installed in the virtualenv")
		}
		return actionNotNeeded()
	})
	task.addAction(builder.Build())

	builder = actionBuilder("install dependencies from the Pipfile", func(ctx *Context) error {
		result := command(ctx, "pipenv", "install", "--system", "--dev").SetEnvVar("PIPENV_QUIET", "1").Run()
		if result.Error != nil {
			return fmt.Errorf("pipenv failed: %s", result.Error)
		}
		return nil
	})
	task.addAction(builder.Build())

	return nil
}
