package context

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

// Context is the interface to the calling execution context
type Context struct {
	Cfg     *config.Config
	Project *project.Project
	Env     *env.Env
	UI      *termui.UI
}

func Load() (*Context, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	proj, err := project.FindCurrent()
	if err != nil {
		return nil, err
	}

	ctx := &Context{
		Cfg:     cfg,
		Project: proj,
		UI:      termui.New(cfg),
		Env:     env.NewFromOS(),
	}

	return ctx, nil
}
