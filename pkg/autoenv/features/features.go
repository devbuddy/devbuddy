package features

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type activateFunc func(string, *config.Config, *project.Project, *env.Env) (bool, error)
type deactivateFunc func(string, *config.Config, *env.Env)

// Feature is the implementation of an environment feature.
type Feature struct {
	Name       string
	Activate   activateFunc
	Deactivate deactivateFunc
}
