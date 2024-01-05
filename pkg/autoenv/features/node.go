package features

import (
	"github.com/devbuddy/devbuddy/pkg/context"
	"github.com/devbuddy/devbuddy/pkg/helpers"
)

func init() {
	register.Register(Node{})
}

type Node struct{}

func (Node) Name() string {
	return "node"
}

func (Node) Activate(ctx *context.Context, param string) (bool, error) {
	node := helpers.NewNode(ctx.Cfg, param)
	if !node.Exists() {
		return true, nil
	}
	ctx.Env.PrependToPath(node.BinPath())
	return false, nil
}

func (Node) Deactivate(ctx *context.Context, param string) {}
