package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/utils"
)

func init() {
	register.Register(ruby{})
}

type ruby struct{}

func (ruby) Name() string {
	return "ruby"
}

func (ruby) Activate(ctx *context.Context, param string) (bool, error) {
	rbEnv, err := helpers.NewRbEnv(ctx)
	if err != nil {
		return true, nil
	}
	binPath := rbEnv.BinPath(param)
	if !utils.PathExists(binPath) {
		return true, nil
	}
	ctx.Env.PrependToPath(binPath)
	return false, nil
}

func (ruby) Deactivate(ctx *context.Context, param string) {}
