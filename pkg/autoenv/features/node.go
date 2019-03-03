package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register("node", nodeActivate, nodeDeactivate)
}

func nodeActivate(ctx *context.Context, version string) (bool, error) {
	node := helpers.NewNode(ctx.Cfg, version)
	if !node.Exists() {
		return true, nil
	}
	ctx.Env.PrependToPath(node.BinPath())
	return false, nil
}

func nodeDeactivate(ctx *context.Context, version string) {
	node := helpers.NewNode(ctx.Cfg, version)
	ctx.Env.RemoveFromPath(node.Path())
}
