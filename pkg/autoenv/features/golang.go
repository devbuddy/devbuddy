package features

import (
	"strings"

	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("golang", golangActivate, golangDeactivate)
}

func golangActivate(ctx *context.Context, version string) (bool, error) {
	golang := helpers.NewGolang(ctx.Cfg, golangExtractVersion(version))

	if !golang.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(golang.BinPath())

	ctx.Env.Set("GOROOT", golang.Path())

	if golangVersionWithModules(version) {
		ctx.Env.Set("GO111MODULE", "on")
	}

	return false, nil
}

func golangDeactivate(ctx *context.Context, version string) {
	// Golang install without version to get the base path
	golang := helpers.NewGolang(ctx.Cfg, "")
	ctx.Env.RemoveFromPath(golang.Path())

	ctx.Env.Unset("GO111MODULE")
	ctx.Env.Unset("GOROOT")
}

const golangModulesSuffix = "+mod"

func golangExtractVersion(version string) string {
	return strings.TrimSuffix(version, golangModulesSuffix)
}

func golangVersionWithModules(version string) bool {
	return strings.HasSuffix(version, golangModulesSuffix)
}
