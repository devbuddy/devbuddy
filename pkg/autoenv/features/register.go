package features

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/autoenv/feature"
	"github.com/devbuddy/devbuddy/pkg/context"
)

type MutableRegister struct {
	nameToFeature map[string]*feature.Feature
}

func NewRegister() *MutableRegister {
	return &MutableRegister{nameToFeature: make(map[string]*feature.Feature)}
}

type Register interface {
	Get(string) (*feature.Feature, error)
	Names() []string
}

var globalRegister *MutableRegister

func GlobalRegister() Register {
	return globalRegister
}

func register(name string, activate feature.ActivateFunc, deactivate feature.DeactivateFunc) {
	if globalRegister == nil {
		globalRegister = NewRegister()
	}
	globalRegister.Register(name, activate, deactivate)
}

func (e *MutableRegister) Register(name string, activate feature.ActivateFunc, deactivate feature.DeactivateFunc) {
	if _, ok := e.nameToFeature[name]; ok {
		panic(fmt.Sprint("Can't re-register a definition:", name))
	}
	if activate == nil {
		panic("activate can't be nil")
	}
	if deactivate == nil {
		deactivate = func(ctx *context.Context, version string) {}
	}

	e.nameToFeature[name] = &feature.Feature{Name: name, Activate: activate, Deactivate: deactivate}
}

func (e *MutableRegister) Get(name string) (*feature.Feature, error) {
	env := e.nameToFeature[name]
	if env == nil {
		return nil, fmt.Errorf("unknown feature: %s", name)
	}
	return env, nil
}

func (e *MutableRegister) Names() (names []string) {
	for name := range e.nameToFeature {
		names = append(names, string(name))
	}
	return
}
