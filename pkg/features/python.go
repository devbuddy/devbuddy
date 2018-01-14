package features

import (
	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

func init() {
	allFeatures["python"] = NewPython
}

type Python struct {
	name string
}

func NewPython(param string) Feature {
	return &Python{name: param}
}

func (p *Python) Enable(cfg *config.Config, proj *project.Project, env *Env, ui *termui.HookUI) error {
	venv := helpers.NewVirtualenv(cfg, p.name)

	if !venv.Exists() {
		return DevUpNeeded
	}

	p.cleanPath(cfg, env)
	env.PrependToPath(venv.BinPath())

	env.Set("VIRTUAL_ENV", venv.Path())

	return nil
}

func (p *Python) Disable(cfg *config.Config, env *Env, ui *termui.HookUI) {
	env.Unset("VIRTUAL_ENV")

	p.cleanPath(cfg, env)
}

// cleanPath removes all virtualenv path, even if multiple of them exists
func (p *Python) cleanPath(cfg *config.Config, env *Env) {
	virtualenvBasePath := helpers.NewVirtualenv(cfg, "").Path()
	env.RemoveFromPath(virtualenvBasePath)
}
