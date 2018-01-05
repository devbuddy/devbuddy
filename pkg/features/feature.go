package features

import (
	"errors"

	"github.com/pior/dad/pkg/config"
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

var DevUpNeeded error

func init() {
	DevUpNeeded = errors.New("dev up needed")
}

type Feature interface {
	Enable(*config.Config, *project.Project, *Env, *termui.HookUI) error
	Disable(*config.Config, *project.Project, *Env, *termui.HookUI)
}

type FeatureBuilder func(param string) Feature

var allFeatures = make(map[string]FeatureBuilder)
