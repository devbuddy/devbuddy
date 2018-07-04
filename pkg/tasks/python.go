package tasks

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	t := registerTask("python")
	t.name = "Python"
	t.parser = parserPython
}

func parserPython(config *taskConfig, task *Task) error {
	version, err := config.getPayloadAsString()
	if err != nil {
		return err
	}

	task.header = version

	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		log.Fatalf("PyEnv helper failed: %s", err)
	}

	name := helpers.VirtualenvName(ctx.proj, version)
	venv := helpers.NewVirtualenv(ctx.cfg, name)

	task.addAction(&pythonPyenv{version: version, pyEnv: pyEnv})
	task.addAction(&pythonInstallVenv{version: version, pyEnv: pyEnv})
	task.addAction(&pythonCreateVenv{version: version, pyEnv: pyEnv, venv: venv})

	return nil
}

func (p *Python) actions(ctx *context) []taskAction {
	pyEnv, err := helpers.NewPyEnv(ctx.cfg)
	if err != nil {
		log.Fatalf("PyEnv helper failed: %s", err)
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

func (p *pythonPyenv) needed(ctx *context) (bool, error) {
	installed, err := p.pyEnv.VersionInstalled(p.version)
	return !installed, err
}

func (p *pythonPyenv) run(ctx *context) error {
	err := command(ctx, "pyenv", "install", p.version).Run()
	if err != nil {
		return fmt.Errorf("failed to install the required python version: %s", err)
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

func (p *pythonInstallVenv) needed(ctx *context) (bool, error) {
	installed := utils.PathExists(p.pyEnv.Which(p.version, "virtualenv"))
	return !installed, nil
}

func (p *pythonInstallVenv) run(ctx *context) error {
	err := command(ctx, p.pyEnv.Which(p.version, "python"), "-m", "pip", "install", "virtualenv").Run()
	if err != nil {
		return fmt.Errorf("failed to install virtualenv: %s", err)
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

func (p *pythonCreateVenv) needed(ctx *context) (bool, error) {
	return !p.venv.Exists(), nil
}

func (p *pythonCreateVenv) run(ctx *context) error {
	err := os.MkdirAll(filepath.Dir(p.venv.Path()), 0750)
	if err != nil {
		return err
	}

	err = command(ctx, p.pyEnv.Which(p.version, "virtualenv"), p.venv.Path()).Run()
	if err != nil {
		return fmt.Errorf("failed to create the virtualenv: %s", err)
	}

	return nil
}

func (p *Python) feature(proj *project.Project) (string, string) {
	name := helpers.VirtualenvName(proj, p.version)
	return "python", name
}
