package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(golang{})
}

type golang struct{}

func (golang) Name() string {
	return "golang"
}

func (golang) Activate(ctx *context.Context, param string) (bool, error) {
	golang := helpers.NewGolang(ctx, param)

	if !golang.Exists() {
		return true, nil
	}

	ctx.Env.PrependToPath(golang.BinPath())

	ctx.Env.Set("GOROOT", golang.Path())

	return false, nil
}

func (golang) Deactivate(ctx *context.Context, param string) {}
