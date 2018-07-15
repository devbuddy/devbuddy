package features

import (
	"errors"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/features/register"
	"github.com/devbuddy/devbuddy/pkg/project"
)

var DevUpNeeded error
var definitions *register.Definitions

func init() {
	DevUpNeeded = errors.New("dev up needed")
	definitions = register.New()
}

func Activate(name string, param string, conf *config.Config, proj *project.Project, env *env.Env) error {
	def := definitions.Get(name)
	if def.Activate == nil {
		panic("123")
	}
	return def.Activate(param, conf, proj, env)
}
