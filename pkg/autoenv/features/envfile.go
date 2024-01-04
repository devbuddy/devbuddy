package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(envfile{})
}

type envfile struct{}

func (envfile) Name() string {
	return "envfile"
}

func (envfile) Activate(ctx *context.Context, param string) (bool, error) {
	envfilePath := param

	err := helpers.LoadEnvfile(ctx.Env, envfilePath)
	if err != nil {
		return true, err
	}

	return false, nil
}

func (envfile) Deactivate(ctx *context.Context, param string) {}
