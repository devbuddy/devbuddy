package tasks

import (
	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

type Context struct {
	proj *project.Project
	ui   *termui.UI
	cfg  *config.Config
}

func RunAll(cfg *config.Config, proj *project.Project, ui *termui.UI) error {
	ctx := Context{cfg: cfg, proj: proj, ui: ui}

	taskList, err := GetTasksFromProject(ctx.proj)
	if err != nil {
		return err
	}

	for _, task := range taskList {
		err = task.Perform(&ctx)
		if err != nil {
			break
		}
	}

	return nil
}
