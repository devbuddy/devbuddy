package tasks

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

const pythonTaskName = "python"

func init() {
	taskapi.RegisterTaskDefinition(pythonTaskName, "Python", parserPython)
}

func parserPython(config *taskapi.TaskConfig, task *taskapi.Task) error {
	version, err := config.GetStringPropertyAllowSingle("version")
	if err != nil {
		return err
	}

	task.Header = version
	task.feature = features.NewFeatureInfo("python", version)

	parserPythonInstallPyenv(task, version)
	parserPythonInstallPythonVersion(task, version)
	parserPythonInstallVirtualenv(task, version)
	parserPythonCreateVirtualenv(task, version)
	return nil
}

func parserPythonInstallPyenv(task *taskapi.Task, version string) {
	needed := func(ctx *taskapi.Context) *taskapi.ActionResult {
		_, err := helpers.NewPyEnv()
		if err != nil {
			return actionNeeded("Pyenv is not installed: %s", err)
		}
		return actionNotNeeded()
	}
	run := func(ctx *taskapi.Context) error {
		result := command(ctx, "brew", "install", "pyenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install pyenv: %s", result.Error)
		}
		return nil
	}
	taskapi.AddActionWithBuilder("install PyEnv", run).OnFunc(needed)
}

func parserPythonInstallPythonVersion(task *taskapi.Task, version string) {
	needed := func(ctx *taskapi.Context) *taskapi.ActionResult {
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
	}
	run := func(ctx *taskapi.Context) error {
		result := command(ctx, "pyenv", "install", version).Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install the required python version: %s", result.Error)
		}
		return nil
	}
	taskapi.AddActionWithBuilder("install Python version with PyEnv", run).OnFunc(needed)
}

func parserPythonInstallVirtualenv(task *taskapi.Task, version string) {
	needed := func(ctx *taskapi.Context) *taskapi.ActionResult {
		pyEnv, err := helpers.NewPyEnv()
		if err != nil {
			return actionFailed("cannot use pyenv: %s", err)
		}
		installed := utils.PathExists(pyEnv.Which(version, "virtualenv"))
		if !installed {
			return actionNeeded("virtualenv is not installed")
		}
		return actionNotNeeded()
	}
	run := func(ctx *taskapi.Context) error {
		pyEnv, err := helpers.NewPyEnv()
		if err != nil {
			return err
		}
		result := command(ctx, pyEnv.Which(version, "python"), "-m", "pip", "install", "virtualenv").Run()
		if result.Error != nil {
			return fmt.Errorf("failed to install virtualenv: %s", result.Error)
		}
		return nil
	}
	taskapi.AddActionWithBuilder("install virtualenv", run).OnFunc(needed)
}

func parserPythonCreateVirtualenv(task *taskapi.Task, version string) {
	needed := func(ctx *taskapi.Context) *taskapi.ActionResult {
		name := helpers.VirtualenvName(ctx.proj, version)
		venv := helpers.NewVirtualenv(ctx.cfg, name)
		if !venv.Exists() {
			return actionNeeded("project virtualenv does not exists")
		}
		return actionNotNeeded()
	}
	run := func(ctx *taskapi.Context) error {
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
	}
	taskapi.AddActionWithBuilder("create virtualenv", run).OnFunc(needed)
}
