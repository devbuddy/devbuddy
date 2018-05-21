package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/utils"
)

func init() {
	allTasks["pipfile"] = newPipfile
}

type Pipfile struct {
}

func newPipfile(config *taskConfig) (Task, error) {
	return &Pipfile{}, nil
}

func (p *Pipfile) name() string {
	return "Pipfile"
}

func (p *Pipfile) header() string {
	return ""
}

func (p *Pipfile) preRunValidation(ctx *context) (err error) {
	_, hasPythonFeature := ctx.features["python"]
	if !hasPythonFeature {
		return fmt.Errorf("You must specify a Python environment to use this task")
	}
	return nil
}

func (p *Pipfile) actions(ctx *context) []taskAction {
	return []taskAction{
		&pipfileInstall{},
		&pipfileRun{},
	}
}

type pipfileInstall struct {
}

func (p *pipfileInstall) description() string {
	return "install pipfile command"
}

func (p *pipfileInstall) needed(ctx *context) (bool, error) {
	pythonParam := ctx.features["python"]
	venv := helpers.NewVirtualenv(ctx.cfg, pythonParam)
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
	err := command(ctx, "pipenv", "install", "--system", "--dev").Run()
	if err != nil {
		return fmt.Errorf("pipenv failed: %s", err)
	}
	p.success = true
	return nil
}
