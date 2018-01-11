package tasks

import (
	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

type Context struct {
	proj     *project.Project
	ui       *termui.UI
	cfg      *config.Config
	features map[string]string
}

func RunAll(cfg *config.Config, proj *project.Project, ui *termui.UI) error {
	taskList, err := GetTasksFromProject(proj)
	if err != nil {
		return err
	}

	features := GetFeaturesFromTasks(proj, taskList)

	ctx := Context{cfg: cfg, proj: proj, ui: ui, features: features}

	for _, task := range taskList {
		err = task.Perform(&ctx)
		if err != nil {
			break
		}
	}

	return nil
}
