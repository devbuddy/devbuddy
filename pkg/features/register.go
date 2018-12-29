package features

import (
	"fmt"

	"github.com/devbuddy/devbuddy/pkg/config"
	"github.com/devbuddy/devbuddy/pkg/env"
	"github.com/devbuddy/devbuddy/pkg/project"
)

type activateFunc func(string, *config.Config, *project.Project, *env.Env) (bool, error)
type deactivateFunc func(string, *config.Config, *env.Env)

type Feature struct {
	Name       string
	Activate   activateFunc
	Deactivate deactivateFunc
}

type featureRegister struct {
	nameToEnv map[string]*Feature
}

func newFeatureRegister() *featureRegister {
	return &featureRegister{nameToEnv: make(map[string]*Feature)}
}

var globalRegister *featureRegister

func register(name string, activate activateFunc, deactivate deactivateFunc) {
	if globalRegister == nil {
		globalRegister = newFeatureRegister()
	}
	globalRegister.register(name, activate, deactivate)
}

func Get(name string) (*Feature, error) {
	return globalRegister.get(name)
}

func (e *featureRegister) register(name string, activate activateFunc, deactivate deactivateFunc) {
	if _, ok := e.nameToEnv[name]; ok {
		panic(fmt.Sprint("Can't re-register a definition:", name))
	}
	if activate == nil {
		panic("activate can't be nil")
	}
	if deactivate == nil {
		panic("deactivate can't be nil")
	}

	e.nameToEnv[name] = &Feature{Name: name, Activate: activate, Deactivate: deactivate}
}

func (e *featureRegister) get(name string) (*Feature, error) {
	env := e.nameToEnv[name]
	if env == nil {
		return nil, fmt.Errorf("unknown feature: %s", name)
	}
	return env, nil
}

func (e *featureRegister) names() (names []string) {
	for name := range e.nameToEnv {
		names = append(names, string(name))
	}
	return
}
