package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(Envfile{})
}

type Envfile struct{}

func (Envfile) Name() string {
	return "envfile"
}

func (Envfile) Activate(ctx *context.Context, param string) (bool, error) {
	envfilePath := param

	err := helpers.LoadEnvfile(ctx.Env, envfilePath)
	if err != nil {
		return true, err
	}

	return false, nil
}

func (Envfile) Deactivate(ctx *context.Context, param string) {}
