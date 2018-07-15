package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func init() {
	f := definitions.Register("golang")
	f.Activate = golangActivate
	f.Deactivate = golangDeactivate
}

func golangActivate(version string, cfg *config.Config, proj *project.Project, env *env.Env) error {
	golang := helpers.NewGolang(cfg, version)

	if !golang.Exists() {
		return DevUpNeeded
	}

	env.PrependToPath(golang.BinPath())

	env.Set("GOROOT", golang.Path())

	// TODO: decide whether we want to enable GO15VENDOREXPERIMENT
	// Introduced in 1.5, enabled by default in 1.7

	return nil
}

func golangDeactivate(version string, cfg *config.Config, env *env.Env) {
	// Golang install without version to get the base path
	golang := helpers.NewGolang(cfg, "")
	env.RemoveFromPath(golang.Path())

	env.Unset("GOROOT")
}
