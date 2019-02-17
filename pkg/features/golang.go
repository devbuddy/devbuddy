package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func init() {
	register("golang", golangActivate, golangDeactivate)
}

func golangActivate(version string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
	golang := helpers.NewGolang(cfg, version)

	if !golang.Exists() {
		return true, nil
	}

	env.PrependToPath(golang.BinPath())

	env.Set("GOROOT", golang.Path())

	// TODO: decide whether we want to enable GO15VENDOREXPERIMENT
	// Introduced in 1.5, enabled by default in 1.7
	if utils.PathExists("go.mod") {
		env.Set("GO111MODULE", "on")
	}

	return false, nil
}

func golangDeactivate(version string, cfg *config.Config, env *env.Env) {
	// Golang install without version to get the base path
	golang := helpers.NewGolang(cfg, "")
	env.RemoveFromPath(golang.Path())

	env.Unset("GOROOT")
}
