package tasks

import (
	"fmt"
	"os"

	"github.com/devbuddy/devbuddy/pkg/features"
)

// Run accepts a list of tasks and check for their requirements and runs them if their conditions are met
func Run(ctx *Context, executor TaskRunner, selector TaskSelector, taskList []*Task) (success bool, err error) {
	for _, task := range taskList {
		if task.requiredTask != "" {
			if _, present := ctx.features[task.requiredTask]; !present {
				ctx.ui.TaskErrorf("You must specify a %s environment to use a %s task", task.requiredTask, task.name)
				return false, nil
			}
		}
	}

	for _, task := range taskList {
		shouldRun, err := selector.ShouldRun(ctx, task)
		if err != nil {
			ctx.ui.TaskError(err)
			return false, nil
		}
		if !shouldRun {
			ctx.ui.TaskHeader(task.name, task.header, "disabled")
			continue
		}

		ctx.ui.TaskHeader(task.name, task.header, "")
		err = executor.Run(ctx, task)
		if err != nil {
			ctx.ui.TaskError(err)
			return false, nil
		}
	}

	return true, nil
}

type TaskRunner interface {
	Run(*Context, *Task) error
}

type TaskRunnerImpl struct{}

func (r *TaskRunnerImpl) Run(ctx *Context, task *Task) (err error) {
	for _, action := range task.actions {
		err = runAction(ctx, action)
		if err != nil {
			return err
		}
	}

	err = r.activateFeature(ctx, task)
	return err
}

func (r *TaskRunnerImpl) activateFeature(ctx *Context, task *Task) error {
	if task.featureName == "" {
		return nil
	}

	def, err := features.Get(task.featureName)
	if err != nil {
		return err
	}

	devUpNeeded, err := def.Activate(task.featureParam, ctx.cfg, ctx.proj, ctx.env)
	if err != nil {
		return err
	}
	if devUpNeeded {
		ctx.ui.TaskWarning(fmt.Sprintf("Something is wrong, the feature %s could not be activated", task.featureName))
	}

	// Special case, we want the bud process to get PATH updates from features to call the right processes.
	// Like the pip process from the newly activated virtualenv.
	// Explanation: exec.Command calls exec.LookPath to find the executable path, which rely on the PATH of
	// the process itself.
	// There is no problem when executing a shell command since the shell process will do the executable lookup
	// itself with the PATH value from the specified Env.
	return os.Setenv("PATH", ctx.env.Get("PATH"))
}
