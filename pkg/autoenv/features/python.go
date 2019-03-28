package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("python", pythonActivate, nil)
}

func pythonActivate(ctx *context.Context, version string) (bool, error) {
	name := helpers.VirtualenvName(ctx.Project, version)
	venv := helpers.NewVirtualenv(ctx.Cfg, name)

	if !venv.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(venv.BinPath())
	ctx.Env.Set("VIRTUAL_ENV", venv.Path())

	return false, nil
}
