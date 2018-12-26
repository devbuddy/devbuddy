package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func init() {
	f := definitions.Register("node")
	f.Activate = nodeActivate
	f.Deactivate = nodeDeactivate
}

func nodeActivate(version string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
	node := helpers.NewNode(cfg, version)
	if !node.Exists() {
		return true, nil
	}
	env.PrependToPath(node.BinPath())
	return false, nil
}

func nodeDeactivate(version string, cfg *config.Config, env *env.Env) {
	node := helpers.NewNode(cfg, version)
	env.RemoveFromPath(node.Path())
}
