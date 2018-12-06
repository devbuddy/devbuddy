package tasks

import (
	"fmt"
	"os"

	"github.com/devbuddy/devbuddy/pkg/features"
)

// RunAll accepts a list of tasks and check for their requirements and runs them if their conditions are met
func RunAll(ctx *Context, taskList []*Task) (success bool, err error) {
	for _, task := range taskList {
		if task.requiredTask != "" {
			if _, present := ctx.features[task.requiredTask]; !present {
				err = fmt.Errorf("You must specify a %s environment to use a %s task", task.requiredTask, task.name)
				ctx.ui.TaskError(err)
				return false, nil
			}
		}
	}

	for _, task := range taskList {
		ctx.ui.TaskHeader(task.name, task.header)
		err = runTask(ctx, task)
		if err != nil {
			ctx.ui.TaskError(err)
			return false, nil
		}
	}

	return true, nil
}

func runTask(ctx *Context, task *Task) (err error) {
	if task.perform != nil {
		err = task.perform(ctx)
		if err != nil {
			return err
		}
	}

	for _, action := range task.actions {
		err = runAction(ctx, action)
		if err != nil {
			return err
		}
	}

	err = activateFeature(ctx, task)
	return err
}

func activateFeature(ctx *Context, task *Task) (err error) {
	if task.featureName == "" {
		return nil
	}

	devUpNeeded, err := features.Activate(task.featureName, task.featureParam, ctx.cfg, ctx.proj, ctx.env)
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
