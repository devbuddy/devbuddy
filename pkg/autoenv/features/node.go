package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(node{})
}

type node struct{}

func (node) Name() string {
	return "node"
}

func (node) Activate(ctx *context.Context, param string) (bool, error) {
	node := helpers.NewNode(ctx.Cfg, param)
	if !node.Exists() {
		return true, nil
	}
	ctx.Env.PrependToPath(node.BinPath())
	return false, nil
}

func (node) Deactivate(ctx *context.Context, param string) {}
