package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	taskapi.Register("pipfile", "Pipfile", parserPipfile).SetRequiredTask(pythonTaskName)
}

func parserPipfile(config *taskapi.TaskConfig, task *taskapi.Task) error {
	installPipfile := func(ctx *context.Context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pipenv: %s", result.Error)
		}
		return nil
	}
	installPipfileNeeded := func(ctx *context.Context) *taskapi.ActionResult {
		version, err := findAutoEnvFeatureParam(ctx, "python")
		if err != nil {
			return taskapi.ActionFailed("missing python feature: %s", err)
		}
		name := helpers.VirtualenvName(ctx.Project, version)
		venv := helpers.NewVirtualenv(ctx.Cfg, name)
		pipenvCmd := venv.Which("pipenv")
		if !utils.PathExists(pipenvCmd) {
			return taskapi.ActionNeeded("Pipenv is not installed in the virtualenv")
		}
		return taskapi.ActionNotNeeded()
	}
	task.AddActionWithBuilder("install pipfile command", installPipfile).
		OnFunc(installPipfileNeeded)

	runPipfileInstall := func(ctx *context.Context) error {
		result := command(ctx, "pipenv", "install", "--system", "--dev").SetEnvVar("PIPENV_QUIET", "1").Run()
		if result.Error != nil {
			return fmt.Errorf("pipenv failed: %s", result.Error)
		}
		return nil
	}
	task.AddActionWithBuilder("install dependencies from the Pipfile", runPipfileInstall).
		OnFileChange("Pipfile").
		OnFileChange("Pipfile.lock")

	return nil
}
