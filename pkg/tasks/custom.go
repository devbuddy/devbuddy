package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.RegisterTaskDefinition("custom", "Custom", parserCustom)
}

func parserCustom(config *taskapi.TaskConfig, task *taskapi.Task) error {
	command, err := config.GetStringProperty("meet")
	if err != nil {
		return err
	}
	condition, err := config.GetStringProperty("met?")
	if err != nil {
		return err
	}
	name, err := config.GetStringPropertyDefault("name", command)
	if err != nil {
		return err
	}

	task.Header = name

	run := func(ctx *taskapi.Context) error {
		result := shell(ctx, command).Run()
		return result.Error
	}
	needed := func(ctx *taskapi.Context) *taskapi.ActionResult {
		result := shellSilent(ctx, condition).Capture()
		if result.LaunchError != nil {
			return taskapi.ActionFailed("failed to run the condition command: %s", result.LaunchError)
		}
		if result.Code != 0 {
			return taskapi.ActionNeeded("the met? command exited with a non-zero code")
		}
		return taskapi.ActionNotNeeded()
	}
	task.AddActionWithBuilder("", run).OnFunc(needed)

	return nil
}
