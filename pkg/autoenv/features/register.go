package features

import (
	"fmt"
)

type Register struct {
	nameToFeature map[string]*Feature
}

func NewRegister() *Register {
	return &Register{nameToFeature: make(map[string]*Feature)}
}

var globalRegister *Register

func GetRegister() *Register {
	return globalRegister
}

func register(name string, activate activateFunc, deactivate deactivateFunc) {
	if globalRegister == nil {
		globalRegister = NewRegister()
	}
	globalRegister.Register(name, activate, deactivate)
}

func (e *Register) Register(name string, activate activateFunc, deactivate deactivateFunc) {
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

func (e *Register) Get(name string) (*Feature, error) {
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
