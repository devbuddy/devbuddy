package features

import (
	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/env"
	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

func init() {
	allFeatures["python"] = newPython
}

type Python struct {
	name string
}

func newPython(param string) Feature {
	return &Python{name: param}
}

func (p *Python) activate(cfg *config.Config, proj *project.Project, env *env.Env) error {
	venv := helpers.NewVirtualenv(cfg, p.name)

	if !venv.Exists() {
		return DevUpNeeded
	}

	p.cleanPath(cfg, env)
	env.PrependToPath(venv.BinPath())

	env.Set("VIRTUAL_ENV", venv.Path())

	return nil
}

func (p *Python) deactivate(cfg *config.Config, env *env.Env) {
	env.Unset("VIRTUAL_ENV")

	p.cleanPath(cfg, env)
}

// cleanPath removes all virtualenv path, even if multiple of them exists
func (p *Python) cleanPath(cfg *config.Config, env *env.Env) {
	virtualenvBasePath := helpers.NewVirtualenv(cfg, "").Path()
	env.RemoveFromPath(virtualenvBasePath)
}
