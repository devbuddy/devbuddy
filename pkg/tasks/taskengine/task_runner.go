package taskengine

import (
	"fmt"
	"os"

	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/tasks/taskapi"
)

// Run accepts a list of tasks and check for their requirements and runs them if their conditions are met
func Run(ctx *taskapi.Context, executor TaskRunner, selector TaskSelector, taskList []*taskapi.Task) (bool, error) {
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

type TaskRunner interface {
	Run(*taskapi.Context, *taskapi.Task) error
}

type TaskRunnerImpl struct{}

func (r *TaskRunnerImpl) Run(ctx *taskapi.Context, task *taskapi.Task) (err error) {
	for _, action := range task.Actions {
		err = runAction(ctx, action)
		if err != nil {
			return err
		}
	}

	err = r.activateFeature(ctx, task)
	return err
}

func (r *TaskRunnerImpl) activateFeature(ctx *taskapi.Context, task *taskapi.Task) error {
	if task.Feature.Name == "" {
		return nil
	}

	def, err := features.Get(task.Feature)
	if err != nil {
		return err
	}

	devUpNeeded, err := def.Activate(task.Feature.Param, ctx.Cfg, ctx.Project, ctx.Env)
	if err != nil {
		return err
	}
	if devUpNeeded {
		ctx.UI.TaskWarning(fmt.Sprintf("Something is wrong, the feature %s could not be activated", task.Feature))
	}

	// Special case, we want the bud process to get PATH updates from features to call the right processes.
	// Like the pip process from the newly activated virtualenv.
	// Explanation: exec.Command calls exec.LookPath to find the executable path, which rely on the PATH of
	// the process itself.
	// There is no problem when executing a shell command since the shell process will do the executable lookup
	// itself with the PATH value from the specified Env.
	return os.Setenv("PATH", ctx.Env.Get("PATH"))
}

func runAction(ctx *taskapi.Context, action taskapi.TaskAction) error {
	desc := action.Description()

	result := action.Needed(ctx)
	if result.Error != nil {
		return fmt.Errorf("The task action (%s) failed to detect whether it need to run: %s", desc, result.Error)
	}

	if result.Needed {
		if desc != "" {
			ctx.UI.TaskActionHeader(desc)
		}
		ctx.UI.Debug("Reason: \"%s\"", result.Reason)

		err := action.Run(ctx)
		if err != nil {
			return fmt.Errorf("The task action failed to run: %s", err)
		}

		result = action.Needed(ctx)
		if result.Error != nil {
			return fmt.Errorf("The task action failed to detect if it is resolved: %s", result.Error)
		}

		if result.Needed {
			return fmt.Errorf("The task action did not produce the expected result: %s", result.Reason)
		}
	}

	return nil
}
