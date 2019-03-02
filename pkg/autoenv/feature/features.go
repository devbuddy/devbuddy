package feature

import (
	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type ActivateFunc func(string, *config.Config, *project.Project, *env.Env) (bool, error)
type DeactivateFunc func(string, *config.Config, *env.Env)

// Feature is the implementation of an environment feature.
type Feature struct {
	Name       string
	Activate   ActivateFunc
	Deactivate DeactivateFunc
}
