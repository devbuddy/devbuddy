package features

import (
	"errors"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
)

var DevUpNeeded error

func init() {
	DevUpNeeded = errors.New("dev up needed")
}

type Feature interface {
	activate(*config.Config, *project.Project, *Env) error
	deactivate(*config.Config, *Env)
}

type featureBuilder func(param string) Feature

var allFeatures = make(map[string]featureBuilder)
