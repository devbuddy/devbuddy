package taskengine

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

// Run accepts a list of tasks and check for their requirements and runs them if their conditions are met
func Run(ctx *context.Context, executor TaskRunner, selector TaskSelector, taskList []*taskapi.Task) (bool, error) {
	err := checkRequiredTasks(taskList)
	if err != nil {
		return false, err
	}

	for _, task := range taskList {
		shouldRun, err := selector.ShouldRun(ctx, task)
		if err != nil {
			ctx.UI.TaskError(err)
			return false, nil
		}
		if !shouldRun {
			ctx.UI.TaskHeader(task.Name, task.Info, "disabled")
			continue
		}

		ctx.UI.TaskHeader(task.Name, task.Info, "")
		err = executor.Run(ctx, task)
		if err != nil {
			ctx.UI.TaskError(err)
			return false, nil
		}
	}

	return true, nil
}

func checkRequiredTasks(taskList []*taskapi.Task) error {
	seen := map[string]bool{}
	for _, task := range taskList {
		if task.RequiredTask != "" && !seen[task.RequiredTask] {
			return fmt.Errorf("You must specify a %s task before a %s task", task.RequiredTask, task.Key)
		}
		seen[task.Key] = true
	}
	return nil
}
