package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

func init() {
	allTasks["pipfile"] = NewPipfile
}

type Pipfile struct {
}

func NewPipfile() Task {
	return &Pipfile{}
}

func (p *Pipfile) Load(definition interface{}) (bool, error) {
	// def, ok := definition.(string)
	// if !ok {
	// 	return false, nil
	// }
	// if def == "pipfile" {
	// 	return true, nil
	// }

	// return false, nil

	return true, nil
}

func (p *Pipfile) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader("Pipfile", p.version)

	pipenvInstalled, err := p.InstallPipenv(ctx)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	InstallRan, err := p.RunInstall(ctx)
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

func (p *Pipfile) InstallPipenv(ctx *Context) (acted bool, err error) {
	pyEnv := helpers.NewPyEnv(ctx.cfg, ctx.proj)

	if config.PathExists(pyEnv.Which(p.version, "virtualenv")) {
		return false, nil
	}

	code, err := executor.Run(pyEnv.Which(p.version, "pip"), "install", "virtualenv")
	if err != nil {
		return
	}
	if code != 0 {
		return false, fmt.Errorf("failed to install virtualenv for the required python version. exit code: %d", code)
	}

	return true, nil
}

func (p *Pipfile) RunInstall(ctx *Context) (acted bool, err error) {
	name := helpers.VirtualenvName(ctx.proj, p.version)
	venv := helpers.NewVirtualenv(ctx.cfg, name)
	pyEnv := helpers.NewPyEnv(ctx.cfg, ctx.proj)

	if venv.Exists() {
		return false, nil
	}

	err = os.MkdirAll(filepath.Dir(venv.Path()), 0750)
	if err != nil {
		return
	}

	code, err := executor.Run(pyEnv.Which(p.version, "virtualenv"), venv.Path())
	if err != nil {
		return
	}
	if code != 0 {
		return false, fmt.Errorf("failed to create the virtualenv. exit code: %d", code)
	}

	return true, nil
}
