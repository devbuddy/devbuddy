package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.Register("python_develop", "Python develop", parserPythonDevelop).SetRequiredTask(pythonTaskName)
}

func parserPythonDevelop(config *taskapi.TaskConfig, task *taskapi.Task) error {
	extras, err := config.GetListOfStringsPropertyDefault("extras", []string{})
	if err != nil {
		return err
	}

	pipTarget := "."
	if len(extras) > 0 {
		pipTarget = fmt.Sprintf(".[%s]", strings.Join(extras, ","))
	}
	pipArgs := []string{"install", "--require-virtualenv", "-e", pipTarget}

	pipInstall := func(ctx *context.Context) error {
		result := command(ctx, "pip", pipArgs...).AddOutputFilter("already satisfied").Run()
		if result.Error != nil {
			return fmt.Errorf("Pip failed: %w", result.Error)
		}

		return nil
	}
	task.AddActionBuilder("install python package in develop mode", pipInstall).
		OnFileChange("setup.py").
		OnFileChange("setup.cfg")

	return nil
}
