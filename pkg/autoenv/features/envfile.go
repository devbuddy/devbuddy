package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"

	"github.com/joho/godotenv"
)

func init() {
	register("envfile", envfileActivate, nil)
}

func envfileActivate(ctx *context.Context, param string) (bool, error) {
	envfilePath := param

	loadedEnv, err := godotenv.Read(envfilePath)
	if err != nil {
		return true, err
	}

	for name, value := range loadedEnv {
		ctx.Env.Set(name, value)
	}

	return false, nil
}
