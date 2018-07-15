package features

import (
	"errors"
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

var DevUpNeeded error
var definitions *register

func init() {
	DevUpNeeded = errors.New("dev up needed")
	definitions = newRegister()
}

func Activate(name string, param string, conf *config.Config, proj *project.Project, env *env.Env) error {
	def := definitions.Get(name)
	if def == nil {
		panic(fmt.Sprintf("unknown feature: %s", name))
	}
	if def.activate == nil {
		panic(fmt.Sprintf("no activate method for feature %s", name))
	}
	return def.activate(param, conf, proj, env)
}

func Deactivate(name string, param string, conf *config.Config, env *env.Env) {
	def := definitions.Get(name)
	if def == nil {
		panic(fmt.Sprintf("unknown feature: %s", name))
	}
	if def.deactivate == nil {
		panic(fmt.Sprintf("no deactivate method for feature %s", name))
	}
	def.deactivate(param, conf, env)
}
