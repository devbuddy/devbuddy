package tasks

import (
	"fmt"
	"os"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

type context struct {
	proj     *project.Project
	ui       *termui.UI
	cfg      *config.Config
	env      *env.Env
	features map[string]string
}

// RunAll builds and execute all tasks found in the project
func RunAll(cfg *config.Config, proj *project.Project, ui *termui.UI) (success bool, err error) {
	taskList, err := GetTasksFromProject(proj)
	if err != nil {
		return false, err
	}

	ctx := &context{
		cfg:      cfg,
		proj:     proj,
		ui:       ui,
		env:      env.NewFromOS(),
		features: GetFeaturesFromTasks(proj, taskList),
	}

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
			return false, err
		}
	}

	return true, nil
}

func runTask(ctx *context, task *Task) (err error) {
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

func runAction(ctx *context, action taskAction) error {
	desc := action.description()

	needed, err := action.needed(ctx)
	if err != nil {
		return fmt.Errorf("The task action (%s) failed to detect whether it need to run: %s", desc, err)
	}

	if needed {
		if desc != "" {
			ctx.ui.TaskActionHeader(desc)
		}

		err = action.run(ctx)
		if err != nil {
			return fmt.Errorf("The task action failed to run: %s", err)
		}
	}

	stillNeeded, err := action.needed(ctx)
	if err != nil {
		return fmt.Errorf("The task action failed to detect if it is resolved: %s", err)
	}

	if stillNeeded {
		return fmt.Errorf("The task action is not resolved after running it")
	}

	return nil
}

func activateFeature(ctx *context, task *Task) (err error) {
	if task.featureName == "" {
		return nil
	}

	err = features.Activate(task.featureName, task.featureParam, ctx.cfg, ctx.proj, ctx.env)
	if err != nil {
		if err == features.DevUpNeeded {
			ctx.ui.TaskWarning(
				fmt.Sprintf("Something is wrong, the feature %s could not be activated", task.featureName))
		} else {
			return err
		}
	}

	// Special case, we want the bud process to get PATH updates from features to call the right processes.
	// Like the pip process from the newly activated virtualenv.
	// Explanation: exec.Command calls exec.LookPath to find the executable path, which rely on the PATH of
	// the process itself.
	// There is no problem when executing a shell command since the shell process will do the executable lookup
	// itself with the PATH value from the specified Env.
	return os.Setenv("PATH", ctx.env.Get("PATH"))
}
