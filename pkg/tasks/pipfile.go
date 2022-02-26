package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	api.Register("pipfile", "Pipfile", parserPipfile).SetRequiredTask(pythonTaskName)
}

func parserPipfile(config *api.TaskConfig, task *api.Task) error {
	installPipfile := func(ctx *context.Context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pipenv: %w", result.Error)
		}
		return nil
	}
	installPipfileNeeded := func(ctx *context.Context) *api.ActionResult {
		version, err := findAutoEnvFeatureParam(ctx, "python")
		if err != nil {
			return api.Failed("missing python feature: %s", err)
		}
		name := helpers.VirtualenvName(ctx.Project, version)
		venv := helpers.NewVirtualenv(ctx.Cfg, name)
		pipenvCmd := venv.Which("pipenv")
		if !utils.PathExists(pipenvCmd) {
			return api.Needed("Pipenv is not installed in the virtualenv")
		}
		return api.NotNeeded()
	}
	task.AddActionBuilder("install pipfile command", installPipfile).
		On(api.FuncCondition(installPipfileNeeded))

	runPipfileInstall := func(ctx *context.Context) error {
		result := command(ctx, "pipenv", "install", "--system", "--dev").SetEnvVar("PIPENV_QUIET", "1").Run()
		if result.Error != nil {
			return fmt.Errorf("pipenv failed: %w", result.Error)
		}
		return nil
	}
	task.AddActionBuilder("install dependencies from the Pipfile", runPipfileInstall).
		On(api.FileCondition("Pipfile")).
		On(api.FileCondition("Pipfile.lock"))

	return nil
}
