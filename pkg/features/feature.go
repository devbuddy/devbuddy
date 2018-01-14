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
	Enable(*config.Config, *project.Project, *Env) error
	Disable(*config.Config, *Env)
}

type FeatureBuilder func(param string) Feature

var allFeatures = make(map[string]FeatureBuilder)
