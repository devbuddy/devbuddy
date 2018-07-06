package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func init() {
	allFeatures["python"] = newPython
}

type Python struct {
	version string
}

func newPython(param string) Feature {
	return &Python{version: param}
}

func (p *Python) Activate(cfg *config.Config, proj *project.Project, env *env.Env) error {
	name := helpers.VirtualenvName(proj, p.version)
	venv := helpers.NewVirtualenv(cfg, name)

	if !venv.Exists() {
		return DevUpNeeded
	}

	p.cleanPath(cfg, env)
	env.PrependToPath(venv.BinPath())

	env.Set("VIRTUAL_ENV", venv.Path())

	return nil
}

func (p *Python) Deactivate(cfg *config.Config, env *env.Env) {
	env.Unset("VIRTUAL_ENV")

	p.cleanPath(cfg, env)
}

// cleanPath removes all virtualenv path, even if multiple of them exists
func (p *Python) cleanPath(cfg *config.Config, env *env.Env) {
	virtualenvBasePath := helpers.NewVirtualenv(cfg, "").Path()
	env.RemoveFromPath(virtualenvBasePath)
}
