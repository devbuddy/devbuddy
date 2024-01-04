package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(python{})
}

type python struct{}

func (python) Name() string {
	return "python"
}

func (python) Activate(ctx *context.Context, param string) (bool, error) {
	name := helpers.VirtualenvName(ctx.Project, param)
	venv := helpers.NewVirtualenv(ctx.Cfg, name)

	if !venv.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(venv.BinPath())
	ctx.Env.Set("VIRTUAL_ENV", venv.Path())

	return false, nil
}

func (python) Deactivate(ctx *context.Context, param string) {}
