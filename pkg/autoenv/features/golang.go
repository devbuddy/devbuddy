package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("golang", golangActivate, golangDeactivate)
}

func golangActivate(ctx *context.Context, version string) (bool, error) {
	golang := helpers.NewGolang(ctx.Cfg, version)

	if !golang.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(golang.BinPath())

	ctx.Env.Set("GOROOT", golang.Path())

	// TODO: decide whether we want to enable GO15VENDOREXPERIMENT
	// Introduced in 1.5, enabled by default in 1.7

	return false, nil
}

func golangDeactivate(ctx *context.Context, version string) {
	// Golang install without version to get the base path
	golang := helpers.NewGolang(ctx.Cfg, "")
	ctx.Env.RemoveFromPath(golang.Path())

	ctx.Env.Unset("GOROOT")
}
