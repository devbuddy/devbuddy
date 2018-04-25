package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/utils"
)

func init() {
	allTasks["python"] = newPython
}

// Python task: setup a virtualenv with a specified Python version
type Python struct {
	version string
}

func newPython(config *taskConfig) (Task, error) {
	task := &Python{}

	var err error
	task.version, err = config.getPayloadAsString()
	if err != nil {
		return nil, err
	}

	return task, nil
}

func (p *Python) name() string {
	return "Python"
}

func (p *Python) header() string {
	return p.version
}

func (p *Python) actions(ctx *Context) []taskAction {
	pyEnv, err := helpers.NewPyEnv(ctx.cfg)
	if err != nil {
	}

	name := helpers.VirtualenvName(ctx.proj, p.version)
	venv := helpers.NewVirtualenv(ctx.cfg, name)

	return []taskAction{
		&pythonPyenv{version: p.version, pyEnv: pyEnv},
		&pythonInstallVenv{version: p.version, pyEnv: pyEnv},
		&pythonCreateVenv{version: p.version, pyEnv: pyEnv, venv: venv},
	}
}

type pythonPyenv struct {
	version string
	pyEnv   *helpers.PyEnv
}

func (p *pythonPyenv) description() string {
	return "install Python version with PyEnv"
}

func (p *pythonPyenv) needed(ctx *Context) (bool, error) {
	installed, err := p.pyEnv.VersionInstalled(p.version)
	return !installed, err
}

func (p *pythonPyenv) run(ctx *Context) error {
	code, err := runCommand(ctx, "pyenv", "install", p.version)
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("failed to install the required python version. exit code: %d", code)
	}
	return nil
}

type pythonInstallVenv struct {
	version string
	pyEnv   *helpers.PyEnv
}

func (p *pythonInstallVenv) description() string {
	return "install virtualenv"
}

func (p *pythonInstallVenv) needed(ctx *Context) (bool, error) {
	installed := utils.PathExists(p.pyEnv.Which(p.version, "virtualenv"))
	return !installed, nil
}

func (p *pythonInstallVenv) run(ctx *Context) error {
	code, err := runCommand(ctx, p.pyEnv.Which(p.version, "python"), "-m", "pip", "install", "virtualenv")
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("failed to install virtualenv. exit code: %d", code)
	}

	return nil
}

type pythonCreateVenv struct {
	version string
	pyEnv   *helpers.PyEnv
	venv    *helpers.Virtualenv
}

func (p *pythonCreateVenv) description() string {
	return "create virtualenv"
}

func (p *pythonCreateVenv) needed(ctx *Context) (bool, error) {
	return !p.venv.Exists(), nil
}

func (p *pythonCreateVenv) run(ctx *Context) error {
	err := os.MkdirAll(filepath.Dir(p.venv.Path()), 0750)
	if err != nil {
		return err
	}

	code, err := runCommand(ctx, p.pyEnv.Which(p.version, "virtualenv"), p.venv.Path())
	if err != nil {
		return err
	}
	if code != 0 {
		return fmt.Errorf("failed to create the virtualenv. exit code: %d", code)
	}

	return nil
}

func (p *Python) feature(proj *project.Project) (string, string) {
	name := helpers.VirtualenvName(proj, p.version)
	return "python", name
}
