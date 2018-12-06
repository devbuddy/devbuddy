package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

type Context struct {
	proj     *project.Project
	ui       *termui.UI
	cfg      *config.Config
	env      *env.Env
	features map[string]string
}

func NewContextFromEnv(cfg *config.Config, proj *project.Project, ui *termui.UI, taskList []*Task) *Context {
	return &Context{
		cfg:      cfg,
		proj:     proj,
		ui:       ui,
		env:      env.NewFromOS(),
		features: GetFeaturesFromTasks(proj, taskList),
	}
}
