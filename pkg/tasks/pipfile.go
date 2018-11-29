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
	task.addAction(
		newAction("install pipfile command", pipfileInstallRun).
			onFunc(pipfileInstallNeeded),
	)
	task.addAction(
		newAction("install dependencies from the Pipfile", pipfileRun))
	return nil
}

func pipfileInstallNeeded(ctx *context) *actionResult {
	pythonParam := ctx.features["python"]
	name := helpers.VirtualenvName(ctx.proj, pythonParam)
	venv := helpers.NewVirtualenv(ctx.cfg, name)
	pipenvCmd := venv.Which("pipenv")

	if !utils.PathExists(pipenvCmd) {
		return actionNeeded("Pipenv is not installed in the virtualenv")
	}
	return actionNotNeeded()
}

func pipfileInstallRun(ctx *context) error {
	result := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to install pipenv: %s", result.Error)
	}
	return nil
}

func pipfileRun(ctx *context) error {
	result := command(ctx, "pipenv", "install", "--system", "--dev").SetEnvVar("PIPENV_QUIET", "1").Run()
	if result.Error != nil {
		return fmt.Errorf("pipenv failed: %s", result.Error)
	}
	return nil
}
