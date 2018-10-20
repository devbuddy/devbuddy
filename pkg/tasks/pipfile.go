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
	task.addAction(&pipfileInstall{})
	task.addAction(&pipfileRun{})
	return nil
}

type pipfileInstall struct {
}

func (p *pipfileInstall) description() string {
	return "install pipfile command"
}

func (p *pipfileInstall) needed(ctx *context) (bool, error) {
	pythonParam := ctx.features["python"]
	name := helpers.VirtualenvName(ctx.proj, pythonParam)
	venv := helpers.NewVirtualenv(ctx.cfg, name)
	pipenvCmd := venv.Which("pipenv")
	return !utils.PathExists(pipenvCmd), nil
}

func (p *pipfileInstall) run(ctx *context) error {
	err := command(ctx, "pip", "install", "--require-virtualenv", "pipenv").Run()
	if err != nil {
		return fmt.Errorf("failed to install pipenv: %s", err)
	}
	return nil
}

type pipfileRun struct {
	success bool
}

func (p *pipfileRun) description() string {
	return "install dependencies from the Pipfile"
}

func (p *pipfileRun) needed(ctx *context) (bool, error) {
	return !p.success, nil
}

func (p *pipfileRun) run(ctx *context) error {
	err := command(ctx, "pipenv", "install", "--system", "--dev").SetEnvVar("PIPENV_QUIET", "1").Run()
	if err != nil {
		return fmt.Errorf("pipenv failed: %s", err)
	}
	p.success = true
	return nil
}
