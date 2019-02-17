package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"

	"github.com/joho/godotenv"
)

func init() {
	register("envfile", envfileActivate, envfileDeactivate)
}

func envfileActivate(version string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
	loadedEnv, err := godotenv.Read(".env")
	if err != nil {
		return true, err
	}

	for name, value := range loadedEnv {
		if !env.Has(name) {
			env.Set(name, value)
		}
	}

	return false, nil
}

func envfileDeactivate(version string, cfg *config.Config, env *env.Env) {
}
