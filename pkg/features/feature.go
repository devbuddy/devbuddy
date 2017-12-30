package features

import (
	"errors"

	"github.com/pior/dad/pkg/project"
	"github.com/pior/dad/pkg/termui"
)

var DevUpNeeded error

func init() {
	DevUpNeeded = errors.New("dev up needed")
}

type Feature interface {
	Enable(*project.Project, *Env, *termui.UI) error
	Disable(*project.Project, *Env, *termui.UI)
}

type FeatureBuilder func(param string) Feature

var FeatureMap = make(map[string]FeatureBuilder)
