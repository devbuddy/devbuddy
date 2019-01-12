package tasks

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

// Context provides the context in which a Task run
type Context struct {
	Project  *project.Project
	UI       *termui.UI
	Cfg      *config.Config
	Env      *env.Env
	Features features.FeatureSet
}

// NewContext returns a *Context for a project, using the environment
func NewContext(cfg *config.Config, proj *project.Project, ui *termui.UI, taskList []*Task) *Context {
	return &Context{
		Cfg:      cfg,
		Project:  proj,
		UI:       ui,
		Env:      env.NewFromOS(),
		Features: GetFeaturesFromTasks(taskList),
	}
}
