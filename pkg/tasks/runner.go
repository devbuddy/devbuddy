package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/env"
	"github.com/pior/dad/pkg/features"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
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
		if t, ok := task.(taskWithPreRunValidation); ok {
			err = t.preRunValidation(ctx)
			if err != nil {
				ctx.ui.TaskError(err)
				return false, nil
			}
		}
	}

	for _, task := range taskList {
		ctx.ui.TaskHeader(task.name(), task.header())

		if t, ok := task.(taskWithPerform); ok {
			err = t.perform(ctx)
			if err != nil {
				ctx.ui.TaskError(err)
				return false, nil
			}
		}

		for _, action := range task.actions(ctx) {
			err = runAction(ctx, action)
			if err != nil {
				ctx.ui.TaskError(err)
				return false, nil
			}
		}

		err = activateFeature(ctx, task)
		if err != nil {
			ctx.ui.TaskError(err)
			return false, nil
		}

	}

	return true, nil
}

func runAction(ctx *context, action taskAction) error {
	desc := action.description()

	needed, err := action.needed(ctx)
	if err != nil {
		return fmt.Errorf("The task action (%s) failed to detect whether it need to run: %s", desc, err)
	}

	if desc != "" {
		ctx.ui.TaskActionHeader(desc)
	}

	if needed {
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

func activateFeature(ctx *context, task Task) (err error) {
	t, ok := task.(TaskWithFeature)
	if !ok {
		return nil
	}

	name, param := t.feature(ctx.proj)
	err = features.New(name, param).Activate(ctx.cfg, ctx.proj, ctx.env)
	if err != nil {
		if err == features.DevUpNeeded {
			ctx.ui.TaskWarning(fmt.Sprintf("Something is wrong, the feature %s could not be activated", name))
		} else {
			return err
		}
	}

	return nil
}
