package tasks

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"

	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

const pythonTaskName = "python"

func init() {
	t := registerTaskDefinition(pythonTaskName)
	t.name = "Python"
	t.parser = parserPython
}

func parserPython(config *taskConfig, task *Task) error {
	version, err := config.getStringProperty("version", true)
	if err != nil {
		return err
	}

	task.header = version
	task.featureName = "python"
	task.featureParam = version

	task.addAction(&pyenvPath{})
	task.addAction(&pyenv{})
	task.addAction(&pythonPyenv{version: version})
	task.addAction(&pythonInstallVenv{version: version})
	task.addAction(&pythonCreateVenv{version: version})

	return nil
}

type pyenvPath struct{}

func (p *pyenvPath) description() string {
	return ""
}

func (p *pyenvPath) needed(ctx *context) (bool, error) {
	match, _ := regexp.MatchString("\\.pyenv/bin", ctx.env.Get("PATH"))
	return !match, nil
}

func (p *pyenvPath) run(ctx *context) error {
	ctx.ui.TaskWarning("The PATH environment variable should have 'pyenv' in it. Try adding the following to ~/.bashrc\n\nexport PATH=\"~/.pyenv/bin:$PATH\"")
	return nil
}

type pyenv struct{}

func (p *pyenv) description() string {
	return "install PyEnv"
}

func (p *pyenv) needed(ctx *context) (bool, error) {
	_, err := helpers.NewPyEnv()
	if err == nil {
		return false, nil
	}

	return true, err
}

func (p *pyenv) run(ctx *context) error {
	err := command(ctx, "bash", "-c", "curl -L https://github.com/pyenv/pyenv-installer/raw/master/bin/pyenv-installer | bash").Run()
	if err != nil {
		return fmt.Errorf("failed to install pyenv: %s", err)
	}
	return nil
}

type pythonPyenv struct {
	version string
}

func (p *pythonPyenv) description() string {
	return "install Python version with PyEnv"
}

func (p *pythonPyenv) needed(ctx *context) (bool, error) {
	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		return false, err
	}

	installed, err := pyEnv.VersionInstalled(p.version)
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
}

func (p *pythonInstallVenv) description() string {
	return "install virtualenv"
}

func (p *pythonInstallVenv) needed(ctx *context) (bool, error) {
	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		return false, err
	}
	installed := utils.PathExists(pyEnv.Which(p.version, "virtualenv"))
	return !installed, nil
}

func (p *pythonInstallVenv) run(ctx *context) error {
	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		return err
	}
	err = command(ctx, pyEnv.Which(p.version, "python"), "-m", "pip", "install", "virtualenv").Run()
	if err != nil {
		return fmt.Errorf("failed to install virtualenv: %s", err)
	}

	return nil
}

type pythonCreateVenv struct {
	version string
}

func (p *pythonCreateVenv) description() string {
	return "create virtualenv"
}

func (p *pythonCreateVenv) needed(ctx *context) (bool, error) {
	name := helpers.VirtualenvName(ctx.proj, p.version)
	venv := helpers.NewVirtualenv(ctx.cfg, name)
	return !venv.Exists(), nil
}

func (p *pythonCreateVenv) run(ctx *context) error {
	name := helpers.VirtualenvName(ctx.proj, p.version)
	venv := helpers.NewVirtualenv(ctx.cfg, name)

	err := os.MkdirAll(filepath.Dir(venv.Path()), 0750)
	if err != nil {
		return err
	}

	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		return err
	}

	err = command(ctx, pyEnv.Which(p.version, "virtualenv"), venv.Path()).Run()
	if err != nil {
		return fmt.Errorf("failed to create the virtualenv: %s", err)
	}

	return nil
}
