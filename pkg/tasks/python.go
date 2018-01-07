package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/helpers"
)

func init() {
	allTasks["python"] = NewPython
}

type Python struct {
	version string
}

func NewPython() Task {
	return &Python{}
}

func (p *Python) Load(definition interface{}) (bool, error) {
	def, ok := definition.(map[interface{}]interface{})
	if !ok {
		return false, nil
	}
	if version, ok := def["python"]; ok {
		p.version, ok = version.(string)
		if !ok {
			return false, nil
		}
		return true, nil
	}

	return false, nil
}

func (p *Python) Perform(ctx *Context) (err error) {
	ctx.ui.TaskHeader("Python", p.version)

	installed, err := p.InstallPython(ctx)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	venvInstalled, err := p.InstallVirtualEnv(ctx)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	venvCreated, err := p.CreateVirtualEnv(ctx)
	if err != nil {
		ctx.ui.TaskError(err)
		return err
	}

	if installed || venvInstalled || venvCreated {
		ctx.ui.TaskActed()
	} else {
		ctx.ui.TaskAlreadyOk()
	}

	return nil
}

func (p *Python) InstallPython(ctx *Context) (acted bool, err error) {
	pyEnv := helpers.NewPyEnv(ctx.cfg, ctx.proj)

	installed, err := pyEnv.VersionInstalled(p.version)
	if err != nil {
		return
	}
	if installed {
		return
	}

	code, err := executor.Run("pyenv", "install", p.version)
	if err != nil {
		return
	}
	if code != 0 {
		return false, fmt.Errorf("failed to install the required python version. exit code: %d", code)
	}

	return true, nil
}

func (p *Python) InstallVirtualEnv(ctx *Context) (acted bool, err error) {
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

func (p *Python) CreateVirtualEnv(ctx *Context) (acted bool, err error) {
	venv := helpers.NewVirtualenv(ctx.cfg, ctx.proj, p.version)
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

func (p *Python) Features() map[string]string {
	return map[string]string{"python": p.version}
}
