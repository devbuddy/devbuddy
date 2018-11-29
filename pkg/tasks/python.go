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

	action := newAction("install Pyenv", func(ctx *context) error {
		result := command(ctx, "brew", "install", "pyenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pyenv: %s", result.Error)
		}
		return nil
	})
	action.onFunc(func(ctx *context) *actionResult {
		_, err := helpers.NewPyEnv()
		if err != nil {
			return actionNeeded("Pyenv is not installed: %s", err)
		}
		return actionNotNeeded()
	})
	task.addAction(action)

	action = newAction("install Python version with PyEnv", func(ctx *context) error {
		result := command(ctx, "pyenv", "install", version).Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install the required python version: %s", result.Error)
		}
		return nil
	})
	action.onFunc(func(ctx *context) *actionResult {
		pyEnv, err := helpers.NewPyEnv()
		if err != nil {
			return actionFailed("cannot use pyenv: %s", err)
		}

		installed, err := pyEnv.VersionInstalled(version)
		if err != nil {
			return actionFailed("failed to check if python version is installed: %s", err)
		}

		if !installed {
			return actionNeeded("python version is not installed")
		}

		return actionNotNeeded()
	})
	task.addAction(action)

	action = newAction("install virtualenv", func(ctx *context) error {
		pyEnv, err := helpers.NewPyEnv()
		if err != nil {
			return err
		}
		result := command(ctx, pyEnv.Which(version, "python"), "-m", "pip", "install", "virtualenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install virtualenv: %s", result.Error)
		}

		return nil
	})
	action.onFunc(func(ctx *context) *actionResult {
		pyEnv, err := helpers.NewPyEnv()
		if err != nil {
			return actionFailed("cannot use pyenv: %s", err)
		}

		installed := utils.PathExists(pyEnv.Which(version, "virtualenv"))
		if !installed {
			return actionNeeded("virtualenv is not installed")
		}

		return actionNotNeeded()
	})
	task.addAction(action)

	action = newAction("create virtualenv", func(ctx *context) error {
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
	})
	action.onFunc(func(ctx *context) *actionResult {
		name := helpers.VirtualenvName(ctx.proj, version)
		venv := helpers.NewVirtualenv(ctx.cfg, name)

		if !venv.Exists() {
			return actionNeeded("project virtualenv does not exists")
		}

		return actionNotNeeded()
	})
	task.addAction(action)

	return nil
}
