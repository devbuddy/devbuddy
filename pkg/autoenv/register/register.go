package register

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/autoenv/feature"
)

type ImmutableRegister interface {
	Get(string) (*feature.Feature, error)
	Names() []string
}

type Register struct {
	nameToFeature map[string]*feature.Feature
}

func NewRegister() *Register {
	return &Register{nameToFeature: make(map[string]*feature.Feature)}
}

var globalRegister *Register

func Global() ImmutableRegister {
	return globalRegister
}

func RegisterFeature(name string, activate feature.ActivateFunc, deactivate feature.DeactivateFunc) {
	if globalRegister == nil {
		globalRegister = NewRegister()
	}
	globalRegister.Register(name, activate, deactivate)
}

func (e *Register) Register(name string, activate feature.ActivateFunc, deactivate feature.DeactivateFunc) {
	if _, ok := e.nameToFeature[name]; ok {
		panic(fmt.Sprint("Can't re-register a definition:", name))
	}
	if activate == nil {
		panic("activate can't be nil")
	}
	if deactivate == nil {
		panic("deactivate can't be nil")
	}

	e.nameToFeature[name] = &feature.Feature{Name: name, Activate: activate, Deactivate: deactivate}
}

func (e *Register) Get(name string) (*feature.Feature, error) {
	env := e.nameToFeature[name]
	if env == nil {
		return nil, fmt.Errorf("unknown feature: %s", name)
	}
	return env, nil
}

func (e *Register) Names() (names []string) {
	for name := range e.nameToFeature {
		names = append(names, string(name))
	}
	return
}
