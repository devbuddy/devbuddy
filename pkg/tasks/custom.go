package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("custom", "Custom", parserCustom)
}

func parserCustom(config *api.TaskConfig, task *api.Task) error {
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

	runNeeded := func(ctx *context.Context) *api.ActionResult {
		result := shellSilent(ctx, condition).Capture()
		if result.LaunchError != nil {
			return api.Failed("failed to run the condition command: %s", result.LaunchError)
		}
		if result.Code != 0 {
			return api.Needed("the met? command exited with a non-zero code")
		}
		return api.NotNeeded()
	}

	task.AddActionBuilder("", runCommand).On(api.FuncCondition(runNeeded))

	return nil
}
