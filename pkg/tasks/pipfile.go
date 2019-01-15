package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	taskapi.Register("pipfile", "Pipfile", parserPipfile).SetRequiredTask(pythonTaskName)
}

func parserPipfile(config *taskapi.TaskConfig, task *taskapi.Task) error {
	installPipfile := func(ctx *taskapi.Context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pipenv: %s", result.Error)
		}
		return nil
	}
	installPipfileNeeded := func(ctx *taskapi.Context) *taskapi.ActionResult {
		featureInfo := ctx.Features["python"]
		name := helpers.VirtualenvName(ctx.Project, featureInfo.Param)
		venv := helpers.NewVirtualenv(ctx.Cfg, name)
		pipenvCmd := venv.Which("pipenv")
		if !utils.PathExists(pipenvCmd) {
			return taskapi.ActionNeeded("Pipenv is not installed in the virtualenv")
		}
		return taskapi.ActionNotNeeded()
	}
	task.AddActionWithBuilder("install pipfile command", installPipfile).
		OnFunc(installPipfileNeeded)

	runPipfileInstall := func(ctx *taskapi.Context) error {
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
