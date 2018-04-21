package features

import (
	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/helpers"
	"github.com/pior/dad/pkg/project"
)

func init() {
	allFeatures["golang"] = NewGolang
}

type Golang struct {
	version string
}

func NewGolang(param string) Feature {
	return &Golang{version: param}
}

func (g *Golang) activate(cfg *config.Config, proj *project.Project, env *Env) error {
	golang := helpers.NewGolang(cfg, g.version)

	if !golang.Exists() {
		return DevUpNeeded
	}

	env.PrependToPath(golang.BinPath())

	env.Set("GOROOT", golang.Path())

	// TODO: decide whether we want to enable GO15VENDOREXPERIMENT
	// Introduced in 1.5, enabled by default in 1.7

	return nil
}

func (g *Golang) deactivate(cfg *config.Config, env *Env) {
	// Golang install without version to get the base path
	golang := helpers.NewGolang(cfg, "")
	env.RemoveFromPath(golang.Path())

	env.Unset("GOROOT")
}
