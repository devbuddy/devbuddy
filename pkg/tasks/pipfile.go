package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/executor"
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

func (p *Pipfile) perform(ctx *Context) (err error) {
	// We should also check that the python task is executed before this one
	pythonParam, hasPythonFeature := ctx.features["python"]
	if !hasPythonFeature {
		return fmt.Errorf("You must specify a Python environment to use this task")
	}
	venv := helpers.NewVirtualenv(ctx.cfg, pythonParam)

	pipenvInstalled, err := p.installPipenv(ctx, venv)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	InstallRan, err := p.runInstall(ctx, venv)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	if pipenvInstalled || InstallRan {
		ctx.ui.TaskActed()
	} else {
		ctx.ui.TaskAlreadyOk()
	}

	return nil
}

func (p *Pipfile) installPipenv(ctx *Context, venv *helpers.Virtualenv) (acted bool, err error) {
	pipCmd := venv.Which("pip")
	pipenvCmd := venv.Which("pipenv")

	if utils.PathExists(pipenvCmd) {
		return false, nil
	}

	code, err := executor.Run(pipCmd, "install", "--require-virtualenv", "pipenv")
	if err != nil {
		return false, err
	}
	if code != 0 {
		return false, fmt.Errorf("failed to install pipenv for the required python version. exit code: %d", code)
	}

	return true, nil
}

func (p *Pipfile) runInstall(ctx *Context, venv *helpers.Virtualenv) (acted bool, err error) {
	pipenvCmd := venv.Which("pipenv")

	code, err := executor.Run(pipenvCmd, "install", "--system", "--dev")
	if err != nil {
		return false, err
	}
	if code != 0 {
		return false, fmt.Errorf("failed to run `pipenv install`. exit code: %d", code)
	}

	return true, nil
}
