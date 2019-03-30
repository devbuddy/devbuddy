package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("envfile", envfileActivate, nil)
}

func envfileActivate(ctx *context.Context, param string) (bool, error) {
	envfilePath := param

	err := helpers.LoadEnvfile(ctx.Env, envfilePath)
	if err != nil {
		return true, err
	}

	return false, nil
}
