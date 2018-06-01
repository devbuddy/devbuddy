package features

import (
	"errors"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

var DevUpNeeded error

func init() {
	DevUpNeeded = errors.New("dev up needed")
}

type Feature interface {
	Activate(*config.Config, *project.Project, *env.Env) error
	Deactivate(*config.Config, *env.Env)
}

type featureBuilder func(param string) Feature

var allFeatures = make(map[string]featureBuilder)

func New(name string, param string) Feature {
	return allFeatures[name](param)
}
