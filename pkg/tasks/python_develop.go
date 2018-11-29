package tasks

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/helpers/store"
)

func init() {
	t := registerTaskDefinition("python_develop")
	t.name = "Python develop"
	t.requiredTask = pythonTaskName
	t.parser = parserPythonDevelop
}

func parserPythonDevelop(config *taskConfig, task *Task) error {
	action := newAction("install python package in develop mode", func(ctx *context) error {
		result := command(ctx, "pip", "install", "--require-virtualenv", "-e", ".").
			AddOutputFilter("already satisfied").Run()

		if result.Error != nil {
			return fmt.Errorf("Pip failed: %s", result.Error)
		}

		return store.New(ctx.proj.Path).RecordFileChange("setup.py")
	})
	action.onFunc(func(ctx *context) *actionResult {
		changed, err := store.New(ctx.proj.Path).HasFileChanged("setup.py")
		if err != nil {
			return actionFailed("failed to check if setup.py has changed: %s", err)
		}
		if changed {
			return actionNeeded("setup.py was modified")
		}
		return actionNotNeeded()
	})
	task.addAction(action)
	return nil
}
