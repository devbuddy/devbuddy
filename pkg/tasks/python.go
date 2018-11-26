package tasks

import (
	"fmt"
	"os"
	"path/filepath"

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

	task.addActionWithBuilder("install Pyenv", func(ctx *context) error {
		result := command(ctx, "brew", "install", "pyenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pyenv: %s", result.Error)
		}
		return nil
	}).addConditionFunc(func(ctx *context) *actionResult {
		_, err := helpers.NewPyEnv()
		if err != nil {
			return actionNeeded("Pyenv is not installed: %s", err)
		}
		return actionNotNeeded()
	})

	task.addAction(&pythonPyenv{version: version})
	task.addAction(&pythonInstallVenv{version: version})

	task.addActionWithBuilder("create virtualenv", func(ctx *context) error {
		name := helpers.VirtualenvName(ctx.proj, version)
		venv := helpers.NewVirtualenv(ctx.cfg, name)

		err := os.MkdirAll(filepath.Dir(venv.Path()), 0750)
		if err != nil {
			return err
		}

		pyEnv, err := helpers.NewPyEnv()
		if err != nil {
			return err
		}

		result := command(ctx, pyEnv.Which(version, "virtualenv"), venv.Path()).Run()
		if result.Error != nil {
			return fmt.Errorf("failed to create the virtualenv: %s", result.Error)
		}

		return nil
	}).addConditionFunc(func(ctx *context) *actionResult {
		name := helpers.VirtualenvName(ctx.proj, version)
		venv := helpers.NewVirtualenv(ctx.cfg, name)

		if !venv.Exists() {
			return actionNeeded("project virtualenv does not exists")
		}

		return actionNotNeeded()
	})

	return nil
}

type pythonPyenv struct {
	version string
}

func (p *pythonPyenv) description() string {
	return "install Python version with PyEnv"
}

func (p *pythonPyenv) needed(ctx *context) *actionResult {
	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		return actionFailed("cannot use pyenv: %s", err)
	}

	installed, err := pyEnv.VersionInstalled(p.version)
	if err != nil {
		return actionFailed("failed to check if python version is installed: %s", err)
	}

	if !installed {
		return actionNeeded("python version is not installed")
	}

	return actionNotNeeded()
}

func (p *pythonPyenv) run(ctx *context) error {
	result := command(ctx, "pyenv", "install", p.version).Run()
	if result.Error != nil {
		return fmt.Errorf("failed to install the required python version: %s", result.Error)
	}
	return nil
}

type pythonInstallVenv struct {
	version string
}

func (p *pythonInstallVenv) description() string {
	return "install virtualenv"
}

func (p *pythonInstallVenv) needed(ctx *context) *actionResult {
	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		return actionFailed("cannot use pyenv: %s", err)
	}

	installed := utils.PathExists(pyEnv.Which(p.version, "virtualenv"))
	if !installed {
		return actionNeeded("virtualenv is not installed")
	}

	return actionNotNeeded()
}

func (p *pythonInstallVenv) run(ctx *context) error {
	pyEnv, err := helpers.NewPyEnv()
	if err != nil {
		return err
	}
	result := command(ctx, pyEnv.Which(p.version, "python"), "-m", "pip", "install", "virtualenv").Run()
	if result.Error != nil {
		return fmt.Errorf("failed to install virtualenv: %s", result.Error)
	}

	return nil
}
