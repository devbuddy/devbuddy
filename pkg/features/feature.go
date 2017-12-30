package features

import (
	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

type Feature interface {
	Enable(*project.Project, *Env, *termui.UI)
	Disable(*project.Project, *Env, *termui.UI)
}

type FeatureBuilder func(param string) Feature

var FeatureMap = make(map[string]FeatureBuilder)
