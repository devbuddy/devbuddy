package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
	"github.com/devbuddy/devbuddy/pkg/termui"
)

// Get returns a feature selected by the name argument.
func Get(name string) (*Feature, error) {
	return globalRegister.get(name)
}

type activateFunc func(string, *config.Config, *project.Project, *env.Env) (bool, error)
type deactivateFunc func(string, *config.Config, *env.Env)

// Feature is the implementation of an environment feature.
type Feature struct {
	Name       string
	Activate   activateFunc
	Deactivate deactivateFunc
}

// Sync activates / deactivates the features in the instance of env.Env.
// When a feature is already active but unknown, it will be ignored completely.
// When a param changes, the feature is deactivated with the current param then activated with the new param.
func Sync(cfg *config.Config, proj *project.Project, ui *termui.UI, env *env.Env, features map[string]string) {
	runner := &runner{cfg: cfg, proj: proj, ui: ui, env: env, reg: globalRegister}
	runner.sync(features)
}
