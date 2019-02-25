package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func init() {
	register("python", pythonActivate, pythonDeactivate)
}

func pythonActivate(version string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
	name := helpers.VirtualenvName(proj, version)
	venv := helpers.NewVirtualenv(cfg, name)

	if !venv.Exists() {
		return true, nil
	}

	pythonCleanPath(cfg, env)
	env.PrependToPath(venv.BinPath())

	env.Set("VIRTUAL_ENV", venv.Path())

	return false, nil
}

func pythonDeactivate(version string, cfg *config.Config, env *env.Env) {
	env.Unset("VIRTUAL_ENV")

	pythonCleanPath(cfg, env)
}

// pythonCleanPath removes all virtualenv path, even if multiple of them exists
func pythonCleanPath(cfg *config.Config, env *env.Env) {
	virtualenvBasePath := helpers.NewVirtualenv(cfg, "").Path()
	env.RemoveFromPath(virtualenvBasePath)
}
