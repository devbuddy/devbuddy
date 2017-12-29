package features

import (
	"github.com/pior/dad/pkg/termui"
)

type Feature interface {
	Enable(*Env, *termui.UI)
	Disable(*Env, *termui.UI)
}

type FeatureBuilder func(param string) Feature

var FeatureMap = make(map[string]FeatureBuilder)
