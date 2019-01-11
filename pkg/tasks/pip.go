package tasks

import (
	"fmt"
	"strings"

	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

func init() {
	taskapi.RegisterTaskDefinition("pip", "Pip", parserPip).AddRequiredTask(pythonTaskName)
}

func parserPip(config *taskapi.TaskConfig, task *taskapi.Task) error {
	files, err := config.GetListOfStrings()
	if err != nil {
		return err
	}
	if len(files) == 0 {
		return fmt.Errorf("no pip files specified")
	}

	task.Header = strings.Join(files, ", ")

	for _, file := range files {
		builder := actionBuilder(fmt.Sprintf("install %s", file), func(ctx *taskapi.Context) error {
			pipArgs := []string{"install", "--require-virtualenv", "-r", file}
			result := command(ctx, "pip", pipArgs...).AddOutputFilter("already satisfied").Run()
			if result.Error != nil {
				return fmt.Errorf("Pip failed: %s", result.Error)
			}
			return nil
		})
		builder.OnFileChange(file)
		task.AddAction(builder.Build())
	}

	return nil
}
