package tasks

import (
	"fmt"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/env"
	"github.com/pior/dad/pkg/executor"
	"github.com/pior/dad/pkg/features"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

type Context struct {
	proj     *project.Project
	ui       *termui.UI
	cfg      *config.Config
	env      *env.Env
	features map[string]string
}

func RunAll(cfg *config.Config, proj *project.Project, ui *termui.UI) error {
	taskList, err := GetTasksFromProject(proj)
	if err != nil {
		return err
	}

	ctx := &Context{
		cfg:      cfg,
		proj:     proj,
		ui:       ui,
		env:      env.NewFromOS(),
		features: GetFeaturesFromTasks(proj, taskList),
	}

	for _, task := range taskList {
		ctx.ui.TaskHeader(task.name(), task.header())

		err = task.perform(ctx)
		if err != nil {
			return err
		}

		err = activateFeature(ctx, task)
		if err != nil {
			ctx.ui.TaskError(err)
			return err
		}
	}

	return nil
}

func activateFeature(ctx *Context, task Task) (err error) {
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
