package context

import (
	"errors"

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

// Load returns a Context, even if the project was not found
func Load(hookMode bool) (*Context, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}

	ui := termui.New(cfg)
	if hookMode {
		ui.SetOutputToStderr()
	}

	proj, err := project.FindCurrent()
	if err != nil {
		if errors.Is(err, project.ErrProjectNotFound) {
			return nil, err
		}
		ui.Debug("Project not found")
	} else {
		ui.Debug("Project found in %s", proj.Path)
	}

	ctx := &Context{
		Cfg:     cfg,
		Project: proj,
		UI:      ui,
		Env:     env.NewFromOS(),
	}

	return ctx, nil
}

// LoadWithProject returns a Context, fails if the project is not found
func LoadWithProject() (*Context, error) {
	ctx, err := Load(false)
	if err != nil {
		return nil, err
	}
	if ctx.Project == nil {
		return nil, project.ErrProjectNotFound
	}
	return ctx, nil
}
