package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("custom", "Custom", parserCustom)
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

	task.Info = name

	runCommand := func(ctx *context.Context) error {
		result := shell(ctx, command).Run()
		return result.Error
	}

	runNeeded := func(ctx *context.Context) *taskapi.ActionResult {
		result := shellSilent(ctx, condition).Capture()
		if result.LaunchError != nil {
			return taskapi.ActionFailed("failed to run the condition command: %s", result.LaunchError)
		}
		if result.Code != 0 {
			return taskapi.ActionNeeded("the met? command exited with a non-zero code")
		}
		return taskapi.ActionNotNeeded()
	}

	task.AddActionBuilder("", runCommand).OnFunc(runNeeded)

	return nil
}
