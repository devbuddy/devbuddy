package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

// Get returns a feature selected by the name argument
func Get(name string) (*Feature, error) {
	return globalRegister.get(name)
}

type activateFunc func(string, *config.Config, *project.Project, *env.Env) (bool, error)
type deactivateFunc func(string, *config.Config, *env.Env)

type Feature struct {
	Name       string
	Activate   activateFunc
	Deactivate deactivateFunc
}
