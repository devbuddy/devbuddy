package features

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/definitions"
	"github.com/devbuddy/devbuddy/pkg/project"
)

func Activate(name string, param string, conf *config.Config, proj *project.Project, env *env.Env) (bool, error) {
	def := definitions.Get(name)
	if def == nil {
		panic(fmt.Sprintf("unknown feature: %s", name))
	}
	if def.Activate != nil {
		return def.Activate(param, conf, proj, env)
	}
	return false, nil
}

func Deactivate(name string, param string, conf *config.Config, env *env.Env) {
	def := definitions.Get(name)
	if def == nil {
		panic(fmt.Sprintf("unknown feature: %s", name))
	}
	if def.Deactivate != nil {
		def.Deactivate(param, conf, env)
	}
}
