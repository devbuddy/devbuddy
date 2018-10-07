package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/devbuddy/devbuddy/pkg/helpers"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func init() {
	f := definitions.Register("dotenv")
	f.Activate = dotenvRefresh
	f.Refresh = dotenvRefresh
}

func dotenvRefresh(path string, cfg *config.Config, proj *project.Project, env *env.Env) (bool, error) {
	vars, err := helpers.DotenvRead(path)
	if err != nil {
		return false, err
	}
	for name, value := range vars {
		env.Set(name, value)
	}
	return false, nil
}
