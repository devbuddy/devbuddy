package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func init() {
	allFeatures["golang"] = newGolang
}

type Golang struct {
	version string
}

func newGolang(param string) Feature {
	return &Golang{version: param}
}

func (g *Golang) Activate(cfg *config.Config, proj *project.Project, env *env.Env) error {
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

func (g *Golang) Deactivate(cfg *config.Config, env *env.Env) {
	// Golang install without version to get the base path
	golang := helpers.NewGolang(cfg, "")
	env.RemoveFromPath(golang.Path())

	env.Unset("GOROOT")
}
