package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/manifest"
)

func init() {
	register("env", envActivate, nil)
}

func envActivate(ctx *context.Context, param string) (bool, error) {
	man, err := manifest.Load(ctx.Project.Path)
	if err != nil {
		return false, err
	}

	for name, value := range man.Env {
		ctx.Env.Set(name, value)
	}

	return false, nil
}
