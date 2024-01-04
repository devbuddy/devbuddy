package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/manifest"
)

func init() {
	register.Register(env{})
}

type env struct{}

func (env) Name() string {
	return "env"
}

func (env) Activate(ctx *context.Context, param string) (bool, error) {
	man, err := manifest.Load(ctx.Project.Path)
	if err != nil {
		return false, err
	}

	for name, value := range man.Env {
		ctx.Env.Set(name, value)
	}

	return false, nil
}

func (env) Deactivate(ctx *context.Context, param string) {}
