package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/api"
)

func init() {
	api.Register("pip", "Pip", parserPip).SetRequiredTask(pythonTaskName)
}

func parserPip(config *api.TaskConfig, task *api.Task) error {
	files, err := config.GetListOfStrings()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no pip files specified")
	}

	task.Info = strings.Join(files, ", ")

	for _, file := range files {
		pipInstall := func(ctx *context.Context) error {
			pipArgs := []string{"install", "--require-virtualenv", "-r", file}
			result := command(ctx, "pip", pipArgs...).AddOutputFilter("already satisfied").Run()
			if result.Error != nil {
				return fmt.Errorf("Pip failed: %w", result.Error)
			}
			return nil
		}
		task.AddActionBuilder(fmt.Sprintf("install %s", file), pipInstall).
			On(api.FileCondition(file))
	}

	return nil
}
