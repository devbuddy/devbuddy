package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("python", pythonActivate, pythonDeactivate)
}

func pythonActivate(ctx *context.Context, version string) (bool, error) {
	name := helpers.VirtualenvName(ctx.Project, version)
	venv := helpers.NewVirtualenv(ctx.Cfg, name)

	if !venv.Exists() {
		return true, nil
	}

	pythonCleanPath(ctx)
	ctx.Env.PrependToPath(venv.BinPath())

	ctx.Env.Set("VIRTUAL_ENV", venv.Path())

	return false, nil
}

func pythonDeactivate(ctx *context.Context, version string) {
	ctx.Env.Unset("VIRTUAL_ENV")

	pythonCleanPath(ctx)
}

// pythonCleanPath removes all virtualenv path, even if multiple of them exists
func pythonCleanPath(ctx *context.Context) {
	virtualenvBasePath := helpers.NewVirtualenv(ctx.Cfg, "").Path()
	ctx.Env.RemoveFromPath(virtualenvBasePath)
}
