package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(Python{})
}

type Python struct{}

func (Python) Name() string {
	return "python"
}

func (Python) Activate(ctx *context.Context, param string) (bool, error) {
	name := helpers.VirtualenvName(ctx.Project, param)
	venv := helpers.NewVirtualenv(ctx.Cfg, name)

	if !venv.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(venv.BinPath())
	ctx.Env.Set("VIRTUAL_ENV", venv.Path())

	return false, nil
}

func (Python) Deactivate(ctx *context.Context, param string) {}
