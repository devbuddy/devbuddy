package features

import (
	"fmt"
)

type featureRegister struct {
	nameToFeature map[string]*Feature
}

func newFeatureRegister() *featureRegister {
	return &featureRegister{nameToFeature: make(map[string]*Feature)}
}

var globalRegister *featureRegister

func register(name string, activate activateFunc, deactivate deactivateFunc) {
	if globalRegister == nil {
		globalRegister = newFeatureRegister()
	}
	globalRegister.register(name, activate, deactivate)
}

func (e *featureRegister) register(name string, activate activateFunc, deactivate deactivateFunc) {
	if _, ok := e.nameToFeature[name]; ok {
		panic(fmt.Sprint("Can't re-register a definition:", name))
	}
	if activate == nil {
		panic("activate can't be nil")
	}
	if deactivate == nil {
		panic("deactivate can't be nil")
	}

	e.nameToFeature[name] = &Feature{Name: name, Activate: activate, Deactivate: deactivate}
}

func (e *featureRegister) get(name string) (*Feature, error) {
	env := e.nameToFeature[name]
	if env == nil {
		return nil, fmt.Errorf("unknown feature: %s", name)
	}
	return env, nil
}

func (e *featureRegister) names() (names []string) {
	for name := range e.nameToFeature {
		names = append(names, string(name))
	}
	return
}
